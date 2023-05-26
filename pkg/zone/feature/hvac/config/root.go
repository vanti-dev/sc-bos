package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	Thermostat                             // announced as {zone}
	ThermostatGroups map[string]Thermostat `json:"thermostatGroups,omitempty"` // announced as {zone}/{key}
}

type Thermostat struct {
	ReadOnlyThermostat bool     `json:"thermostatReadOnly,omitempty"`
	Thermostats        []string `json:"thermostats,omitempty"`
}
