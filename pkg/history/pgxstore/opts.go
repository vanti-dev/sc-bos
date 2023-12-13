package pgxstore

import (
	"time"

	"go.uber.org/zap"
)

type Option func(*Store)

// WithMaxAge is an option to set the maximum age of records in the store.
func WithMaxAge(maxAge time.Duration) Option {
	return func(s *Store) {
		s.maxAge = maxAge
	}
}

// WithMaxCount is an option to set the maximum number of records in the store.
func WithMaxCount(maxCount int64) Option {
	return func(s *Store) {
		s.maxCount = maxCount
	}
}

// WithLogger is an option to set the logger used by the store.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Store) {
		s.logger = logger
	}
}
