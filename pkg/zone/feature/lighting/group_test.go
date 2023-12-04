package lighting

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/vanti-dev/sc-bos/pkg/util/chans"
)

func TestGroup_GetBrightness(t *testing.T) {
	r := light.NewApiRouter()
	client := light.WrapApi(r)
	group := &Group{
		client: client,
		names: []string{
			"L1", "L2", "L3", "L4", "L5",
		},
		readOnly: false,
		logger:   zap.NewNop(),
	}

	for _, name := range group.names {
		r.Add(name, light.WrapApi(light.NewMemoryDevice()))
	}

	type response struct {
		R *traits.PullBrightnessResponse
		E error
	}
	responses := make(chan response)
	pullCtx, pullCancel := context.WithCancel(context.Background())
	defer pullCancel()
	go func() {
		defer close(responses)
		groupClient := light.WrapApi(group)
		stream, err := groupClient.PullBrightness(pullCtx, &traits.PullBrightnessRequest{
			Name: "anything will do", // we're using a direct client call, not routed
		})
		if err != nil {
			responses <- response{E: err}
			return
		}
		for {
			r, err := stream.Recv()
			if err != nil {
				responses <- response{E: err}
				return
			}
			responses <- response{R: r}
		}
	}()

	chanWait := 10 * time.Millisecond
	time.Sleep(500 * time.Millisecond) // wait for the pull calls to start

	// test initial updates, should all be 0, but we should only get one as we omit duplicates
	res, err := chans.RecvWithin(responses, chanWait)
	if err != nil {
		t.Fatalf("initial update %v", err)
	}
	if res.E != nil {
		t.Fatalf("initial update %v", res.E)
	}
	if v := len(res.R.GetChanges()); v != 1 {
		t.Fatalf("initial update; want len(changes)==1, got %v", v)
	}
	change := res.R.GetChanges()[0]
	if v := change.GetBrightness().GetLevelPercent(); v != 0 {
		t.Fatalf("initial update; want level_percent==0, got %v", v)
	}

	expectedAverage := func() float32 {
		var total float32
		for _, name := range group.names {
			c, err := r.GetBrightness(context.Background(), &traits.GetBrightnessRequest{Name: name})
			if err != nil {
				t.Fatalf("get brightness %v: %v", name, err)
			}
			total += c.GetLevelPercent()
		}
		return total / float32(len(group.names))
	}
	testUpdate := func(name string, level float32) {
		_, err := r.UpdateBrightness(context.Background(), &traits.UpdateBrightnessRequest{
			Name: name,
			Brightness: &traits.Brightness{
				LevelPercent: level,
			},
		})
		if err != nil {
			t.Fatalf("update brightness %s: %v", name, err)
		}
		want := expectedAverage()
		got, err := chans.RecvWithin(responses, chanWait)
		if err != nil {
			t.Fatalf("update brightness group %s: %v", name, err)
		}
		if got.E != nil {
			t.Fatalf("update brightness group %s: %v", name, got.E)
		}
		if v := len(got.R.GetChanges()); v != 1 {
			t.Fatalf("update brightness group %s; want len(changes)==1, got %v", name, v)
		}
		change := got.R.GetChanges()[0]
		if change.Brightness.LevelPercent != want {
			t.Fatalf("update brightness group %s; want level_percent==%v, got %v", name, want, change.Brightness.LevelPercent)
		}
	}
	// start sending updates and waiting for results
	testUpdate("L1", 100)
	testUpdate("L2", 50)
	testUpdate("L2", 0)
	testUpdate("L1", 0)
}
