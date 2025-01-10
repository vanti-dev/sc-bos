package rx

import (
	"context"
	"sync"

	"github.com/vanti-dev/sc-bos/pkg/minibus"
	"github.com/vanti-dev/sc-bos/pkg/util/slices"
)

// Set is a reactive set.
// Set is safe for concurrent use.
// Use NewSet to create a new Set.
type Set[T any] struct {
	m sync.Mutex
	v *slices.Sorted[T]
	b *minibus.Bus[Change[T]]
	c <-chan struct{} // nil when not notifying, blocks when notifying, closed when done notifying
}

func NewSet[T any](v *slices.Sorted[T]) *Set[T] {
	s := &Set[T]{
		v: v,
		b: &minibus.Bus[Change[T]]{},
	}
	return s
}

// Set adds item to the set, replacing any existing item.
func (s *Set[T]) Set(item T) (i int, old T, replaced bool, sent <-chan struct{}) {
	s.m.Lock()
	defer s.m.Unlock()
	i, old, replaced = s.v.Set(item)
	e := Change[T]{
		Type: Add,
		New:  item,
	}
	if replaced {
		e.Type = Update
		e.Old = old
	}

	s.c = send(s.b, e, s.c)
	return i, old, replaced, s.c
}

func (s *Set[T]) Remove(item T) (i int, old T, removed bool, sent <-chan struct{}) {
	s.m.Lock()
	defer s.m.Unlock()
	i, old, removed = s.v.Remove(item)
	if removed {
		e := Change[T]{
			Type: Remove,
			Old:  old,
		}
		s.c = send(s.b, e, s.c)
	}
	return i, old, removed, s.c
}

// Replace replaces the set with items.
func (s *Set[T]) Replace(items []T) (added, deleted, updated *slices.Sorted[T], sent <-chan struct{}) {
	s.m.Lock()
	defer s.m.Unlock()

	var sends []<-chan struct{}

	deleted = s.v.Copy()
	added = slices.NewSortedFunc(s.v.Cmp())
	updated = slices.NewSortedFunc(s.v.Cmp())
	for _, item := range items {
		deleted.Remove(item)
		_, old, replaced := s.v.Set(item)

		e := Change[T]{
			Type: Add,
			New:  item,
		}
		if replaced {
			updated.Set(old)
			e.Type = Update
			e.Old = old
		} else {
			added.Set(item)
		}
		sends = append(sends, send(s.b, e, s.c))
	}

	deleted.All(func(i int, item T) bool {
		s.v.Remove(item)
		e := Change[T]{
			Type: Remove,
			Old:  item,
		}
		sends = append(sends, send(s.b, e, s.c))
		return true
	})

	return added, deleted, updated, allDone(sends)
}

func (s *Set[T]) Len() int {
	s.m.Lock()
	defer s.m.Unlock()
	return s.v.Len()
}

// Get returns the item at index i.
func (s *Set[T]) Get(i int) T {
	s.m.Lock()
	defer s.m.Unlock()
	return s.v.Get(i)
}

// Find returns the index, item, and true if k is in the set.
func (s *Set[T]) Find(k T) (int, T, bool) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.v.Find(k)
}

// All calls the given function for each item in the set, breaking early if the function returns false.
func (s *Set[T]) All(yield func(int, T) bool) {
	s.m.Lock()
	defer s.m.Unlock()
	s.v.All(yield)
}

// Sub returns a copy of the set and a channel that receives changes.
func (s *Set[T]) Sub(ctx context.Context) (*slices.Sorted[T], <-chan Change[T]) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.v.Copy(), s.b.Listen(ctx)
}

type Change[T any] struct {
	Type ChangeType
	Old  T // non-zero during update and remove
	New  T // non-zero during add and update
}

type ChangeType int

const (
	Add ChangeType = iota
	Remove
	Update
)
