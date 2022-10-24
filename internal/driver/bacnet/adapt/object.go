package adapt

import (
	"fmt"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/config"
)

// DeviceName returns the smart core name we should use for the configured object.
func DeviceName(o config.Device) string {
	if o.Name != "" {
		return o.Name
	}
	return fmt.Sprintf("%d", o.ID)
}

// ObjectName returns the smart core name we should use for the configured object.
func ObjectName(o config.Object) string {
	if o.Name != "" {
		return o.Name
	}
	return o.ID.String()
}
