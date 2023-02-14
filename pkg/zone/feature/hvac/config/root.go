package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	ReadOnlyThermostats bool     `json:"readOnlyThermostats,omitempty"`
	Thermostats         []string `json:"thermostats,omitempty"`
}
