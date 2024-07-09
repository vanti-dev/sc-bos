package config

import (
	"encoding/json"
	"errors"
	"time"
)

type Duration struct {
	// We embed to give config.Duration all the methods of time.Duration. Mostly String().

	time.Duration
}

//goland:noinspection GoMixedReceiverTypes
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

//goland:noinspection GoMixedReceiverTypes
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
