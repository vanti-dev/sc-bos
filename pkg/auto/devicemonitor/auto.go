// Package devicemonitor is used to monitor the measurements and statuses of devices and look for any abnormal behaviour.
// If abnormal behaviour is detected, it will raise a status notification.
package devicemonitor

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/devicemonitor/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const AutoName = "devicemonitor"

var Factory auto.Factory = factory{}

type factory struct{}

type deviceMonitorAuto struct {
	*service.Service[config.Root]
	auto.Services

	alertAdminClient gen.AlertAdminApiClient
	announcer        *node.ReplaceAnnouncer
}

func (f factory) New(services auto.Services) service.Lifecycle {
	a := &deviceMonitorAuto{Services: services}
	a.Service = service.New(service.MonoApply(a.applyConfig), service.WithParser(config.ReadBytes))
	a.Logger = a.Logger.Named(AutoName)
	return a
}

func (d *deviceMonitorAuto) applyConfig(ctx context.Context, cfg config.Root) error {

	now := cfg.Now
	if now == nil {
		now = time.Now
	}

	d.announcer = node.NewReplaceAnnouncer(d.Node)
	// todo check we can actually do it like this
	d.alertAdminClient = gen.NewAlertAdminApiClient(d.Node.ClientConn())
	if cfg.AirTempConfig != nil && len(cfg.AirTempConfig.Devices) > 0 {
		fcuClient := traits.NewAirTemperatureApiClient(d.Node.ClientConn())

		atConfig := cfg.AirTempConfig
		go func() {
			t := now()
			for {
				next := atConfig.MonitorSchedule.Next(t)
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Until(next)):
					t = next
					d.runAirTemperatureMonitor(ctx, cfg, fcuClient)
				}
			}
		}()
	}

	return nil
}
