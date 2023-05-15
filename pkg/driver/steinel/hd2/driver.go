package hd2

import (
	"context"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"

	"github.com/vanti-dev/sc-bos-drivers/pkg/driver/steinel/hd2/config"
)

const DriverName = "steinel-hd2"

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
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announcer := node.AnnounceContext(ctx, d.announcer)

	d.client = NewInsecureClient(cfg.IpAddress)

	d.airQualitySensor = NewAirQualitySensor(d.client, d.logger.Named("AirQuality"), 0)
	announcer.Announce(cfg.Name+"/airQuality", node.HasTrait(trait.AirQualitySensor, node.WithClients(airqualitysensor.WrapApi(&d.airQualitySensor))))

	d.occupancy = Occupancy{client: d.client}
	announcer.Announce(cfg.Name+"/occupancy", node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensor.WrapApi(&d.occupancy))))

	return nil
}
