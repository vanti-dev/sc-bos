package times

import "time"

// StartOfDay returns t with zero hour, min, sec, and nsec components.
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// NextWeekday returns a time.Time at local time t but on the next weekday where test returns true.
// If test never returns true then a zero time will be returned.
func NextWeekday(t time.Time, test func(wd time.Weekday) bool) time.Time {
	day := t.Weekday()
	for i := day + 1; i < day+7; i++ {
		if test(i % time.Saturday) {
			if day == i {
				return t
			}
			return t.AddDate(0, 0, int(i-day))
		}
	}

	return time.Time{}
}
