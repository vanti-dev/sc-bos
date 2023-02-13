package pull

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Changes calls Pull on poller unless it's not supported, in which case it polls.
func Changes[C any](ctx context.Context, poller Fetcher[C], changes chan<- C, opts ...Option) error {
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
			err := poller.Pull(ctx, changes)
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
						conf.logger.Warn("updates are failing, will keep retrying", zap.Error(err))
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
