package lights

import (
	"context"
	"errors"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/auto/lights/config"
	"go.uber.org/zap"
)

// Patcher represents a single patch that adjusts ReadState.
type Patcher interface {
	Patch(s *ReadState)
}

type PatchFunc func(s *ReadState)

func (p PatchFunc) Patch(s *ReadState) {
	p(s)
}

// setupReadSources configures the automation to pull events from the configured devices and emit patches into changes.
// Any configuration changes, via Configure, are recognised and the event sources updated.
//
// Blocks until fatal errors in the subscriptions or ctx is done.
func (b *BrightnessAutomation) setupReadSources(ctx context.Context, changes chan<- Patcher) error {
	// eagerly fetch the clients we might be using.
	// While the config might mean we don't use them, better to have the system fail early just in case
	var occupancySensorClient traits.OccupancySensorApiClient
	if err := b.clients.Client(&occupancySensorClient); err != nil {
		return fmt.Errorf("%w traits.OccupancySensorApiClient", err)
	}

	// runningSources, keyed by device name, tracks which sources are currently running.
	// The value can be called to cancel the context used to start that source.
	runningSources := make(map[string]context.CancelFunc)
	defer func() {
		// cancel everything if we're returning.
		for _, cancelFunc := range runningSources {
			cancelFunc()
		}
	}()

	processConfig := func(cfg config.Root) {
		sourcesToStop := shallowCopyMap(runningSources)
		for _, name := range cfg.OccupancySensors {
			// are we already watching this name?
			if _, ok := sourcesToStop[name]; ok {
				delete(sourcesToStop, name)
				continue
			}
			// I guess not, lets start watching
			ctx, stop := context.WithCancel(ctx)
			runningSources[name] = stop
			source := &OccupancySensorPatches{name: name, client: occupancySensorClient}
			go func() {
				err := source.Subscribe(ctx, changes)
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}
				if err != nil {
					// todo: handle error, the subscription has failed without us asking it to stop. Retry?
					b.logger.Warn("Subscription ended before it should", zap.String("source", source.name), zap.Error(err))
				}
			}()
		}

		// stop any sources that are no longer in the config
		for _, cancelFunc := range sourcesToStop {
			cancelFunc()
		}

		// update the config in the ReadState too
		changes <- PatchFunc(func(s *ReadState) {
			s.Config = cfg
		})
	}

	configChanged := b.bus.On("config")
	defer b.bus.Off("config", configChanged)

	processConfig(b.config)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-configChanged:
			processConfig(e.Args[0].(config.Root))
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
