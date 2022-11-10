package lights

import (
	"context"
	"errors"
	"fmt"
	"github.com/olebedev/emitter"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/auto/lights/config"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Patcher represents a single patch that adjusts ReadState.
type Patcher interface {
	Patch(s *ReadState)
}

type PatchFunc func(s *ReadState)

func (p PatchFunc) Patch(s *ReadState) {
	p(s)
}

type subscriber interface {
	Subscribe(ctx context.Context, changes chan<- Patcher) error
}

// setupReadSources configures the automation to pull events from the configured devices and emit patches into changes.
// Any configuration changes, via Configure, should be sent via configChanged chan which matches the return type for Emitter.On.
//
// Blocks until fatal errors in the subscriptions or ctx is done.
func (b *BrightnessAutomation) setupReadSources(ctx context.Context, configChanged <-chan emitter.Event, changes chan<- Patcher) error {
	// eagerly fetch the clients we might be using.
	// While the config might mean we don't use them, better to have the system fail early just in case
	var occupancySensorClient traits.OccupancySensorApiClient
	if err := b.clients.Client(&occupancySensorClient); err != nil {
		return fmt.Errorf("%w traits.OccupancySensorApiClient", err)
	}
	var brightnessSensorClient traits.BrightnessSensorApiClient
	if err := b.clients.Client(&brightnessSensorClient); err != nil {
		return fmt.Errorf("%w traits.BrightnessSensorApiClient", err)
	}

	// Setup the sources that we can pull patches from.
	type source struct {
		new   func(name string, logger *zap.Logger) subscriber
		names func(cfg config.Root) []string
		// runningSources, keyed by device name, tracks which sources are currently running.
		// The value can be called to cancel the context used to start that source.
		runningSources map[string]context.CancelFunc
	}
	sources := []source{
		{
			names: func(cfg config.Root) []string { return cfg.OccupancySensors },
			new: func(name string, logger *zap.Logger) subscriber {
				return &OccupancySensorPatches{name: name, client: occupancySensorClient, logger: logger}
			},
		},
		{
			names: func(cfg config.Root) []string { return cfg.BrightnessSensors },
			new: func(name string, logger *zap.Logger) subscriber {
				return &BrightnessSensorPatches{name: name, client: brightnessSensorClient, logger: logger}
			},
		},
	}

	// cancel everything if we're returning.
	defer func() {
		for _, source := range sources {
			for _, cancelFunc := range source.runningSources {
				cancelFunc()
			}
		}
	}()

	processConfig := func(cfg config.Root) (sourceCount int) {
		for _, source := range sources {
			names := source.names(cfg)
			if source.runningSources == nil && len(names) > 0 {
				source.runningSources = make(map[string]context.CancelFunc, len(names))
			}
			sourcesToStop := shallowCopyMap(source.runningSources)
			for _, name := range names {
				sourceCount++

				// are we already watching this name?
				if _, ok := sourcesToStop[name]; ok {
					delete(sourcesToStop, name)
					continue
				}
				// I guess not, lets start watching
				ctx, stop := context.WithCancel(ctx)
				source.runningSources[name] = stop
				logger := b.logger.With(zap.String("name", name))
				impl := source.new(name, logger)
				go func() {
					err := impl.Subscribe(ctx, changes)
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					if s, ok := status.FromError(err); ok {
						if s.Code() == codes.Unimplemented {
							logger.Warn(fmt.Sprintf("Subscription does not implement required features: %v", s.Message()))
							return
						}
					}
					if err != nil {
						// todo: handle error, the subscription has failed without us asking it to stop.
						logger.Warn("Subscription ended before it should", zap.Error(err))
					}
				}()
			}

			// stop any sources that are no longer in the config
			for _, cancelFunc := range sourcesToStop {
				cancelFunc()
			}
		}

		// update the config in the ReadState too
		changes <- PatchFunc(func(s *ReadState) {
			s.Config = cfg
		})

		return sourceCount
	}

	processConfig(b.config)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-configChanged:
			if sc := processConfig(e.Args[0].(config.Root)); sc == 0 {
				b.logger.Debug("no sources configured, automation will do nothing")
			}
		}
	}
}

func shallowCopyMap[K comparable, V any](m map[K]V) map[K]V {
	n := make(map[K]V, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}
