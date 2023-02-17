package history

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/util/chans"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_automation_collectOccupancyChanges(t *testing.T) {
	model := occupancysensor.NewModel(&traits.Occupancy{})
	// n is used as the clienter and announcer in the automation
	n := node.New("test")
	n.Support(node.Clients(occupancysensor.WrapApi(occupancysensor.NewModelServer(model))))

	collector := &automation{
		clients:  n,
		announce: n,
		logger:   zap.NewNop(),
	}

	payloads := make(chan []byte, 5)
	ctx, stop := context.WithCancel(context.Background())
	t.Cleanup(stop)
	go func() {
		collector.collectOccupancyChanges(ctx, "anything", payloads)
	}()

	if err := chans.IsEmptyWithin(payloads, time.Second); err != nil {
		t.Fatal(err)
	}

	if _, err := model.SetOccupancy(&traits.Occupancy{State: traits.Occupancy_OCCUPIED}); err != nil {
		t.Fatal(err)
	}
	payload, err := chans.RecvWithin(payloads, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	msg := &traits.Occupancy{}
	err = proto.Unmarshal(payload, msg)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(&traits.Occupancy{State: traits.Occupancy_OCCUPIED}, msg, protocmp.Transform()); diff != "" {
		t.Fatalf("payload (-want,+got)\n%s", diff)
	}
}
