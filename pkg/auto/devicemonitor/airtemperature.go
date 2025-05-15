package devicemonitor

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto/devicemonitor/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func (d *deviceMonitorAuto) checkReturnTemperaturesAreNormal(ctx context.Context, client traits.AirTemperatureApiClient, config *config.AirTempConfig) bool {
	for _, device := range config.Devices {
		resp, err := client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{
			Name: device.Name,
		})
		if err != nil {
			d.Logger.Error("failed to get air temperature", zap.Error(err))
			return false
		}
		if resp.AmbientTemperature != nil {
			if resp.AmbientTemperature.ValueCelsius < *config.OkRtLowerBound {
				alert, err := d.alertAdminClient.CreateAlert(ctx, &gen.CreateAlertRequest{
					Alert: &gen.Alert{
						Id:          uuid.New().String(),
						CreateTime:  timestamppb.Now(),
						Description: "Ambient temperature is abnormally low",
						Severity:    gen.Alert_WARNING,
						Source:      device.Name,
					},
				})
				if err != nil {
					d.Logger.Error("failed to create alert", zap.Error(err))
				} else {
					device.AbnormalLowId = alert.Id
				}
			} else {
				if device.AbnormalLowId != "" {
					_, err := d.alertAdminClient.ResolveAlert(ctx, &gen.ResolveAlertRequest{
						Alert: &gen.Alert{
							Id:          device.AbnormalLowId,
							ResolveTime: timestamppb.Now(),
						},
					})
					if err != nil {
						d.Logger.Error("failed to resolve alert", zap.Error(err))
					} else {
						device.AbnormalLowId = ""
					}
				}
			}
			if resp.AmbientTemperature.ValueCelsius > *config.OkRtUpperBound {
				alert, err := d.alertAdminClient.CreateAlert(ctx, &gen.CreateAlertRequest{
					Alert: &gen.Alert{
						Id:          uuid.New().String(),
						CreateTime:  timestamppb.Now(),
						Description: "Ambient temperature is abnormally high",
						Severity:    gen.Alert_WARNING,
						Source:      device.Name,
					},
				})
				if err != nil {
					d.Logger.Error("failed to create alert", zap.Error(err))
				} else {
					device.AbnormalHighId = alert.Id
				}
			} else {
				if device.AbnormalHighId != "" {
					_, err := d.alertAdminClient.ResolveAlert(ctx, &gen.ResolveAlertRequest{
						Alert: &gen.Alert{
							Id:          device.AbnormalHighId,
							ResolveTime: timestamppb.Now(),
						},
					})
					if err != nil {
						d.Logger.Error("failed to resolve alert", zap.Error(err))
					} else {
						device.AbnormalHighId = ""
					}
				}
			}
		}
	}
	return true
}

func (d *deviceMonitorAuto) runAirTemperatureMonitor(ctx context.Context, cfg config.Root, client traits.AirTemperatureApiClient) {
	d.checkReturnTemperaturesAreNormal(ctx, client, cfg.AirTempConfig)
}
