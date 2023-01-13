package config

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type Root struct {
	driver.BaseConfig
	Devices []Device `json:"devices,omitempty"`
}

type Device struct {
	*traits.Metadata
}
