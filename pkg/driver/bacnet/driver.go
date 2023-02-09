package bacnet

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/types/objecttype"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/adapt"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/merge"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/rpc"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

const DriverName = "bacnet"

var Factory driver.Factory = factory{}

type factory struct{}

func (_ factory) New(services driver.Services) service.Lifecycle {
	return NewDriver(services)
}

func (_ factory) AddSupport(supporter node.Supporter) {
	Register(supporter)
}

// Register makes sure this driver and its device apis are available in the given node.
func Register(supporter node.Supporter) {
	r := rpc.NewBacnetDriverServiceRouter()
	supporter.Support(
		node.Routing(r),
		node.Clients(rpc.WrapBacnetDriverService(r)),
	)
}

// Driver brings BACnet devices into Smart Core.
type Driver struct {
	announcer node.Announcer // Any device we setup gets announced here
	logger    *zap.Logger

	*service.Service[config.Root]
	client *gobacnet.Client // How we interact with bacnet systems

	devices *known.Map
}

func NewDriver(services driver.Services) *Driver {
	announcer := node.AnnounceWithNamePrefix("bacnet/", services.Node)
	d := &Driver{
		announcer: announcer,
		devices:   known.NewMap(),
		logger:    services.Logger.Named("bacnet"),
	}
	d.Service = service.New(service.MonoApply(d.applyConfig), service.WithParser(config.ReadBytes))
	return d
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	// AnnounceContext only makes sense if using MonoApply, which we are in NewDriver
	rootAnnouncer := node.AnnounceContext(ctx, d.announcer)
	go func() {
		<-ctx.Done() // ctx is cancelled on stop or before another applyConfig call, thanks to MonoApply above.
		// clear resources setup during the last applyConfig call
		d.devices.Clear()
		if d.client != nil {
			d.client.Close()
			d.client = nil
		}
	}()

	client, err := gobacnet.NewClient(cfg.LocalInterface, int(cfg.LocalPort))
	if err != nil {
		return err
	}
	client.Log.SetLevel(logrus.InfoLevel)
	d.client = client
	if address, err := client.LocalUDPAddress(); err == nil {
		d.logger.Debug("bacnet client configured", zap.Stringer("local", address),
			zap.String("localInterface", cfg.LocalInterface), zap.Uint16("localPort", cfg.LocalPort))
	}

	// setup all our devices and objects...
	for _, device := range cfg.Devices {
		logger := d.logger.With(zap.Uint32("deviceId", uint32(device.ID)))
		bacDevice, e := d.findDevice(device)
		if e != nil {
			err = multierr.Append(err, e)
			continue
		}

		deviceName := adapt.DeviceName(device)
		d.devices.StoreDevice(deviceName, bacDevice)

		announcer := node.AnnounceWithNamePrefix("device/", rootAnnouncer)
		adapt.Device(deviceName, d.client, bacDevice, d.devices).AnnounceSelf(announcer)

		prefix := fmt.Sprintf("device/%v/obj/", deviceName)
		announcer = node.AnnounceWithNamePrefix(prefix, rootAnnouncer)

		// Collect all the object that we will be announcing.
		// This will be a combination of configured objects and those we discover on the device.
		objects, e := d.fetchObjects(cfg, device, bacDevice)
		if e != nil {
			logger.Warn("Failed discovering objects", zap.Error(e))
		}

		for _, object := range objects {
			co, bo := object.co, object.bo
			logger := logger.With(zap.Stringer("object", co))
			// Device types are handled separately
			if bo.ID.Type == objecttype.Device {
				// We're assuming that devices in the wild follow the spec
				// which says each network device has exactly one bacnet device.
				// We check for this explicitly to make sure our assumptions hold
				if bo.ID != bacDevice.ID {
					logger.Error("BACnet device with multiple advertised devices!")
				}
				continue
			}

			// no error, we added the device before we entered the loop so it should exist
			_ = d.devices.StoreObject(bacDevice, adapt.ObjectName(co), *bo)

			impl, err := adapt.Object(d.client, bacDevice, co, logger.Named(co.Name))
			if errors.Is(err, adapt.ErrNoDefault) {
				// logger.Debug("No default adaptation trait for object")
				continue
			}
			if errors.Is(err, adapt.ErrNoAdaptation) {
				logger.Error("No adaptation from object to trait", zap.Stringer("trait", co.Trait))
				continue
			}
			if err != nil {
				logger.Error("Error adapting object", zap.Error(err))
				continue
			}
			impl.AnnounceSelf(announcer)
		}
	}

	// Combine objects together into traits...
	announcer := node.AnnounceWithNamePrefix("trait/", rootAnnouncer)
	for _, trait := range cfg.Traits {
		logger := d.logger.With(zap.Stringer("trait", trait.Kind), zap.String("name", trait.Name))
		impl, err := merge.IntoTrait(d.client, d.devices, trait, logger)
		if errors.Is(err, merge.ErrTraitNotSupported) {
			logger.Error("Cannot combine into trait, not supported")
			continue
		}
		if err != nil {
			logger.Error("Cannot combine into trait", zap.Error(err))
			continue
		}
		impl.AnnounceSelf(announcer)
	}

	return err
}
