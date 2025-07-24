package minirx

import (
	"context"
	"iter"
	"testing"
	"testing/synctest"
)

func TestValue(t *testing.T) {
	v := NewValue(1)
	if v.Get() != 1 {
		t.Errorf("Expected initial value to be 1, got %v", v.Get())
	}
	v.Set(2)
	if v.Get() != 2 {
		t.Errorf("Expected value to be updated to 2, got %v", v.Get())
	}
}

func TestValue_Pull(t *testing.T) {
	synctest.Run(func() {
		v := NewValue(1)
		ctx, cancel := context.WithCancel(context.Background())
		initial, changes := v.Pull(ctx)
		if initial != 1 {
			t.Errorf("Expected initial value to be 1, got %v", initial)
		}
		next, stop := iter.Pull(changes)
		defer stop()
		// tests that reader (which is not being consumed) does not block the write
		v.Set(2)
		v.Set(3)
		synctest.Wait()

		current, more := next()
		if !more || current != 3 {
			t.Errorf("Expected value to be updated to 3, got %v", current)
		}
		v.Set(4)
		synctest.Wait()
		current, more = next()
		if !more || current != 4 {
			t.Errorf("Expected value to be updated to 4, got %v", current)
		}
		cancel()
		synctest.Wait()
		current, more = next()
		if more {
			t.Errorf("Expected no more values after cancel, got %v", current)
		}
	})
}
