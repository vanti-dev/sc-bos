package allautos

import (
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/export"
	"github.com/vanti-dev/sc-bos/pkg/auto/history"
	"github.com/vanti-dev/sc-bos/pkg/auto/lights"
	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
)

// Factories returns a new map containing all known auto factories.
func Factories() map[string]auto.Factory {
	return map[string]auto.Factory{
		"export-mqtt":   export.MQTTFactory,
		"history":       history.Factory,
		lights.AutoType: lights.Factory,
		udmi.AutoType:   udmi.Factory,
	}
}
