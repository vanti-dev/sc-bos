package memstore

import (
	"time"
)

type Option func(*Store)

// WithNow is an option to set where the store gets the current time from.
func WithNow(now func() time.Time) Option {
	return func(s *Store) {
		s.now = now
	}
}

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
