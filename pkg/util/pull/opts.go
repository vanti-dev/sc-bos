package pull

import (
	"time"

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
	WithPollDelay(DefaultPollDelay),
}

type changeOpts struct {
	logger    *zap.Logger
	pollDelay time.Duration

	fallbackInitialDelay time.Duration
	fallbackMaxDelay     time.Duration

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
