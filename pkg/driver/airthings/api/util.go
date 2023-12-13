// Package api allows interaction with the AirThings API.
// Some files were generated and copied from openapi-generator.
// Additional utils and fixes were added manually.
package api

import (
	"bytes"
	"encoding/json"
	"time"
)

// This file contains manual utilities to fix some of the generated code.

// A wrapper for strict JSON decoding
func newStrictDecoder(data []byte) *json.Decoder {
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.DisallowUnknownFields()
	return dec
}

type Time struct {
	time.Time
}

const timeFormat = "2006-01-02T15:04:05"

func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	tt, err := time.Parse(timeFormat, s)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(timeFormat))
}
