package comm

import (
	"fmt"
)

type LiftCarDirection uint16 // use uint16 since values up to 65535 are valid

const (
	DirectionUnknown   LiftCarDirection = 0
	DirectionNone      LiftCarDirection = 1
	DirectionStopped   LiftCarDirection = 2
	DirectionUp        LiftCarDirection = 3
	DirectionDown      LiftCarDirection = 4
	DirectionUpAndDown LiftCarDirection = 5

	// 6..1023 reserved by ASHRAE
	// 1024..65535 vendor specific
)

func (d LiftCarDirection) String() string {
	switch d {
	case DirectionUnknown:
		return "unknown"
	case DirectionNone:
		return "none"
	case DirectionStopped:
		return "stopped"
	case DirectionUp:
		return "up"
	case DirectionDown:
		return "down"
	case DirectionUpAndDown:
		return "up-and-down"
	default:
		if d <= 1023 {
			return fmt.Sprintf("reserved(%d)", d)
		}
		return fmt.Sprintf("vendor(%d)", d)
	}
}

type LandingCall struct {
	Floor     uint8
	Direction LiftCarDirection
}

type MakingCarCall []uint8

type DoorStatus uint16 // up to 65535

const (
	DoorClosed        DoorStatus = 0
	DoorOpened        DoorStatus = 1
	DoorUnknown       DoorStatus = 2
	DoorFault         DoorStatus = 3
	DoorUnused        DoorStatus = 4
	DoorNone          DoorStatus = 5
	DoorClosing       DoorStatus = 6
	DoorOpening       DoorStatus = 7
	DoorSafetyLocked  DoorStatus = 8
	DoorLimitedOpened DoorStatus = 9

	// 10..1023 reserved by ASHRAE
	// 1024..65535 vendor-specific
)

func (d DoorStatus) String() string {
	switch d {
	case DoorClosed:
		return "closed"
	case DoorOpened:
		return "opened"
	case DoorUnknown:
		return "unknown"
	case DoorFault:
		return "door-fault"
	case DoorUnused:
		return "unused"
	case DoorNone:
		return "none"
	case DoorClosing:
		return "closing"
	case DoorOpening:
		return "opening"
	case DoorSafetyLocked:
		return "safety-locked"
	case DoorLimitedOpened:
		return "limited-opened"
	default:
		if d <= 1023 {
			return fmt.Sprintf("reserved(%d)", d)
		}
		return fmt.Sprintf("vendor(%d)", d)
	}
}

type LiftCarMode int64

const (
	LiftCarModeUnknown            LiftCarMode = 0
	LiftCarModeNormal             LiftCarMode = 1 // in service
	LiftCarModeVIP                LiftCarMode = 2
	LiftCarModeHoming             LiftCarMode = 3
	LiftCarModeParking            LiftCarMode = 4
	LiftCarModeAttendantControl   LiftCarMode = 5
	LiftCarModeFirefighterControl LiftCarMode = 6
	LiftCarModeEmergencyPower     LiftCarMode = 7
	LiftCarModeInspection         LiftCarMode = 8
	LiftCarModeCabinetRecall      LiftCarMode = 9
	LiftCarModeEarthquakeOp       LiftCarMode = 10
	LiftCarModeFireOp             LiftCarMode = 11
	LiftCarModeOutOfService       LiftCarMode = 12
	LiftCarModeOccupantEvac       LiftCarMode = 13
)

func (m LiftCarMode) String() string {
	switch m {
	case LiftCarModeUnknown:
		return "unknown"
	case LiftCarModeNormal:
		return "normal"
	case LiftCarModeVIP:
		return "vip"
	case LiftCarModeHoming:
		return "homing"
	case LiftCarModeParking:
		return "parking"
	case LiftCarModeAttendantControl:
		return "attendant-control"
	case LiftCarModeFirefighterControl:
		return "firefighter-control"
	case LiftCarModeEmergencyPower:
		return "emergency-power"
	case LiftCarModeInspection:
		return "inspection"
	case LiftCarModeCabinetRecall:
		return "cabinet-recall"
	case LiftCarModeEarthquakeOp:
		return "earthquake-operation"
	case LiftCarModeFireOp:
		return "fire-operation"
	case LiftCarModeOutOfService:
		return "out-of-service"
	case LiftCarModeOccupantEvac:
		return "occupant-evacuation"
	default:
		if m <= 1023 {
			return fmt.Sprintf("reserved(%d)", m)
		}
		return fmt.Sprintf("vendor(%d)", m)
	}
}

type LiftFault uint16

const (
	LiftFaultControllerFault            LiftFault = 0
	LiftFaultDriveAndMotorFault         LiftFault = 1
	LiftFaultGovernorAndSafetyGearFault LiftFault = 2
	LiftFaultLiftShaftDeviceFault       LiftFault = 3
	LiftFaultPowerSupplyFault           LiftFault = 4
	LiftFaultSafetyInterlockFault       LiftFault = 5
	LiftFaultDoorClosingFault           LiftFault = 6
	LiftFaultDoorOpeningFault           LiftFault = 7
	LiftFaultCarStoppedOutsideLanding   LiftFault = 8
	LiftFaultCallButtonStuck            LiftFault = 9
	LiftFaultStartFailure               LiftFault = 10
	LiftFaultControllerSupplyFault      LiftFault = 11
	LiftFaultSelfTestFailure            LiftFault = 12
	LiftFaultRuntimeLimitExceeded       LiftFault = 13
	LiftFaultPositionLost               LiftFault = 14
	LiftFaultDriveTempExceeded          LiftFault = 15
	LiftFaultLoadMeasurementFault       LiftFault = 16
)

func (f LiftFault) String() string {
	switch f {
	case LiftFaultControllerFault:
		return "controller-fault"
	case LiftFaultDriveAndMotorFault:
		return "drive-and-motor-fault"
	case LiftFaultGovernorAndSafetyGearFault:
		return "governor-and-safety-gear-fault"
	case LiftFaultLiftShaftDeviceFault:
		return "lift-shaft-device-fault"
	case LiftFaultPowerSupplyFault:
		return "power-supply-fault"
	case LiftFaultSafetyInterlockFault:
		return "safety-interlock-fault"
	case LiftFaultDoorClosingFault:
		return "door-closing-fault"
	case LiftFaultDoorOpeningFault:
		return "door-opening-fault"
	case LiftFaultCarStoppedOutsideLanding:
		return "car-stopped-outside-landing-zone"
	case LiftFaultCallButtonStuck:
		return "call-button-stuck"
	case LiftFaultStartFailure:
		return "start-failure"
	case LiftFaultControllerSupplyFault:
		return "controller-supply-fault"
	case LiftFaultSelfTestFailure:
		return "self-test-failure"
	case LiftFaultRuntimeLimitExceeded:
		return "runtime-limit-exceeded"
	case LiftFaultPositionLost:
		return "position-lost"
	case LiftFaultDriveTempExceeded:
		return "drive-temperature-exceeded"
	case LiftFaultLoadMeasurementFault:
		return "load-measurement-fault"
	default:
		if f <= 1023 {
			return fmt.Sprintf("reserved(%d)", f)
		}
		return fmt.Sprintf("vendor(%d)", f)
	}
}
