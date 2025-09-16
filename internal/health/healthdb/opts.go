package healthdb

import (
	"go.uber.org/zap"
)

type opts struct {
	logger *zap.Logger
}

type Option func(*opts)

// WithLogger is an option to set the logger used by the store.
func WithLogger(logger *zap.Logger) Option {
	return func(s *opts) {
		s.logger = logger
	}
}
