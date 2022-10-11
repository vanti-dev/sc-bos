package minibus

import (
	"container/list"
	"context"
	"sync"
)

type Bus[T any] struct {
	listenerM sync.RWMutex
	listeners list.List // use linked list to make edits easier - i.e. removing from the middle
}

// Send sends a message to all active listeners.
// The returned chan will emit how many listeners were notified before ctx is done.
func (b *Bus[T]) Send(ctx context.Context, event T) <-chan int {
	sentToAll := make(chan int, 1)

	b.listenerM.RLock()
	defer b.listenerM.RUnlock()

	size := b.listeners.Len()
	if size == 0 {
		sentToAll <- 0
		return sentToAll
	}
	sentToOne := make(chan bool)
	for el := b.listeners.Front(); el != nil; el = el.Next() {
		lis := el.Value.(*listener[T])
		go func() {
			sentToOne <- lis.send(ctx, event)
		}()
	}

	go func() {
		defer close(sentToAll)

		remaining := size
		var sent int
		for remaining > 0 {
			select {
			case <-ctx.Done():
				sentToAll <- sent
				return
			case success := <-sentToOne:
				remaining--
				if success {
					sent++
				}

			}
		}
		sentToAll <- sent
	}()

	return sentToAll
}

func (b *Bus[T]) Listen(ctx context.Context) <-chan T {
	ch := make(chan T)

	l := &listener[T]{
		ch:  ch,
		ctx: ctx,
	}

	b.listenerM.Lock()
	el := b.listeners.PushBack(l)
	b.listenerM.Unlock()

	go func() {
		<-ctx.Done()
		b.listenerM.Lock()
		defer b.listenerM.Unlock()
		b.listeners.Remove(el)
		l.stop()
	}()

	return ch
}

type listener[T any] struct {
	m   sync.RWMutex
	ch  chan T
	ctx context.Context
}

func (l *listener[T]) send(ctx context.Context, event T) (sent bool) {
	l.m.RLock()
	defer l.m.RUnlock()

	select {
	case <-ctx.Done():
		// send context cancelled
		return false

	case <-l.ctx.Done():
		// listen context cancelled
		return false

	case l.ch <- event:
		// event sent successfully
		return true
	}
}

func (l *listener[T]) stop() {
	l.m.Lock()
	defer l.m.Unlock()
	if l.ch != nil {
		close(l.ch)
		l.ch = nil
	}
}

func (l *listener[T]) alive() bool {
	return l.ctx.Err() == nil
}
