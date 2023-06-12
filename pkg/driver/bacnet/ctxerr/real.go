package ctxerr

import (
	"context"
	"errors"

	"go.uber.org/multierr"
)

// Cause returns err or ctx.Err depending on if ctx is done.
// A lot of the bacnet client code returns unhelpful, non-cause errors which make detecting context cancellation difficult.
// This is being fixed, so if err is a context error, err is returned as-is.
func Cause(ctx context.Context, err error) error {
	if Is(err) {
		return err
	}

	select {
	case <-ctx.Done():
		return multierr.Combine(ctx.Err(), err)
	default:
		return err
	}
}

// Is returns true if err is a context error.
func Is(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
