package config

import (
	"encoding/json"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

func ReadBytes(data []byte) (cfg Root, err error) {
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	cfg.ModeSource.ApplyDefaults(DefaultModeSource)
	for i := range cfg.OccupancyModeTargets {
		cfg.OccupancyModeTargets[i].ApplyDefaults(DefaultOccupancyModeTarget)
	}
	for i := range cfg.DeadbandModeTargets {
		cfg.DeadbandModeTargets[i].ApplyDefaults(DefaultDeadbandModeTarget)
	}
	return
}

type Root struct {
	auto.Config

	DryRun              bool `json:"dryRun,omitempty"`              // Don't actually write to devices
	LogDeviceWrites     bool `json:"logDeviceWrites,omitempty"`     // Log each write to a device as they happen
	LogDuplicateChanges bool `json:"logDuplicateChanges,omitempty"` // Log even if all writes were blocked by caching
	LogTTLDelays        bool `json:"logTTLDelays,omitempty"`        // Log the TTL when no writes were performed
	LogReads            bool `json:"logReads,omitempty"`            // Log each read from a device as they happen. Warning this is noisy

	// How long should a past write affect future writes.
	WriteCacheExpiry *jsontypes.Duration `json:"writeCacheExpiry,omitempty"` // Defaults to no expiry.
	// If processing fails, how long to wait before trying again.
	WriteRetryDelay *jsontypes.Duration `json:"WriteRetryDelay,omitempty"` // Defaults to 1m.
	// Reprocess state at least this often.
	WriteEvery *jsontypes.Duration `json:"writeEvery,omitempty"` // Defaults to never.
	// How long do we allow for writes to be picked up as reads.
	// For devices that we both read from and write to, we cache the read state in the write state to avoid writing the same value back to the device.
	// We only cache the read value if it is newer than the write, this property configures what "newer" means.
	// The read must be at least WriteReadPropagation newer than the write to be cached.
	WriteReadPropagation *jsontypes.Duration `json:"writeReadPropagation,omitempty"` // Defaults to 5s.

	// The device that we read the automation mode from.
	// See DefaultModeSource for the default values.
	ModeSource SwitchMode `json:"modeSource,omitempty"`
	// How long after ModeSource changes away from "auto" do we change it back to "auto".
	// Defaults to 4h.
	ResetModeSourceDelay *jsontypes.Duration `json:"resetModeSourceDelay,omitempty"`
	AutoModeSetPoint     *float32            `json:"autoModeSetPoint,omitempty"` // Defaults to 21.0
	AutoThermostats      []string            `json:"autoThermostats,omitempty"`  // Thermostats that we control in auto mode.

	OccupancySensors     []string            `json:"occupancySensors,omitempty"`     // Sensors whose occupancy is linked with OccupancyModeTargets On mode.
	OccupancyModeTargets []SwitchMode        `json:"occupancyModeTargets,omitempty"` // Defaults: on=occupied, off=unoccupied
	UnoccupiedDelay      *jsontypes.Duration `json:"unoccupiedDelay,omitempty"`      // Defaults to 15m.

	DeadbandSchedule    []Range      `json:"deadbandSchedule,omitempty"`    // Periods of time when DeadbandModeTargets should be On
	DeadbandModeTargets []SwitchMode `json:"deadbandModeTargets,omitempty"` // Defaults: on=comfort, off=eco
}

type Range struct {
	Start Schedule `json:"start,omitempty"`
	End   Schedule `json:"end,omitempty"`
}

var (
	DefaultModeSource           = SwitchMode{Key: "hvac.mode", On: "auto", Off: "manual"}
	DefaultOccupancyModeTarget  = SwitchMode{Key: "occupancy", On: "occupied", Off: "unoccupied"}
	DefaultDeadbandModeTarget   = SwitchMode{Key: "deadband", On: "comfort", Off: "eco"}
	DefaultWriteCacheExpiry     = 0 * time.Second
	DefaultWriteRetryDelay      = time.Minute
	DefaultWriteEvery           = 0 * time.Second
	DefaultWriteReadPropagation = 5 * time.Second
	DefaultUnoccupiedDelay      = 15 * time.Minute
	DefaultResetModeSourceDelay = 4 * time.Hour
	DefaultAutoModeSetPoint     = float32(21.0)
)

// SwitchMode represents a mode option that we switch between two values.
type SwitchMode struct {
	Name string `json:"name,omitempty"` // name of the device that implements Mode
	Key  string `json:"key,omitempty"`  // mode key to set
	On   string `json:"on,omitempty"`   // value when on
	Off  string `json:"off,omitempty"`  // value when off
}

func (s *SwitchMode) ApplyDefaults(d SwitchMode) {
	if s.Key == "" {
		s.Key = d.Key
	}
	if s.On == "" {
		s.On = d.On
	}
	if s.Off == "" {
		s.Off = d.Off
	}
}

func PtrOr[T any](p *T, or T) T {
	if p == nil {
		return or
	}
	return *p
}
