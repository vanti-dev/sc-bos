// Package pull provides a reliable way to subscribe to changes from a device.
package pull

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Getter = func(context.Context) error

// OrPoll calls pull, falling back to repeated calls to get if pull returns a codes.Unimplemented error.
// Argument pull should subscribe to changes, blocking until the subscription fails.
// Argument get should get the current state, blocking until the state is got or an error occurs.
// OrPoll blocks until the context is cancelled, retrying either pull or get on error.
func OrPoll(ctx context.Context, pull, get Getter, opts ...Option) error {
	conf := calcOpts(opts...)

	return runBlocking(ctx, "blocking tasks", func(ctx context.Context) error {
		err := runPull(ctx, pull, conf)
		if shouldReturn(err) {
			return err
		}
		// the runPull err is never nil, so if we're here then we need to try polling
		err = get(ctx) // get one to see if this is unimplemented too
		if shouldReturn(err) {
			return err
		}
		if !isUnimplemented(err) {
			// pull is unimplemented but get is not, there's a chance things have changed (i.e. during boot)
			// so try pull again just to be sure.
			pullErr := runPull(ctx, pull, conf)
			if shouldReturn(pullErr) {
				return pullErr
			}
			// pull is still unimplemented, let's start the poll loop for real
			err = runPoll(ctx, err, get, conf)
			if shouldReturn(err) {
				return err
			}
			// else get is unimplemented
		}

		// if we're here both pull and get returned unimplemented, so retry everything after a delay
		return err
	}, conf, func(error) bool { return false })
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

// runPull retries calling pull returning on context cancellation or when pull returns an error that has codes.Unimplemented.
// runPull will never return a nil error.
func runPull(ctx context.Context, pull Getter, conf changeOpts) error {
	return runBlocking(ctx, "pulls", pull, conf, isUnimplemented)
}

// runPoll blocks repeatedly calling get after some delay until context is cancelled or get returns an error that has codes.Unimplemented.
// runPoll never returns a nil error.
func runPoll(ctx context.Context, err error, get Getter, conf changeOpts) error {
	pollDelay := conf.pollDelay
	retry := conf.backoff(backoff.WithInitialInterval(pollDelay + pollDelay/5)) // pollDelay * 1.2 (without casting)
	errCount := 0
	ticker := conf.clock.NewTicker(pollDelay)
	defer ticker.Stop()

	incErr := func() {
		errCount++
		pollDelay = retry.NextBackOff()
		ticker.Reset(pollDelay)
		if errCount == 5 {
			conf.logger.Warn("poll is failing, will try keep retrying", zap.Error(err))
		}
	}

	if err != nil {
		incErr()
	}

	for {
		// delay before first poll as the caller will have called get for us already.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.Chan():
		}

		err := get(ctx)
		if shouldReturn(err) || isUnimplemented(err) {
			return err
		}
		if err != nil {
			incErr()
		} else {
			// get worked, reset the poll delay back to the configured one (if we have to)
			if pollDelay != conf.pollDelay {
				pollDelay = conf.pollDelay
				retry.Reset()
				ticker.Reset(conf.pollDelay)
			}
		}
	}
}

// runBlocking reties task until ctx is cancelled or fatal returns true for returned errors.
func runBlocking(ctx context.Context, name string, task Getter, conf changeOpts, fatal func(error) bool) error {
	var mu sync.Mutex
	retry := conf.backoff()
	var errCount int
	resetErr := func() int {
		mu.Lock()
		defer mu.Unlock()
		old := errCount
		errCount = 0
		retry.Reset()
		return old
	}
	incErr := func() int {
		mu.Lock()
		defer mu.Unlock()
		errCount++
		return errCount
	}
	incDelay := func() time.Duration {
		mu.Lock()
		defer mu.Unlock()
		return retry.NextBackOff()
	}

	// These are used to track successful pulls.
	// The task method itself blocks until error so we have to track separately for success,
	// mostly for peace of mind and logging.
	var taskDuration time.Duration
	const successfulMultiplier = 4
	successTimer := conf.clock.AfterFunc(math.MaxInt64, func() {
		attempts := resetErr()
		if attempts > 5 { // we only log failure after 5 attempts
			conf.logger.Debug(name+" are now succeeding", zap.Int("attempts", attempts))
		}
	})
	successTimer.Stop() // avoid the timer actually running, the above was just to avoid nil timers
	defer successTimer.Stop()

	for {
		t0 := conf.clock.Now()
		err := task(ctx) // blocks
		taskDuration = conf.clock.Since(t0)
		successTimer.Stop()

		if shouldReturn(err) || fatal(err) {
			return err
		}
		if err != nil {
			errCount := incErr()
			successTimer.Reset(taskDuration * successfulMultiplier)
			if errCount == 5 {
				conf.logger.Warn(name+" are failing, will keep retrying", zap.Error(err))
			}
		} else {
			// A nil error means the task stopped successfully, this is unusual as it should block forever, but not technically an error.
			// Let's make sure we run it again and hope that the next time it stays alive for longer.
			resetErr()
			continue // skip the wait
		}

		delay := incDelay()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-conf.clock.After(delay):
		}
	}
}

func shouldReturn(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func isUnimplemented(err error) bool {
	if grpcError, ok := status.FromError(err); ok {
		if grpcError.Code() == codes.Unimplemented {
			return true
		}
	}
	return false
}
