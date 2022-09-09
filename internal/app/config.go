package app

import (
	"encoding/json"
)

type ControllerConfig struct {
	Drivers    []RawDriverConfig     `json:"drivers"`
	Automation []RawAutomationConfig `json:"automation"`
	Spaces     []RawSpaceConfig      `json:"spaces"`
}

type BaseDriverConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type RawDriverConfig struct {
	BaseDriverConfig
	Raw json.RawMessage `json:"-"`
}

func (c *RawDriverConfig) UnmarshalJSON(buf []byte) error {
	c.Raw = buf
	return json.Unmarshal(buf, &c.BaseDriverConfig)
}

type BaseAutomationConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type RawAutomationConfig struct {
	BaseAutomationConfig
	Raw json.RawMessage `json:"-"`
}

func (c *RawAutomationConfig) UnmarshalJSON(buf []byte) error {
	c.Raw = buf
	return json.Unmarshal(buf, &c.BaseAutomationConfig)
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
