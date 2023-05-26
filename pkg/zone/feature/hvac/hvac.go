package hvac

import (
	"context"
	"path"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/hvac/config"
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
	publish := func(name string, t config.Thermostat) error {
		var client traits.AirTemperatureApiClient
		if err := f.clients.Client(&client); err != nil {
			return err
		}

		group := &Group{
			client:   client,
			names:    t.Thermostats,
			readOnly: t.ReadOnlyThermostat,
			logger:   logger,
		}
		announce.Announce(name, node.HasTrait(trait.AirTemperature, node.WithClients(airtemperature.WrapApi(group))))
		return nil
	}

	if len(cfg.Thermostats) > 0 {
		if err := publish(cfg.Name, cfg.Thermostat); err != nil {
			return err
		}
	}

	for k, t := range cfg.ThermostatGroups {
		name := path.Join(cfg.Name, k)
		if err := publish(name, t); err != nil {
			return err
		}
	}

	return nil
}
