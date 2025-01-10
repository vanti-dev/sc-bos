package meter

import (
	"context"
	"path"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/meter/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("meter")
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
	clients   node.Clienter
	logger    *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := f.announcer.Replace(ctx)
	logger := f.logger

	var apiClient gen.MeterApiClient
	var infoClient gen.MeterInfoClient
	if len(cfg.Meters) > 0 || len(cfg.MeterGroups) > 0 {
		if err := f.clients.Client(&apiClient); err != nil {
			return err
		}
		if err := f.clients.Client(&infoClient); err != nil {
			return err
		}
	}
	announceGroup := func(name string, devices []string) {
		if len(devices) == 0 {
			return
		}

		group := &Group{
			apiClient:  apiClient,
			infoClient: infoClient,
			names:      devices,
			logger:     logger,
		}
		f.devices.Add(devices...)
		announce.Announce(name, node.HasTrait(meter.TraitName, node.WithClients(gen.WrapMeterApi(group), gen.WrapMeterInfo(group))))
	}

	announceGroup(cfg.Name, cfg.Meters)
	for name, meters := range cfg.MeterGroups {
		announceGroup(path.Join(cfg.Name, name), meters)
	}

	return nil
}
