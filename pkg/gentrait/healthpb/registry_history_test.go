package healthpb

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// ExampleRegistry_history shows how to connect a Registry to a historical datastore.
// Checks are initially loaded from the datastore to populate their initial state,
// updates are recorded to the datastore as they occur.
func ExampleRegistry_history() {
	// history is a sample in-memory store holding the last known state of a check.
	history := store{
		// our store has a previous record for device1's "example" check
		// indicating it was previously in an error state.
		key{"device1", "example"}: &gen.HealthCheck{
			Id: "example",
			Check: &gen.HealthCheck_Check{
				State:        gen.HealthCheck_Check_ABNORMAL,
				AbnormalTime: timeAt(12, 42),
				LastError:    ErrorToProto(errors.New("out of paper")),
			},
		},
	}

	registry := &Registry{
		onCheckCreate: func(name string, c *gen.HealthCheck) *gen.HealthCheck {
			old, err := history.LoadLastCheck(name, c.Id)
			if err != nil {
				return nil
			}
			if old == nil {
				// don't store the initial check as it likely has no state
				return nil
			}
			// copy over any updates from the new check into old,
			// anything not set on c (like timestamps or state) will be preserved
			mergeCheck(proto.Merge, old, c)
			return old
		},
		onCheckUpdate: func(name string, c *gen.HealthCheck) {
			err := history.SaveCheck(name, c)
			if err != nil {
				panic(err)
			}
		},
	}

	// create the check for device1 owned by "example"
	exampleChecks := registry.ForOwner("example")
	dev1Check, err := exampleChecks.NewErrorCheck("device1", &gen.HealthCheck{
		DisplayName: "Paper level",
	})
	if err != nil {
		panic(err)
	}
	defer dev1Check.Dispose()

	// check the initial state was loaded from history
	storeCheck1 := registry.GetCheck("device1", "example")
	fmt.Printf("Before update: %q=%v\n", storeCheck1.GetDisplayName(), storeCheck1.GetCheck().GetState())

	// perform a check
	dev1Check.UpdateError(nil) // all good now
	storeCheck2 := registry.GetCheck("device1", "example")
	fmt.Printf("After update: %q=%v\n", storeCheck2.GetDisplayName(), storeCheck2.GetCheck().GetState())
	// Output:
	// Before update: "Paper level"=ABNORMAL
	// After update: "Paper level"=NORMAL
}

func timeAt(h, m int) *timestamppb.Timestamp {
	return timestamppb.New(time.Date(2025, 9, 4, h, m, 0, 0, time.UTC))
}

type store map[key]*gen.HealthCheck

type key struct {
	name, id string
}

func (s store) LoadLastCheck(name, id string) (*gen.HealthCheck, error) {
	c := s[key{name, id}]
	if c == nil {
		return nil, nil
	}
	return proto.Clone(c).(*gen.HealthCheck), nil
}

func (s store) SaveCheck(name string, c *gen.HealthCheck) error {
	s[key{name, c.Id}] = c
	return nil
}
