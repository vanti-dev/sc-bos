package lights

import (
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// captureButtonActions consumes button clicks returning the intended actions.
// This will modify writeState to record the last time a button was pressed to avoid an old button press triggering multiple times.
func captureButtonActions(readState *ReadState, writeState *WriteState, logger *zap.Logger) (onButtonClicked bool, offButtonClicked bool) {
	mostRecentButtonName, mostRecentButtonState, buttonActionRequired := getMostRecentButtonPress(readState)

	if buttonActionRequired {
		buttonActionRequired = isButtonActionRequired(mostRecentButtonState, writeState)

		logger.Debug("Checking if button action required for button",
			zap.String("button", mostRecentButtonName),
			zap.Bool("action required", buttonActionRequired),
			zap.Time("state change time", mostRecentButtonState.StateChangeTime.AsTime()),
			zap.Time("last action time", writeState.LastButtonAction),
			zap.Stringer("button state", mostRecentButtonState.State),
			zap.Stringer("last gesture", mostRecentButtonState.GetMostRecentGesture().GetKind()),
		)
	}

	if buttonActionRequired {
		// Either onButtonClicked or offButtonClicked should be set
		buttonType := getButtonType(readState, mostRecentButtonName)
		switch buttonType {
		case OnButton:
			onButtonClicked = true
		case OffButton:
			offButtonClicked = true
		case ToggleButton:
			// decide if the toggle button should be treated as an on or off button based on updates we've written
			if brightnessAllOff(writeState) {
				onButtonClicked = true
			} else {
				offButtonClicked = true
			}
		}
		buttonClickTime := mostRecentButtonState.StateChangeTime.AsTime()
		logger.Debug("Button action required",
			zap.Stringer("Button type", buttonType),
			zap.Bool("onButtonClicked", onButtonClicked),
			zap.Bool("offButtonClicked", offButtonClicked),
			zap.Time("buttonClickTime", buttonClickTime),
			zap.Time("last button action", writeState.LastButtonAction),
		)

		// Update the last time a button action happened
		if onButtonClicked {
			writeState.LastButtonOnTime = buttonClickTime
		}
		writeState.LastButtonAction = buttonClickTime
	}
	return onButtonClicked, offButtonClicked
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

func getMostRecentButtonPress(readState *ReadState) (name string, state *gen.ButtonState, ok bool) {
	mostRecentTime := time.Time{}
	for n, button := range readState.Buttons {
		if button.StateChangeTime.AsTime().After(mostRecentTime) {
			mostRecentTime = button.StateChangeTime.AsTime()
			name = n
			state = button
		}
	}
	return name, state, !mostRecentTime.IsZero()
}

// isButtonActionRequired returns true if state is unpressed and change time is more recent than last button action
func isButtonActionRequired(button *gen.ButtonState, writeState *WriteState) bool {
	if button.State == gen.ButtonState_UNPRESSED && button.StateChangeTime.AsTime().After(writeState.LastButtonAction) {
		return true
	}
	return false
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
