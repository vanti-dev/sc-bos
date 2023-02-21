package prioritypb

import (
	"context"

	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/priority"
	"github.com/vanti-dev/sc-bos/pkg/util/chans"
	"golang.org/x/sync/errgroup"
)

type indexedChange[T any] struct {
	i   int
	val T
}

// collect inserts any values from pullers into dst at the corresponding index.
// dst should have a len greater or equal to pullers.
func collect[T any](ctx context.Context, dst *priority.List[T], pullers ...pull.Fetcher[T]) error {
	group, ctx := errgroup.WithContext(ctx)
	changes := make(chan indexedChange[T])
	defer close(changes)
	for i, puller := range pullers {
		localChanges := make(chan T)
		puller := puller
		group.Go(func() error {
			defer close(localChanges)
			return pull.Changes[T](ctx, puller, localChanges)
		})
		i := i
		group.Go(func() error {
			for localChange := range localChanges {
				err := chans.SendContext(ctx, changes, indexedChange[T]{
					i:   i,
					val: localChange,
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	group.Go(func() error {
		for change := range changes {
			dst.Set(change.i, change.val)
		}
		return nil
	})

	return group.Wait()
}
