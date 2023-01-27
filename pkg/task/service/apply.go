package service

import (
	"context"
)

// MonoApply wraps apply to ensure that consecutive calls to the returned ApplyFunc will cancel the context passed to apply.
func MonoApply[C any](apply ApplyFunc[C]) ApplyFunc[C] {
	var lastCtx context.Context
	var stopLast context.CancelFunc
	return func(ctx context.Context, config C) error {
		if stopLast != nil {
			stopLast()
		}
		lastCtx, stopLast = context.WithCancel(ctx)
		return apply(lastCtx, config)
	}
}
