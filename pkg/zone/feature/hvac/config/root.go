package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	Thermostats []string `json:"thermostats,omitempty"`
}
