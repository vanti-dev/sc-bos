package hpd3

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverName = "steinel/hpd3"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	d := &Driver{
		services: services,
	}
	d.Service = service.New(
		service.MonoApply(d.applyConfig),
		service.WithParser(ParseDriverConfig),
	)
	return d
}

type Driver struct {
	*service.Service[DriverConfig]
	services driver.Services

	m            sync.Mutex
	unannouncers []node.Undo
}

func (d *Driver) applyConfig(_ context.Context, conf DriverConfig) error {
	d.m.Lock()
	defer d.m.Unlock()

	for _, unannounce := range d.unannouncers {
		unannounce()
	}
	d.unannouncers = nil

	for _, devConf := range conf.Devices {
		logger := d.services.Logger.With(zap.String("deviceName", devConf.Name))
		dev, err := newDevice(devConf, logger)
		if err != nil {
			d.services.Logger.Error("failed to initialise a hpd3 device",
				zap.String("name", devConf.Name),
				zap.Error(err))
			continue
		}

		unannounce := d.services.Node.Announce(devConf.Name, dev.features()...)
		d.unannouncers = append(d.unannouncers, unannounce)
	}

	return nil
}
