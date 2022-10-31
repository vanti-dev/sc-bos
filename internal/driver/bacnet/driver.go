package bacnet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/adapt"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/config"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/known"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/merge"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/rpc"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/util/state"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/types/objecttype"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const DriverName = "bacnet"

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

	status *state.Manager[driver.Status] // allows us to implement Stateful
	config config.Root                   // The config that was used to setup the device into its current state
	client *gobacnet.Client              // How we interact with bacnet systems

	configC chan config.Root
	stopCtx context.Context
	stop    context.CancelFunc

	devices *known.Map
}

func NewDriver(services driver.Services) *Driver {
	announcer := node.AnnounceWithNamePrefix("bacnet/", services.Node)
	return &Driver{
		announcer: announcer,
		logger:    services.Logger.Named("bacnet"),
		status:    state.NewManager(driver.StatusInactive),
		devices:   known.NewMap(),
	}
}

// Factory creates a new Driver and calls Start then Configure on it.
func Factory(ctx context.Context, services driver.Services, rawConfig json.RawMessage) (out driver.Driver, err error) {
	d := NewDriver(services)

	err = d.Start(ctx)
	if err != nil {
		return nil, err
	}

	// now the driver is started, make sure we stop it if we happen to fail while setting up the driver.
	defer func() {
		if err != nil {
			_ = d.Stop()
		}
	}()

	// assume that people are using the ctx to stop the driver
	go func() {
		<-ctx.Done()
		_ = d.Stop()
	}()

	err = d.Configure(rawConfig)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Driver) Name() string {
	return d.config.Name
}

// Start makes this driver available to be configured.
// Call Stop when you're done with the driver to free up resources.
//
// Start must be called before Configure.
// Once started Configure and Stop may be called from any go routine.
func (d *Driver) Start(_ context.Context) error {
	// We implement a main loop pattern to avoid locks.
	// Methods like Configure and Stop both push messages onto channels
	// which we select on in a loop using a single go routine, avoiding any
	// locking issues (that we have to deal with ourselves).

	d.configC = make(chan config.Root, 5)
	d.stopCtx, d.stop = context.WithCancel(context.Background())
	go func() {
		d.status.Update(driver.StatusActive)
		defer d.status.Update(driver.StatusInactive)

		for {
			select {
			case cfg := <-d.configC:
				d.status.Update(driver.StatusLoading)
				err := d.applyConfig(cfg)
				d.status.Update(driver.StatusActive)
				if err != nil {
					d.logger.Error("failed to apply config update", zap.Error(err))
					continue
				}
				d.config = cfg
			case <-d.stopCtx.Done():
				return
			}
		}
	}()

	return nil
}

// Configure instructs the driver to setup and announce any devices found in configData.
// configData should be an encoded JSON object matching config.Root.
//
// Configure must not be called before Start, but once Started can be called concurrently.
func (d *Driver) Configure(configData []byte) error {
	if d.configC == nil {
		return errors.New("not started")
	}
	c, err := config.Read(bytes.NewReader(configData))
	if err != nil {
		return err
	}
	d.configC <- c
	return nil
}

// Stop stops the driver and releases resources.
// Stop races with Start before Start has completed, but can be called concurrently once started.
func (d *Driver) Stop() error {
	if d.stop == nil {
		// not started
		return nil
	}
	d.stop()
	return nil
}

func (d *Driver) WaitForStateChange(ctx context.Context, sourceState driver.Status) error {
	return d.status.WaitForStateChange(ctx, sourceState)
}

func (d *Driver) CurrentState() driver.Status {
	return d.status.CurrentState()
}

func (d *Driver) applyConfig(cfg config.Root) error {
	// todo: make this process atomic
	// todo: allow more than one config change, i.e. Undo announcements we need to remove on config change

	var err error
	// todo: support re-setting up the client if config changes
	if d.client == nil {
		client, err := gobacnet.NewClient(cfg.LocalInterface, int(cfg.LocalPort))
		if err != nil {
			return err
		}
		d.client = client
		if address, err := client.LocalUDPAddress(); err == nil {
			d.logger.Debug("bacnet client configured", zap.Stringer("local", address),
				zap.String("localInterface", cfg.LocalInterface), zap.Uint16("localPort", cfg.LocalPort))
		}
	}

	d.devices.Clear()

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

		announcer := node.AnnounceWithNamePrefix("device/", d.announcer)
		adapt.Device(deviceName, d.client, bacDevice).AnnounceSelf(announcer)

		prefix := fmt.Sprintf("device/%v/obj/", deviceName)
		announcer = node.AnnounceWithNamePrefix(prefix, d.announcer)

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

			impl, err := adapt.Object(d.client, bacDevice, co)
			if errors.Is(err, adapt.ErrNoDefault) {
				logger.Debug("No default adaptation trait for object")
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
	announcer := node.AnnounceWithNamePrefix("trait/", d.announcer)
	for _, trait := range cfg.Traits {
		logger := d.logger.With(zap.Stringer("trait", trait.Kind), zap.String("name", trait.Name))
		impl, err := merge.IntoTrait(d.client, d.devices, trait)
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
