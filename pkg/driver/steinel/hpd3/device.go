package hpd3

import (
	"errors"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Device struct {
	client Client

	occupancyServer *occupancyServer
}

func newDevice(conf DeviceConfig, logger *zap.Logger) (*Device, error) {
	if conf.Host == "" {
		return nil, errors.New("host not specified")
	}
	password, err := conf.Password.Read()
	if err != nil {
		return nil, err
	}
	client := &HTTPClient{
		Password: password,
		Host:     conf.Host,
	}
	dev := &Device{
		client: client,
		occupancyServer: &occupancyServer{
			client: client,
			logger: logger.With(zap.String("trait", string(trait.OccupancySensor))),
		},
	}
	return dev, nil
}

func (d *Device) features() []node.Feature {
	return []node.Feature{
		node.HasTrait(trait.OccupancySensor, node.WithClients(
			occupancysensor.WrapApi(d.occupancyServer),
			occupancysensor.WrapInfo(d.occupancyServer),
		)),
	}
}
