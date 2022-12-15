package jsontypes

import (
	"encoding/json"
	"time"
)

// Duration wraps a time.Duration.
// It marshals / unmarshals to a string with the time.ParseDuration format.
type Duration struct {
	time.Duration
}

//goland:noinspection GoMixedReceiverTypes
func (d *Duration) UnmarshalJSON(raw []byte) error {
	var str string
	err := json.Unmarshal(raw, &str)
	if err != nil {
		return err
	}
	parsed, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
