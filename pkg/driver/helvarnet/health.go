package helvarnet

import (
	"context"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/healthpb"
)

const (
	DeviceOfflineCode     = -1
	BadResponseCode       = -2
	UnrecognisedErrorCode = -100
	SystemName            = "HelvarNet Lighting"
)

type State struct {
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
	FlagValue   int64  `json:"flagValue,omitempty"`
}

// DeviceStatuses lists all possible device states and their associated metadata.
// These states, descriptions and codes are copied from the HelvarNet documentation.
var DeviceStatuses = []State{
	{"Disabled", "Device or subdevice has been disabled, usually an IR subdevice or a DMX channel", 0x00000001},
	{"LampFailure", "Unspecified lamp problem", 0x00000002},
	{"Missing", "The device previously existed but is not currently present", 0x00000004},
	{"Faulty", "Ran out of addresses (DALI subnet) / unknown Digidim control device / DALI load that keeps responding with multi-replies", 0x00000008},
	{"Refreshing", "DALI subnet, DALI load or Digidim control device is being discovered", 0x00000010},
	{"Resting", "The load is intentionally off whilst the control gear is being powered by the emergency supply", 0x00000100},
	{"InEmergency", "No mains power is being supplied", 0x00000400},
	{"InProlong", "Mains has been restored but device is still using the emergency supply", 0x00000800},
	{"FTInProgress", "The Functional Test is in progress (brief test where the control gear is being powered by the emergency supply)", 0x00001000},
	{"DTInProgress", "The Duration Test is in progress. This test involves operating the control gear using the battery until the battery is completely discharged. The duration that the control gear was operational for is recorded, and then the battery recharges itself from the mains supply", 0x00002000},
	{"DTPending", "The Duration Test has been requested but has not yet commenced. The test can be delayed if the battery is not fully charged", 0x00010000},
	{"FTPending", "The Functional Test has been requested but has not yet commenced. The test can be delayed if there is not enough charge in the battery", 0x00020000},
	{"BatteryFail", "Battery has failed", 0x00040000},
	{"Inhibit", "Prevents an emergency fitting from going into emergency mode", 0x00200000},
	{"FTRequested", "Emergency Function Test has been requested", 0x00400000},
	{"DTRequested", "Emergency Duration Test has been requested", 0x00800000},
	{"Unknown", "Initial state of an emergency fitting", 0x01000000},
	{"OverTemperature", "Load is over temperature/heating", 0x02000000},
	{"OverCurrent", "Too much current is being drawn by the load", 0x04000000},
	{"CommsError", "Communications error", 0x08000000},
	{"SevereError", "Indicates that a load is either over temperature or drawing too much current, or both", 0x10000000},
	{"BadReply", "Indicates that a reply to a query was malformed", 0x20000000},
	{"DeviceMismatch", "The actual load type does not match the expected type", 0x80000000},
}

// GetStatusListFromFlag takes a bitwise flag integer and returns a list of State structs
// corresponding to the set bits in the flag.
func GetStatusListFromFlag(flag int64) []State {

	if flag == 0 {
		return []State{}
	}

	var statusList []State
	for _, ds := range DeviceStatuses {
		if flag&ds.FlagValue != 0 {
			statusList = append(statusList, ds)
		}
	}

	return statusList
}

func updateDeviceFaults(ctx context.Context, status int64, fc *healthpb.FaultCheck, raisedFaults map[int64]bool) {

	// Handle negative status codes for special conditions, mainly comms issues
	if status < 0 {
		rel := &gen.HealthCheck_Reliability{
			UnreliableTime: timestamppb.Now(),
		}
		switch status {
		case DeviceOfflineCode:
			rel.State = gen.HealthCheck_Reliability_NO_RESPONSE
			fc.SetFault(&gen.HealthCheck_Error{
				SummaryText: "Device Offline",
				DetailsText: "No communication received from device since the last Smart Core restart",
				Code: &gen.HealthCheck_Error_Code{
					Code:   strconv.Itoa(DeviceOfflineCode),
					System: SystemName,
				},
			})
		case BadResponseCode:
			rel.State = gen.HealthCheck_Reliability_BAD_RESPONSE
			fc.SetFault(&gen.HealthCheck_Error{
				SummaryText: "Bad Response",
				DetailsText: "The device has sent an invalid response to a command",
				Code: &gen.HealthCheck_Error_Code{
					Code:   strconv.Itoa(BadResponseCode),
					System: SystemName,
				},
			})
		default:
			// this should really never happen, but if it does, then it is a problem with the driver
			// and it should be flagged
			rel.State = gen.HealthCheck_Reliability_UNRELIABLE
			fc.SetFault(&gen.HealthCheck_Error{
				SummaryText: "Internal Driver Error",
				DetailsText: "The device has an unrecognised internal status code",
				Code: &gen.HealthCheck_Error_Code{
					Code:   strconv.Itoa(UnrecognisedErrorCode),
					System: SystemName,
				},
			})
		}
		fc.UpdateReliability(ctx, rel)

	} else {

		statuses := GetStatusListFromFlag(status)

		if len(statuses) == 0 {
			fc.ClearFaults()

			for code := range raisedFaults {
				raisedFaults[code] = false
			}
		} else {
			setDeviceFaults(statuses, fc, raisedFaults)
		}
	}
}

func setDeviceFaults(statuses []State, fc *healthpb.FaultCheck, raisedFaults map[int64]bool) {
	for _, s := range statuses {
		fc.AddOrUpdateFault(&gen.HealthCheck_Error{
			SummaryText: s.State,
			DetailsText: s.Description,
			Code: &gen.HealthCheck_Error_Code{
				Code:   strconv.Itoa(int(s.FlagValue)),
				System: SystemName,
			},
		})
		raisedFaults[s.FlagValue] = true
	}

	for code, raised := range raisedFaults {
		if raised {
			// if we have raised the fault in smart core but it is no longer being reported by the device, it needs to be removed in sc
			found := false
			for _, s := range statuses {
				if s.FlagValue == code {
					found = true
					break
				}
			}
			if !found {
				fc.RemoveFault(&gen.HealthCheck_Error{
					Code: &gen.HealthCheck_Error_Code{
						Code:   strconv.Itoa(int(code)),
						System: SystemName,
					},
				})
				raisedFaults[code] = false
			}
		}
	}
}
