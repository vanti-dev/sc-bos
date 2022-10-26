package known

import (
	"errors"
	"fmt"
	bactypes "github.com/vanti-dev/gobacnet/types"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrUnknownDevice = fmt.Errorf("device %w", ErrNotFound)
)

// Map holds information about devices and objects.
// Map implements Context.
type Map struct {
	devicesByID   map[bactypes.ObjectInstance]*device
	devicesByName map[string]*device
}

func NewMap() *Map {
	return &Map{}
}

func (m *Map) StoreDevice(name string, d bactypes.Device) {
	if m.devicesByID == nil {
		m.devicesByID = make(map[bactypes.ObjectInstance]*device)
		m.devicesByName = make(map[string]*device)
	}

	item := &device{
		bacDevice:     d,
		name:          name,
		objectsByID:   make(map[bactypes.ObjectID]*object),
		objectsByName: make(map[string]*object),
	}
	m.devicesByID[d.ID.Instance] = item
	m.devicesByName[name] = item
}

func (m *Map) StoreObject(d bactypes.Device, name string, o bactypes.Object) error {
	item, ok := m.devicesByID[d.ID.Instance]
	if !ok {
		return ErrUnknownDevice
	}
	obj := &object{
		bacObject: o,
		name:      name,
	}
	item.objectsByID[o.ID] = obj
	item.objectsByName[name] = obj
	return nil
}

func (m *Map) DeleteDevice(device bactypes.Device) {
	m.DeleteDeviceByID(device.ID.Instance)
}

func (m *Map) DeleteDeviceByID(id bactypes.ObjectInstance) {
	d, ok := m.devicesByID[id]
	if !ok {
		return
	}
	delete(m.devicesByID, d.bacDevice.ID.Instance)
	delete(m.devicesByName, d.name)
}

func (m *Map) DeleteDeviceByName(name string) {
	d, ok := m.devicesByName[name]
	if !ok {
		return
	}
	delete(m.devicesByID, d.bacDevice.ID.Instance)
	delete(m.devicesByName, d.name)
}

func (m *Map) DeleteObject(device bactypes.Device, object bactypes.Object) {
	m.DeleteObjectByID(device, object.ID)
}

func (m *Map) DeleteObjectByID(device bactypes.Device, id bactypes.ObjectID) {
	d, ok := m.devicesByID[device.ID.Instance]
	if !ok {
		return
	}
	o, ok := d.objectsByID[id]
	if !ok {
		return
	}
	delete(d.objectsByID, o.bacObject.ID)
	delete(d.objectsByName, o.name)
}

func (m *Map) DeleteObjectByName(device bactypes.Device, name string) {
	d, ok := m.devicesByID[device.ID.Instance]
	if !ok {
		return
	}
	o, ok := d.objectsByName[name]
	if !ok {
		return
	}
	delete(d.objectsByID, o.bacObject.ID)
	delete(d.objectsByName, o.name)
}

func (m *Map) Clear() {
	m.devicesByID = nil
	m.devicesByName = nil
}

func (m *Map) LookupDeviceByID(id bactypes.ObjectInstance) (bactypes.Device, error) {
	d, ok := m.devicesByID[id]
	if !ok {
		return bactypes.Device{}, ErrNotFound
	}
	return d.bacDevice, nil
}

func (m *Map) LookupDeviceByName(name string) (bactypes.Device, error) {
	d, ok := m.devicesByName[name]
	if !ok {
		return bactypes.Device{}, ErrNotFound
	}
	return d.bacDevice, nil
}

func (m *Map) LookupObjectByID(device bactypes.Device, id bactypes.ObjectID) (bactypes.Object, error) {
	d, ok := m.devicesByID[device.ID.Instance]
	if !ok {
		return bactypes.Object{}, ErrUnknownDevice
	}
	o, ok := d.objectsByID[id]
	if !ok {
		return bactypes.Object{}, ErrNotFound
	}
	return o.bacObject, nil
}

func (m *Map) LookupObjectByName(device bactypes.Device, name string) (bactypes.Object, error) {
	d, ok := m.devicesByID[device.ID.Instance]
	if !ok {
		return bactypes.Object{}, ErrUnknownDevice
	}
	o, ok := d.objectsByName[name]
	if !ok {
		return bactypes.Object{}, ErrNotFound
	}
	return o.bacObject, nil
}

type device struct {
	bacDevice bactypes.Device
	name      string

	objectsByID   map[bactypes.ObjectID]*object
	objectsByName map[string]*object
}

type object struct {
	bacObject bactypes.Object
	name      string
}
