package once

import (
	"context"
	"errors"
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once *RetryError, o *one, c chan bool) {
	once.Do(context.Background(), func() error {
		o.Increment()
		return nil
	})
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestRetryError(t *testing.T) {
	t.Run("run once", func(t *testing.T) {
		o := new(one)
		once := new(RetryError)
		c := make(chan bool)
		const N = 10
		for range N {
			go run(t, once, o, c)
		}
		for range N {
			<-c
		}
		if *o != 1 {
			t.Errorf("once failed outside run: %d is not 1", *o)
		}
	})

	t.Run("retry errors", func(t *testing.T) {
		o := new(one)
		once := new(RetryError)
		knownErr := errors.New("known")
		err := once.Do(context.Background(), func() error {
			o.Increment()
			return knownErr
		})
		if err != knownErr {
			t.Fatal("Expecting known error")
		}
		err = once.Do(context.Background(), func() error {
			o.Increment()
			return knownErr
		})
		if err != knownErr {
			t.Fatal("Expecting known error")
		}
		if v := *o; v != 2 {
			t.Fatalf("Expecting multiple calls on error")
		}
	})

	t.Run("panic", func(t *testing.T) {
		once := RetryError{recoverPanic: true}
		err := once.Do(context.Background(), func() error {
			panic("expected test panic")
		})
		if err != ErrPanicSuppressed {
			t.Fatalf("Expecting %v, got %v", ErrPanicSuppressed, err)
		}

		// make sure panicking didn't settle the once
		var ran bool
		once.Do(context.Background(), func() error {
			ran = true
			return nil
		})
		if !ran {
			t.Fatalf("Panic caused once to be settled")
		}
	})
}
