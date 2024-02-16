// Package lighting implements group lighting control for a zone.
// This package provides a single LightApi endpoint for controlling multiple underlying fixtures.
package lighting

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/lighting/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("lighting")
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

	if len(cfg.Lights) > 0 {
		var lightClient traits.LightApiClient
		if err := f.clients.Client(&lightClient); err != nil {
			return err
		}

		group := &Group{
			client:   lightClient,
			names:    cfg.Lights,
			readOnly: cfg.ReadOnlyLights,
			logger:   logger,
		}
		f.devices.Add(cfg.Lights...)
		announce.Announce(cfg.Name, node.HasTrait(trait.Light, node.WithClients(light.WrapApi(group))))
	}
	for key, lights := range cfg.LightGroups {
		var lightClient traits.LightApiClient
		if err := f.clients.Client(&lightClient); err != nil {
			return err
		}
		group := &Group{
			client:   lightClient,
			names:    lights,
			readOnly: cfg.ReadOnlyLights,
			logger:   logger.With(zap.String("lightGroup", key)),
		}
		name := fmt.Sprintf("%s/%s", cfg.Name, key)
		f.devices.Add(lights...)
		announce.Announce(name, node.HasTrait(trait.Light, node.WithClients(light.WrapApi(group))))
	}

	return nil
}
