package chans

import (
	"context"
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
