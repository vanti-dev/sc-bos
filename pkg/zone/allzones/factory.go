package allzones

import (
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/area"
)

// Factories returns a new map containing all known zone factories.
func Factories() map[string]zone.Factory {
	return map[string]zone.Factory{
		"area": area.Factory,
	}
}
