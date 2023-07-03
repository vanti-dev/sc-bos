package node

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAnnounceContext(t *testing.T) {
	var names []string
	an := AnnouncerFunc(func(name string, features ...Feature) Undo {
		i := len(names)
		names = append(names, name)
		return func() {
			names[i] += ":undo"
		}
	})

	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	a := AnnounceContext(ctx, an)
	a.Announce("a")
	ub := a.Announce("b")
	a.Announce("c")

	if diff := cmp.Diff(names, []string{"a", "b", "c"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}

	ub()
	if diff := cmp.Diff(names, []string{"a", "b:undo", "c"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}

	a.Announce("d")
	if diff := cmp.Diff(names, []string{"a", "b:undo", "c", "d"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}

	stop()
	cont := make(chan struct{})
	go func() {
		<-ctx.Done()
		time.Sleep(100 * time.Millisecond)
		close(cont)
	}()
	<-cont
	if diff := cmp.Diff(names, []string{"a:undo", "b:undo", "c:undo", "d:undo"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}

	a.Announce("e")
	if diff := cmp.Diff(names, []string{"a:undo", "b:undo", "c:undo", "d:undo"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}
}
