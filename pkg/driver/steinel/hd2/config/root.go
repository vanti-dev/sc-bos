package config

import (
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type Root struct {
	driver.BaseConfig

	IpAddress string `json:"ipAddress"`
	Password  string `json:"password, omitempty"`
}
