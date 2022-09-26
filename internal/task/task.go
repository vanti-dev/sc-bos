package task

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

// A Task is a unit of work which can be restarted if it fails.
type Task func(ctx context.Context) (Next, error)

// Next allows a task to specify a preferred retry behaviour.
type Next int

const (
	// Normal mode of operation - Task is retried if it returns a non-nil error.
	Normal Next = iota
	// StopNow will prevent the task from being restarted, even if it returns a non-nil error.
	StopNow
	// RetryNow will restart the task immediately, without any delay. The backoff will also be reset.
	RetryNow
	// ResetBackoff will reset the Task restart delay to its starting value, and then continue as normal.
	ResetBackoff
)

func (n Next) String() string {
	switch n {
	case Normal:
		return "Normal"
	case StopNow:
		return "StopNow"
	case RetryNow:
		return "RetryNow"
	case ResetBackoff:
		return "ResetBackoff"
	}
	return fmt.Sprintf("Next(%d)", int(n))
}

type Option func(o *Runner)

// RetryUnlimited allows a Task to be restarted forever. Used as an argument to WithRetry.
const RetryUnlimited int = math.MaxInt

// WithRetry places a limit on the number of times a Task may be restarted before we give up.
// By default, a Task will only be run once. Pass RetryUnlimited to retry forever.
func WithRetry(attempts int) Option {
	return func(o *Runner) {
		o.attemptsRemaining = attempts
	}
}

// WithRetryDelay adds a fixed delay between a Task returning and the next attempt starting.
// By default, there is no delay and retries happen immediately.
func WithRetryDelay(delay time.Duration) Option {
	return WithBackoff(delay, delay)
}

// WithBackoff adds exponential backoff when retrying tasks. The delay begins at start, and is capped at max.
// After each attempt, the delay increases by a factor of 1.5.
func WithBackoff(start time.Duration, max time.Duration) Option {
	return func(o *Runner) {
		o.delay = start
		o.minDelay = start
		o.maxDelay = max
	}
}

// WithTimeout imposes a time limit on each individual invocation of the Task.
func WithTimeout(timeout time.Duration) Option {
	return func(o *Runner) {
		o.timeout = timeout
	}
}

// WithErrorLogger will log to the provided logger every time the Task returns a non-nil error.
func WithErrorLogger(logger *zap.Logger) Option {
	return func(o *Runner) {
		o.logger = logger
	}
}

// Run will run a Task to completion in a blocking fashion.
// By default, the Task is run once. Pass Options in order to add automatic retry with backoff, logging etc.
//
// This is a convenience function that constructs a Runner and calls Runner.Step in a loop until the runner completes.
func Run(ctx context.Context, task Task, options ...Option) error {
	r := NewRunner(task, options...)
	for {
		err, again, delay := r.Step(ctx)
		if !again {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func NewRunner(task Task, options ...Option) *Runner {
	r := &Runner{
		task:              task,
		attemptsRemaining: 1,
	}
	for _, opt := range options {
		opt(r)
	}
	return r
}

// A Runner allows the caller control over the invocation of a Task.
// Call Step to run the Task - the return values will tell you if and when to call Step again.
type Runner struct {
	task              Task
	attemptsRemaining int
	delay             time.Duration
	minDelay          time.Duration
	maxDelay          time.Duration
	timeout           time.Duration
	logger            *zap.Logger
}

// Step will run this Runner's Task once.
// The error from the task is returned as err.
// If the task should be run again due to the applicable retry options, then the return value again will be true
// and the required delay before the next invocation is returned in delay.
//
// Callers will generally want to call Step in a loop until again=false.
func (r *Runner) Step(ctx context.Context) (err error, again bool, delay time.Duration) {
	var (
		subctx context.Context
		cancel context.CancelFunc
	)
	if r.timeout != 0 {
		subctx, cancel = context.WithTimeout(ctx, r.timeout)
	} else {
		subctx, cancel = context.WithCancel(ctx)
	}
	var next Next
	next, err = r.task(subctx)
	cancel()

	if err != nil && r.logger != nil {
		r.logger.Error("task returned an error", zap.Error(err), zap.String("next", next.String()))
	}

	if r.attemptsRemaining != RetryUnlimited {
		r.attemptsRemaining--
	}
	if r.attemptsRemaining <= 0 || err == nil {
		again = false
		return
	}

	switch next {
	case Normal:
		again = true
		delay = r.delay
		r.delay = backoff(r.delay, r.maxDelay)
		return
	case StopNow:
		again = false
		return
	case RetryNow:
		again = true
		delay = 0
		r.delay = r.minDelay
		return
	case ResetBackoff:
		again = true
		delay = r.minDelay
		r.delay = backoff(r.minDelay, r.maxDelay)
		return
	default:
		panic(fmt.Sprintf("invalid task.Next value %v", next))
	}
}

func backoff(delay time.Duration, max time.Duration) time.Duration {
	delay += delay / 2
	if delay > max {
		delay = max
	}
	return delay
}
