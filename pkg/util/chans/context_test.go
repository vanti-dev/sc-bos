package chans

import (
	"context"
	"testing"
	"time"
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
