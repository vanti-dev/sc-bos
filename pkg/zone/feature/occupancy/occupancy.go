package occupancy

import (
	"context"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("occupancy")
	f := &feature{
		announcer: node.NewReplaceAnnouncer(services.Node),
		devices:   services.Devices,
		clients:   services.Node,
		logger:    services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig))
	return f
})

type feature struct {
	*service.Service[config.Root]
	announcer *node.ReplaceAnnouncer
	devices   *zone.Devices
	clients   node.ClientConner
	logger    *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := f.announcer.Replace(ctx)
	logger := f.logger

	if len(cfg.OccupancySensors) > 0 || len(cfg.EnterLeaveOccupancySensors) > 0 {
		group := &Group{logger: logger}
		conn := f.clients.ClientConn()

		if len(cfg.OccupancySensors) > 0 {
			group.client = traits.NewOccupancySensorApiClient(conn)
			group.names = cfg.OccupancySensors
		}
		if len(cfg.EnterLeaveOccupancySensors) > 0 {
			elServer := &enterLeave{
				model:  occupancysensorpb.NewModel(),
				client: traits.NewEnterLeaveSensorApiClient(conn),
				names:  cfg.EnterLeaveOccupancySensors,
				logger: logger,
			}
			group.clients = append(group.clients, occupancysensorpb.WrapApi(elServer))
		}

		f.devices.Add(cfg.OccupancySensors...)
		f.devices.Add(cfg.EnterLeaveOccupancySensors...)
		announce.Announce(cfg.Name, node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensorpb.WrapApi(group))))
	}

	return nil
}
