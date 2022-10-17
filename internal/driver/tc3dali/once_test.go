package tc3dali

import (
	"context"
	"errors"
	"testing"
)

func TestOnce_Do(t *testing.T) {
	// Check that the callback won't be run again after a success
	t.Run("Serial", func(t *testing.T) {
		ctx := context.Background()
		o := newOnce()

		n := 0
		f := func() error {
			n++
			return nil
		}

		// this should call f
		err, done := o.Do(ctx, f)
		if !done {
			t.Error("expected first call's done value to be true")
		}
		if err != nil {
			t.Errorf("expected first error value to be nil but got %v", err)
		}
		if n != 1 {
			t.Errorf("expected callback to be called once - it has been called %d times", n)
		}

		// this should not call f
		err, done = o.Do(ctx, f)
		if !done {
			t.Error("expected second call's done value to be true")
		}
		if err != nil {
			t.Errorf("expected second error value to be nil but got %v", err)
		}
		if n != 1 {
			t.Errorf("expected callback to be called once - it has been called %d times", n)
		}
	})

	t.Run("Retry", func(t *testing.T) {
		ctx := context.Background()
		o := newOnce()

		fail := errors.New("intentional failure")
		n := 0
		f := func() error {
			n++
			if n == 1 {
				return fail
			} else {
				return nil
			}
		}

		// this should call f
		err, done := o.Do(ctx, f)
		if !done {
			t.Error("expected first call's done value to be true")
		}
		if !errors.Is(err, fail) {
			t.Errorf("expected first error value to be fail but got %v", err)
		}
		if n != 1 {
			t.Errorf("expected callback to be called once - it has been called %d times", n)
		}

		// this should call f again, and succeed
		err, done = o.Do(ctx, f)
		if !done {
			t.Error("expected second call's done value to be true")
		}
		if err != nil {
			t.Errorf("expected second error value to be nil but got %v", err)
		}
		if n != 2 {
			t.Errorf("expected callback to be called twice - it has been called %d times", n)
		}
	})

	// check that a Do call which is waiting for an existing call to finish can be interrupted by cancelling the context
	t.Run("Wait_Expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		o := newOnce()

		running := make(chan struct{}) // closed when the callback is known to be running
		go func() {
			err, done := o.Do(ctx, func() error {
				close(running)
				<-ctx.Done()
				return ctx.Err()
			})
			if !done {
				t.Error("expected goroutine Do to return done=true")
			}
			if !errors.Is(err, context.Canceled) {
				t.Errorf("expected goroutine Do to return err=context.Cancelled, got %v", err)
			}
		}()

		// force concurrency by waiting until the above callback invocation is definitely running
		<-running
		// create an expired context
		expired, cancelExpired := context.WithCancel(context.Background())
		cancelExpired()

		ran := false
		err, done := o.Do(expired, func() error {
			ran = true
			return nil
		})

		if ran {
			t.Error("callbacks ran concurrently!")
		}
		if done {
			t.Error("Do reported that it completed")
		}
		if !errors.Is(err, expired.Err()) {
			t.Errorf("error mismatch: expected %v, got %v", expired.Err(), err)
		}
	})
}
