package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	Lights      []string            `json:"lights,omitempty"`      // Announces as {zone}
	LightGroups map[string][]string `json:"lightGroups,omitempty"` // Announced as {zone}/lights/{key}
}
