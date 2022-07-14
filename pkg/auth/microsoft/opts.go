package microsoft

import (
	"time"
)

type Option func(o *opts)

func WithClock(now func() time.Time) Option {
	return func(o *opts) {
		o.now = now
	}
}

type opts struct {
	now func() time.Time
}

func defaultOpts() opts {
	return opts{now: time.Now}
}

func resolveOpts(options ...Option) opts {
	resolved := defaultOpts()
	for _, opt := range options {
		opt(&resolved)
	}
	return resolved
}
