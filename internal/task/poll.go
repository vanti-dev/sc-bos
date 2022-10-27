package task

import (
	"context"
	"time"
)

// Poll creates a task that calls a function at a regular interval while it is running.
// The action will not be run until the returned Intermittent.Attach is called.
func Poll(action func(context.Context), interval time.Duration) *Intermittent {
	if interval <= 0 {
		panic("invalid interval")
	}

	// we don't need to do anything to initialise, so we don't use the initialisation context passed in
	start := func(_ context.Context) (StopFn, error) {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					action(ctx)
				}
			}
		}()

		return cancel, nil
	}

	return NewIntermittent(start)
}
