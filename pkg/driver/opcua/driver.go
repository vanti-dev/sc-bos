package opcua

import (
	"context"
	"fmt"
	"time"

	"github.com/gopcua/opcua"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/smart-core-os/sc-bos/pkg/driver"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/meter"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/transport"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/udmipb"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-bos/pkg/task/service"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/electricpb"

	"github.com/smart-core-os/sc-bos/pkg/driver/opcua/config"
)

const DriverName = "opcua"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	logger := services.Logger.Named(DriverName)

	d := &Driver{
		announcer: node.NewReplaceAnnouncer(services.Node),
	}
	d.Service = service.New(
		service.MonoApply(d.applyConfig),
		service.WithParser(config.ReadBytes),
		service.WithRetry[config.Root](service.RetryWithLogger(func(logContext service.RetryContext) {
			logContext.LogTo("applyConfig", logger)
		}), service.RetryWithMinDelay(5*time.Second), service.RetryWithInitialDelay(5*time.Second)),
	)
	d.logger = services.Logger.Named(DriverName)
	return d
}

type Driver struct {
	*service.Service[config.Root]
	logger    *zap.Logger
	announcer *node.ReplaceAnnouncer
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {

	a := d.announcer.Replace(ctx)

	opcClient, err := opcua.NewClient(cfg.Conn.Endpoint)
	if err != nil {
		d.logger.Warn("NewClient error", zap.String("error", err.Error()))
		return err
	}

	err = opcClient.Connect(ctx)
	if err != nil {
		d.logger.Warn("Connect error", zap.String("error", err.Error()))
		return err
	}

	client := NewClient(opcClient, d.logger, cfg.Conn.SubscriptionInterval.Duration, cfg.Conn.ClientId)

	if cfg.Meta != nil {
		a.Announce(cfg.Name, node.HasMetadata(cfg.Meta))
	}

	grp, ctx := errgroup.WithContext(ctx)
	var errs error
	for _, dev := range cfg.Devices {
		var allFeatures []node.Feature
		opcDev := NewDevice(&dev, d.logger, client)

		for _, t := range dev.Traits {
			switch t.Kind {
			case meter.TraitName:
				opcDev.meter, err = newMeter(dev.Name, t, d.logger)
				if err != nil {
					errs = fmt.Errorf("failed to add trait for device %s: %w", dev.Name, err)
				} else {
					allFeatures = append(allFeatures, node.HasTrait(meter.TraitName, node.WithClients(gen.WrapMeterApi(opcDev.meter), gen.WrapMeterInfo(opcDev.meter))))
				}
			case transport.TraitName:
				opcDev.transport, err = newTransport(dev.Name, t, d.logger)
				if err != nil {
					errs = fmt.Errorf("failed to add trait for device %s: %w", dev.Name, err)
				} else {
					allFeatures = append(allFeatures, node.HasTrait(transport.TraitName, node.WithClients(gen.WrapTransportApi(opcDev.transport), gen.WrapTransportInfo(opcDev.transport))))
				}
			case udmipb.TraitName:
				opcDev.udmi, err = newUdmi(dev.Name, t, d.logger)
				if err != nil {
					errs = fmt.Errorf("failed to add trait for device %s: %w", dev.Name, err)
				} else {
					allFeatures = append(allFeatures, node.HasTrait(udmipb.TraitName, node.WithClients(gen.WrapUdmiService(opcDev.udmi))))
				}
			case trait.Electric:
				opcDev.electric, err = newElectric(dev.Name, t, d.logger)
				if err != nil {
					errs = fmt.Errorf("failed to add trait for device %s: %w", dev.Name, err)
				} else {
					allFeatures = append(allFeatures, node.HasTrait(trait.Electric, node.WithClients(electricpb.WrapApi(opcDev.electric))))
				}
			default:
				d.logger.Error("unknown trait", zap.String("trait", t.Name))
			}
		}

		if dev.Meta != nil {
			allFeatures = append(allFeatures, node.HasMetadata(dev.Meta))
		}

		if errs != nil {
			d.logger.Error("errors encountered whilst loading driver", zap.String("device", dev.Name), zap.Error(errs))
		}

		a.Announce(dev.Name, allFeatures...)
		grp.Go(func() error {
			return opcDev.run(ctx)
		})
	}

	go func() {
		err := grp.Wait()
		d.logger.Error("run error", zap.String("error", err.Error()))
		_ = opcClient.Close(ctx)
	}()
	return nil
}
