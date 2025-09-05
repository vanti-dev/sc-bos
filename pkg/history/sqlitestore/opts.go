package sqlitestore

import (
	"time"

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

type writeOpts struct {
	enableMaxCount bool
	maxCount       int64
	trimTime       time.Time
	trimAge        time.Duration
}

type WriteOption func(*writeOpts)

// WithMaxCount will automatically delete old records after the write operation is complete, so that there at most
// maxCount records for each source involved in the write operation.
// Like calling Database.TrimTime after each write, but happens within the same transaction.
func WithMaxCount(maxCount int64) WriteOption {
	if maxCount < 0 {
		panic("maxCount must be >= 0")
	}
	return func(o *writeOpts) {
		o.maxCount = maxCount
		o.enableMaxCount = true
	}
}

// WithNoMaxCount disables automatic deletion, cancelling a previous WithMaxCount option.
func WithNoMaxCount() WriteOption {
	return func(o *writeOpts) {
		o.enableMaxCount = false
	}
}

// WithEarliestTime will delete all records older than the given time after the write operation is complete.
// Only sources involved in the write operation are affected.
// Like calling Database.TrimTime after each write, but happens within the same transaction.
// Passing a zero time disables automatic deletion, cancelling a previous WithEarliestTime or WithMaxAge option.
func WithEarliestTime(t time.Time) WriteOption {
	return func(o *writeOpts) {
		o.trimAge = 0
		o.trimTime = t
	}
}

// WithMaxAge will delete all records older than the given duration after the write operation is complete.
// Works the same as WithEarliestTime, but the time is calculated as time.Now().Add(-d) at the time of each write operation.
// Overrides any previous WithEarliestTime or WithMaxAge option.
// Passing a zero duration disables automatic deletion, cancelling a previous WithEarliestTime or WithMaxAge option.
func WithMaxAge(d time.Duration) WriteOption {
	return func(o *writeOpts) {
		o.trimTime = time.Time{}
		o.trimAge = d
	}
}
