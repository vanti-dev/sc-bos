package driver

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var ErrDriverFactoryNotFound = errors.New("driver factory not found")

type BuildResult struct {
	Type   string
	Driver Driver
	Err    error
}

func Build(ctx context.Context, services Services, factories map[string]Factory, configs []RawConfig) map[string]BuildResult {
	results := make(map[string]BuildResult)

	for _, config := range configs {
		logger := services.Logger.With(
			zap.Namespace("driver"),
			zap.String("type", config.Type),
			zap.String("name", config.Name),
		)
		if _, ok := results[config.Name]; ok {
			logger.Warn("skipping duplicate driver entry")
			continue
		}

		factory, ok := factories[config.Type]
		if !ok {
			logger.Error("unknown driver type")
			results[config.Name] = BuildResult{
				Type: config.Type,
				Err:  ErrDriverFactoryNotFound,
			}
			continue
		}

		driverServices := services
		driverServices.Logger = logger

		start := time.Now()
		logger.Debug("starting driver initialisation")
		driver, err := factory(ctx, driverServices, config.Raw)
		duration := time.Now().Sub(start)

		if err != nil {
			logger.Error("driver failed to initialise", zap.Error(err), zap.Duration("duration", duration))
		} else {
			logger.Debug("driver finished initialisation", zap.Duration("duration", duration))
		}
		results[config.Name] = BuildResult{
			Type:   config.Type,
			Driver: driver,
			Err:    err,
		}
	}

	return results
}
