package config

import (
	"encoding/json"
	"time"
)

// Root represent the configuration parameters available for the lighting automation.
// This should be convertable to/from json.
type Root struct {
	OccupancySensors []string `json:"occupancySensors,omitempty"`
	Lights           []string `json:"lights,omitempty"`

	// UnoccupiedOffDelay configures how long we wait after the most recent occupancy sensor reported unoccupied before
	// we turn the light off.
	UnoccupiedOffDelay time.Duration `json:"unoccupiedOffDelay,omitempty"`

	// Now returns the current time. It's configurable for testing purposes, typically for testing the logic.
	Now func() time.Time `json:"-"`
}

func Read(data []byte) (Root, error) {
	root := Default()
	err := json.Unmarshal(data, &root)
	return root, err
}

func Default() Root {
	return Root{
		Now:                time.Now,
		UnoccupiedOffDelay: 10 * time.Minute,
	}
}
