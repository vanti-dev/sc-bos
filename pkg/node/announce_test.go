package node

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func ExampleAnnounceScope() {
	rootAnnouncer := New("test")

	announcer, undo := AnnounceScope(rootAnnouncer)
	// announce names in the new scope
	announcer.Announce("a")

	// undo must be called in the same goroutine as the new scope,
	// ideally just before the new scope is created.
	undo()
	announcer, undo = AnnounceScope(rootAnnouncer)
	announcer.Announce("b")
}

func TestAnnounceScope(t *testing.T) {
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
	check := func(want ...string) {
		t.Helper()
		m.Lock()
		defer m.Unlock()
		if diff := cmp.Diff(want, names); diff != "" {
			t.Errorf("unexpected names (-want +got):\n%s", diff)
		}
	}

	a, undo := AnnounceScope(an)
	a.Announce("a")
	ub := a.Announce("b")
	a.Announce("c")

	check("a", "b", "c")

	ub()
	check("a", "b:undo", "c")

	a.Announce("d")
	check("a", "b:undo", "c", "d")

	undo()
	check("a:undo", "b:undo", "c:undo", "d:undo")

	// new name should not be announced once the scope is finished
	a.Announce("e")
	check("a:undo", "b:undo", "c:undo", "d:undo")
}

func TestReplaceAnnouncer(t *testing.T) {
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
	check := func(want ...string) {
		t.Helper()
		m.Lock()
		defer m.Unlock()
		if diff := cmp.Diff(want, names); diff != "" {
			t.Errorf("unexpected names (-want +got):\n%s", diff)
		}
	}

	ra := NewReplaceAnnouncer(an)

	ctx1, cancel1 := context.WithCancel(context.Background())
	a1 := ra.Replace(ctx1)
	a1.Announce("a")
	ub := a1.Announce("b")
	a1.Announce("c")

	check("a", "b", "c")

	ub()
	check("a", "b:undo", "c")

	a1.Announce("d")
	check("a", "b:undo", "c", "d")

	// cancelling the context should undo all the announcements asynchronously
	cancel1()
	time.Sleep(50 * time.Millisecond)

	check("a:undo", "b:undo", "c:undo", "d:undo")

	// new name should not be announced once the context has been cancelled
	a1.Announce("e")
	check("a:undo", "b:undo", "c:undo", "d:undo")

	a2 := ra.Replace(context.Background())
	a2.Announce("f")

	check("a:undo", "b:undo", "c:undo", "d:undo", "f")

	a3 := ra.Replace(context.Background())
	// should have undone everything announced by a2
	check("a:undo", "b:undo", "c:undo", "d:undo", "f:undo")

	a3.Announce("g")
	check("a:undo", "b:undo", "c:undo", "d:undo", "f:undo", "g")
}
