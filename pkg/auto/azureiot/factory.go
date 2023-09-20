package azureiot

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/azureiot/auth"
	"github.com/vanti-dev/sc-bos/pkg/auto/azureiot/iothub"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverType = "azureiot"

const minPollInterval = 5 * time.Second

var Factory auto.Factory = factory{}

type factory struct{}

func (f factory) New(services auto.Services) service.Lifecycle {
	a := &Auto{services: services}
	a.Service = service.New(service.MonoApply(a.applyConfig))
	return a
}

type Auto struct {
	*service.Service[Config]
	services auto.Services
}

func (a *Auto) applyConfig(ctx context.Context, cfg Config) error {
	var source gen.PointApiClient
	err := a.services.Node.Client(&source)
	if err != nil {
		return fmt.Errorf("PointApiClient unavailable: %w", err)
	}

	if cfg.PollInterval.Duration < minPollInterval {
		return fmt.Errorf("pollInterval must be at least %v", minPollInterval)
	}
	if len(cfg.Devices) == 0 {
		a.services.Logger.Warn("no devices configured; no polling will happen")
		return nil
	}

	// load the group key from string or disk, but only if a device will need it
	var needsGroupKey bool // group key is required if a device lacks a connection string
	for _, deviceCfg := range cfg.Devices {
		if deviceCfg.ConnectionString == "" {
			needsGroupKey = true
		}
	}
	var groupKey auth.SASKey
	if needsGroupKey {
		var err error
		groupKey, err = loadGroupKey(cfg)
		if err != nil {
			return fmt.Errorf("failed to load group key: %w", err)
		}

		// if the group key is used, then an ID scope is also mandatory
		if cfg.IDScope == "" {
			return fmt.Errorf("id scope is required when using group keys")
		}
	}

	// construct a poller for each device registered
	pollers := make([]*poller, 0, len(cfg.Devices))
	for _, deviceConf := range cfg.Devices {
		logger := a.services.Logger.With(zap.String("device", deviceConf.Name))
		devDialler, err := diallerFromConfig(deviceConf, cfg.IDScope, groupKey)
		if err != nil {
			logger.Error("device poller not initialised due to invalid configuration", zap.Error(err))
			continue
		}

		d := &poller{
			source:  source,
			dialler: devDialler,
			logger:  logger,
			config:  deviceConf,
		}
		pollers = append(pollers, d)
	}

	go a.run(ctx, pollers, cfg.PollInterval.Duration)

	return nil
}

func (a *Auto) run(ctx context.Context, pollers []*poller, interval time.Duration) {
	defer func() {
		var err error
		for _, p := range pollers {
			err = multierr.Append(err, p.Close())
		}
		if err != nil {
			a.services.Logger.Error("one or more Azure IoT Hub connections could not be closed",
				zap.Errors("errors", multierr.Errors(err)))
		}
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, p := range pollers {
				p.poll(ctx)
			}
		}
	}
}

func loadGroupKey(cfg Config) (auth.SASKey, error) {
	if cfg.GroupKey != "" {
		return auth.ParseSASKey(cfg.GroupKey)
	}

	raw, err := os.ReadFile(cfg.GroupKeyFile)
	if err != nil {
		return nil, err
	}
	return auth.ParseSASKey(string(raw))
}

func diallerFromConfig(devCfg DeviceConfig, idScope string, grpKey auth.SASKey) (dialler, error) {
	if devCfg.ConnectionString != "" {
		// the device specifies its own connection string, no need to use the DPS
		params, err := iothub.ParseConnectionString(devCfg.ConnectionString)
		if err != nil {
			return nil, fmt.Errorf("invalid connection string for device %q: %w", devCfg.Name, err)
		}

		return &directDialler{params: params}, nil
	}

	regId := devCfg.RegistrationID
	if regId == "" {
		return nil, fmt.Errorf("device %q is missing a registration ID", devCfg.Name)
	}

	return &dpsDialler{
		idScope: idScope,
		regID:   regId,
		key:     grpKey,
	}, nil
}
