package pull

import (
	"context"
)

// Fetcher can pull and poll for changes.
type Fetcher[C any] interface {
	// Pull blocks subscribing to changes and sends them on the given channel.
	// Retries and reconnects are managed by the caller.
	Pull(ctx context.Context, changes chan<- C) error
	// Poll queries for the current value sending it on the given channel.
	// Delays between polls are managed by the caller.
	Poll(ctx context.Context, changes chan<- C) error
}

// NewFetcher returns a Fetcher backed by the given pull and poll functions.
func NewFetcher[C any](pull, poll func(ctx context.Context, changes chan<- C) error) Fetcher[C] {
	return fetcher[C]{pull: pull, poll: poll}
}

type fetcher[C any] struct {
	pull, poll func(ctx context.Context, changes chan<- C) error
}

func (f fetcher[C]) Pull(ctx context.Context, changes chan<- C) error {
	return f.pull(ctx, changes)
}

func (f fetcher[C]) Poll(ctx context.Context, changes chan<- C) error {
	return f.poll(ctx, changes)
}
