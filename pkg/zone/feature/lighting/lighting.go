// Package lighting implements group lighting control for a zone.
// This package provides a single LightApi endpoint for controlling multiple underlying fixtures.
package lighting

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/lightpb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/lighting/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	services.Logger = services.Logger.Named("lighting")
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

	announceGroup := func(name string, lights []string, logger *zap.Logger) error {
		var apiClient traits.LightApiClient
		if err := f.clients.Client(&apiClient); err != nil {
			return err
		}
		var infoClient traits.LightInfoClient
		if err := f.clients.Client(&infoClient); err != nil {
			// we don't support info api, Group can handle this so just continue
		}
		group := &Group{
			client:   apiClient,
			info:     infoClient,
			names:    lights,
			readOnly: cfg.ReadOnlyLights,
			logger:   logger,
		}
		f.devices.Add(lights...)
		announce.Announce(name, node.HasTrait(trait.Light, node.WithClients(lightpb.WrapApi(group), lightpb.WrapInfo(group))))
		return nil
	}

	if len(cfg.Lights) > 0 {
		if err := announceGroup(cfg.Name, cfg.Lights, logger); err != nil {
			return err
		}
	}
	for key, lights := range cfg.LightGroups {
		name := fmt.Sprintf("%s/%s", cfg.Name, key)
		logger := logger.With(zap.String("lightGroup", key))
		if err := announceGroup(name, lights, logger); err != nil {
			return err
		}
	}

	return nil
}
