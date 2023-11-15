package config

import (
	"github.com/vanti-dev/sc-bos/pkg/auto"
)

type Root struct {
	auto.Config
	// Schedule defines when to reset the enter/leave counter.
	// Defaults to midnight every day.
	Schedule *Schedule `json:"schedule,omitempty"`
	// Devices is the name of each device to reset.
	Devices []string `json:"devices,omitempty"`
}
