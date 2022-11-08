package testlight

import (
	"encoding/json"
	"time"
)

type Config struct {
	Devices  []string `json:"devices"`
	Interval time.Duration
}

func DefaultConfig() Config {
	return Config{
		Interval: 2 * time.Minute,
	}
}

func DecodeConfig(configJSON []byte) (decoded Config, err error) {
	decoded = DefaultConfig()
	err = json.Unmarshal(configJSON, &decoded)
	return
}
