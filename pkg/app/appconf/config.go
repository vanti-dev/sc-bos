package appconf

import (
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Config struct {
	Name       string             `json:"name,omitempty"`
	Drivers    []driver.RawConfig `json:"drivers,omitempty"`
	Automation []auto.RawConfig   `json:"automation,omitempty"`
	Zones      []zone.RawConfig   `json:"zones,omitempty"`
}
