package mode

import (
	"context"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/mode"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/mode/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("mode")
	f := &feature{
		announcer: node.NewReplaceAnnouncer(services.Node),
		devices:   services.Devices,
		clients:   services.Node,
		logger:    services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig), service.WithParser(config.ReadConfigBytes))
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
	if len(cfg.Modes) == 0 {
		return nil
	}
	announce := f.announcer.Replace(ctx)
	logger := f.logger

	var api traits.ModeApiClient
	if err := f.clients.Client(&api); err != nil {
		return err
	}
	f.devices.Add(cfg.AllDeviceNames()...)
	group := &Group{
		client: api,
		cfg:    cfg,
		logger: logger,
	}
	announce.Announce(cfg.Name, node.HasTrait(trait.Mode, node.WithClients(mode.WrapApi(group), mode.WrapInfo(group))))

	return nil
}
