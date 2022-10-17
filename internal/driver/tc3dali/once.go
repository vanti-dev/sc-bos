package tc3dali

import (
	"context"
)

// once works like sync.Once but supports callbacks that return errors and callers that use contexts.
// The callback will be allowed to succeed only once - once the callback returns a nil error, no further attempts will
// be made. If it returns a non-nil error, then that call to Do will fail, but other calls will run the call
type once struct {
	sem       chan struct{} // can't select on a mutex, so we use a channel as a binary semaphore instead
	completed bool          // set when the callback has run with a non-nil error - won't run again after that
}

func newOnce() *once {
	o := &once{
		sem: make(chan struct{}, 1),
	}
	// prefill sem to indicate that it's available
	o.sem <- struct{}{}
	return o
}

// Do will run f if no previous call to Do has resulted in a nil err.
// At most once invocation will run concurrently - if a function is already running, then we wait for it to finish.
// If the wait context is cancelled before the concurrent invocation completes, then done=false is returned
// along with wait.Err().
// If we are waiting for a concurrent invocation, and it returns a non-nil error, then f will be invoked immediately
// to try again. This means that the error returned by f will only be returned to a single caller of Do.
// When no concurrent invocation is running, f will be invoked synchronously. The wait context will not be used while
// f runs, so the caller is responsible for cancelling it if required.
func (o *once) Do(wait context.Context, f func() error) (err error, done bool) {
	select {
	case <-wait.Done():
		return wait.Err(), false
	case <-o.sem: // acquire the lock
	}
	defer func() {
		// release the lock
		// as this is a buffered channel which we have emptied above, this will never block
		o.sem <- struct{}{}
	}()

	if o.completed {
		return nil, true
	}

	err = f()
	if err == nil {
		o.completed = true
	}
	return err, true
}
