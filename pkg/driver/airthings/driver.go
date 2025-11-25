// Package airthings integrates AirThings devices into Smart Core.
// AirThings manufacture sensors, typically air quality sensors, that connect directly to their cloud.
// The cloud API provides access to various information about the devices and sites.
// The primary api used by this driver is the "location latest samples" api.
// See https://developer.airthings.com/api-docs#tag/Locations/paths/~1v1~1locations~1%7BlocationId%7D~1latest-samples/get
//
// The driver pulls all data into a local model, then translates that local model into Smart Core traits.
// Package [local] defines the local model.
// The code that pulls the data from the AirThings cloud API into local is in [client.go].
// The supported traits are defined in [traits.go].
package airthings

import (
	"context"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/smart-core-os/sc-bos/pkg/block"
	"github.com/smart-core-os/sc-bos/pkg/driver"
	"github.com/smart-core-os/sc-bos/pkg/driver/airthings/api"
	"github.com/smart-core-os/sc-bos/pkg/driver/airthings/local"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/statuspb"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-bos/pkg/task/service"
)

const DriverName = "airthings"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	services.Logger = services.Logger.Named(DriverName)
	d := &Driver{
		Services:  services,
		announcer: node.NewReplaceAnnouncer(services.Node),
	}
	d.Service = service.New(service.MonoApply(d.applyConfig))
	return d
}

func (_ factory) ConfigBlocks() []block.Block {
	return Blocks
}

type Driver struct {
	*service.Service[Config]
	driver.Services
	announcer *node.ReplaceAnnouncer

	cfg    Config
	client *http.Client

	listLocationsOnce sync.Once
	locationsErr      error
	locations         api.GetLocationsResponse
}

func (d *Driver) applyConfig(ctx context.Context, cfg Config) error {
	announcer := d.announcer.Replace(ctx)
	d.listLocationsOnce = sync.Once{}
	d.cfg = cfg

	ccConfig, err := d.cfg.Auth.ClientCredentialsConfig()
	if err != nil {
		return err
	}
	d.client = ccConfig.Client(ctx)

	status := statuspb.NewMap(announcer)

	grp, ctx := errgroup.WithContext(ctx)
	for _, location := range cfg.Locations {
		location := location
		ll := local.NewLocation()
		grp.Go(func() error {
			return d.pollLocationLatestSamples(ctx, location, ll)
		})

		for _, device := range location.Devices {
			n := device.Name
			announcer.Announce(n, node.HasMetadata(device.Metadata))
			status.UpdateProblem(n, &gen.StatusLog_Problem{
				Level:       gen.StatusLog_NOMINAL,
				Description: "Device configured successfully",
				Name:        n + ":setup",
			})
			err = d.announceDevice(ctx, announcer, device, ll, status)
			if err != nil {
				return err // failure of configuration, not runtime
			}
		}
	}
	return nil
}
