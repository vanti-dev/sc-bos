package service

import (
	"encoding/json"
	"time"
)

type Option[T any] interface {
	apply(l *Service[T])
}

func DefaultOpts[C any]() []Option[C] {
	return []Option[C]{
		WithNow[C](time.Now),
		WithParser(func(data []byte) (C, error) {
			var c C
			err := json.Unmarshal(data, &c)
			return c, err
		}),
	}
}

// OptionFunc adapts a func of the correct signature to implement Option.
type OptionFunc[T any] func(l *Service[T])

func (o OptionFunc[T]) apply(l *Service[T]) {
	o(l)
}

// WithParser configures a Service to use the given parse func instead of the default json.Unmarshaler.
func WithParser[T any](parse ParseFunc[T]) Option[T] {
	return OptionFunc[T](func(l *Service[T]) {
		l.parse = parse
	})
}

// WithNow configures a service with a custom time functions instead of the default time.Now.
// Useful for testing.
func WithNow[T any](now func() time.Time) Option[T] {
	return OptionFunc[T](func(l *Service[T]) {
		l.now = now
	})
}
