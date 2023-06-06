package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	Electrics      []string            `json:"electrics,omitempty"`
	ElectricGroups map[string][]string `json:"electricGroups,omitempty"`
}
