package pull

import (
	"time"

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

func WithPollDelay(delay time.Duration) Option {
	return func(opts *changeOpts) {
		opts.pollDelay = delay
	}
}

var defaultChangeOpts = []Option{
	WithLogger(zap.NewNop()),
	WithPullFallback(100*time.Millisecond, 10*time.Second),
	WithPollDelay(time.Second),
}

type changeOpts struct {
	logger    *zap.Logger
	pollDelay time.Duration

	fallbackInitialDelay time.Duration
	fallbackMaxDelay     time.Duration
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
