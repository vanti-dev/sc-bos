package electric

import (
	"context"
	"path"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/electric/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
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
	if len(cfg.Electrics) == 0 && len(cfg.ElectricGroups) == 0 {
		return nil
	}
	var client traits.ElectricApiClient
	if err := f.clients.Client(&client); err != nil {
		return err
	}

	announce := node.AnnounceContext(ctx, f.announce)
	logger := f.logger.With(zap.String("zone", cfg.Name))

	announceGroup := func(name string, devices []string) {
		if len(devices) == 0 {
			return
		}
		group := &Group{
			client: client,
			names:  devices,
			logger: logger,
		}
		f.devices.Add(devices...)
		announce.Announce(name, node.HasTrait(trait.Electric, node.WithClients(electric.WrapApi(group))))
	}

	announceGroup(cfg.Name, cfg.Electrics)
	for name, group := range cfg.ElectricGroups {
		announceGroup(path.Join(cfg.Name, name), group)
	}

	return nil
}
