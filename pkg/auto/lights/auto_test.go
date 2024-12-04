package lights

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/olebedev/emitter"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/brightnesssensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/mode"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/auto/lights/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

var errFailedBrightnessUpdate = errors.New("failed to update brightness this time")

func TestPirsTurnLightsOn(t *testing.T) {
	// we update this to send messages to the automation
	pir01 := occupancysensor.NewModel()
	pir02 := occupancysensor.NewModel()

	clients := node.ClientFunc(func(p any) error {
		switch v := p.(type) {
		case *traits.OccupancySensorApiClient:
			r := occupancysensor.NewApiRouter()
			r.Add("pir01", occupancysensor.WrapApi(occupancysensor.NewModelServer(pir01)))
			r.Add("pir02", occupancysensor.WrapApi(occupancysensor.NewModelServer(pir02)))
			*v = occupancysensor.WrapApi(r)
		case *traits.LightApiClient:
			*v = light.WrapApi(light.NewApiRouter())
		case *traits.BrightnessSensorApiClient:
			*v = brightnesssensor.WrapApi(brightnesssensor.NewApiRouter())
		case *gen.ButtonApiClient:
			*v = gen.WrapButtonApi(gen.NewButtonApiRouter())
		case *traits.ModeApiClient:
			*v = mode.WrapApi(mode.NewApiRouter())
		default:
			return errors.New("unsupported lightClient type")
		}
		return nil
	})

	testActions := newTestActions(t)

	now := time.Unix(0, 0)

	automation := PirsTurnLightsOn(clients, zap.NewNop())
	automation.makeActions = func(_ node.Clienter) (actions, error) { return testActions, nil }
	automation.autoStartTime = now

	cfg := config.Default()
	cfg.Now = func() time.Time { return now }
	cfg.OccupancySensors = []string{"pir01", "pir02"}
	cfg.Lights = []string{"light01", "light02"}
	cfg.UnoccupiedOffDelay = jsontypes.Duration{Duration: 10 * time.Minute}
	cfg.RefreshEvery = jsontypes.Duration{Duration: 8 * time.Minute}

	tickChan := make(chan time.Time, 1)
	automation.newTimer = func(d time.Duration) (<-chan time.Time, func() bool) {
		return tickChan, func() bool { return true }
	}

	if err := automation.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if err := automation.configure(cfg); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	processComplete := automation.bus.On("process-complete", emitter.Sync)
	waitForState := func(wait time.Duration, test func(state *ReadState) bool) (time.Duration, error) {
		t.Helper()
		timeout := time.NewTimer(wait)
		for {
			select {
			case <-timeout.C:
				t.Fatalf("timeout waiting for state")
				return 0, nil
			case e := <-processComplete:
				state := e.Args[2].(*ReadState)
				if test(state) {
					if e.Args[1] == nil {
						return e.Args[0].(time.Duration), nil
					}
					return e.Args[0].(time.Duration), e.Args[1].(error)
				}
			}
		}
	}

	// check setting occupied on one PIR causes the lights to come on
	_, _ = pir01.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_OCCUPIED})
	// DEAR FUTURE DEV, if you get test failures here it's probably because you're missing a case in the
	// clients switch statement above
	ttl, err := waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o, ok := state.Occupancy["pir01"]
		if !ok {
			return false
		}
		return o.State == traits.Occupancy_OCCUPIED
	})
	assertNoErrAndTtl(t, ttl, err, 0)

	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light02",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})

	// check that setting occupied on the other PIR does nothing
	_, _ = pir02.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_OCCUPIED})
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o, ok := state.Occupancy["pir02"]
		if !ok {
			return false
		}
		return o.State == traits.Occupancy_OCCUPIED
	})
	assertNoErrAndTtl(t, ttl, err, 0)
	testActions.assertNoMoreCalls()

	// check that making both PIRs unoccupied doesn't do anything, but then does
	_, _ = pir01.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED, StateChangeTime: timestamppb.New(now.Add(-3 * time.Minute))})
	_, _ = pir02.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED, StateChangeTime: timestamppb.New(now.Add(-8 * time.Minute))})
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		o02, ok02 := state.Occupancy["pir02"]
		if !ok01 || !ok02 {
			return false
		}
		return o01.State == traits.Occupancy_UNOCCUPIED && o02.State == traits.Occupancy_UNOCCUPIED
	})
	if want := 7 * time.Minute; ttl != want {
		t.Fatalf("TTL want %v, got %v", want, ttl)
	}
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	// trigger the timer
	now = now.Add(7 * time.Minute)
	tickChan <- now
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		return true // no state change, only time change
	})
	assertNoErrAndTtl(t, ttl, err, 0)
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 0,
		},
	})
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light02",
		Brightness: &traits.Brightness{
			LevelPercent: 0,
		},
	})

	// test 1 retry
	testActions.err = errFailedBrightnessUpdate
	_, _ = pir01.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_OCCUPIED})
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		if !ok01 {
			return false
		}
		return o01.State == traits.Occupancy_OCCUPIED
	})
	// jitter is set to ±0.2
	assertErrorAndTtl(t, ttl, err, cfg.OnProcessError.BackOffMultiplier.Duration*8/10, errFailedBrightnessUpdate)
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})
	// since newTimer is intercepted by this test, we force a replay here
	now = now.Add(time.Millisecond * 500)
	tickChan <- now
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		if !ok01 {
			return false
		}
		return o01.State == traits.Occupancy_OCCUPIED
	})
	assertNoErrAndTtl(t, ttl, err, 0)
	// it works after the retry
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light02",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})

	// testing retries getting cancelled after max attempts
	_, _ = pir01.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED, StateChangeTime: timestamppb.New(time.Unix(0, 0).Add(-3 * time.Minute))})
	_, _ = pir02.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED, StateChangeTime: timestamppb.New(time.Unix(0, 0).Add(-3 * time.Minute))})
	tickChan <- now
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		o02, ok02 := state.Occupancy["pir02"]
		if !ok01 || !ok02 {
			return false
		}
		return o01.State == traits.Occupancy_UNOCCUPIED && o02.State == traits.Occupancy_UNOCCUPIED
	})
	assertNoErrAndTtl(t, ttl, err, 0)
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 0,
		},
	})
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light02",
		Brightness: &traits.Brightness{
			LevelPercent: 0,
		},
	})

	testActions.err = errFailedBrightnessUpdate
	_, _ = pir01.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_OCCUPIED})
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		if !ok01 {
			return false
		}
		return o01.State == traits.Occupancy_OCCUPIED
	})
	assertErrorAndTtl(t, ttl, err, cfg.OnProcessError.BackOffMultiplier.Duration*8/10, errFailedBrightnessUpdate)
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})

	// second try
	testActions.err = errFailedBrightnessUpdate
	// since newTimer is intercepted by this test, we force a replay here
	now = now.Add(time.Millisecond * 500)
	tickChan <- now
	ttl, err = waitForState(time.Millisecond*500, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		if !ok01 {
			return false
		}
		return o01.State == traits.Occupancy_OCCUPIED
	})
	assertErrorAndTtl(t, ttl, err, 2*cfg.OnProcessError.BackOffMultiplier.Duration*(8/10), errFailedBrightnessUpdate)
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})

	// third try and is cancelled
	testActions.err = errFailedBrightnessUpdate
	// since newTimer is intercepted by this test, we force a replay here
	now = now.Add(time.Millisecond * 500)
	tickChan <- now
	ttl, err = waitForState(time.Millisecond*502, func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		if !ok01 {
			return false
		}
		return o01.State == traits.Occupancy_OCCUPIED
	})
	assertErrorAndTtl(t, ttl, err, -time.Nanosecond, errFailedBrightnessUpdate)
	testActions.assertNextCall(&traits.UpdateBrightnessRequest{
		Name: "light01",
		Brightness: &traits.Brightness{
			LevelPercent: 100,
		},
	})
	// ensure we have effectively cancelled reprocessing
	testActions.assertNoMoreCalls()
}
