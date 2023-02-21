package priority

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/priority/config"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/prioritypb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"
)

const DriverName = "priority"

var Factory driver.Factory = factory{}

type factory struct{}

func (_ factory) New(services driver.Services) service.Lifecycle {
	return NewDriver(services)
}

func NewDriver(services driver.Services) *Driver {
	d := &Driver{
		announcer: services.Node,
		clients:   services.Node,
		logger:    services.Logger.Named(DriverName),
	}
	d.Service = service.New(service.MonoApply(d.applyConfig))
	return d
}

type Driver struct {
	*service.Service[config.Root]

	logger    *zap.Logger
	announcer node.Announcer
	clients   node.Clienter
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := node.AnnounceContext(ctx, d.announcer)

	for _, deviceMd := range cfg.Devices {
		for _, traitMd := range deviceMd.Traits {
			switch trait.Name(traitMd.Name) {
			case trait.Light:
				var client traits.LightApiClient
				if err := d.clients.Client(&client); err != nil {
					return err
				}
				opts := cfg.Names.Options(
					prioritypb.WithLogger(d.logger),
					prioritypb.WithMetadata(deviceMd.Metadata),
				)
				prioritypb.NewLightPriority(client, deviceMd.Name, opts...).AnnounceSelf(announce)
			default:
				d.logger.Warn("Unsupported priority trait", zap.String("trait", traitMd.Name))
			}
		}
	}
	return nil
}
