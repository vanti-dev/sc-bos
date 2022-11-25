package messageports

import "time"

const TimePattern = "02/01/2006 15:04:05"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalText(data []byte) error {
	tt, err := time.Parse(TimePattern, string(data))
	if err != nil {
		return err
	}
	*t = Time{tt}
	return nil
}

func (t *Time) UnmarshalBinary(data []byte) error {
	return t.UnmarshalText(data)
}
