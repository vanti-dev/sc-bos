package lights

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"

	"github.com/vanti-dev/sc-bos/pkg/auto/lights/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// ReadState models everything we have read from the system.
// For example if we PullBrightness, then the responses will be recoded here.
type ReadState struct {
	Config config.Root

	Occupancy map[string]*traits.Occupancy
	// used for daylight dimming
	AmbientBrightness map[string]*traits.AmbientBrightness
	Buttons           map[string]*gen.ButtonState
	Force             *ForceState
}

type ForceState struct {
	On   bool
	Time time.Time
}

func (s *ForceState) Clone() *ForceState {
	if s != nil {
		return &ForceState{
			On:   s.On,
			Time: s.Time,
		}
	} else {
		return nil
	}
}

func NewReadState() *ReadState {
	return &ReadState{
		Occupancy:         make(map[string]*traits.Occupancy),
		AmbientBrightness: make(map[string]*traits.AmbientBrightness),
		Buttons:           make(map[string]*gen.ButtonState),
	}
}

func (s *ReadState) Clone() *ReadState {
	clone := NewReadState()
	clone.Config = s.Config
	// assume values in the map are immutable!
	for name, val := range s.Occupancy {
		clone.Occupancy[name] = val
	}
	for name, val := range s.AmbientBrightness {
		clone.AmbientBrightness[name] = val
	}
	for name, val := range s.Buttons {
		clone.Buttons[name] = val
	}
	clone.Force = s.Force.Clone()
	return clone
}

// WriteState models the state of the system based on the changes we've made to it.
// For example if we UpdateBrightness, then the response to that call is recorded in this state.
type WriteState struct {
	Brightness map[string]*traits.Brightness
}

func NewWriteState() *WriteState {
	return &WriteState{
		Brightness: make(map[string]*traits.Brightness),
	}
}

func (s *WriteState) MergeFrom(other *WriteState) {
	for name, brightness := range other.Brightness {
		if brightness == nil {
			delete(s.Brightness, name)
		} else {
			s.Brightness[name] = brightness
		}
	}
}

// readStateChanges collates changes and emits *ReadState.
func readStateChanges(ctx context.Context, workingState *ReadState, changes <-chan Patcher, newStateAvailable chan<- *ReadState) error {
	var readyToNotify chan struct{}

	// applyChange updates workingState by applying change to it.
	// This signals that we have some new state ready for someone to process.
	applyChange := func(change Patcher) {
		change.Patch(workingState)

		// let the loop know we have something to broadcast
		if readyToNotify == nil {
			readyToNotify = make(chan struct{})
			close(readyToNotify) // notify without sending anything
		}
	}

	// The following code can be summarised as:
	//   Apply changes to the state, notify when the state has updated
	//
	// As with most things it's a little more complicated than that though.
	// We _really_ want to prioritise updating the state over notifying that the state has updated,
	// if the thing that is processing the state misses some state updates then this isn't a bad thing.
	//
	// To accomplish this we make sure that at each step that could block we're checking for state changes again,
	// in practice this means each select has a case that checks for changes.
	// Unfortunately this doesn't guarantee that changes will be processed before the state is processed because
	// Go doesn't provide any guarantees or mechanisms for prioritising select cases.
	// We try our best though by introducing a 'drain loop' that will empty the changes chan of all items before
	// we consider even looking to see if we should be notifying of state changes.
	for {
	drainChanges:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				applyChange(change)
			default:
				break drainChanges
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case change := <-changes:
			applyChange(change)
		case <-readyToNotify:
			clonedState := workingState.Clone()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case newStateAvailable <- clonedState:
				readyToNotify = nil
			case change := <-changes:
				applyChange(change)
			}
		}
	}
}
