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
