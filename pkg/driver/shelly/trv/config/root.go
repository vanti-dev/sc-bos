package config

import (
	"time"

	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type Root struct {
	driver.BaseConfig

	Devices []TRVConfig `json:"devices,omitempty"`
}

type TRVConfig struct {
	Name         string        `json:"name"`
	Address      string        `json:"address"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	PollInterval time.Duration `json:"poll-interval"`
}
