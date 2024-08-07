package task

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPollErr(t *testing.T) {
	t.Run("repeats", func(t *testing.T) {
		runCount := 0
		action := func(ctx context.Context) error {
			runCount++
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

		if runCount < 2 {
			t.Errorf("expected at least 2 runs, got %d", runCount)
		}
	})

	t.Run("stops", func(t *testing.T) {
		var mu sync.Mutex
		runCount := 0
		action := func(ctx context.Context) error {
			mu.Lock()
			runCount++
			mu.Unlock()
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

		mu.Lock()
		runCount2 := runCount
		mu.Unlock()

		time.Sleep(10 * time.Millisecond)
		if runCount != runCount2 {
			t.Errorf("expected no more runs, got %d extra", runCount-runCount2)
		}
	})

	t.Run("uses err backoff", func(t *testing.T) {
		runCount := 0
		action := func(ctx context.Context) error {
			runCount++
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

		if runCount < 2 {
			t.Errorf("expected at least 2 runs, got %d", runCount)
		}
	})
}
