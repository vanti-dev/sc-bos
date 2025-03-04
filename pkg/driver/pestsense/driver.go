package pestsense

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/pestsense/config"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverName = "pestsense"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	d := &Driver{
		announcer: node.NewReplaceAnnouncer(services.Node),
	}
	d.Service = service.New(service.MonoApply(d.applyConfig))
	d.logger = services.Logger.Named(DriverName)
	return d
}

type Driver struct {
	*service.Service[config.Root]
	logger    *zap.Logger
	announcer *node.ReplaceAnnouncer

	devices map[string]*PestSensor

	client mqtt.Client
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announcer := d.announcer.Replace(ctx)

	d.devices = make(map[string]*PestSensor)
	// Add devices and apis
	for _, device := range cfg.Devices {
		sensor := NewPestSensor(device.Id)
		d.devices[device.Id] = sensor

		announcer.Announce(device.Name, node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensorpb.WrapApi(sensor))))
	}

	// Connect to MQTT
	var err error
	d.client, err = newMqttClient(cfg)
	if err != nil {
		return err
	}
	connected := d.client.Connect()
	connected.Wait()
	if connected.Error() != nil {
		return connected.Error()
	}
	d.logger.Debug("connected")

	var responseHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		go handleResponse(msg.Payload(), d.devices, d.logger)
	}

	token := d.client.Subscribe(cfg.Broker.Topic, 0, responseHandler)
	token.Wait()

	return nil
}

func newMqttClient(cfg config.Root) (mqtt.Client, error) {
	options, err := cfg.Broker.ClientOptions()
	if err != nil {
		return nil, err
	}
	return mqtt.NewClient(options), nil
}
