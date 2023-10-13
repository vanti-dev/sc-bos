package config

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type Root struct {
	driver.BaseConfig

	// Smart core metadata associated with this device.
	Metadata *traits.Metadata `json:"metadata,omitempty"`

	IpAddress string `json:"ipAddress"`
	Password  string `json:"password, omitempty"`
}
