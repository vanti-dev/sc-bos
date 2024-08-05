package node

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAnnounceContext(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	var names []string
	var m sync.Mutex
	an := AnnouncerFunc(func(name string, features ...Feature) Undo {
		m.Lock()
		i := len(names)
		names = append(names, name)
		m.Unlock()
		return func() {
			m.Lock()
			names[i] += ":undo"
			m.Unlock()
		}
	})

	a := AnnounceContext(ctx, an)
	a.Announce("a")
	ub := a.Announce("b")
	a.Announce("c")

	m.Lock()
	if diff := cmp.Diff(names, []string{"a", "b", "c"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}
	m.Unlock()

	ub()
	m.Lock()
	if diff := cmp.Diff(names, []string{"a", "b:undo", "c"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}
	m.Unlock()

	a.Announce("d")
	m.Lock()
	if diff := cmp.Diff(names, []string{"a", "b:undo", "c", "d"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}
	m.Unlock()

	stop()
	cont := make(chan struct{})
	go func() {
		<-ctx.Done()
		time.Sleep(100 * time.Millisecond)
		close(cont)
	}()
	<-cont
	m.Lock()
	if diff := cmp.Diff(names, []string{"a:undo", "b:undo", "c:undo", "d:undo"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}
	m.Unlock()

	a.Announce("e")
	m.Lock()
	if diff := cmp.Diff(names, []string{"a:undo", "b:undo", "c:undo", "d:undo"}); diff != "" {
		t.Errorf("unexpected names (-want +got):\n%s", diff)
	}
	m.Unlock()
}
