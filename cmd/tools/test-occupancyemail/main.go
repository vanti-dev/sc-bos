package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/occupancyemail"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/historypb"
	"github.com/vanti-dev/sc-bos/pkg/history/memstore"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	root := node.New("test")
	alltraits.AddSupport(root)

	now, _ := time.Parse(time.DateTime, "2023-11-15 11:36:00")
	now = now.Round(time.Second) // get rid of millis, etc

	oc := func(age time.Duration, pc int) *traits.PullOccupancyResponse_Change {
		return &traits.PullOccupancyResponse_Change{
			ChangeTime: timestamppb.New(now.Add(-age)),
			Occupancy:  &traits.Occupancy{PeopleCount: int32(pc)},
		}
	}
	testData := []*traits.PullOccupancyResponse_Change{
		// note: these _must_ be in chronological order
		oc(7*24*time.Hour+time.Second, 20), // before the 7-day window
		oc(7*24*time.Hour-2*time.Second, 6),
		oc(7*24*time.Hour-2*time.Hour, 0),
		oc(7*24*time.Hour-3*time.Hour, 7),
		oc(3*24*time.Hour, 4),
		oc(-time.Second, 22), // in the future, just in case
	}
	store := memstore.New()
	for _, td := range testData {
		td := td
		memstore.SetNow(store, td.ChangeTime.AsTime)
		payload, _ := proto.Marshal(td.Occupancy)
		_, err := store.Append(nil, payload)
		if err != nil {
			panic(err)
		}
	}
	device := historypb.NewOccupancySensorServer(store)
	client := gen.WrapOccupancySensorHistory(device)
	root.Announce("test", node.HasTrait(trait.OccupancySensor, node.WithClients(client)))

	serv := auto.Services{
		Logger: logger,
		Node:   root,
		Now: func() time.Time {
			return now.Add(-2 * time.Second)
		},
	}
	lifecycle := occupancyemail.Factory.New(serv)
	defer lifecycle.Stop()
	cfg := `{
  "name": "emails", "type": "occupancyemail",
  "source": {
    "name": "test",
    "title": "Enterprise Wharf"
  },
  "destination": {
    "host": "smtp.gmail.com",
    "from": "Enterprise Wharf <no-reply@enterprisewharf.co.uk>",
    "to": ["Matt Nathan <matt.nathan@vanti.co.uk>"],
    "passwordFile": ".secrets/ew-email-pass",
    "sendTime": "36 11 * * *"
  }
}`
	_, err = lifecycle.Configure([]byte(cfg))
	if err != nil {
		panic(err)
	}
	_, err = lifecycle.Start()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
