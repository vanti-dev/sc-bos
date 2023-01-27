package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

var (
	// ErrAlreadyStarted is returned from Service.Start if the service state is active.
	ErrAlreadyStarted = errors.New("already started")
	// ErrAlreadyStopped is returned from Service.Stop if the service state is not active.
	ErrAlreadyStopped = errors.New("already stopped")
	// ErrAlreadyLoading is returned from Service.Configure if the service state is loading.
	ErrAlreadyLoading = errors.New("already loading")
)

// State contains all the properties of the service.
// A services state can be fetched on demand or you can be notified of changes via the Service.State,
// Service.StateChanges, and Service.StateAndChanges depending on your needs.
type State struct {
	// Active records whether a Service is performing its job or not.
	// If true then the Service is doing what it is designed to do, it potentially has running background tasks, and is responding to stimulus.
	// If false then the Service is stopped, is not doing any work, and has no background processes active.
	// Active is updated via Service.Start and Service.Stop.
	// Additionally if a service fails to load config then this can also cause the Service to become inactive with an error.
	Active bool
	// Config is the raw configuration bytes last successfully used to configure the Service.
	// The contents of Config should not be modified after calling Configure
	Config []byte
	// Loading indicates that the service is Active and is processing an update to config.
	Loading bool
	// Err holds the error returned by the ApplyFunc.
	Err error

	// Times active was last set to false or true respectively.
	LastInactiveTime, LastActiveTime time.Time
	// Times these fields were last filled.
	// Setting err to nil does not update the time.
	LastErrTime, LastConfigTime time.Time
	// Times when loading was set to true or false respectively.
	// lastLoadingEndTime will be set to the zero time when loading is set to true.
	LastLoadingStartTime, LastLoadingEndTime time.Time
}

type Lifecycle interface {
	Start() (State, error)
	Configure(data []byte) (State, error)
	Stop() (State, error)

	State() State
	StateChanges(ctx context.Context) <-chan State
	StateAndChanges(ctx context.Context) (State, <-chan State)
}

// ApplyFunc is called each time an active service has its config updated.
// The func should block for as long as config is being read but no longer.
// Background tasks, like opening connections to other network devices, should not block.
// The given context will be cancelled if the Service is stopped.
//
// Only one call to ApplyFunc will happen at the same time, but it may be called more than once if Service.Configure is
// called on an active service.
// The implementation is responsible for cleaning up any outdated resources as part of the second ApplyFunc call.
type ApplyFunc[C any] func(ctx context.Context, config C) error

// ParseFunc is called to convert raw []byte config into a known config structure.
// The func is called in line with calls to Service.Configure and any errors are reported immediately.
// The func should perform any validation on the configuration in order to avoid errors once the parsed config is passed
// to ApplyFunc.
type ParseFunc[C any] func(data []byte) (C, error)

// Service manages the lifecycle of a background task.
type Service[C any] struct {
	mu       sync.Mutex // guards the next few fields
	state    State
	bus      *minibus.Bus[State]
	config   *C
	stopCtx  context.Context
	stopFunc context.CancelFunc

	parse  ParseFunc[C]
	apply  ApplyFunc[C]
	onStop func()

	now func() time.Time
}

// New creates a new Service using apply to spin up background tasks based on config of type C.
// The default ParserFunc uses json.Unmarshal to convert []byte to C.
// The service is created inactive and without config.
func New[C any](apply ApplyFunc[C], opts ...Option[C]) *Service[C] {
	s := &Service[C]{
		bus:   &minibus.Bus[State]{},
		apply: apply,
	}
	for _, opt := range DefaultOpts[C]() {
		opt.apply(s)
	}
	for _, opt := range opts {
		opt.apply(s)
	}
	s.state.LastInactiveTime = s.now()
	return s
}

// State returns the current state of the service.
func (l *Service[C]) State() State {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.state
}

// StateChanges returns a chan that will emit each time the service state changes.
func (l *Service[C]) StateChanges(ctx context.Context) <-chan State {
	return l.bus.Listen(ctx)
}

// StateAndChanges atomically returns the current state and a chan that emits future state changes.
func (l *Service[C]) StateAndChanges(ctx context.Context) (State, <-chan State) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.state, l.bus.Listen(ctx)
}

// Start transitions the Service to the active state.
// If the service has config then that config will be loaded as part of this call without blocking.
// Starting an active service returns ErrAlreadyStarted.
func (l *Service[C]) Start() (State, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	state := l.state
	if state.Active {
		return state, ErrAlreadyStarted
	}

	state.Active = true
	state.Loading = false
	state.Err = nil
	state.LastActiveTime = l.now()

	l.stopCtx, l.stopFunc = context.WithCancel(context.Background())

	if l.config == nil {
		// start without config
		return l.saveLocked(state)
	}

	return l.applyConfig(state, *l.config)
}

// Configure updates the config associated with this service.
// If the service is active then the config will be loaded as part of this call without blocking.
// If ParseFunc returns an error parsing data, that error will be returned and no state transition will be applied.
func (l *Service[C]) Configure(data []byte) (State, error) {
	// parse outside of holding the lock
	config, err := l.parse(data)
	if err != nil {
		return State{}, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	state := l.state
	if state.Loading {
		return state, ErrAlreadyLoading
	}

	state.LastConfigTime = l.now()
	state.Config = data
	l.config = &config

	if !state.Active {
		// configure without applying
		return l.saveLocked(state)
	}

	return l.applyConfig(state, config)
}

// Stop transitions the service to the inactive state.
// The context for any ApplyFunc calls will be cancelled.
// Config and other state will not be adjusted.
// If the service is inactive calling Stop will return ErrAlreadyStopped.
func (l *Service[C]) Stop() (State, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	state := l.state
	if !state.Active {
		return state, ErrAlreadyStopped
	}

	state.Active = false
	state.LastInactiveTime = l.now()
	l.stopLocked()

	if l.onStop != nil {
		l.onStop()
	}

	return l.saveLocked(state)
}

func (l *Service[C]) applyConfig(state State, config C) (State, error) {
	state.Loading = true
	state.LastLoadingStartTime = l.now()
	ctx := l.stopCtx
	go func() {
		err := l.apply(ctx, config)

		l.mu.Lock()
		defer l.mu.Unlock()
		state := l.state
		state.LastLoadingEndTime = l.now()
		state.Loading = false
		if err != nil && ctx == l.stopCtx {
			// only record the error if we're active, otherwise we've stopped and the error is likely because
			// the context has been cancelled
			state.Err = err
			state.Active = false
			state.LastErrTime = l.now()
			state.LastInactiveTime = l.now()
			l.stopLocked()
		}
		_, _ = l.saveLocked(state)
	}()

	return l.saveLocked(state)
}

func (l *Service[C]) saveLocked(state State) (State, error) {
	l.state = state
	go l.bus.Send(context.Background(), state)
	return state, nil
}

func (l *Service[C]) stopLocked() {
	if stop := l.stopFunc; stop != nil {
		// clear before calling stop to avoid races with go routines that are blocked on ctx.Done
		l.stopFunc = nil
		l.stopCtx = nil
		stop()
	}
}
