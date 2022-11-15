package testlight

import (
	"encoding/json"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/util/jsontypes"
)

type Config struct {
	Devices       []string           `json:"devices"`
	PollInterval  jsontypes.Duration `json:"pollInterval"`  // The minimum interval from polling one light to another.
	CycleInterval jsontypes.Duration `json:"cycleInterval"` // How often to poll all lights
}

func DefaultConfig() Config {
	return Config{
		PollInterval:  jsontypes.Duration{Duration: 2 * time.Minute},
		CycleInterval: jsontypes.Duration{Duration: time.Hour},
	}
}

func DecodeConfig(configJSON []byte) (decoded Config, err error) {
	decoded = DefaultConfig()
	err = json.Unmarshal(configJSON, &decoded)
	return
}
