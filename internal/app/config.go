package app

import (
	"encoding/json"
	"github.com/vanti-dev/sc-bos/internal/auto"

	"github.com/vanti-dev/sc-bos/internal/driver"
)

type ControllerConfig struct {
	Name       string                     `json:"name"`
	Drivers    []driver.RawConfig         `json:"drivers"`
	Automation []auto.RawConfig           `json:"automation"`
	Spaces     []RawSpaceConfig           `json:"spaces"`
	Systems    map[string]json.RawMessage `json:"systems,omitempty"`
}

type BaseSpaceConfig struct {
	Name string `json:"name"`
}

type RawSpaceConfig struct {
	BaseSpaceConfig
	Raw json.RawMessage `json:"-"`
}

func (c *RawSpaceConfig) UnmarshalJSON(buf []byte) error {
	c.Raw = buf
	return json.Unmarshal(buf, &c.BaseSpaceConfig)
}
