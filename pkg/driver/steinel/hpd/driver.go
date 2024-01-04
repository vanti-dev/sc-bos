package hpd

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/steinel/hpd/config"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverName = "steinel-hpd"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	d := &Driver{
		announcer: services.Node,
	}
	d.Service = service.New(service.MonoApply(d.applyConfig))
	d.logger = services.Logger.Named(DriverName)
	return d
}

type Driver struct {
	*service.Service[config.Root]
	logger    *zap.Logger
	announcer node.Announcer

	client *Client

	airQualitySensor AirQualitySensor
	occupancy        Occupancy
	temperature      TemperatureSensor
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announcer := node.AnnounceContext(ctx, d.announcer)

	if cfg.Metadata != nil {
		announcer.Announce(cfg.Name, node.HasMetadata(cfg.Metadata))
	}

	p, err := cfg.LoadPassword()
	if err != nil {
		return err
	}

	if cfg.IpAddress == "" {
		return fmt.Errorf("ipAddress is required")
	}
	d.client = NewInsecureClient(cfg.IpAddress, p)

	d.airQualitySensor = NewAirQualitySensor(d.client, d.logger.Named("AirQuality").With(zap.String("ipAddress", cfg.IpAddress)), 0)
	announcer.Announce(cfg.Name, node.HasTrait(trait.AirQualitySensor, node.WithClients(airqualitysensor.WrapApi(&d.airQualitySensor))))

	d.occupancy = NewOccupancySensor(d.client, d.logger.Named("Occupancy").With(zap.String("ipAddress", cfg.IpAddress)), 0)
	announcer.Announce(cfg.Name, node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensor.WrapApi(&d.occupancy))))

	d.temperature = NewTemperatureSensor(d.client, d.logger.Named("Temperature").With(zap.String("ipAddress", cfg.IpAddress)), 0)
	announcer.Announce(cfg.Name, node.HasTrait(trait.AirTemperature, node.WithClients(airtemperature.WrapApi(&d.temperature))))

	return nil
}
