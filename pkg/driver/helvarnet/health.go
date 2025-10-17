package helvarnet

import (
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

const (
	DeviceOfflineCode = "-1"
	BadResponseCode   = "-2"
	SystemName        = "HelvarNet Lighting"
)

var deviceStatusCheck = &gen.HealthCheck{
	Id:              "device_status_check",
	DisplayName:     "Light Status Check",
	Description:     "Checks the status from the device itself and also if communication is healthy",
	Reliability:     &gen.HealthCheck_Reliability{},
	OccupantImpact:  gen.HealthCheck_COMFORT,
	EquipmentImpact: gen.HealthCheck_FUNCTION,
}

var allHealthChecks = [1]*gen.HealthCheck{deviceStatusCheck}

// this health check is checks the helvarnet status field on the device itself and reports any faults
// based on that. If we cannot get the status it means we have no comms with the device
func getDeviceStatusCheck(status int64) *gen.HealthCheck {

	check := deviceStatusCheck

	var faults []*gen.HealthCheck_Error
	if status < 0 {
		check.Normality = gen.HealthCheck_ABNORMAL
		check.AbnormalTime = timestamppb.Now()
		check.Reliability.UnreliableTime = timestamppb.Now()
		if status == NoResponse {
			check.Reliability.State = gen.HealthCheck_Reliability_NO_RESPONSE
			faults = append(faults, &gen.HealthCheck_Error{
				SummaryText: "Device Offline",
				DetailsText: "The device has not responded to commands since the last restart of this driver",
				Code: &gen.HealthCheck_Error_Code{
					Code:   DeviceOfflineCode,
					System: SystemName,
				},
			})
		} else if status == BadResponse {
			check.Reliability.State = gen.HealthCheck_Reliability_BAD_RESPONSE
			faults = append(faults, &gen.HealthCheck_Error{
				SummaryText: "Bad Response",
				DetailsText: "The device has sent an invalid response to a command",
				Code: &gen.HealthCheck_Error_Code{
					Code:   BadResponseCode,
					System: SystemName,
				},
			})
		}
	} else {
		check.Reliability.State = gen.HealthCheck_Reliability_RELIABLE
		check.Reliability.ReliableTime = timestamppb.Now()

		statuses := config.GetStatusListFromFlag(status)

		if len(statuses) == 0 {
			check.Normality = gen.HealthCheck_NORMAL
			check.NormalTime = timestamppb.Now()
		} else {
			check.Normality = gen.HealthCheck_ABNORMAL
			check.AbnormalTime = timestamppb.Now()
			for _, s := range statuses {
				faults = append(faults, &gen.HealthCheck_Error{
					SummaryText: s.State,
					DetailsText: s.Description,
					Code: &gen.HealthCheck_Error_Code{
						Code:   strconv.Itoa(int(s.FlagValue)),
						System: SystemName,
					},
				})
			}
		}
	}

	check.Check = &gen.HealthCheck_Faults_{
		Faults: &gen.HealthCheck_Faults{
			CurrentFaults: faults,
		},
	}

	return check
}
