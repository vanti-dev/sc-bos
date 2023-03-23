package config

import (
	"github.com/vanti-dev/sc-bos/pkg/system"
)

type Root struct {
	system.Config

	// Ignore contains a list of enrolled host:port that we should not proxy.
	// Useful if you setup more than one proxy enrolled with the same hub.
	// This controller will always ignore it's own endpoint.
	Ignore []string `json:"ignore,omitempty"`
}
