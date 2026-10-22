package healthpb

import (
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// TestRegistry_recreateWhileDisposing recreates a check while the delete callbacks for its disposal are running,
// as healthbounds does when a device is reconfigured.
// The callbacks must see the recreate after the delete, or the mirror loses the new check (SCB-1508).
func TestRegistry_recreateWhileDisposing(t *testing.T) {
	tests := map[string]struct {
		sibling bool // whether another check keeps the name alive
	}{
		"last check for name":          {},
		"same id with a sibling check": {sibling: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				mirrorMu sync.Mutex
				mirror   = make(map[string]map[string]bool) // name -> ids
				reg      *Registry
				once     sync.Once
				done     = make(chan struct{})
			)
			recreate := func() {
				defer close(done)
				if _, err := reg.ForOwner("owner").NewFaultCheck("dev", &HealthCheck{Id: "comms"}); err != nil {
					t.Errorf("recreate: %v", err)
				}
			}
			reg = NewRegistry(
				WithOnNameCreate(func(name string) {
					mirrorMu.Lock()
					defer mirrorMu.Unlock()
					if _, ok := mirror[name]; ok {
						return // as pkg/app does, keep the existing entry
					}
					mirror[name] = make(map[string]bool)
				}),
				WithOnCheckCreate(func(name string, c *HealthCheck) *HealthCheck {
					mirrorMu.Lock()
					defer mirrorMu.Unlock()
					if ids, ok := mirror[name]; ok {
						ids[c.Id] = true
					}
					return nil
				}),
				WithOnCheckDelete(func(name, id string) {
					once.Do(func() {
						// give the recreate the chance to complete before this callback acts
						go recreate()
						select {
						case <-done:
						case <-time.After(50 * time.Millisecond):
						}
					})
					mirrorMu.Lock()
					defer mirrorMu.Unlock()
					delete(mirror[name], id)
				}),
				WithOnNameDelete(func(name string) {
					mirrorMu.Lock()
					defer mirrorMu.Unlock()
					delete(mirror, name)
				}),
			)

			fc, err := reg.ForOwner("owner").NewFaultCheck("dev", &HealthCheck{Id: "comms"})
			if err != nil {
				t.Fatal(err)
			}
			want := map[string]map[string]bool{"dev": {"owner:comms": true}}
			if tt.sibling {
				if _, err := reg.ForOwner("other").NewFaultCheck("dev", &HealthCheck{Id: "comms"}); err != nil {
					t.Fatal(err)
				}
				want["dev"]["other:comms"] = true
			}

			fc.Dispose()
			<-done

			if reg.GetCheck("dev", "owner:comms") == nil {
				t.Fatal("recreated check not in the registry")
			}
			mirrorMu.Lock()
			defer mirrorMu.Unlock()
			if diff := cmp.Diff(want, mirror); diff != "" {
				t.Errorf("callbacks lost the recreated check (-want +got):\n%s", diff)
			}
		})
	}
}
