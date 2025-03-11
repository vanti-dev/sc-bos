package modbus

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

/** TODO:
This should be configurable per modbus device implementation
since we don't have **many** modbus device specifications
I will assume this is compatible with **all** modbus devices, for now
*/

func convertFuel(res []byte, scale float32) float32 {
	if len(res) != 2 {
		return 0
	}

	return float32(uint16(res[0])<<8|uint16(res[1])) * scale
}

func checkFaults(res []byte) traits.Emergency_Level {
	if len(res) != 2 {
		return traits.Emergency_OK
	}

	status := uint16(res[0])<<8 | uint16(res[1])

	// shutdown alarm active
	if status&(1<<13) != 0 {
		return traits.Emergency_EMERGENCY
	}

	// control unit failure
	if status&(1<<14) != 0 {
		return traits.Emergency_WARNING
	}

	// electrical trip / controlled shutdown
	if status&(1<<12) != 0 {
		return traits.Emergency_WARNING
	}

	//  warning alarm
	if status&(1<<11) != 0 {
		return traits.Emergency_WARNING
	}

	// telemetry alarm active
	if status&(1<<10) != 0 {
		return traits.Emergency_WARNING
	}

	// satellite telemetry alarm
	if status&(1<<9) != 0 {
		return traits.Emergency_WARNING
	}

	return traits.Emergency_OK
}

func checkEngineState(res []byte) (gen.StatusLog_Level, string) {
	if len(res) != 2 {
		return gen.StatusLog_OFFLINE, ""
	}

	state := uint16(res[0])<<8 | uint16(res[1])

	if desc, ok := engineStates[state]; ok {
		switch desc {
		case "Engine Stopped":
			return gen.StatusLog_NON_FUNCTIONAL, desc
		case "Running":
			return gen.StatusLog_NOMINAL, desc
		case "Cooling Down":
			return gen.StatusLog_REDUCED_FUNCTION, desc
		case "Pre-Start":
			return gen.StatusLog_REDUCED_FUNCTION, desc
		case "Post Run":
			return gen.StatusLog_NOMINAL, desc
		default:
			return gen.StatusLog_NOMINAL, ""
		}
	}

	return gen.StatusLog_LEVEL_UNDEFINED, ""

}

var engineStates = map[uint16]string{
	0:  "Engine Stopped",
	1:  "Pre-Start",
	2:  "Warming Up",
	3:  "Running",
	4:  "Cooling Down",
	5:  "Engine Stopped",
	6:  "Post Run",
	7:  "Reserved",
	8:  "Available for SAE Assignment",
	9:  "Available for SAE Assignment",
	10: "Available for SAE Assignment",
	11: "Available for SAE Assignment",
	12: "Available for SAE Assignment",
	13: "Available for SAE Assignment",
	14: "Reserved",
	15: "Not Available",
}
