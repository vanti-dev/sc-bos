package scale

import (
	"time"
)

// TimeFunc returns a scale factor based on a time.
// The scale return value will be in the range [0,1].
// If t is before the period represented by this scale func, tense will be -1, 0 when during, and 1 when t is after.
// A tense of 0 should be returned if the func applies to the entire timeline.
type TimeFunc func(t time.Time) (scale float64, tense int)

// Time scales based on time.Time values.
// The first entry that returns a non-positive tense return value will have its value returned.
type Time []TimeFunc

// Now returns a scale factor between [0, 1] for the current time.
func (s Time) Now() float64 {
	return s.At(time.Now())
}

// At returns a scale factor between [0, 1] for the given time.
func (s Time) At(t time.Time) float64 {
	var v float64
	for _, timeFunc := range s {
		var cmp int
		v, cmp = timeFunc(t)
		if cmp <= 0 {
			return v
		}
	}
	return v
}

// NineToFive is a Time scaler that reports higher values during working days/hours.
var NineToFive = Time{
	weekendRamp(0.1, 1),            // no work at weekends, must be first
	linearRampHour(8, 10, 0.1, 1),  // start of day
	linearRampHour(12, 13, 1, 0.5), // start of lunch
	linearRampHour(13, 14, 0.5, 1), // end of lunch
	linearRampHour(16, 18, 1, 0.1), // end of day
}

func weekendRamp(min, max float64) TimeFunc {
	return func(t time.Time) (float64, int) {
		if isWeekend(t) {
			return min, 0
		}
		return max, 1 // after the weekend means we'll check the next TimeFunc
	}
}

func linearRampHour(from, to, min, max float64) TimeFunc {
	return func(t time.Time) (float64, int) {
		hour := float64(t.Hour()) + float64(t.Minute())/60.0
		if hour < from {
			return min, -1
		}
		if hour > to {
			return max, 1
		}
		return mapRange(hour, from, to, min, max), 0
	}
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// mapRange maps a value from one range to another
func mapRange[T float32 | float64](value, fromLow, fromHigh, toLow, toHigh T) T {
	return (value-fromLow)/(fromHigh-fromLow)*(toHigh-toLow) + toLow
}
