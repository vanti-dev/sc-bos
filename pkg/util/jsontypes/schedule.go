package jsontypes

import (
	"encoding/json"

	"github.com/robfig/cron/v3"
)

// Schedule represent a cron formatted schedule suitable for including in json types.
type Schedule struct {
	cron.Schedule
	Raw string
}

// MustParseSchedule parses a cron formatted schedule and panics if it fails.
func MustParseSchedule(raw string) *Schedule {
	schedule, err := cron.ParseStandard(raw)
	if err != nil {
		panic(err)
	}
	return &Schedule{Schedule: schedule, Raw: raw}
}

func (s *Schedule) String() string {
	return s.Raw
}

func (s *Schedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Raw)
}

func (s *Schedule) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	schedule, err := cron.ParseStandard(str)
	if err != nil {
		return err
	}
	*s = Schedule{Schedule: schedule, Raw: str}
	return nil
}
