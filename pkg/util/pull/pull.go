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

type Getter = func(context.Context) error

// OrPoll will attempt to call pull, falling back to poll if it's not supported.
// pull subscribe to changes, blocking until the subscription fails, if possible failures will be retried.
// poll should query for the current state, this function manages delays between calls.
func OrPoll(ctx context.Context, pull, poll Getter, opts ...Option) error {
	conf := calcOpts(opts...)

	shouldPoll := false

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

		if shouldPoll {
			return runPoll(ctx, poll, conf)
		} else {
			if errCount > 0 {
				pullSuccessTimer.Reset(pullDuration * successfulPullMultiplier)
			}

			t0 := time.Now()
			err := pull(ctx)
			pullDuration = time.Since(t0)

			if err != nil {
				pullSuccessTimer.Stop()
				if shouldReturn(err) {
					return err
				}
				if fallBackToPolling(err) {
					conf.logger.Debug("pull not supported, polling instead")
					shouldPoll = true
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

// Changes calls Pull on poller unless it's not supported, in which case it polls.
// It will retry on error, backing off exponentially up to a maximum delay.
// It will return if the context is cancelled or a non-recoverable error occurs.
func Changes[C any](ctx context.Context, poller Fetcher[C], changes chan<- C, opts ...Option) error {
	return OrPoll(ctx,
		func(ctx context.Context) error {
			return poller.Pull(ctx, changes)
		},
		func(ctx context.Context) error {
			return poller.Poll(ctx, changes)
		},
		opts...)
}

func runPoll(ctx context.Context, get Getter, conf changeOpts) error {
	pollDelay := conf.pollDelay
	errCount := 0
	ticker := time.NewTicker(conf.pollDelay)
	defer ticker.Stop()
	for {
		err := get(ctx)
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
