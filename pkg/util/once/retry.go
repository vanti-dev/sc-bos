package once

import (
	"context"
	"errors"
	"sync"
)

// RetryError is like sync.Once except if f returns an err Do will try again the next time it's called.
type RetryError struct {
	mu  sync.Mutex
	c   chan struct{} // nil if needs init, blocked if init-ing, closed if successfully done
	err error

	recoverPanic bool // for testing, when true panics will be not be left to the system to log and fail tests.
}

var (
	ErrPanicSuppressed = errors.New("panic suppressed")
)

// Do calls f if Do has not successfully been called before.
// Do will return before f completes if ctx is done, f will continue to run until it returns.
// Concurrent calls to Do will wait until the first f completes, returning that error.
// Subsequent calls to Do after an error will invoke f as if Do had not been called before.
func (o *RetryError) Do(ctx context.Context, f func() error) error {
	o.mu.Lock()
	if o.c != nil {
		// some other routine is working on it
		c := o.c
		o.mu.Unlock()
		return o.wait(ctx, c)
	}

	// we should work on it
	c := make(chan struct{})
	o.err = nil
	o.c = c
	o.mu.Unlock()

	// we call f in a go routine to allow this call to return when ctx is done, even if f hasn't completed yet.
	go func() {
		var ran bool
		defer func() {
			// by default we don't actually want to recover from panics as this looses the stack,
			// but it's helpful to avoid incorrect test failures
			if o.recoverPanic {
				recover()
			}
			o.mu.Lock()
			defer o.mu.Unlock()
			if !ran && o.err == nil {
				// there was a panic
				o.err = ErrPanicSuppressed
			}
			if o.err != nil {
				o.c = nil // try again next time Do is called
			}
			close(c) // release any of those waiting for Do to complete
		}()
		// no lock needed here as err is only read when c is closed
		o.err = f()
		ran = true // only not set if f() panics
	}()

	return o.wait(ctx, c)
}

func (o *RetryError) wait(ctx context.Context, c chan struct{}) error {
	select {
	case <-c:
		return o.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
