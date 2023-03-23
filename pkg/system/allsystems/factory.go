package allsystems

import (
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/alerts"
	"github.com/vanti-dev/sc-bos/pkg/system/authn"
	"github.com/vanti-dev/sc-bos/pkg/system/hub"
	"github.com/vanti-dev/sc-bos/pkg/system/publications"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants"
)

// Factories returns a new map containing all known system factories.
func Factories() map[string]system.Factory {
	return map[string]system.Factory{
		"alerts":       alerts.Factory,
		"authn":        authn.Factory(),
		"hub":          hub.Factory(),
		"publications": publications.Factory,
		"tenants":      tenants.Factory,
	}
}
