package bridge

import (
	"context"
	"io"
)

type Command int

const (
	DirectArcPowerControl Command = 1 + iota
	QueryActualLevel
	QueryStatus
	SetEventScheme
	QueryEventScheme
	SetEventFilter
	QueryEventFilter
	IdentifyDevice102
	IdentifyDevice103
	EnableInstance
	GoToScene
	SetFadeTime
	QueryInputValue

	QueryBatteryCharge
	QueryDurationTestResult
	QueryEmergencyMode
	QueryEmergencyStatus
	QueryRatedDuration
	QueryTestTiming
	ResetDurationTestDoneFlag
	ResetFunctionTestDoneFlag
	StartDurationTest
	StartFunctionTest
	StopTest
	SetDurationTestInterval
	SetFunctionTestInterval

	TestCommand Command = 100
)

type AddressType byte

const (
	Short AddressType = iota
	Group
	Broadcast
	BroadcastUnaddr
)

type InstanceAddressType byte

const (
	IATInstanceNumber InstanceAddressType = iota
	IATInstanceGroup
	IATInstanceType
	IATFeatureNumber
	IATFeatureGroup
	IATFeatureType
	IATFeatureBroadcast
	IATInstanceBroadcast
	IATFeatureDevice
	IATDevice
)

// Request is an instruction to the DALI implementation to perform a DALI command.
type Request struct {
	// Command specifies which command will be sent to the DALI bus
	Command Command `tc3ads:"command"`
	// AddressType specifies which style of addressing will be used to select the target DALI control gear/devices
	AddressType AddressType `tc3ads:"addressType"`
	// Address is the address used to target the command. If the AddressType is a broadcast type, when Address is
	// ignored.
	Address             byte                `tc3ads:"address"`
	InstanceAddressType InstanceAddressType `tc3ads:"instanceAddressType"`
	InstanceAddress     byte                `tc3ads:"instanceAddress"`
	// Data is the payload of the command. For example, when performing a DirectArcPowerControl, then Data contains
	// the level to write. If the Command selected does not need a data input, then Data is ignored.
	Data byte `tc3ads:"data"`
}

type InputEventHandler func(event InputEvent, err error)

type Dali interface {
	ExecuteCommand(ctx context.Context, request Request) (data uint32, err error)
	EnableInputEventListener(params InputEventParameters) error
	OnInputEvent(handler InputEventHandler) error

	io.Closer
}
