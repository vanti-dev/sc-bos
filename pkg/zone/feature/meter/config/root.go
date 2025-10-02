package config

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.Config

	Meters                []string            `json:"meters,omitempty"`
	MeterGroups           map[string][]string `json:"meterGroups,omitempty"`
	UseHistoryBackupOnErr bool                `json:"useHistoryBackupOnErr,omitempty"`
}
