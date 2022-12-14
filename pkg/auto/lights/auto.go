package lights

import (
	"context"
	"github.com/olebedev/emitter"
	"github.com/vanti-dev/sc-bos/pkg/auto/lights/config"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

// BrightnessAutomation implements turning lights on or off based on occupancy readings from PIRs and other devices.
type BrightnessAutomation struct {
	logger  *zap.Logger
	clients node.Clienter // clients are not got until Start

	// bus emits "stop" and "config" events triggered by Stop and Configure.
	bus *emitter.Emitter

	makeActions func(clienter node.Clienter) (actions, error)                // override for testing
	newTimer    func(duration time.Duration) (<-chan time.Time, func() bool) // override for testing
}

// PirsTurnLightsOn creates an automation that controls light brightness based on PIR occupancy status.
func PirsTurnLightsOn(clients node.Clienter, logger *zap.Logger) *BrightnessAutomation {
	return &BrightnessAutomation{
		logger:      logger,
		clients:     clients,
		makeActions: newActions,
		newTimer: func(duration time.Duration) (<-chan time.Time, func() bool) {
			t := time.NewTimer(duration)
			return t.C, t.Stop
		},
	}
}

// Start implements Starter and initialises this automation.
// Start may be called before or after Configure.
func (b *BrightnessAutomation) Start(_ context.Context) error {
	if b.bus != nil {
		b.bus.Off("*")
	} else {
		b.bus = emitter.New(1)
	}
	// We make the actions impl here so that we can create the automation before clients are available,
	// so long as they're available before Start is called.
	actions, err := b.makeActions(b.clients)
	if err != nil {
		return err
	}

	ctx, stop := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	// make sure we stop the group when Stop is called
	stopCalled := b.bus.On("stop")
	group.Go(func() error {
		defer b.bus.Off("stop", stopCalled)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-stopCalled:
				stop()
				return context.Canceled
			}
		}
	})

	changes := make(chan Patcher, 10)

	// "actions" the subscriptions. This causes any Pull or Poll or other even sources to start running and placing
	// state changes into the changes chan.
	configChanged := b.bus.On("config")
	group.Go(func() error {
		defer b.bus.Off("config", configChanged)
		return b.setupReadSources(ctx, configChanged, changes)
	})

	// readStates receives state that should be processed, for example to work out if lights should be turned on.
	// readStates is how state changes are communicates to the func that processes those state changes.
	readStates := make(chan *ReadState)
	initialState := NewReadState()

	// collect and collate the state changes
	group.Go(func() error {
		return readStateChanges(ctx, initialState, changes, readStates)
	})

	// process state changes and execute RPC calls to update the real world
	group.Go(func() error {
		return b.processStateChanges(ctx, readStates, actions)
	})

	// note: we don't wait for the group to complete!
	return nil
}

// Configure updates the configured devices and settings this automation uses.
func (b *BrightnessAutomation) Configure(configData []byte) error {
	cfg, err := config.Read(configData)
	if err != nil {
		return err
	}
	return b.configure(cfg)
}

func (b *BrightnessAutomation) configure(cfg config.Root) error {
	<-b.bus.Emit("config", cfg) // wait for anyone who is listening to apply that config
	return nil
}

// Stop stops the automation from running.
// You must call Start to have automated action occur again.
func (b *BrightnessAutomation) Stop() error {
	b.bus.Emit("stop") // don't wait
	return nil
}

// processStateChanges reads ReadState from a channel and analyses each entry deciding if light levels should be changed.
// This function backs off to processState which has the actual logic for what to do given a certain state,
// this function handles the channel management, retry logic, TTL on decisions, and all that type of thing.
func (b *BrightnessAutomation) processStateChanges(ctx context.Context, readStates <-chan *ReadState, actions actions) error {
	// the below are the innards of time.Timer, but expanded so we can stop/select on them even if
	// we don't have a timer active right now

	var ttlExpired <-chan time.Time
	cancelTtlTimer := func() bool { return false }

	var retryFailedProcessing <-chan time.Time
	cancelRetryTimer := func() bool { return false }

	// writeState is only accessed from this go routine.
	writeState := NewWriteState()

	processStateFn := func(readState *ReadState) error {
		cancelTtlTimer()
		cancelRetryTimer()

		ttl, err := processState(ctx, readState, writeState, actions)
		b.bus.Emit("process-complete", ttl, err, readState, writeState) // used only for testing, notify that processing has completed
		if err != nil {
			// todo: setup retries for processing the state
			return err
		}

		// Setup ttl for the transformed model.
		// After this time it should be recalculated.
		if ttl > 0 {
			ttlExpired, cancelTtlTimer = b.newTimer(ttl)
		}

		return nil
	}

	var lastReadState *ReadState

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case readState := <-readStates:
			lastReadState = readState
			err := processStateFn(readState)
			if err != nil {
				return err
			}
		case <-ttlExpired:
			err := processStateFn(lastReadState)
			if err != nil {
				return err
			}
		case <-retryFailedProcessing:
			err := processStateFn(lastReadState)
			if err != nil {
				return err
			}
		}
	}
}

type notify interface {
	Emit(topic string, args ...any) chan struct{}
}
