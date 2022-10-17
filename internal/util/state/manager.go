package state

import (
	"context"
	"errors"
	"sync"
)

// Manager keeps track of changes to state.
type Manager[T comparable] struct {
	// This code is "heavily inspired" by grpc.connectivityStateManager and grpc.ClientConn.WaitForStateChange

	mu         sync.Mutex
	state      T
	notifyChan chan struct{}
}

// NewManager returns a new Manager with an initial state.
func NewManager[T comparable](initialState T) *Manager[T] {
	return &Manager[T]{state: initialState}
}

// Update updates the state if the current state isn't terminal.
// If there's a change it notifies goroutines waiting on state change to happen.
func (sm *Manager[T]) Update(state T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if IsTerminal(sm.state) {
		return
	}
	if sm.state == state {
		return
	}
	sm.state = state
	if sm.notifyChan != nil {
		// There are other goroutines waiting on this channel.
		close(sm.notifyChan)
		sm.notifyChan = nil
	}
}

// GetNotifyChan returns a chan that can be used to be notified when the state changes.
func (sm *Manager[T]) GetNotifyChan() <-chan struct{} {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.notifyChan == nil {
		sm.notifyChan = make(chan struct{})
	}
	return sm.notifyChan
}

// CurrentState returns the current state.
func (sm *Manager[T]) CurrentState() T {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.state
}

// WaitForStateChange implements Stateful.WaitForStateChange.
func (sm *Manager[S]) WaitForStateChange(ctx context.Context, sourceState S) error {
	ch := sm.GetNotifyChan()
	s := sm.CurrentState()
	if s != sourceState {
		return nil
	}
	if IsTerminal(s) {
		return ErrTerminalState
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}

// IsTerminal returns whether state implements Terminal and returns true for Terminal.IsTerminal.
func IsTerminal(state any) bool {
	if ts, ok := state.(Terminal); ok {
		return ts.IsTerminal()
	}
	return false
}

// Stateful describes types that have a single state.
type Stateful[S comparable] interface {
	// WaitForStateChange waits until the state changes from sourceState or ctx expires.
	// An error will be returned if ctx expires or the state will never change.
	WaitForStateChange(ctx context.Context, sourceState S) error
	// CurrentState returns the current state.
	CurrentState() S
}

// ErrTerminalState is returned from Stateful.WaitForStateChange if the currentState is a terminal state.
var ErrTerminalState = errors.New("terminal state")

// Terminal allows a state to indicate that it is a terminal state or not.
type Terminal interface {
	IsTerminal() bool
}
