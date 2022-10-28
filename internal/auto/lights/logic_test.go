package lights

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func Test_processState(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		readState := NewReadState()
		writeState := NewWriteState()
		actions := newTestActions(t)
		ttl, err := processState(context.Background(), readState, writeState, actions)
		assertNoTTLOrErr(t, ttl, err)
		actions.assertNoMoreCalls()
	})

	t.Run("turn on when occupied", func(t *testing.T) {
		readState := NewReadState()
		writeState := NewWriteState()
		actions := newTestActions(t)

		readState.Config.OccupancySensors = []string{"pir01"}
		readState.Config.Lights = []string{"light01"}
		readState.Occupancy["pir01"] = &traits.Occupancy{State: traits.Occupancy_OCCUPIED}

		ttl, err := processState(context.Background(), readState, writeState, actions)
		assertNoTTLOrErr(t, ttl, err)
		actions.assertNextCall(&traits.UpdateBrightnessRequest{
			Name: "light01",
			Brightness: &traits.Brightness{
				LevelPercent: 100,
			},
		})
	})

	t.Run("ignore non-relevant occupancy", func(t *testing.T) {
		readState := NewReadState()
		writeState := NewWriteState()
		actions := newTestActions(t)

		readState.Config.OccupancySensors = []string{"pir02"}
		readState.Config.Lights = []string{"light01"}
		readState.Occupancy["pir01"] = &traits.Occupancy{State: traits.Occupancy_OCCUPIED}

		ttl, err := processState(context.Background(), readState, writeState, actions)
		assertNoTTLOrErr(t, ttl, err)
		actions.assertNoMoreCalls()
	})

	t.Run("turns lights off when unoccupied", func(t *testing.T) {
		readState := NewReadState()
		writeState := NewWriteState()
		actions := newTestActions(t)
		now := time.Unix(0, 0)

		readState.Config.Now = func() time.Time { return now }
		readState.Config.UnoccupiedOffDelay = 10 * time.Minute
		readState.Config.OccupancySensors = []string{"pir01"}
		readState.Config.Lights = []string{"light01"}
		readState.Occupancy["pir01"] = &traits.Occupancy{
			State:           traits.Occupancy_UNOCCUPIED,
			StateChangeTime: timestamppb.New(now.Add(-20 * time.Minute)),
		}

		ttl, err := processState(context.Background(), readState, writeState, actions)

		assertNoTTLOrErr(t, ttl, err)
		actions.assertNextCall(&traits.UpdateBrightnessRequest{
			Name: "light01",
			Brightness: &traits.Brightness{
				LevelPercent: 0,
			},
		})
	})
	t.Run("ttl returned when lights should change", func(t *testing.T) {
		readState := NewReadState()
		writeState := NewWriteState()
		actions := newTestActions(t)
		now := time.Unix(0, 0)

		readState.Config.Now = func() time.Time { return now }
		readState.Config.UnoccupiedOffDelay = 10 * time.Minute
		readState.Config.OccupancySensors = []string{"pir01"}
		readState.Config.Lights = []string{"light01"}
		readState.Occupancy["pir01"] = &traits.Occupancy{
			State:           traits.Occupancy_UNOCCUPIED,
			StateChangeTime: timestamppb.New(now.Add(-5 * time.Minute)),
		}

		ttl, err := processState(context.Background(), readState, writeState, actions)
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		if ttl != 5*time.Minute {
			t.Fatalf("TTL want %v, got %v", 5*time.Minute, ttl)
		}
		actions.assertNoMoreCalls()
	})
}

func assertNoTTLOrErr(t *testing.T, ttl time.Duration, err error) {
	if ttl != 0 {
		t.Fatalf("TTL want 0, got %v", ttl)
	}
	if err != nil {
		t.Fatalf("Error want <nil>, got %v", err)
	}
}

func newTestActions(t *testing.T) *testActions {
	return &testActions{t: t}
}

type testActions struct {
	t *testing.T

	calls    []any
	nextCall int // updated via assertNextCall

	brightnessCalls []*traits.UpdateBrightnessRequest
}

func (ta *testActions) assertNoMoreCalls() {
	ta.t.Helper()

	if len(ta.calls) > ta.nextCall {
		callStr := ""
		for i, call := range ta.calls[ta.nextCall:] {
			callStr += fmt.Sprintf("  [%d] %+v\n", i, call)
		}
		ta.t.Fatalf("Call count want 0, got %d\n%s", len(ta.calls)-ta.nextCall, callStr)
	}
}

func (ta *testActions) assertNextCall(req any) {
	ta.t.Helper()

	if len(ta.calls) <= ta.nextCall {
		ta.t.Fatalf("Call count want >%d, got %d", ta.nextCall, len(ta.calls))
	}
	call := ta.calls[ta.nextCall]
	ta.nextCall++

	if diff := cmp.Diff(req, call, protocmp.Transform()); diff != "" {
		ta.t.Fatalf("Next call (+want, -got)\n%s", diff)
	}
}

func (ta *testActions) UpdateBrightness(ctx context.Context, req *traits.UpdateBrightnessRequest, state *WriteState) error {
	ta.calls = append(ta.calls, req)
	ta.brightnessCalls = append(ta.brightnessCalls, req)
	state.Brightness[req.Name] = req.Brightness
	return nil
}
