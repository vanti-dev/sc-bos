package service

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// AbortRetry returns a new error wrapping err that when returned from an ApplyFunc passed to RetryApply will cause it to not attempt a retry.
func AbortRetry(err error) error {
	return abortRetry{err: err}
}

type abortRetry struct {
	err error
}

func (a abortRetry) Error() string {
	return a.err.Error()
}

func (a abortRetry) Unwrap() error {
	return a.err
}

type retryOptions struct {
	InitialDelay time.Duration // time between the first and second attempt
	MaxDelay     time.Duration // maximum time between attempts
	MinDelay     time.Duration // minimum time between attempts
	Factor       float64       // factor by which to increase the delay between attempts
	MaxAttempts  int           // maximum number of attempts, 0 means unlimited

	Logger func(logContext RetryContext)
}

type RetryContext struct {
	Attempt int
	Err     error
	Delay   time.Duration // until next attempt, 0 if no more attempts will be performed
	T0      time.Time     // first attempt time
}

// LogTo logs using a reasonable default format at a reasonable frequency to logger.
func (ctx RetryContext) LogTo(desc string, logger *zap.Logger) {
	success := ctx.Err == nil
	switch {
	case success && ctx.Attempt == 0: // worked first try, don't log anything
	case success: // worked after at least one retry
		logger.Info(fmt.Sprintf("%s has now succeeded", desc), zap.Int("attempts", ctx.Attempt), zap.Duration("duration", time.Since(ctx.T0)))
	case ctx.Attempt == 5: // has failed enough times to be worth logging
		logger.Warn(fmt.Sprintf("%s is failing, will retry", desc), zap.Int("attempts", ctx.Attempt), zap.Duration("duration", time.Since(ctx.T0)), zap.Error(ctx.Err))
	case (ctx.Attempt-5)%20 == 0:
		logger.Warn(fmt.Sprintf("%s is still failing, will retry", desc), zap.Int("attempts", ctx.Attempt), zap.Duration("duration", time.Since(ctx.T0)), zap.Error(ctx.Err))
	}
}

type RetryOption func(*retryOptions)

var defaultRetryOptions = retryOptions{
	InitialDelay: 500 * time.Millisecond,
	MaxDelay:     30 * time.Second,
	MinDelay:     500 * time.Millisecond,
	Factor:       1.5,
	MaxAttempts:  0,
	Logger:       func(logContext RetryContext) {},
}

func RetryWithInitialDelay(d time.Duration) RetryOption {
	return func(o *retryOptions) {
		o.InitialDelay = d
	}
}

func RetryWithMaxDelay(d time.Duration) RetryOption {
	return func(o *retryOptions) {
		o.MaxDelay = d
	}
}

func RetryWithFactor(f float64) RetryOption {
	return func(o *retryOptions) {
		o.Factor = f
	}
}

func RetryWithMaxAttempts(n int) RetryOption {
	return func(o *retryOptions) {
		o.MaxAttempts = n
	}
}

func RetryWithLogger(l func(logContext RetryContext)) RetryOption {
	if l == nil {
		l = func(logContext RetryContext) {}
	}
	return func(o *retryOptions) {
		o.Logger = l
	}
}
