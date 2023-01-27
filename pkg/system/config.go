package system

import (
	"encoding/json"
)

type Config struct {
	Disabled bool `json:"disabled,omitempty"`
}

type RawConfig struct {
	Config
	Raw json.RawMessage `json:"-"`
}

func (c *RawConfig) UnmarshalJSON(buf []byte) error {
	c.Raw = buf
	return json.Unmarshal(buf, &c.Config)
}
