package lights

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/brightnesssensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/bsp-ew/internal/auto/lights/config"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPirsTurnLightsOn(t *testing.T) {
	// we update this to send messages to the automation
	pir01 := occupancysensor.NewModel(&traits.Occupancy{})
	pir02 := occupancysensor.NewModel(&traits.Occupancy{})

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
		default:
			return errors.New("unsupported lightClient type")
		}
		return nil
	})

	testActions := newTestActions(t)

	automation := PirsTurnLightsOn(clients, zap.NewNop())
	automation.makeActions = func(_ node.Clienter) (actions, error) { return testActions, nil }

	now := time.Unix(0, 0)

	cfg := config.Default()
	cfg.Now = func() time.Time { return now }
	cfg.OccupancySensors = []string{"pir01", "pir02"}
	cfg.Lights = []string{"light01", "light02"}
	cfg.UnoccupiedOffDelay = config.Duration{Duration: 10 * time.Minute}

	if err := automation.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if err := automation.configure(cfg); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	processComplete := automation.bus.On("process-complete")
	waitForState := func(test func(state *ReadState) bool) (time.Duration, error) {
		t.Helper()
		timeout := time.NewTimer(500 * time.Millisecond)
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
	ttl, err := waitForState(func(state *ReadState) bool {
		o, ok := state.Occupancy["pir01"]
		if !ok {
			return false
		}
		return o.State == traits.Occupancy_OCCUPIED
	})
	assertNoTTLOrErr(t, ttl, err)

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
	ttl, err = waitForState(func(state *ReadState) bool {
		o, ok := state.Occupancy["pir02"]
		if !ok {
			return false
		}
		return o.State == traits.Occupancy_OCCUPIED
	})
	assertNoTTLOrErr(t, ttl, err)
	testActions.assertNoMoreCalls()

	// check that making both PIRs unoccupied doesn't do anything, but then does
	tickChan := make(chan time.Time, 1)
	automation.newTimer = func(d time.Duration) (<-chan time.Time, func() bool) {
		return tickChan, func() bool {
			return false
		}
	}
	_, _ = pir01.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED, StateChangeTime: timestamppb.New(now.Add(-3 * time.Minute))})
	_, _ = pir02.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED, StateChangeTime: timestamppb.New(now.Add(-8 * time.Minute))})
	ttl, err = waitForState(func(state *ReadState) bool {
		o01, ok01 := state.Occupancy["pir01"]
		o02, ok02 := state.Occupancy["pir02"]
		if !ok01 || !ok02 {
			return false
		}
		return o01.State == traits.Occupancy_UNOCCUPIED && o02.State == traits.Occupancy_UNOCCUPIED
	})
	if ttl != 7*time.Minute {
		t.Fatalf("TTL want %v, got %v", 5*time.Minute, ttl)
	}
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	// trigger the timer
	now = now.Add(7 * time.Minute)
	tickChan <- now
	ttl, err = waitForState(func(state *ReadState) bool {
		return true // no state change, only time change
	})
	assertNoTTLOrErr(t, ttl, err)
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
}
