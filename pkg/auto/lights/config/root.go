package config

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/multierr"

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

	Mode // default mode
	// Modes describe modes of operation and when they should be active by default.
	Modes []ModeOption `json:"modes,omitempty"`

	ModeSource   string `json:"modeSource,omitempty"` // the device name to read the active mode from
	ModeValueKey string `json:"modeName,omitempty"`   // the name of the mode value in ModeSource that represents the active mode. Defaults to lighting.mode

	// Devices implementing the Button trait will be used to switch the lights on and off when clicked once.
	OnButtons     []string `json:"onButtons,omitempty"`
	OffButtons    []string `json:"offButtons,omitempty"`
	ToggleButtons []string `json:"toggleButtons,omitempty"`

	// Now returns the current time. It's configurable for testing purposes, typically for testing the logic.
	Now func() time.Time `json:"-"`
}

type DaylightDimming struct {
	// Segments configures generated Thresholds based on some parameters
	// if Thresholds has any values, this will be ignored
	Segments *ThresholdSegments `json:"segments,omitempty"`
	// Thresholds configures a mapping between measured lux levels and output brightness of lights.
	// With Thresholds you can say "below 300 lux set brightness to 80%, below 700 lux set to 50%".
	// The threshold with the highest BelowLux value below the measured lux level will be selected.
	Thresholds []LevelThreshold `json:"thresholds,omitempty"`
	// PercentageTowardsGoal configures how quickly we reach our goal. If set to 50 then we calculate the desired level
	// from the lookup table and then go half way between current and desired.
	// A new lux reading will come in based on that light level; we will approach the goal, always undershooting.
	PercentageTowardsGoal float32 `json:"percentageTowardsGoal,omitempty"`
}

// process will turn Segments into Thresholds, if configured
func (d *DaylightDimming) process() error {
	if d == nil {
		return nil // nothing to do
	}
	if len(d.Thresholds) > 0 {
		return nil // nothing to do
	}
	if d.Segments == nil {
		return nil // nothing to do
	}
	seg := d.Segments
	if seg.MinLevel < 0 {
		return fmt.Errorf("invalid daylight dimming minLevel %d, expected >0", seg.MinLevel)
	}
	if seg.MaxLevel > 100 {
		return fmt.Errorf("invalid daylight dimming maxLevel %d, expected <100", seg.MaxLevel)
	}
	if seg.MinLux > seg.MaxLux || seg.MinLux == seg.MaxLux {
		return fmt.Errorf("invalid daylight dimming minLux %d maxLux %d, expected minLux < maxLux", seg.MinLux, seg.MaxLux)
	}
	if seg.MaxLevel == 0 {
		seg.MaxLevel = 100
	}
	if seg.Steps == 0 {
		seg.Steps = 100
	}
	// negative lux step, we go from max to min
	luxStep := (seg.MinLux - seg.MaxLux) / seg.Steps
	levelStep := (seg.MaxLevel - seg.MinLevel) / (seg.Steps - 1)
	for i := 0; i < seg.Steps; i++ {
		d.Thresholds = append(d.Thresholds, LevelThreshold{
			BelowLux:     float32(seg.MaxLux + i*luxStep),
			LevelPercent: float32(seg.MinLevel + i*levelStep),
		})
	}
	return nil
}

// ThresholdSegments will generate one LevelThreshold per step, with each threshold LevelPercent being evenly
// spread between MinLevel and MaxLevel, and each threshold BelowLux being evenly spread between MinLux and MaxLux
type ThresholdSegments struct {
	MinLux int `json:"minLux,omitempty"`
	MaxLux int `json:"maxLux,omitempty"`
	// defaults to 0!
	MinLevel int `json:"minLevel,omitempty"`
	// defaults to 100 (if 0), max 100
	MaxLevel int `json:"maxLevel,omitempty"`
	// defaults to 100 (if 0)
	Steps int `json:"steps,omitempty"`
}

type LevelThreshold struct {
	BelowLux     float32 `json:"belowLux,omitempty"`
	LevelPercent float32 `json:"levelPercent,omitempty"`
}

type Mode struct {
	// UnoccupiedOffDelay configures how long we wait after the most recent occupancy sensor reported unoccupied before
	// we turn the light off.
	UnoccupiedOffDelay jsontypes.Duration `json:"unoccupiedOffDelay,omitempty"`
	// DaylightDimming configures how the brightness measured in the space affects the luminosity of the lights that
	// are on.
	DaylightDimming *DaylightDimming `json:"daylightDimming,omitempty"`
	// Levels to use when the lights are on or off. If present overrides daylight dimming.
	OnLevelPercent  *float32 `json:"onLevelPercent,omitempty"`
	OffLevelPercent *float32 `json:"offLevelPercent,omitempty"`
}

type ModeOption struct {
	Name string `json:"name,omitempty"`
	Mode
	Start *Schedule `json:"start,omitempty"`
	End   *Schedule `json:"end,omitempty"`
}

func Read(data []byte) (Root, error) {
	root := Default()
	err := json.Unmarshal(data, &root)
	if err != nil {
		return root, err
	}
	var errs error
	if root.DaylightDimming != nil {
		// err is returned below
		errs = multierr.Append(errs, root.DaylightDimming.process())
	}
	for _, mode := range root.Modes {
		if mode.DaylightDimming != nil {
			errs = multierr.Append(errs, mode.DaylightDimming.process())
		}
	}
	root.Modes = applyModeDefaults(root.Mode, root.Modes)
	return root, err
}

func applyModeDefaults(defaults Mode, modes []ModeOption) []ModeOption {
	for i, mode := range modes {
		if mode.DaylightDimming == nil {
			mode.DaylightDimming = defaults.DaylightDimming
		}
		if mode.UnoccupiedOffDelay.Duration == 0 {
			mode.UnoccupiedOffDelay = defaults.UnoccupiedOffDelay
		}
		if mode.OnLevelPercent == nil {
			mode.OnLevelPercent = defaults.OnLevelPercent
		}
		if mode.OffLevelPercent == nil {
			mode.OffLevelPercent = defaults.OffLevelPercent
		}
		modes[i] = mode
	}
	return modes
}

func Default() Root {
	return Root{
		Now: time.Now,
		Mode: Mode{
			UnoccupiedOffDelay: jsontypes.Duration{Duration: 10 * time.Minute},
		},
	}
}
