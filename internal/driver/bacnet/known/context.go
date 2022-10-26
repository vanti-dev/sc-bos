package known

import bactypes "github.com/vanti-dev/gobacnet/types"

// Context describes what we know about a bacnet system.
type Context interface {
	LookupDeviceByID(id bactypes.ObjectInstance) (bactypes.Device, error)
	LookupDeviceByName(name string) (bactypes.Device, error)
	LookupObjectByID(device bactypes.Device, id bactypes.ObjectID) (bactypes.Object, error)
	LookupObjectByName(device bactypes.Device, name string) (bactypes.Object, error)
}
