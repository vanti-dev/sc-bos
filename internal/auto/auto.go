package auto

import (
	"context"
	"errors"
)

// Starter describes types that can be started.
type Starter interface {
	// Start instructs the automation to start.
	// The ctx represents how long the type can spend starting before it should give up.
	Start(ctx context.Context) error
}

// Stopper describes types that can be stopped.
type Stopper interface {
	// Stop instructs the automation to stop what it is doing.
	Stop() error
}

// ErrCannotBeStopped is returned when attempting to Stop an automation that does not implement Stopper.
var ErrCannotBeStopped = errors.New("cannot be stopped")

// Stop attempts to stop the given automation.
// If the automation does not implement Stopper then ErrCannotBeStopped will be returned.
func Stop(s any) error {
	if impl, ok := s.(Stopper); ok {
		return impl.Stop()
	}
	return ErrCannotBeStopped
}

// Stoppable returns whether the given automation can be stopped.
func Stoppable(s any) bool {
	_, ok := s.(Stopper)
	return ok
}

// Configurer describes types that can be configured.
// Configuration data is represented as a []byte and the automation is expected to decode this data into an internally
// significant data structure.
type Configurer interface {
	Configure(configData []byte) error
}

// ErrCannotBeConfigured is returned when attempting to Configure an automation that does not implement Configurer.
var ErrCannotBeConfigured = errors.New("cannot be configured")

// Configure attempts to configure the given automation.
// If the automation does not implement Configurer then ErrCannotBeConfigured will be returned.
func Configure(s any, configData []byte) error {
	if impl, ok := s.(Configurer); ok {
		return impl.Configure(configData)
	}
	return ErrCannotBeConfigured
}

// Configurable returns whether the given automation can be configured.
func Configurable(s any) bool {
	_, ok := s.(Configurer)
	return ok
}
