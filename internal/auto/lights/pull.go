package lights

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type pullPoller interface {
	pull(ctx context.Context, changes chan<- Patcher) error
	poll(ctx context.Context, changes chan<- Patcher) error
}

type subscribeOpts struct {
	logger    *zap.Logger
	pollDelay time.Duration

	fallbackInitialDelay time.Duration
	fallbackMaxDelay     time.Duration
}

type subscribeOption func(opts *subscribeOpts)

func withLogger(logger *zap.Logger) subscribeOption {
	return func(opts *subscribeOpts) {
		opts.logger = logger
	}
}

func withPullFallback(initial, max time.Duration) subscribeOption {
	return func(opts *subscribeOpts) {
		opts.fallbackInitialDelay = initial
		opts.fallbackMaxDelay = max
	}
}

func withPollDelay(delay time.Duration) subscribeOption {
	return func(opts *subscribeOpts) {
		opts.pollDelay = delay
	}
}

var defaultSubscribeOptions = []subscribeOption{
	withLogger(zap.NewNop()),
	withPullFallback(100*time.Millisecond, 10*time.Second),
	withPollDelay(time.Second),
}

func calcOpts(opts ...subscribeOption) subscribeOpts {
	out := &subscribeOpts{}
	for _, opt := range defaultSubscribeOptions {
		opt(out)
	}
	for _, opt := range opts {
		opt(out)
	}
	return *out
}

func subscribe(ctx context.Context, poller pullPoller, changes chan<- Patcher, opts ...subscribeOption) error {
	conf := calcOpts(opts...)

	poll := false
	var delay time.Duration
	var errCount int

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if poll {
			return runPoll(ctx, poller, changes, conf)
		} else {
			err := poller.pull(ctx, changes)
			if err != nil {
				if shouldReturn(err) {
					return err
				}
				if fallBackToPolling(err) {
					conf.logger.Debug("pull not supported, polling instead")
					poll = true
					delay = 0
					errCount = 0
					continue // skip the wait
				}
				if err != nil {
					errCount++
					if errCount == 5 {
						conf.logger.Warn("subscriptions are failing, will keep retrying", zap.Error(err))
					}
				}
			} else {
				errCount = 0
				delay = 0
			}
		}

		if delay == 0 {
			delay = conf.fallbackInitialDelay
		} else {
			delay = time.Duration(float64(delay) * 1.2)
			if delay > conf.fallbackMaxDelay {
				delay = conf.fallbackMaxDelay
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func runPoll(ctx context.Context, poller pullPoller, changes chan<- Patcher, conf subscribeOpts) error {
	pollDelay := conf.pollDelay
	errCount := 0
	ticker := time.NewTicker(conf.pollDelay)
	defer ticker.Stop()
	for {
		err := poller.poll(ctx, changes)
		if err != nil {
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
