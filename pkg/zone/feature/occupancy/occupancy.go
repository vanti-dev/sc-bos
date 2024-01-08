package occupancy

import (
	"context"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("occupancy")
	f := &feature{
		announce: services.Node,
		devices:  services.Devices,
		clients:  services.Node,
		logger:   services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig))
	return f
})

type feature struct {
	*service.Service[config.Root]
	announce node.Announcer
	devices  *zone.Devices
	clients  node.Clienter
	logger   *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := node.AnnounceContext(ctx, f.announce)
	logger := f.logger

	if len(cfg.OccupancySensors) > 0 || len(cfg.EnterLeaveOccupancySensors) > 0 {
		group := &Group{logger: logger}

		if len(cfg.OccupancySensors) > 0 {
			if err := f.clients.Client(&group.client); err != nil {
				return err
			}
			group.names = cfg.OccupancySensors
		}
		if len(cfg.EnterLeaveOccupancySensors) > 0 {
			elServer := &enterLeave{
				model: occupancysensor.NewModel(),
				names: cfg.EnterLeaveOccupancySensors,
			}
			if err := f.clients.Client(&elServer.client); err != nil {
				return err
			}
			group.clients = append(group.clients, occupancysensor.WrapApi(elServer))
		}

		f.devices.Add(cfg.OccupancySensors...)
		f.devices.Add(cfg.EnterLeaveOccupancySensors...)
		announce.Announce(cfg.Name, node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensor.WrapApi(group))))
	}

	return nil
}
