package dali

import "time"

// EventScheme is an enumeration of the combinations of metadata that DALI control devices can send with their events.
// Only two identifiers can be sent per event, and the EventScheme selects which ones will be present.
// This enum exactly matches the TC3 enum documented at
// https://infosys.beckhoff.com/content/1033/tcplclib_tc3_dali/5643252875.html?id=5542922191822951473
type EventScheme byte

const (
	EventSchemeInstance       EventScheme = 0 // instance type & instance number
	EventSchemeDevice         EventScheme = 1 // device short address & instance type
	EventSchemeDeviceInstance EventScheme = 2 // device short address & instance number
	EventSchemeDeviceGroup    EventScheme = 3 // device group & instance type
	EventSchemeInstanceGroup  EventScheme = 4 // instance group & instance type
	EventSchemeUnknown        EventScheme = 255
)

// FadeTime is an enumeration of how long changing level will take.
// This enum exactly matches the TC3 equivalent documented at
// https://infosys.beckhoff.com/content/1033/tcplclib_tc3_dali/6430577547.html?id=7861596685241704773
type FadeTime byte

const (
	FadeTimeDisabled FadeTime = 0
	FadeTime00707ms  FadeTime = 1
	FadeTime01000ms  FadeTime = 2
	FadeTime01400ms  FadeTime = 3
	FadeTime02000ms  FadeTime = 4
	FadeTime02800ms  FadeTime = 5
	FadeTime04000ms  FadeTime = 6
	FadeTime05700ms  FadeTime = 7
	FadeTime08000ms  FadeTime = 8
	FadeTime11300ms  FadeTime = 9
	FadeTime16000ms  FadeTime = 10
	FadeTime22600ms  FadeTime = 11
	FadeTime32000ms  FadeTime = 12
	FadeTime45300ms  FadeTime = 13
	FadeTime64000ms  FadeTime = 14
	FadeTime90500ms  FadeTime = 15
	FadeTimeUnknown  FadeTime = 255
)

func DurationToFadeTime(d time.Duration) (fadeTime FadeTime, ok bool) {
	if d < 0 || d >= 90500*time.Millisecond {
		return 0, false
	}
	ok = true

	switch {
	case d <= 1000*time.Millisecond:
		fadeTime = FadeTime00707ms
	case d <= 1400*time.Millisecond:
		fadeTime = FadeTime01000ms
	case d <= 2000*time.Millisecond:
		fadeTime = FadeTime01400ms
	case d <= 2800*time.Millisecond:
		fadeTime = FadeTime02000ms
	case d <= 4000*time.Millisecond:
		fadeTime = FadeTime02800ms
	case d <= 5700*time.Millisecond:
		fadeTime = FadeTime04000ms
	case d <= 8000*time.Millisecond:
		fadeTime = FadeTime05700ms
	case d <= 11300*time.Millisecond:
		fadeTime = FadeTime08000ms
	case d <= 16000*time.Millisecond:
		fadeTime = FadeTime11300ms
	case d <= 22600*time.Millisecond:
		fadeTime = FadeTime16000ms
	case d <= 3200*time.Millisecond:
		fadeTime = FadeTime22600ms
	case d <= 45300*time.Millisecond:
		fadeTime = FadeTime32000ms
	case d <= 64000*time.Millisecond:
		fadeTime = FadeTime45300ms
	case d <= 90500*time.Millisecond:
		fadeTime = FadeTime64000ms
	default:
		fadeTime = FadeTime90500ms
	}
	return
}

func (ft FadeTime) AsDuration() (duration time.Duration, ok bool) {
	ok = true
	switch ft {
	case FadeTime00707ms:
		duration = 707 * time.Millisecond
	case FadeTime01000ms:
		duration = 1000 * time.Millisecond
	case FadeTime01400ms:
		duration = 1400 * time.Millisecond
	case FadeTime02000ms:
		duration = 2000 * time.Millisecond
	case FadeTime02800ms:
		duration = 2800 * time.Millisecond
	case FadeTime04000ms:
		duration = 4000 * time.Millisecond
	case FadeTime05700ms:
		duration = 5700 * time.Millisecond
	case FadeTime08000ms:
		duration = 8000 * time.Millisecond
	case FadeTime11300ms:
		duration = 11300 * time.Millisecond
	case FadeTime16000ms:
		duration = 16000 * time.Millisecond
	case FadeTime22600ms:
		duration = 22600 * time.Millisecond
	case FadeTime32000ms:
		duration = 32000 * time.Millisecond
	case FadeTime45300ms:
		duration = 45300 * time.Millisecond
	case FadeTime64000ms:
		duration = 64000 * time.Millisecond
	case FadeTime90500ms:
		duration = 90500 * time.Millisecond
	default:
		ok = false
	}
	return
}

type InstanceType byte

// The constants for the different types of control devices (input devices)
// See https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc3_dali/9185112587.html&id=5259714273337466114
// for a reference
const (
	InstanceTypeGeneric              uint8 = 0
	InstanceTypePushButton           uint8 = 1
	InstanceTypeAbsolute             uint8 = 2
	InstanceTypeOccupancy            uint8 = 3
	InstanceTypeLight                uint8 = 4
	InstanceTypeColour               uint8 = 5
	InstanceTypeGeneralPurposeSensor uint8 = 6
)
