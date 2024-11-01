package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	OpenClose       []string            `json:"openClose,omitempty"`
	OpenCloseGroups map[string][]string `json:"openCloseGroups,omitempty"`
}
