package airtempmonitor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto/airtempmonitor/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// abnormalTemperatureCheck checks if the return air temperatures are within the normal range defined in the config.
// It generates alerts for abnormal temperatures and resolves them when temperatures return to normal.
func (a *deviceMonitorAuto) abnormalTemperatureCheck(ctx context.Context, client traits.AirTemperatureApiClient, config *config.AirTempConfig) {
	for _, device := range config.Devices {
		resp, err := client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{
			Name: device.Name,
		})
		if err != nil {
			a.Logger.Error("failed to get air temperature", zap.Error(err))
			return
		}
		if resp.AmbientTemperature != nil {
			if resp.AmbientTemperature.ValueCelsius < *config.OkRtLowerBound {
				if device.AbnormalLowId == "" { // don't re-raise the alert if it is already there
					alert, err := a.alertAdminClient.CreateAlert(ctx, &gen.CreateAlertRequest{
						Alert: &gen.Alert{
							Id:          uuid.New().String(),
							CreateTime:  timestamppb.Now(),
							Description: "Ambient temperature is abnormally low",
							Severity:    gen.Alert_WARNING,
							Source:      device.Name,
						},
					})
					if err != nil {
						a.Logger.Error("failed to create alert", zap.Error(err))
					} else {
						device.AbnormalLowId = alert.Id
					}
				}
			} else {
				if device.AbnormalLowId != "" {
					_, err := a.alertAdminClient.ResolveAlert(ctx, &gen.ResolveAlertRequest{
						Alert: &gen.Alert{
							Id:          device.AbnormalLowId,
							ResolveTime: timestamppb.Now(),
						},
					})
					if err != nil {
						a.Logger.Error("failed to resolve alert", zap.Error(err))
					} else {
						device.AbnormalLowId = ""
					}
				}
			}
			if resp.AmbientTemperature.ValueCelsius > *config.OkRtUpperBound {
				if device.AbnormalHighId == "" {
					alert, err := a.alertAdminClient.CreateAlert(ctx, &gen.CreateAlertRequest{
						Alert: &gen.Alert{
							Id:          uuid.New().String(),
							CreateTime:  timestamppb.Now(),
							Description: "Ambient temperature is abnormally high",
							Severity:    gen.Alert_WARNING,
							Source:      device.Name,
						},
					})
					if err != nil {
						a.Logger.Error("failed to create alert", zap.Error(err))
					} else {
						device.AbnormalHighId = alert.Id
					}
				}
			} else {
				if device.AbnormalHighId != "" {
					_, err := a.alertAdminClient.ResolveAlert(ctx, &gen.ResolveAlertRequest{
						Alert: &gen.Alert{
							Id:          device.AbnormalHighId,
							ResolveTime: timestamppb.Now(),
						},
					})
					if err != nil {
						a.Logger.Error("failed to resolve alert", zap.Error(err))
					} else {
						device.AbnormalHighId = ""
					}
				}
			}
		}
	}
}

// hasReachedSetPoint checks if the measured temperature is within the acceptable tolerance range of the setPoint value.
// It also respects the okSettlingTime by allowing a grace period after a set point change before evaluating.
func (a *deviceMonitorAuto) hasReachedSetPoint(measured, setPoint float64, cfg *config.AirTempConfig, spChangedTime time.Time) bool {

	if spChangedTime.Add(cfg.OkSettlingTime.Duration).After(a.Now()) {
		// if we aren't expecting the set point to be reached until after Now, just return true
		return true
	}

	tolerance := *cfg.Tolerance
	if measured >= setPoint-tolerance && measured <= setPoint+tolerance {
		return true
	}
	return false
}

// returnAirReachesSetPointCheck verifies if the return air temperature of each device aligns with its set point within tolerance.
// It generates alerts if the set point is not reached within the allowed time and resolves them when conditions normalize.
// If the set point is changed, the timer is reset and the okSettlingTime begins again from Now.
func (a *deviceMonitorAuto) returnAirReachesSetPointCheck(ctx context.Context, client traits.AirTemperatureApiClient, config *config.AirTempConfig) {
	for _, device := range config.Devices {
		resp, err := client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{
			Name: device.Name,
		})
		if err != nil {
			a.Logger.Error("failed to get air temperature", zap.Error(err))
			return
		}
		setPoint := resp.GetTemperatureSetPoint()
		if resp.AmbientTemperature != nil && setPoint != nil {
			if device.PreviousSetPoint == nil {
				// we need to be able to detect when the set point has been changed and reset the timer
				// when the auto first runs, initialise the previous set point to the current set point
				device.PreviousSetPoint = proto.Float64(setPoint.ValueCelsius)
				device.SetPointChangedTime = a.Now()
			}
			if *device.PreviousSetPoint != setPoint.ValueCelsius {
				// if the set point has been changed, we need to reset the timer
				device.SetPointChangedTime = a.Now()
			}
			if !a.hasReachedSetPoint(resp.AmbientTemperature.ValueCelsius, setPoint.ValueCelsius, config, device.SetPointChangedTime) {
				// make sure we aren't re-raising the alert if it is already there
				if device.SetPointNotReachedId == "" {
					alert, err := a.alertAdminClient.CreateAlert(ctx, &gen.CreateAlertRequest{
						Alert: &gen.Alert{
							Id:          uuid.New().String(),
							CreateTime:  timestamppb.Now(),
							Description: "Measured temperature is not reaching set point within expected time",
							Severity:    gen.Alert_WARNING,
							Source:      device.Name,
						},
					})
					if err != nil {
						a.Logger.Error("failed to create alert", zap.Error(err))
					} else {
						device.SetPointNotReachedId = alert.Id
					}
				}
			} else {
				// this case ideally should not happen:
				// the idea of the auto is that it catches devices which don't reach the set point within a
				// generous time period, meaning something is likely wrong with the device and requires a manual check.
				// However, if we see that many alerts are being created and auto resolved here,
				// the parameters might be too strict.
				if device.SetPointNotReachedId != "" {
					_, err := a.alertAdminClient.ResolveAlert(ctx, &gen.ResolveAlertRequest{
						Alert: &gen.Alert{
							Id:          device.SetPointNotReachedId,
							ResolveTime: timestamppb.Now(),
						},
					})
					if err != nil {
						a.Logger.Error("failed to resolve alert", zap.Error(err))
					} else {
						device.SetPointNotReachedId = ""
					}
				}
			}
			device.PreviousSetPoint = proto.Float64(setPoint.ValueCelsius)
		} else {
			a.Logger.Error("failed to get air temperature or set point", zap.Error(err))
		}
	}
}

func (a *deviceMonitorAuto) runAirTemperatureMonitor(ctx context.Context, cfg config.Root, client traits.AirTemperatureApiClient) {

	if cfg.AirTempConfig.OkRtUpperBound != nil && cfg.AirTempConfig.OkRtLowerBound != nil {
		a.abnormalTemperatureCheck(ctx, client, cfg.AirTempConfig)
	}

	if cfg.AirTempConfig.Tolerance != nil && cfg.AirTempConfig.OkSettlingTime != nil {
		a.returnAirReachesSetPointCheck(ctx, client, cfg.AirTempConfig)
	}
}
