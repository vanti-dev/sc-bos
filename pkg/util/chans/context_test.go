package chans

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSendContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := make(chan string, 1)
		v := "hello"
		ctx := context.Background()
		err := SendContext(ctx, c, v)
		if err != nil {
			t.Fatal(err)
		}
		got, err := RecvWithin(c, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		if got != v {
			t.Fatalf("want %v, got %v", v, got)
		}
	})

	t.Run("done", func(t *testing.T) {
		c := make(chan string)
		ctx, stop := context.WithCancel(context.Background())
		stop()
		err := SendContext(ctx, c, "hello")
		if err == nil {
			t.Fatal("expecting error")
		}
	})
}

func TestRecvContextFunc(t *testing.T) {
	t.Run("first success", func(t *testing.T) {
		c := make(chan string, 1)
		c <- "hello"
		ctx := context.Background()
		got, err := RecvContextFunc(ctx, c, func(v string) error {
			if v != "hello" {
				return ErrSkip
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if got != "hello" {
			t.Fatalf("want %v, got %v", "hello", got)
		}
	})
	t.Run("skip", func(t *testing.T) {
		items := []string{"hello", "world", "after"}
		c := make(chan string, len(items))
		for _, v := range items {
			c <- v
		}
		var checked []string
		ctx := context.Background()
		got, err := RecvContextFunc(ctx, c, func(v string) error {
			checked = append(checked, v)
			if v != "world" {
				return ErrSkip
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if got != "world" {
			t.Fatalf("want %v, got %v", "world", got)
		}
		if diff := cmp.Diff(items[:2], checked); diff != "" {
			t.Fatalf("checked mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("fn error", func(t *testing.T) {
		wantErr := errors.New("fn error")
		c := make(chan string, 1)
		c <- "hello"
		ctx := context.Background()
		got, err := RecvContextFunc(ctx, c, func(v string) error {
			return wantErr
		})
		if !errors.Is(err, wantErr) {
			t.Fatalf("want %v, got %v", wantErr, err)
		}
		if got != "" {
			t.Fatalf("want empty string, got %v", got)
		}
	})
	t.Run("context done", func(t *testing.T) {
		c := make(chan string)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // immediately cancel the context
		got, err := RecvContextFunc(ctx, c, func(v string) error {
			t.Fatalf("filter func called when context is done")
			return nil // this should not be called
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("want %v, got %v", context.Canceled, err)
		}
		if got != "" {
			t.Fatalf("want empty string, got %v", got)
		}
	})
}
