package lights

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/auto/lights/config"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		readState.Config.UnoccupiedOffDelay = config.Duration{Duration: 10 * time.Minute}
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
		readState.Config.UnoccupiedOffDelay = config.Duration{Duration: 10 * time.Minute}
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

	t.Run("threshold brightness", func(t *testing.T) {
		dd := config.DaylightDimming{
			Thresholds: []config.LevelThreshold{
				{BelowLux: 10, LevelPercent: 100},
				{BelowLux: 200, LevelPercent: 70},
				{BelowLux: 1000, LevelPercent: 50},
				{BelowLux: 10_000, LevelPercent: 30},
				{BelowLux: 30_000, LevelPercent: 1},
				{LevelPercent: 0},
			},
		}
		tests := []struct {
			name string
			want float32
			lux  []float32
		}{
			{"no readings", 100, []float32{}},
			{"0 reading", 100, []float32{0}},
			{"average", 100, []float32{3, 4, 5}},
			{"just below threshold", 100, []float32{9.999}},
			{"on threshold", 70, []float32{10}},
			{"50%", 50, []float32{1000 - 1}},
			{"30%", 30, []float32{10_000 - 1}},
			{"1%", 1, []float32{30_000 - 1}},
			{"off", 0, []float32{30_000}},
			{"very bright", 0, []float32{100_000}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				readState := NewReadState()
				writeState := NewWriteState()
				actions := newTestActions(t)

				readState.Config.DaylightDimming = &dd
				readState.Config.OccupancySensors = []string{"pir01"}
				readState.Config.Lights = []string{"light01"}
				readState.Occupancy["pir01"] = &traits.Occupancy{State: traits.Occupancy_OCCUPIED}

				readState.Config.BrightnessSensors = make([]string, len(tt.lux))
				for i, lux := range tt.lux {
					name := fmt.Sprintf("bri%02d", i)
					readState.Config.BrightnessSensors[i] = name
					readState.AmbientBrightness[name] = &traits.AmbientBrightness{BrightnessLux: lux}
				}

				ttl, err := processState(context.Background(), readState, writeState, actions)
				assertNoTTLOrErr(t, ttl, err)
				actions.assertNextCall(&traits.UpdateBrightnessRequest{
					Name: "light01",
					Brightness: &traits.Brightness{
						LevelPercent: tt.want,
					},
				})
			})
		}
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
