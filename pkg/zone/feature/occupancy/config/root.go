package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	OccupancySensors []string `json:"occupancySensors,omitempty"`
	// These EnterLeave devices get converted to occupancy using `EnterTotal - LeaveTotal`
	EnterLeaveOccupancySensors []string `json:"enterLeaveOccupancySensors,omitempty"`
}
