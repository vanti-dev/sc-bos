package xovis

import (
	"context"
	"errors"
	"sync"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/enterleavesensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

const DriverName = "xovis"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) task.Starter {
	d := &Driver{
		Services: services,
	}
	d.Lifecycle = task.NewLifecycle(d.applyConfig)
	return d
}

type Driver struct {
	driver.Services
	*task.Lifecycle[DriverConfig]

	m                 sync.Mutex
	config            DriverConfig
	client            *Client
	unannounceDevices []node.Undo
}

func (d *Driver) applyConfig(_ context.Context, conf DriverConfig) error {
	d.m.Lock()
	defer d.m.Unlock()

	// A route can't be removed from an HTTP ServeMux, so if it's been changed or removed then we can't support the
	// new configuration. This is likely to be rare in practice. Adding a route is fine.
	var oldWebhook, newWebhook string
	if d.config.DataPush != nil {
		oldWebhook = d.config.DataPush.WebhookPath
	}
	if conf.DataPush != nil {
		newWebhook = d.config.DataPush.WebhookPath
	}
	if oldWebhook != "" && newWebhook != oldWebhook {
		return errors.New("can't change webhook path once service is running")
	}

	// create a new client to communicate with the Xovis sensor
	d.client = NewInsecureClient(conf.Host, conf.Username, conf.Password)
	// unannounce any devices left over from a previous configuration
	for _, unannounce := range d.unannounceDevices {
		unannounce()
	}
	d.unannounceDevices = nil
	// annouce new devices
	for _, dev := range conf.Devices {
		var features []node.Feature
		if dev.Occupancy != nil {
			features = append(features, node.HasTrait(trait.OccupancySensor,
				node.WithClients(occupancysensor.WrapApi(&occupancyServer{
					client:      d.client,
					multiSensor: conf.MultiSensor,
					logicID:     dev.Occupancy.ID,
				})),
			))
		}
		if dev.EnterLeave != nil {
			features = append(features, node.HasTrait(trait.OccupancySensor,
				node.WithClients(enterleavesensor.WrapApi(&enterLeaveServer{
					client:      d.client,
					logicID:     dev.EnterLeave.ID,
					multiSensor: conf.MultiSensor,
				})),
			))
		}

		d.unannounceDevices = append(d.unannounceDevices, d.Node.Announce(dev.Name, features...))
	}
	// register data push webhook
	if dp := conf.DataPush; dp != nil && dp.WebhookPath != "" {

	}

	d.config = conf
}
