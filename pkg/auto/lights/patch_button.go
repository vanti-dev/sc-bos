package lights

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ButtonPatches struct {
	name       string
	client     gen.ButtonApiClient
	logger     *zap.Logger
	isOnButton bool // true for on switch, false for off switch
}

func (p *ButtonPatches) Subscribe(ctx context.Context, changes chan<- Patcher) error {
	defer func() {
		changes <- clearButtonStatePatcher(p.name)
	}()
	return pull.Changes[Patcher](ctx, p, changes, pull.WithLogger(p.logger.Named("button")))
}

func (p *ButtonPatches) Pull(ctx context.Context, changes chan<- Patcher) error {
	stream, err := p.client.PullButtonState(ctx, &gen.PullButtonStateRequest{Name: p.name})
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}
		patcher := pullButtonStatePatcher{
			response:   res,
			isOnSwitch: p.isOnButton,
		}
		select {
		case changes <- patcher:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *ButtonPatches) Poll(ctx context.Context, changes chan<- Patcher) error {
	res, err := p.client.GetButtonState(ctx, &gen.GetButtonStateRequest{Name: p.name})
	if err != nil {
		return err
	}
	patcher := getButtonStatePatcher{
		name:        p.name,
		buttonState: res,
		isOnSwitch:  p.isOnButton,
	}
	select {
	case changes <- patcher:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type pullButtonStatePatcher struct {
	response   *gen.PullButtonStateResponse
	isOnSwitch bool
}

func (p pullButtonStatePatcher) Patch(state *ReadState) {
	for _, change := range p.response.Changes {
		updateButtonState(state, change.Name, p.isOnSwitch, change.ButtonState)
	}
}

type getButtonStatePatcher struct {
	name        string
	buttonState *gen.ButtonState
	isOnSwitch  bool
}

func (p getButtonStatePatcher) Patch(state *ReadState) {
	updateButtonState(state, p.name, p.isOnSwitch, p.buttonState)
}

type clearButtonStatePatcher string

func (name clearButtonStatePatcher) Patch(state *ReadState) {
	delete(state.Buttons, string(name))
}

func updateButtonState(state *ReadState, name string, isOnButton bool, newState *gen.ButtonState) {
	oldButtonState, ok := state.Buttons[name]
	// we only want to apply the force if the button gesture represents a new change we haven't seen before
	// this prevents picking up old button gestures when the automation starts
	if ok {
		if t, ok := isNewSingleClick(oldButtonState, newState); ok {
			state.Force = &ForceState{
				On:   isOnButton,
				Time: t,
			}
		}
	}
	state.Buttons[name] = newState
}

// does the button state contain a single click gesture that we haven't processed before?
func isNewSingleClick(oldState, newState *gen.ButtonState) (t time.Time, ok bool) {
	oldGesture := oldState.GetMostRecentGesture()
	newGesture := newState.GetMostRecentGesture()

	if newGesture == nil {
		return time.Time{}, false
	}
	hasNewID := oldGesture.GetId() != newGesture.GetId()
	isSingleClick := newGesture.Kind == gen.ButtonState_Gesture_CLICK && newGesture.Count == 1
	if hasNewID && isSingleClick {
		ok = true
		t = newGesture.GetEndTime().AsTime()
	}
	return
}
