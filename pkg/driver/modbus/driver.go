package modbus

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/emergencypb"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystoragepb"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/modbus/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const (
	DriverName = "modbus"

	defaultTimeout      = 5 * time.Second
	defaultPollInterval = 30 * time.Second
)

type factory struct{}

var Factory driver.Factory = factory{}

func (f factory) New(services driver.Services) service.Lifecycle {
	d := &Driver{
		logger: services.Logger.Named(DriverName),
		node:   services.Node,
	}

	d.Service = service.New(
		service.MonoApply(d.applyConfig),
		service.WithRetry[config.Root](service.RetryWithLogger(func(logCtx service.RetryContext) {
			logCtx.LogTo("applyConfig", d.logger)
		})),
		service.WithOnStop[config.Root](d.Clean),
	)

	return d
}

type Driver struct {
	*service.Service[config.Root]
	node   node.Announcer
	logger *zap.Logger

	clients sync.Map
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	announcer := node.NewReplaceAnnouncer(d.node).Replace(ctx)

	for _, device := range cfg.Devices {
		var client *Client
		if device.TcpHandle != nil {
			handler := NewTCPClientHandler(
				fmt.Sprintf("%s:%d", device.TcpHandle.Address, device.TcpHandle.Port),
				WithTCPTimeout(device.TcpHandle.Timeout.Or(defaultTimeout)),
				WithTCPSlaveId(device.TcpHandle.SlaveId),
				WithTCPLogger(d.logger.Named(device.Name)),
			)
			client = NewClient(ctx, handler)
		}

		if device.RTUHandle != nil {
			handler := NewRTUClientHandler(
				device.RTUHandle.Address,
				WithRTUBaudRate(device.RTUHandle.BaudRate),
				WithRTUDataBits(device.RTUHandle.DataBits),
				WithRTUStopBits(device.RTUHandle.StopBits),
				WithRTUParity(device.RTUHandle.Parity),
				WithRTUTimeout(device.RTUHandle.Timeout.Or(defaultTimeout)),
				WithRTUSlaveId(device.RTUHandle.SlaveId),
				WithRTULogger(d.logger.Named(device.Name)),
			)
			client = NewClient(ctx, handler)
		}

		if client == nil {
			continue
		}

		if err := client.Connect(); err != nil {
			d.logger.Error("failed to connect modbus client", zap.Error(err), zap.String("device", device.Name))
			return err
		}

		d.clients.Store(device.Name, client)

		for _, deviceTrait := range device.Traits {
			if deviceTrait.Name == trait.EnergyStorage {
				fuelName, err := url.JoinPath(cfg.ScNamePrefix, device.Name, fuel)
				if err != nil {
					return err
				}
				announcer.Announce(fuelName,
					node.HasTrait(trait.EnergyStorage, node.WithClients(energystoragepb.WrapApi(client))),
					node.HasMetadata(device.Metadata),
				)

				// poll the client
				if err := client.Poll(*deviceTrait.PDU, fuel, deviceTrait.Address, deviceTrait.Quantity, deviceTrait.ScaleFactor, deviceTrait.PollInterval.Or(defaultPollInterval)); err != nil {
					d.logger.Error("connecting client", zap.String("device", device.Name), zap.Error(err))
					return err
				}
			}

			if deviceTrait.Name == trait.Emergency {
				faultsName, err := url.JoinPath(cfg.ScNamePrefix, device.Name, faults)
				if err != nil {
					return err
				}
				announcer.Announce(faultsName,
					node.HasTrait(trait.Emergency, node.WithClients(emergencypb.WrapApi(client))),
					node.HasMetadata(device.Metadata),
				)

				// poll the client
				if err := client.Poll(*deviceTrait.PDU, faults, deviceTrait.Address, deviceTrait.Quantity, deviceTrait.ScaleFactor, deviceTrait.PollInterval.Or(defaultPollInterval)); err != nil {
					d.logger.Error("connecting client", zap.String("device", device.Name), zap.Error(err))
					return err
				}
			}

			if deviceTrait.Name == statuspb.TraitName {
				monitorName, err := url.JoinPath(cfg.ScNamePrefix, device.Name, monitor)
				if err != nil {
					return err
				}
				announcer.Announce(monitorName,
					node.HasTrait(statuspb.TraitName, node.WithClients(gen.WrapStatusApi(client))),
					node.HasMetadata(device.Metadata),
				)

				// poll the client
				if err := client.Poll(*deviceTrait.PDU, monitor, deviceTrait.Address, deviceTrait.Quantity, deviceTrait.ScaleFactor, deviceTrait.PollInterval.Or(defaultPollInterval)); err != nil {
					d.logger.Error("connecting client", zap.String("device", device.Name), zap.Error(err))
					return err
				}
			}
		}
	}

	return nil
}

func (d *Driver) Clean() {
	d.clients.Range(func(key, value any) bool {
		client := value.(*Client)
		if err := client.Close(); err != nil {
			d.logger.Error("failed to close client", zap.Error(err))
			return true
		}
		// remove the client
		d.clients.Delete(key)
		return true
	})
}
