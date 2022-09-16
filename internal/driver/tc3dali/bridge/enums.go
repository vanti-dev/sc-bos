package bridge

import "time"

// EventScheme is an enumeration of the combinations of metadata that DALI control devices can send with their events.
// Only two identifiers can be sent per event, and the EventScheme selects which ones will be present.
type EventScheme byte

const (
	// EventSchemeInstance - instance type & instance number
	EventSchemeInstance EventScheme = iota
	// EventSchemeDevice - device short address & instance type
	EventSchemeDevice
	// EventSchemeDeviceInstance - device short address & instance number
	EventSchemeDeviceInstance
	// EventSchemeDeviceGroup - device group & instance type
	EventSchemeDeviceGroup
	// EventSchemeInstanceGroup - instance group & instance type
	EventSchemeInstanceGroup
	EventSchemeUnknown EventScheme = 255
)

type FadeTime byte

const (
	FadeTimeDisabled FadeTime = iota
	FadeTime00707ms
	FadeTime01000ms
	FadeTime01400ms
	FadeTime02000ms
	FadeTime02800ms
	FadeTime04000ms
	FadeTime05700ms
	FadeTime08000ms
	FadeTime11300ms
	FadeTime16000ms
	FadeTime22600ms
	FadeTime32000ms
	FadeTime45300ms
	FadeTime64000ms
	FadeTime90500ms
	FadeTimeUnknown FadeTime = 255
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

const (
	InstanceTypeGeneric              uint8 = 0
	InstanceTypePushButton           uint8 = 1
	InstanceTypeAbsolute             uint8 = 2
	InstanceTypeOccupancy            uint8 = 3
	InstanceTypeLight                uint8 = 4
	InstanceTypeColour               uint8 = 5
	InstanceTypeGeneralPurposeSensor uint8 = 6
)
