package pull

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jonboulle/clockwork"
	"go.uber.org/zap"
)

type Option func(opts *changeOpts)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *changeOpts) {
		opts.logger = logger
	}
}

func WithPullFallback(initial, max time.Duration) Option {
	return func(opts *changeOpts) {
		opts.fallbackInitialDelay = initial
		opts.fallbackMaxDelay = max
	}
}

func WithPullFallbackJitter(jitter float64) Option {
	return func(opts *changeOpts) {
		opts.fallbackJitter = jitter
	}
}

func WithPollDelay(delay time.Duration) Option {
	return func(opts *changeOpts) {
		opts.pollDelay = delay
	}
}

func withClock(clock clockwork.Clock) Option {
	return func(opts *changeOpts) {
		opts.clock = clock
	}
}

const (
	DefaultPollDelay = time.Second
	DefaultRetryInit = 100 * time.Millisecond
	DefaultRetryMax  = 10 * time.Second
)

var defaultChangeOpts = []Option{
	withClock(clockwork.NewRealClock()),
	WithLogger(zap.NewNop()),
	WithPullFallback(DefaultRetryInit, DefaultRetryMax),
	WithPullFallbackJitter(backoff.DefaultRandomizationFactor),
	WithPollDelay(DefaultPollDelay),
}

type changeOpts struct {
	logger    *zap.Logger
	pollDelay time.Duration

	fallbackInitialDelay time.Duration
	fallbackMaxDelay     time.Duration
	fallbackJitter       float64

	clock clockwork.Clock
}

func calcOpts(opts ...Option) changeOpts {
	out := &changeOpts{}
	for _, opt := range defaultChangeOpts {
		opt(out)
	}
	for _, opt := range opts {
		opt(out)
	}
	return *out
}

func (co changeOpts) backoff(opts ...backoff.ExponentialBackOffOpts) *backoff.ExponentialBackOff {
	opts = append([]backoff.ExponentialBackOffOpts{
		backoff.WithClockProvider(co.clock),
		backoff.WithInitialInterval(co.fallbackInitialDelay),
		backoff.WithMaxInterval(co.fallbackMaxDelay),
		backoff.WithRandomizationFactor(co.fallbackJitter),
		backoff.WithMaxElapsedTime(0), // no max time
	}, opts...)
	return backoff.NewExponentialBackOff(opts...)
}
