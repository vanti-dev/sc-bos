package minirx

import (
	"context"
	"iter"
	"sync"
)

type Value[T any] struct {
	m         sync.RWMutex
	v         T
	listeners map[*valueListener[T]]struct{}
}

func NewValue[T any](v T) *Value[T] {
	return &Value[T]{
		v:         v,
		listeners: make(map[*valueListener[T]]struct{}),
	}
}

func (v *Value[T]) Get() T {
	v.m.RLock()
	defer v.m.RUnlock()
	return v.v
}

func (v *Value[T]) Set(val T) {
	v.m.Lock()
	defer v.m.Unlock()
	v.v = val

	for listener := range v.listeners {
		if !listener.Dispatch(val) {
			delete(v.listeners, listener)
		}
	}
}

func (v *Value[T]) Pull(ctx context.Context) (initial T, changes iter.Seq[T]) {
	v.m.Lock()
	defer v.m.Unlock()
	initial = v.v
	l := &valueListener[T]{}
	l.c = sync.Cond{L: &l.m}
	l.cleanupAfterFunc = context.AfterFunc(ctx, l.Stop)
	v.listeners[l] = struct{}{}
	return initial, l.Iter
}

type valueListener[T any] struct {
	cleanupAfterFunc func() bool

	m           sync.Mutex
	c           sync.Cond
	stopped     bool
	hasNewValue bool
	value       T
}

func (l *valueListener[T]) Iter(yield func(T) bool) {
	l.m.Lock()
	defer l.m.Unlock()
	for {
		if l.stopped {
			return
		}
		if l.hasNewValue {
			l.hasNewValue = false
			value := l.value

			l.m.Unlock()
			keepIterating := yield(value)
			l.m.Lock()

			if !keepIterating {
				l.stopped = true
				return
			}
		} else {
			l.c.Wait()
		}
	}
}

func (l *valueListener[T]) Stop() {
	_ = l.cleanupAfterFunc()
	l.m.Lock()
	l.stopped = true
	l.m.Unlock()
	l.c.Broadcast()
}

func (l *valueListener[T]) Dispatch(v T) bool {
	l.m.Lock()
	defer l.m.Unlock()

	if l.stopped {
		return false
	}

	l.value = v
	l.hasNewValue = true
	l.c.Broadcast()

	return true
}
