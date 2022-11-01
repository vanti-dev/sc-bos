package bridge

import (
	"context"
	"io"
)

type Command int

const (
	DirectArcPowerControl     Command = 1
	QueryActualLevel          Command = 2
	QueryStatus               Command = 3
	SetEventScheme            Command = 4
	QueryEventScheme          Command = 5
	SetEventFilter            Command = 6
	QueryEventFilter          Command = 7
	IdentifyDevice102         Command = 8
	IdentifyDevice103         Command = 9
	EnableInstance            Command = 10
	GoToScene                 Command = 11
	SetFadeTime               Command = 12
	QueryInputValue           Command = 13
	QueryBatteryCharge        Command = 14
	QueryDurationTestResult   Command = 15
	QueryEmergencyMode        Command = 16
	QueryEmergencyStatus      Command = 17
	QueryFailureStatus        Command = 18
	QueryRatedDuration        Command = 19
	QueryTestTiming           Command = 20
	ResetDurationTestDoneFlag Command = 21
	ResetFunctionTestDoneFlag Command = 22
	StartDurationTest         Command = 23
	StartFunctionTest         Command = 24
	StopTest                  Command = 25
	SetDurationTestInterval   Command = 26
	SetFunctionTestInterval   Command = 27
	StartIdentification202    Command = 28
	QueryGroups               Command = 29
	AddToGroup                Command = 30
	RemoveFromGroup           Command = 31

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
