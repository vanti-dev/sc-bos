// Package pull provides a reliable way to subscribe to changes from a device.
package pull

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Changes calls Pull on poller unless it's not supported, in which case it polls.
// It will retry on error, backing off exponentially up to a maximum delay.
// It will return if the context is cancelled or a non-recoverable error occurs.
func Changes[C any](ctx context.Context, poller Fetcher[C], changes chan<- C, opts ...Option) error {
	conf := calcOpts(opts...)

	poll := false

	var mu sync.Mutex
	var delay time.Duration
	var errCount int
	resetErr := func() int {
		mu.Lock()
		defer mu.Unlock()
		old := errCount
		errCount = 0
		delay = 0
		return old
	}
	incErr := func() int {
		mu.Lock()
		defer mu.Unlock()
		errCount++
		return errCount
	}
	incDelay := func() {
		mu.Lock()
		defer mu.Unlock()
		if delay == 0 {
			delay = conf.fallbackInitialDelay
		} else {
			delay = time.Duration(float64(delay) * 1.2)
			if delay > conf.fallbackMaxDelay {
				delay = conf.fallbackMaxDelay
			}
		}
	}

	// These are used to track successful pulls.
	// The Pull method itself blocks until error so we have to track separately for success,
	// mostly for peace of mind and logging.
	var pullDuration time.Duration
	const successfulPullMultiplier = 4
	pullSuccessTimer := time.AfterFunc(math.MaxInt64, func() {
		attempts := resetErr()
		if attempts > 5 { // we only log failure after 5 attempts
			conf.logger.Debug("pulls are now succeeding", zap.Int("attempts", attempts))
		}
	})
	pullSuccessTimer.Stop() // avoid the timer actually running, the above was just to avoid nil timers

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if poll {
			return runPoll(ctx, poller, changes, conf)
		} else {
			if errCount > 0 {
				pullSuccessTimer.Reset(pullDuration * successfulPullMultiplier)
			}

			t0 := time.Now()
			err := poller.Pull(ctx, changes)
			pullDuration = time.Since(t0)

			if err != nil {
				pullSuccessTimer.Stop()
				if shouldReturn(err) {
					return err
				}
				if fallBackToPolling(err) {
					conf.logger.Debug("pull not supported, polling instead")
					poll = true
					resetErr()
					continue // skip the wait
				}
				if err != nil {
					if incErr() == 5 {
						conf.logger.Warn("updates are failing, will keep retrying", zap.Error(err))
					}
				}
			} else {
				if errCount > 0 {
					conf.logger.Debug("updates are now succeeding")
				}
				resetErr()
			}
		}

		incDelay()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func runPoll[C any](ctx context.Context, poller Fetcher[C], changes chan<- C, conf changeOpts) error {
	pollDelay := conf.pollDelay
	errCount := 0
	ticker := time.NewTicker(conf.pollDelay)
	defer ticker.Stop()
	for {
		err := poller.Poll(ctx, changes)
		if err != nil {
			if status.Code(err) == codes.Unimplemented {
				return err
			}

			errCount++
			pollDelay = time.Duration(float64(pollDelay) * 1.2)
			if pollDelay > conf.fallbackMaxDelay {
				pollDelay = conf.fallbackMaxDelay
			}
			if errCount == 5 {
				conf.logger.Warn("poll is failing, will try keep retrying", zap.Error(err))
			}
		} else {
			if pollDelay != conf.pollDelay {
				pollDelay = conf.pollDelay
				ticker.Reset(conf.pollDelay)
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func shouldReturn(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func fallBackToPolling(err error) bool {
	if grpcError, ok := status.FromError(err); ok {
		if grpcError.Code() == codes.Unimplemented {
			return true
		}
	}
	return false
}
