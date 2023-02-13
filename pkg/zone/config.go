package zone

import (
	"encoding/json"
)

type Config struct {
	// Name is the announced name of the zone, for example `"zones/Room2"`.
	Name string `json:"name"`
	// Type distinguishes between different types of zone. For example "meeting room" or "lobby".
	// Type is used to identify the controller used for the zone.
	Type string `json:"type"`
	// Disabled zones do nothing.
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
