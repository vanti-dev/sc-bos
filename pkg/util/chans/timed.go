package chans

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimeout = errors.New("timeout waiting for chan")
	ErrClosed  = errors.New("chan closed")
	ErrOpen    = errors.New("chan open")
)

func SendWithin[T any](dest chan<- T, value T, wait time.Duration) error {
	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case dest <- value:
	case <-timer.C:
		return ErrTimeout
	}
	return nil
}

func RecvWithin[T any](src <-chan T, wait time.Duration) (T, error) {
	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case value, ok := <-src:
		if !ok {
			var zero T
			return zero, ErrClosed
		}
		return value, nil
	case <-timer.C:
		var zero T
		return zero, ErrTimeout
	}
}

func IsEmptyWithin[T any](src <-chan T, wait time.Duration) error {
	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case _, ok := <-src:
		if !ok {
			return ErrClosed
		}
		return ErrTimeout
	case <-timer.C:
		return nil
	}
}

func IsClosedWithin[T any](ch <-chan T, wait time.Duration) error {
	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case value, ok := <-ch:
		if ok {
			return fmt.Errorf("%w got val %v", ErrOpen, value)
		}
		return nil
	case <-timer.C:
		return ErrTimeout
	}
}
