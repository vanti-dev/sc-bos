package config

import (
	"encoding/json"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

// Root represent the configuration parameters available for the lighting automation.
// This should be convertable to/from json.
type Root struct {
	auto.Config

	OccupancySensors  []string `json:"occupancySensors,omitempty"`
	Lights            []string `json:"lights,omitempty"`
	BrightnessSensors []string `json:"brightnessSensors,omitempty"`

	// UnoccupiedOffDelay configures how long we wait after the most recent occupancy sensor reported unoccupied before
	// we turn the light off.
	UnoccupiedOffDelay jsontypes.Duration `json:"unoccupiedOffDelay,omitempty"`
	// DaylightDimming configures how the brightness measured in the space affects the luminosity of the lights that
	// are on.
	DaylightDimming *DaylightDimming `json:"daylightDimming,omitempty"`

	// Devices implementing the Button trait will be used to switch the lights on and off when clicked once.
	OnButtons     []string `json:"onButtons,omitempty"`
	OffButtons    []string `json:"offButtons,omitempty"`
	ToggleButtons []string `json:"toggleButtons, omitempty"`

	// Now returns the current time. It's configurable for testing purposes, typically for testing the logic.
	Now func() time.Time `json:"-"`
}

type DaylightDimming struct {
	// Thresholds configures a mapping between measured lux levels and output brightness of lights.
	// With Thresholds you can say "below 300 lux set brightness to 80%, below 700 lux set to 50%".
	// The threshold with the highest BelowLux value below the measured lux level will be selected.
	Thresholds []LevelThreshold `json:"thresholds,omitempty"`
}

type LevelThreshold struct {
	BelowLux     float32 `json:"belowLux,omitempty"`
	LevelPercent float32 `json:"levelPercent,omitempty"`
}

func Read(data []byte) (Root, error) {
	root := Default()
	err := json.Unmarshal(data, &root)
	return root, err
}

func Default() Root {
	return Root{
		Now:                time.Now,
		UnoccupiedOffDelay: jsontypes.Duration{Duration: 10 * time.Minute},
	}
}
