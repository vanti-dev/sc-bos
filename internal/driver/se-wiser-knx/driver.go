package se_wiser_knx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	scLight "github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/vanti-dev/inf-sc-bos/internal/trait/light"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverName = "se-wiser-knx"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	services.Logger = services.Logger.Named(DriverName)
	d := &Driver{
		Services: services,
	}
	d.Service = service.New(
		service.MonoApply(d.applyConfig),
		service.WithParser(ParseConfig),
	)
	d.logger = services.Logger.Named(DriverName)
	d.devicesByName = make(map[string]*light.Model)
	d.devicesByAddress = make(map[string]*light.Model)
	return d
}

type Driver struct {
	*service.Service[Config]
	driver.Services
	logger *zap.Logger

	cfg              Config
	client           *Client
	devicesByName    map[string]*light.Model
	devicesByAddress map[string]*light.Model
}

func (d *Driver) applyConfig(ctx context.Context, cfg Config) error {
	announcer := node.AnnounceContext(ctx, d.Node)
	d.cfg = cfg

	// create a new client to communicate with the Wiser controller
	pass, err := cfg.LoadPassword()
	if err != nil {
		return err
	}
	d.client = NewInsecureClient(cfg.Host, cfg.Username, pass)

	for _, dev := range cfg.Devices {
		if dev.Metadata != nil {
			announcer.Announce(dev.Name, node.HasMetadata(dev.Metadata))
		}

		l := light.NewModel(&traits.Brightness{})
		c := scLight.WrapApi(lightServer{LightApiServer: light.NewModelServer(l), client: d.client, device: &dev, logger: d.logger.Named(dev.Name)})
		announcer.Announce(dev.Name, node.HasTrait(trait.Light, node.WithClients(c)))

		d.devicesByName[dev.Name] = l
		d.devicesByAddress[dev.Address] = l
	}

	go d.poll(ctx)

	return nil
}

func (d *Driver) poll(ctx context.Context) {
	ticker := time.NewTicker(d.cfg.Poll)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// query all objects
			objects, err := QueryObjects(d.client)
			if err != nil {
				d.logger.Error("Error querying objects", zap.Error(err))
			}
			// loop through response objects
			for _, obj := range objects {
				// if matching device address
				if dev, ok := d.devicesByAddress[obj.Address]; ok {
					// update model brightness value for that device
					lvl, err := strconv.ParseFloat(obj.Data, 32)
					if err != nil {
						d.logger.Error("Error parsing brightness", zap.Error(err))
					}
					b := &traits.Brightness{
						LevelPercent: float32(lvl),
					}
					_, err = dev.UpdateBrightness(b)
					if err != nil {
						d.logger.Error("Error updating brightness", zap.Error(err))
					}
				}
			}
		}
	}
}

type lightServer struct {
	traits.LightApiServer
	client *Client
	device *Device
	logger *zap.Logger
}

func (l lightServer) UpdateBrightness(ctx context.Context, req *traits.UpdateBrightnessRequest) (*traits.Brightness, error) {
	err := SetValue(l.client, l.device.Address, fmt.Sprintf("%f", req.Brightness.LevelPercent))
	if err != nil {
		return nil, err
	}

	brightness, err := l.LightApiServer.UpdateBrightness(ctx, req)
	if err != nil {
		return nil, err
	}
	return brightness, nil
}
