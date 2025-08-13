package chans

import (
	"context"
	"errors"
)

// SendContext sends v on c unless ctx is Done.
func SendContext[T any](ctx context.Context, c chan<- T, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c <- v:
		return nil
	}
}

// RecvContext receives from c unless ctx is Done.
func RecvContext[T any](ctx context.Context, c <-chan T) (T, error) {
	select {
	case <-ctx.Done():
		var t T
		return t, ctx.Err()
	case v := <-c:
		return v, nil
	}
}

var (
	// ErrSkip is used to indicate that a value should be skipped when processing a channel.
	ErrSkip = errors.New("skip")
)

// RecvContextFunc applies fn to each value received from c until fn returns nil, ctx is done, or an error occurs.
// fn should return ErrSkip to continue receiving values without returning an error.
func RecvContextFunc[T any](ctx context.Context, c <-chan T, fn func(T) error) (accepted T, _ error) {
	for {
		select {
		case <-ctx.Done():
			return accepted, ctx.Err()
		case v, ok := <-c:
			if !ok {
				return v, ErrClosed
			}
			err := fn(v)
			switch {
			case err == nil:
				return v, nil
			case errors.Is(err, ErrSkip): // continue for loop
			default:
				return accepted, err
			}
		}
	}
}
