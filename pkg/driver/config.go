package driver

import (
	"encoding/json"
)

type BaseConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Disabled bool   `json:"disabled,omitempty"`
}

type RawConfig struct {
	BaseConfig
	Raw json.RawMessage `json:"-"`
}

func (c *RawConfig) UnmarshalJSON(buf []byte) error {
	c.Raw = buf
	return json.Unmarshal(buf, &c.BaseConfig)
}
