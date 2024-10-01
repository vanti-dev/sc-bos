package main

import (
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/historypb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/history/memstore"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

// announceMeter with events in order
func announceMeter(root node.Announcer, name, unit string, events []float32) error {

	model := meter.NewModel()

	modelInfoServer := &meter.InfoServer{
		UnimplementedMeterInfoServer: gen.UnimplementedMeterInfoServer{},
		MeterReading:                 &gen.MeterReadingSupport{Unit: unit},
	}

	client := node.WithClients(meter.NewModelServer(model), gen.WrapMeterInfo(modelInfoServer))
	root.Announce(name, node.HasTrait(meter.TraitName, client))

	store := memstore.New()

	now := time.Now()

	for idx, event := range events {
		randomInt := rand.Int() + idx
		rec, err := proto.Marshal(&gen.MeterReading{
			Usage:     event,
			StartTime: timestamppb.New(now.Add(-time.Duration(randomInt) * time.Minute)),
			EndTime:   timestamppb.New(now.Add(-time.Duration(randomInt) * time.Duration(rand.Int()) * time.Minute)),
		})
		if err != nil {
			return err
		}
		_, err = store.Append(nil, rec)
		if err != nil {
			return err
		}
	}

	root.Announce(name, node.HasClient(gen.WrapMeterHistory(historypb.NewMeterServer(store))))

	return nil
}
