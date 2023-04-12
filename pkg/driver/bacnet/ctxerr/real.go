package ctxerr

import (
	"context"
)

// Cause returns err or ctx.Err depending on if ctx is done.
// A lot of the bacnet client code returns unhelpful, non-cause errors which make detecting context cancellation difficult.
func Cause(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}
