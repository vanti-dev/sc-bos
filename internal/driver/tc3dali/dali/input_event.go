package dali

import (
	"fmt"
)

// InputEventParameters determines which identification data is bundled with event notifications from control devices.
// A DALI input device will only ever send two of these identifiers, in addition to the Scheme. The other identifiers
// will not contain useful data, and must be ignored.
// See EventScheme documentation for information on which fields will be valid.
type InputEventParameters struct {
	Scheme       EventScheme `tc3ads:"eventScheme"`
	AddressInfo1 byte        `tc3ads:"addressInfo1"`
	AddressInfo2 byte        `tc3ads:"addressInfo2"`
}

func InputEventParametersForInstance(deviceShortAddress byte, instance byte) InputEventParameters {
	return InputEventParameters{
		Scheme:       EventSchemeDeviceInstance,
		AddressInfo1: deviceShortAddress,
		AddressInfo2: instance,
	}
}

func (p InputEventParameters) InstanceType() InstanceType {
	switch p.Scheme {
	case EventSchemeInstance:
		return InstanceType(p.AddressInfo1)
	case EventSchemeDevice, EventSchemeDeviceGroup, EventSchemeInstanceGroup:
		return InstanceType(p.AddressInfo2)
	}
	panic(fmt.Sprintf("EventScheme %v does not contain an InstanceType", p.Scheme))
}

func (p InputEventParameters) InstanceNumber() byte {
	switch p.Scheme {
	case EventSchemeInstance, EventSchemeDeviceInstance:
		return p.AddressInfo2
	}
	panic(fmt.Sprintf("EventScheme %v does not contain an InstanceNumber", p.Scheme))
}

func (p InputEventParameters) DeviceShortAddress() byte {
	switch p.Scheme {
	case EventSchemeDevice, EventSchemeDeviceInstance:
		return p.AddressInfo1
	}
	panic(fmt.Sprintf("EventScheme %v does not contain a DeviceShortAddress", p.Scheme))
}

func (p InputEventParameters) DeviceGroup() byte {
	switch p.Scheme {
	case EventSchemeDeviceGroup:
		return p.AddressInfo1
	}
	panic(fmt.Sprintf("EventScheme %v does not contain a DeviceGroup", p.Scheme))
}

func (p InputEventParameters) InstanceGroup() byte {
	switch p.Scheme {
	case EventSchemeInstanceGroup:
		return p.AddressInfo1
	}
	panic(fmt.Sprintf("EventScheme %v does not contain an InstanceGroup", p.Scheme))
}

type InputEvent struct {
	InputEventParameters
	Err  error
	Data uint16
}
