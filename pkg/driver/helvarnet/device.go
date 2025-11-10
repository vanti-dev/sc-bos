package helvarnet

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/healthpb"
)

const (
	DeviceOfflineCode = "-1"
	BadResponseCode   = "-2"
	SystemName        = "HelvarNet Lighting"
)

type device struct {
	client     *tcpClient
	conf       *config.Device
	logger     *zap.Logger
	faultCheck *healthpb.FaultCheck
}

// this health check is checks the helvarnet status field on the device itself and reports any faults
// based on that. If we cannot get the status it means we have no comms with the device
func getDeviceHealthCheck() *gen.HealthCheck {
	return &gen.HealthCheck{
		Id:              "device_status_check",
		DisplayName:     "Device Status Check",
		Description:     "Checks the status from the device itself and also if communication is healthy",
		OccupantImpact:  gen.HealthCheck_COMFORT,
		EquipmentImpact: gen.HealthCheck_FUNCTION,
	}
}

func updateDeviceFaults(ctx context.Context, status int64, fc *healthpb.FaultCheck) {
	if status < 0 {
		rel := &gen.HealthCheck_Reliability{
			UnreliableTime: timestamppb.Now(),
		}
		if status == NoResponse {
			rel.State = gen.HealthCheck_Reliability_NO_RESPONSE
			fc.SetFault(&gen.HealthCheck_Error{
				SummaryText: "Device Offline",
				DetailsText: "No communication received from device since the last smart core restart",
				Code: &gen.HealthCheck_Error_Code{
					Code:   DeviceOfflineCode,
					System: SystemName,
				},
			})
		} else if status == BadResponse {
			rel.State = gen.HealthCheck_Reliability_BAD_RESPONSE
			fc.SetFault(&gen.HealthCheck_Error{
				SummaryText: "Bad Response",
				DetailsText: "The device has sent an invalid response to a command",
				Code: &gen.HealthCheck_Error_Code{
					Code:   BadResponseCode,
					System: SystemName,
				},
			})
		}
		fc.UpdateReliability(ctx, rel)

	} else {

		statuses := config.GetStatusListFromFlag(status)

		if len(statuses) == 0 {
			fc.ClearFaults()
		} else {
			for _, s := range statuses {
				fc.AddOrUpdateFault(&gen.HealthCheck_Error{
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
}

// getHelvarnetStatus queries the device and updates the status value
func (d *device) getHelvarnetStatus() (int64, error) {
	command := queryDeviceState(d.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := d.client.sendAndReceive(command, want)
	if err != nil {
		return NoResponse, err
	}

	split := strings.Split(r, "=")
	if len(split) < 2 {
		return BadResponse, fmt.Errorf("invalid response in getHelvarnetStatus: %s", r)
	}

	s := strings.TrimSuffix(split[1], "#")
	statusInt, err := strconv.Atoi(s)
	if err != nil {
		return BadResponse, err
	}

	return int64(statusInt), nil
}
