package concurrent

import (
	"context"
	"sync"
	"time"
)

type Value[T any] struct {
	m        sync.RWMutex
	value    T
	modified time.Time
	bus      Bus[ValueEvent[T]]
}

func NewValue[T any](initial T) *Value[T] {
	return &Value[T]{
		value:    initial,
		modified: time.Now(),
	}
}

func (v *Value[T]) Get() (value T, modified time.Time) {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.value, v.modified
}

func (v *Value[T]) Set(ctx context.Context, value T) (old T, ok bool) {
	v.m.Lock()
	old = v.value
	v.value = value
	timestamp := time.Now()
	v.modified = timestamp
	v.m.Unlock()

	ok = v.bus.Send(ctx, ValueEvent[T]{
		Timestamp: timestamp,
		Old:       old,
		New:       value,
	})
	return
}

func (v *Value[T]) Changes(ctx context.Context, backpressure bool, bufferSize int) (value T, changes <-chan ValueEvent[T]) {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.value, v.bus.Listen(ctx, backpressure, bufferSize)
}

type ValueEvent[T any] struct {
	Timestamp time.Time
	Old       T
	New       T
}
