package task

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestPollErr(t *testing.T) {
	t.Run("repeats", func(t *testing.T) {
		var runCount atomic.Int32
		action := func(ctx context.Context) error {
			runCount.Add(1)
			return nil
		}
		ctx, stop := context.WithCancel(context.Background())
		defer stop()
		err := PollErr(action, WithPollInterval(time.Millisecond/10)).Attach(ctx)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(10 * time.Millisecond)
		stop()

		if c := runCount.Load(); c < 2 {
			t.Errorf("expected at least 2 runs, got %d", c)
		}
	})

	t.Run("stops", func(t *testing.T) {
		var runCount atomic.Int32
		action := func(ctx context.Context) error {
			runCount.Add(1)
			return nil
		}
		ctx, stop := context.WithCancel(context.Background())
		defer stop()

		runner := PollErr(action, WithPollInterval(time.Millisecond/10))
		err := runner.Attach(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// override the stop function to signal when it is called
		stopCalled := make(chan struct{})
		stopFn := runner.stop
		runner.stop = func() {
			stopFn()
			close(stopCalled)
		}

		time.Sleep(10 * time.Millisecond)
		stop()

		<-stopCalled // wait for the stop to be called as this is a short-running task

		runCount2 := runCount.Load()

		time.Sleep(10 * time.Millisecond)
		if delta := runCount.Load() - runCount2; delta > 0 {
			t.Errorf("expected no more runs, got %d extra", delta)
		}
	})

	t.Run("uses err backoff", func(t *testing.T) {
		var runCount atomic.Int32
		action := func(ctx context.Context) error {
			runCount.Add(1)
			return errors.New("expected test error")
		}
		ctx, stop := context.WithCancel(context.Background())
		defer stop()
		err := PollErr(action, WithPollInterval(30*time.Second), WithPollErrBackoff(time.Millisecond/10, time.Millisecond, 2)).Attach(ctx)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(10 * time.Millisecond)
		stop()

		if c := runCount.Load(); c < 2 {
			t.Errorf("expected at least 2 runs, got %d", c)
		}
	})
}
