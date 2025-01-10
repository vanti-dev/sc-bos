package rx

import (
	"context"

	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

// send sends e on b after receiving on wait, returning a channel that closes when b.Send is complete.
// If wait is nil, it will not be waited on.
// send does not block.
func send[T any](b *minibus.Bus[T], e T, wait <-chan struct{}) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		defer close(c)
		// wait for the last event to be sent
		if wait != nil {
			<-wait
		}
		// a non-expiring context means we always send to all listeners
		b.Send(context.Background(), e)
	}()
	return c
}

// allDone returns a channel that is closed when all channels in cs are closed.
// allDone does not block.
func allDone(cs []<-chan struct{}) <-chan struct{} {
	if len(cs) == 0 {
		// nothing to wait for
		c := make(chan struct{})
		close(c)
		return c
	}
	c := make(chan struct{})
	go func() {
		defer close(c)
		for _, w := range cs {
			<-w
		}
	}()
	return c
}
