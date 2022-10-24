package bacnet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/adapt"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/config"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/util/state"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/types/objecttype"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"log"
)

const DriverName = "bacnet"

type Driver struct {
	announcer node.Announcer // Any device we setup gets announced here
	logger    *zap.Logger

	status *state.Manager[driver.Status]

	config config.Root // The config that was used to setup the device into its current state

	client *gobacnet.Client // How we interact with bacnet systems

	configC chan config.Root
	stopCtx context.Context
	stop    context.CancelFunc
}

func Factory(ctx context.Context, services driver.Services, rawConfig json.RawMessage) (out driver.Driver, err error) {
	announcer := node.AnnounceWithNamePrefix("bacnet/", services.Node)
	d := &Driver{
		announcer: announcer,
		logger:    services.Logger.Named("bacnet"),
		status:    state.NewManager(driver.StatusActive),
		configC:   make(chan config.Root, 5),
	}

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

	err = d.Configure(rawConfig)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Driver) Name() string {
	return d.config.Name
}

func (d *Driver) Start(ctx context.Context) error {
	d.stopCtx, d.stop = context.WithCancel(context.Background())
	go func() {
		select {
		case cfg := <-d.configC:
			d.config = cfg
			err := d.applyConfig(cfg)
			if err != nil {
				d.logger.Error("failed to apply config update", zap.Error(err))
			}
		case <-d.stopCtx.Done():
			return
		}
	}()

	return nil
}

func (d *Driver) Configure(configData []byte) error {
	c, err := config.Read(bytes.NewReader(configData))
	if err != nil {
		return err
	}
	d.configC <- c
	return nil
}

func (d *Driver) Stop() error {
	d.stop()
	return nil
}

func (d *Driver) applyConfig(cfg config.Root) error {
	var err error
	if d.client == nil {
		client, err := gobacnet.NewClient("bridge100", 0)
		if err != nil {
			return err
		}
		d.client = client
		if address, err := client.LocalUDPAddress(); err == nil {
			d.logger.Debug("bacnet client configured", zap.Stringer("local", address))
		}
	}

	for _, device := range cfg.Devices {
		bacDevice, e := d.findDevice(device)
		if e != nil {
			err = multierr.Append(err, e)
			continue
		}

		prefix := fmt.Sprintf("device/%v/obj/", adapt.DeviceName(device))
		announcer := node.AnnounceWithNamePrefix(prefix, d.announcer)

		for _, object := range device.Objects {
			switch object.ID.Type {
			case objecttype.BinaryValue:
				impl := adapt.BinaryValue(d.client, bacDevice, object)
				impl.AnnounceSelf(announcer)
			default:
				log.Printf("Unsupported object type: %v", object.ID.Type)
			}
		}
	}

	return err
}
