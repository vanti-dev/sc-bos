package driver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/vanti-dev/bsp-ew/internal/util/state"
	"go.uber.org/zap"
)

// RunConfigFunc is called when a drivers config changes and should apply those changes.
// The ctx will be cancelled if the driver stops or the config is replaced.
type RunConfigFunc[C any] func(ctx context.Context, cfg C) error

// Lifecycle manages the lifecycle of a driver as per task.Starter, task.Configurer, and task.Stopper.
// Embed Lifecycle into your own driver type and provide a RunConfigFunc that binds your driver based on config.
type Lifecycle[C any] struct {
	// ApplyConfig is a function that gets called each time config is changed.
	ApplyConfig RunConfigFunc[C]
	Logger      *zap.Logger
	// ReadConfig converts bytes into C.
	// Defaults to json.Unmarshal.
	ReadConfig func(bytes []byte) (C, error)

	status *state.Manager[Status] // allows us to implement Stateful

	// these track our state and config changes
	configC chan C
	stopCtx context.Context
	stop    context.CancelFunc
}

func NewLifecycle[C any](runFunc RunConfigFunc[C]) *Lifecycle[C] {
	return &Lifecycle[C]{
		Logger:      zap.NewNop(),
		ApplyConfig: runFunc,
		ReadConfig: func(bytes []byte) (C, error) {
			var c C
			err := json.Unmarshal(bytes, &c)
			return c, err
		},
		status: state.NewManager(StatusInactive),
	}
}

// Start makes this driver available to be configured.
// Call Stop when you're done with the driver to free up resources.
//
// Start must be called before Configure.
// Once started Configure and Stop may be called from any go routine.
func (s *Lifecycle[C]) Start(_ context.Context) error {
	// We implement a main loop pattern to avoid locks.
	// Methods like Configure and Stop both push messages onto channels
	// which we select on in a loop using a single go routine, avoiding any
	// locking issues (that we have to deal with ourselves).

	s.configC = make(chan C, 5)
	s.stopCtx, s.stop = context.WithCancel(context.Background())
	go func() {
		s.status.Update(StatusActive)
		defer s.status.Update(StatusInactive)

		// allow stopping and re-running based on new config
		var cfgCtx context.Context
		cfgStop := func() {}
		defer func() {
			cfgStop()
		}()

		for {
			select {
			case cfg := <-s.configC:
				s.status.Update(StatusLoading)

				cfgStop()
				// It's called in the defer
				//goland:noinspection GoVetLostCancel
				cfgCtx, cfgStop = context.WithCancel(s.stopCtx)

				err := s.ApplyConfig(cfgCtx, cfg)
				s.status.Update(StatusActive)
				if err != nil {
					s.Logger.Error("failed to apply config update", zap.Error(err))
					continue
				}
			case <-s.stopCtx.Done():
				//goland:noinspection GoVetLostCancel
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
func (s *Lifecycle[C]) Configure(configData []byte) error {
	if s.configC == nil {
		return errors.New("not started")
	}
	var c C
	err := json.Unmarshal(configData, &c)
	if err != nil {
		return err
	}
	s.configC <- c
	return nil
}

// Stop stops the driver and releases resources.
// Stop races with Start before Start has completed, but can be called concurrently once started.
func (s *Lifecycle[C]) Stop() error {
	if s.stop == nil {
		// not started
		return nil
	}
	s.stop()
	return nil
}

func (s *Lifecycle[C]) WaitForStateChange(ctx context.Context, sourceState Status) error {
	return s.status.WaitForStateChange(ctx, sourceState)
}

func (s *Lifecycle[C]) CurrentState() Status {
	return s.status.CurrentState()
}
