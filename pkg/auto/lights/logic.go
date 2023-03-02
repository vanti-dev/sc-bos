package lights

import (
	"context"
	"sort"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/auto/lights/config"

	"github.com/smart-core-os/sc-api/go/traits"
)

// processState executes clientActions based on both read and write states.
// Here is where the logic that says "when PIRs report occupied, turn lights on" lives.
//
// Returning a non-zero duration indicates that processing should be rerun after this delay even if ReadState doesn't
// change.
func processState(ctx context.Context, readState *ReadState, writeState *WriteState, actions actions) (time.Duration, error) {
	var rerunAfter time.Duration

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

	// check if the buttons have been used to override the state recently enough
	var now time.Time
	if readState.Config.Now != nil {
		now = readState.Config.Now()
	} else {
		now = time.Now()
	}
	var (
		isForcedOn  bool
		isForcedOff bool
	)
	if readState.Force != nil {
		sinceButtonPress := now.Sub(readState.Force.Time)
		configuredTimeout := readState.Config.UnoccupiedOffDelay.Duration
		if sinceButtonPress < configuredTimeout {
			isForcedOn = readState.Force.On
			isForcedOff = !readState.Force.On
			rerunAfter = configuredTimeout - sinceButtonPress
		}
	}

	if isForcedOff {
		return rerunAfter, updateBrightnessLevelIfNeeded(ctx, writeState, actions, 0, readState.Config.Lights...)
	}

	// We can do easy checks for occupancy and turn things on if they are occupied
	if anyOccupied || isForcedOn {
		level, ok := computeOnLevelPercent(readState)
		if !ok {
			// todo: here we are in a position where daylight dimming is supposed to be enabled but we don't have enough
			//  info to actually choose the output light level. We should probably not make any changes and wait for
			//  more data to come in, but we'll leave that to future us as part of snagging.
		}
		return rerunAfter, updateBrightnessLevelIfNeeded(ctx, writeState, actions, level, readState.Config.Lights...)
	}

	// We can also delay changes if we need to.
	// This code check when occupancy last reported unoccupied and only turns the lights off
	// if it's been unoccupied for more than 10 minutes.
	// If it hasn't been 10 minutes yet, it waits some time and turns the lights off when it has been
	// 10 minutes.
	var occupancyExpired bool
	if len(readState.Config.OccupancySensors) == 0 {
		// if no occupancy sensors are configured, always consider the space unoccupied.
		occupancyExpired = true
	} else if becameUnoccupied := lastUnoccupiedTime(readState); !becameUnoccupied.IsZero() {
		unoccupiedDelayBeforeDarkness := readState.Config.UnoccupiedOffDelay.Duration

		sinceUnoccupied := now.Sub(becameUnoccupied)
		if sinceUnoccupied >= unoccupiedDelayBeforeDarkness {
			// we've been unoccupied for long enough, turn things off now
			occupancyExpired = true
		} else if !isForcedOff {
			// if the lights are forced off, no point waking up before that expires, as nothing will change before then.

			// we haven't written anything, but in `unoccupiedDelayBeforeDarkness - sinceUnoccupied` time we will, let the
			// caller know
			rerunAfter = unoccupiedDelayBeforeDarkness - sinceUnoccupied
		}
	}

	if occupancyExpired {
		return rerunAfter, updateBrightnessLevelIfNeeded(ctx, writeState, actions, 0, readState.Config.Lights...)
	}

	// no change
	return rerunAfter, nil
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

func computeOnLevelPercent(readState *ReadState) (level float32, ok bool) {
	dd := readState.Config.DaylightDimming
	if dd == nil {
		return 100, true
	}
	if len(readState.AmbientBrightness) == 0 {
		return 100, false
	}

	sensorLux := combinedLuxLevel(readState.AmbientBrightness)
	threshold, ok := closestThresholdBelow(sensorLux, dd.Thresholds)
	if !ok {
		return 100, true
	}

	return threshold.LevelPercent, true
}

func combinedLuxLevel(brightness map[string]*traits.AmbientBrightness) float32 {
	var n, v float32
	n, v = float32(len(brightness)), 0
	for _, ambientBrightness := range brightness {
		v += ambientBrightness.BrightnessLux / n
	}
	return v
}

func closestThresholdBelow(lux float32, thresholds []config.LevelThreshold) (config.LevelThreshold, bool) {
	if len(thresholds) == 0 {
		return config.LevelThreshold{}, false
	}

	// BelowLux 100 now comes before 400 in the slice
	sort.Slice(thresholds, func(i, j int) bool {
		return thresholds[i].BelowLux < thresholds[j].BelowLux
	})
	for _, threshold := range thresholds {
		if lux < threshold.BelowLux {
			return threshold, true
		}
	}
	if thresholds[0].BelowLux == 0 {
		return thresholds[0], true
	}
	return config.LevelThreshold{}, false
}
