package priority

import (
	"context"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

// List represents an ordered list of values where the first (lowest index) value is treated as the active one.
type List[T any] struct {
	entries   []Entry[T]
	notifyTop *minibus.Bus[Entry[T]]
	notifyAll *minibus.Bus[[]Entry[T]]
	Now       func() time.Time
}

func NewList[T any](len int) *List[T] {
	return &List[T]{
		entries:   make([]Entry[T], len),
		notifyTop: &minibus.Bus[Entry[T]]{},
		notifyAll: &minibus.Bus[[]Entry[T]]{},
		Now:       time.Now,
	}
}

// All returns all entries currently set.
func (l *List[T]) All() []Entry[T] {
	return l.entries
}

// Get returns the lowest index entry that is set and that index.
// If no entries are set returns an Entry whose Set property is false and the length of this list.
func (l *List[T]) Get() (Entry[T], int) {
	for i, entry := range l.entries {
		if entry.Set {
			return entry, i
		}
	}
	return Entry[T]{}, len(l.entries)
}

// Len returns how many entries this list can hold.
func (l *List[T]) Len() int {
	return len(l.entries)
}

// Set updates the entries value at position i to v.
func (l *List[T]) Set(i int, v T) {
	_, top := l.Get()
	l.entries[i] = Entry[T]{
		Value:      v,
		Set:        true,
		UpdateTime: l.Now(),
	}
	if i <= top {
		go l.sendTop()
	}
}

// Clear unsets the value at index i.
func (l *List[T]) Clear(i int) {
	_, top := l.Get()
	l.entries[i] = Entry[T]{}
	if top == i {
		go l.sendTop()
	}
}

// Listen returns a chan that emits when Get would return a new value.
// The chan will be closed when ctx is done.
func (l *List[T]) Listen(ctx context.Context) <-chan Entry[T] {
	return l.notifyTop.Listen(ctx)
}

func (l *List[T]) ListenAll(ctx context.Context) <-chan []Entry[T] {
	return l.notifyAll.Listen(ctx)
}

// sentTop emits the top priority item onto the notifyTop bus.
func (l *List[T]) sendTop() {
	top, _ := l.Get()
	l.notifyTop.Send(context.Background(), top)
}

type Entry[T any] struct {
	Value      T
	Set        bool
	UpdateTime time.Time
}
