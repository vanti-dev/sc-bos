package hpd3

import (
	"errors"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/motionsensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Device struct {
	client Client

	occupancyServer *occupancyServer
	motionServer    *motionServer
	airTempServer   *airTemperatureServer
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
		motionServer: &motionServer{
			client: client,
			logger: logger.With(zap.String("trait", string(trait.MotionSensor))),
		},
		airTempServer: &airTemperatureServer{
			client: client,
			logger: logger.With(zap.String("trait", string(trait.AirTemperature))),
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
		node.HasTrait(trait.MotionSensor, node.WithClients(
			motionsensor.WrapApi(d.motionServer),
			motionsensor.WrapSensorInfo(d.motionServer),
		)),
		node.HasTrait(trait.AirTemperature, node.WithClients(
			airtemperature.WrapApi(d.airTempServer),
			airtemperature.WrapInfo(d.airTempServer),
		)),
	}
}
