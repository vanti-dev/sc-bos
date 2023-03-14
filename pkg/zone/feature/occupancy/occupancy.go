package occupancy

import (
	"context"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/occupancy/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	f := &feature{
		announce: services.Node,
		clients:  services.Node,
		logger:   services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig))
	return f
})

type feature struct {
	*service.Service[config.Root]
	announce node.Announcer
	clients  node.Clienter
	logger   *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := node.AnnounceContext(ctx, f.announce)
	logger := f.logger.With(zap.String("zone", cfg.Name))

	if len(cfg.OccupancySensors) > 0 {
		var client traits.OccupancySensorApiClient
		if err := f.clients.Client(&client); err != nil {
			return err
		}

		group := &Group{
			client: client,
			names:  cfg.OccupancySensors,
			logger: logger,
		}
		announce.Announce(cfg.Name, node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensor.WrapApi(group))))
	}

	return nil
}
