package zone

import (
	"context"
	"sort"
	"sync"
)

type Devices struct {
	mu    sync.Mutex
	names []string

	// using Context here as it has the same semantics of Done we want so we don't have to implement it ourselves
	freezeOnce sync.Once
	frozen     context.Context
	freeze     context.CancelFunc
}

func (d *Devices) init() {
	d.freezeOnce.Do(func() {
		d.frozen, d.freeze = context.WithCancel(context.Background())
	})
}

func (d *Devices) Add(names ...string) {
	d.init()
	select {
	case <-d.Frozen():
		panic("cannot add devices after Freeze")
	default:
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, name := range names {
		i := sort.SearchStrings(d.names, name)
		switch {
		case i == len(d.names):
			d.names = append(d.names, name)
		case d.names[i] != name:
			d.names = append(d.names, "")
			copy(d.names[i+1:], d.names[i:])
			d.names[i] = name
			// else d.names[i] == name, so it's already there
		}
	}
}

func (d *Devices) Names() []string {
	return d.names
}

func (d *Devices) Freeze() {
	d.init()
	d.mu.Lock()
	d.freeze()
	d.mu.Unlock()
}

func (d *Devices) Frozen() <-chan struct{} {
	d.init()
	return d.frozen.Done()
}
