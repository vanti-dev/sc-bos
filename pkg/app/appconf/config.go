package appconf

import (
	"encoding/json"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type Config struct {
	Name       string             `json:"name,omitempty"`
	Drivers    []driver.RawConfig `json:"drivers,omitempty"`
	Automation []auto.RawConfig   `json:"automation,omitempty"`
	Spaces     []RawSpaceConfig   `json:"spaces,omitempty"`
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
