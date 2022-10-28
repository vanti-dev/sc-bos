package lights

import (
	"context"
	"github.com/smart-core-os/sc-api/go/traits"
	"time"
)

// processState executes clientActions based on both read and write states.
// Here is where the logic that says "when PIRs report occupied, turn lights on".
func processState(ctx context.Context, readState *ReadState, writeState *WriteState, actions actions) (time.Duration, error) {
	// Work out what we need to do to apply the given writeState and make those changes for as long as ctx is valid

	anyOccupied := false
	for _, name := range readState.Config.OccupancySensors {
		if o, ok := readState.Occupancy[name]; ok {
			if o.State == traits.Occupancy_OCCUPIED {
				anyOccupied = true
				break
			}
		}
	}

	// We can do easy checks for occupancy and turn things on if they are occupied
	if anyOccupied {
		return 0, updateBrightnessLevelIfNeeded(ctx, writeState, actions, 100, readState.Config.Lights...)
	}

	// We can also delay changes if we need to.
	// This code check when occupancy last reported unoccupied and only turns the lights off
	// if it's been unoccupied for more than 10 minutes.
	// If it hasn't been 10 minutes yet, it waits some time and turns the lights off when it has been
	// 10 minutes.
	lastOccupiedTime := lastUnoccupiedTime(readState)
	if !lastOccupiedTime.IsZero() {
		unoccupiedDelayBeforeDarkness := readState.Config.UnoccupiedOffDelay

		now := readState.Config.Now()
		sinceUnoccupied := now.Sub(lastOccupiedTime)
		if sinceUnoccupied >= unoccupiedDelayBeforeDarkness {
			// we've been unoccupied for long enough, turn things off now
			return 0, updateBrightnessLevelIfNeeded(ctx, writeState, actions, 0, readState.Config.Lights...)
		}

		// we haven't written anything, but in `unoccupiedDelayBeforeDarkness - sinceUnoccupied` time we will, let the
		// caller know
		return unoccupiedDelayBeforeDarkness - sinceUnoccupied, nil
	}

	// no change
	return 0, nil
}

// lastUnoccupiedTime returns the most recent Occupancy.StateChangeTime across each Config.OccupancySensors that have an
// unoccupied state.
func lastUnoccupiedTime(state *ReadState) time.Time {
	var mostRecentUnoccupiedTime time.Time
	for _, name := range state.Config.OccupancySensors {
		o, ok := state.Occupancy[name]
		if !ok {
			continue
		}

		if o.State == traits.Occupancy_UNOCCUPIED {
			if o.StateChangeTime == nil {
				continue
			}
			candidate := o.StateChangeTime.AsTime()
			if candidate.After(mostRecentUnoccupiedTime) {
				mostRecentUnoccupiedTime = candidate
			}
		}
	}
	return mostRecentUnoccupiedTime
}
