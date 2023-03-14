package lights

import (
	"context"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/lights/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"

	"github.com/smart-core-os/sc-api/go/traits"
)

// processState executes clientActions based on both read and write states.
// Here is where the logic that says "when PIRs report occupied, turn lights on" lives.
//
// Returning a non-zero duration indicates that processing should be rerun after this delay even if ReadState doesn't
// change.
func processState(ctx context.Context, readState *ReadState, writeState *WriteState, actions actions, logger *zap.Logger) (time.Duration, error) {
	var rerunAfter time.Duration

	// Work out what we need to do to apply the given writeState and make those changes for as long as ctx is valid

	var now time.Time
	if readState.Config.Now != nil {
		now = readState.Config.Now()
	} else {
		now = time.Now()
	}

	var (
		isSwitchedOn  bool
		isSwitchedOff bool
	)

	mostRecentButtonName := getMostRecentButtonPress(readState)

	var buttonActionRequired bool
	if len(mostRecentButtonName) == 0 {
		buttonActionRequired = false
	} else {
		buttonActionRequired = isButtonActionRequired(readState.Buttons[mostRecentButtonName], writeState)
	}

	if buttonActionRequired {
		// Either isSwitchedOn or isSwitchedOff should be set
		buttonType := getButtonType(readState, mostRecentButtonName)
		switch buttonType {
		case OnButton:
			isSwitchedOn = true
			rerunAfter = readState.Config.UnoccupiedOffDelay.Duration
		case OffButton:
			isSwitchedOff = true
		case ToggleButton:
			if getNewToggleState(writeState) {
				isSwitchedOn = true
			} else {
				isSwitchedOff = true
			}
			rerunAfter = readState.Config.UnoccupiedOffDelay.Duration
		}
		// Update the last time a button action happened
		writeState.LastButtonAction = readState.Buttons[mostRecentButtonName].StateChangeTime.AsTime()
	}

	logger.Debug("Is switched on ", zap.Bool("isSwitchedOn", isSwitchedOn))
	logger.Debug("Is switched off ", zap.Bool("isSwitchedOff", isSwitchedOff))

	if isSwitchedOff {
		return rerunAfter, updateBrightnessLevelIfNeeded(ctx, writeState, actions, 0, readState.Config.Lights...)
	}

	anyOccupied := areAnyOccupied(readState.Config.OccupancySensors, readState.Occupancy)

	// We can do easy checks for occupancy and turn things on if they are occupied
	if anyOccupied || isSwitchedOn {
		level, ok := computeOnLevelPercent(readState)
		if !ok {
			logger.Warn("Could not get level for daylight dimming")
			// todo: here we are in a position where daylight dimming is supposed to be enabled but we don't have enough
			//  info to actually choose the output light level. We should probably not make any changes and wait for
			//  more data to come in, but we'll leave that to future us as part of snagging.
		}
		return rerunAfter, updateBrightnessLevelIfNeeded(ctx, writeState, actions, level, readState.Config.Lights...)
	}

	// This code check when occupancy last reported unoccupied and only turns the lights off
	// if it's been unoccupied for more than unoccupied timeout.
	// If it hasn't been 10 minutes yet, it waits some time and turns the lights off when it has been
	// greater than the unoccupied timeout.
	// If a push button hasn't been pressed for the timeout period lights will be switched off too
	var occupancyExpired bool

	mostRecentButtonTime := time.Time{}
	if len(mostRecentButtonName) > 0 {
		mostRecentButtonTime = readState.Buttons[mostRecentButtonName].StateChangeTime.AsTime()
	}
	becameUnoccupied := lastUnoccupiedTime(readState)

	var sinceUnoccupied time.Duration

	if mostRecentButtonTime.After(becameUnoccupied) {
		sinceUnoccupied = now.Sub(mostRecentButtonTime)
	} else {
		sinceUnoccupied = now.Sub(becameUnoccupied)
	}

	if becameUnoccupied.IsZero() && mostRecentButtonTime.IsZero() {
		logger.Debug("Both time last unoccupied and last button press are zero.")
	} else {
		unoccupiedDelayBeforeDarkness := readState.Config.UnoccupiedOffDelay.Duration

		if sinceUnoccupied >= unoccupiedDelayBeforeDarkness {
			// we've been unoccupied for long enough, turn things off now
			occupancyExpired = true
		} else {
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

// getNewToggleState returns the new toggle state based on the last light brightness write
// if any light is on then we switch off, if all lights are off we switch on
func getNewToggleState(state *WriteState) bool {
	for _, brightness := range state.Brightness {
		if brightness.LevelPercent > 0 {
			return false
		}
	}
	return true
}

// getButtonType returns the type of button based on where it appeared in the config
func getButtonType(state *ReadState, buttonName string) ButtonType {
	nFound := 0
	buttonType := UndefinedButton
	for _, name := range state.Config.OnButtons {
		if name == buttonName {
			nFound++
			buttonType = OnButton
		}
	}
	for _, name := range state.Config.OffButtons {
		if name == buttonName {
			nFound++
			buttonType = OffButton
		}
	}
	for _, name := range state.Config.ToggleButtons {
		if name == buttonName {
			nFound++
			buttonType = ToggleButton
		}
	}
	// Todo: Add some logging if nButtons != 1
	return buttonType
}

func getMostRecentButtonPress(readState *ReadState) string {
	mostRecentTime := time.Time{}
	mostRecentName := ""
	for name, button := range readState.Buttons {
		if button.StateChangeTime.AsTime().After(mostRecentTime) {
			mostRecentTime = button.StateChangeTime.AsTime()
			mostRecentName = name
		}
	}
	return mostRecentName
}

// isButtonActionRequired returns true if state is unpressed and change time is more recent than last button action
func isButtonActionRequired(button *gen.ButtonState, writeState *WriteState) bool {
	if button.State == gen.ButtonState_UNPRESSED && button.StateChangeTime.AsTime().After(writeState.LastButtonAction) {
		return true
	}
	return false
}

// areAnyOccupied returns true if any occupancy sensors in the list are occupied
func areAnyOccupied(sensorsPresent []string, occupancyStates map[string]*traits.Occupancy) bool {
	var ret = false
	for _, name := range sensorsPresent {
		if o, ok := occupancyStates[name]; ok {
			if o.State == traits.Occupancy_OCCUPIED {
				ret = true
				break
			}
		}
	}
	return ret
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

type ButtonType int

const (
	UndefinedButton ButtonType = iota
	OnButton
	OffButton
	ToggleButton
)

func (s ButtonType) String() string {
	return [...]string{"Undefined button", "On button", "Off button", "Toggle button"}[s]
}
