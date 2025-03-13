package main

import (
	"google.golang.org/protobuf/proto"

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
		MeterReading:                 &gen.MeterReadingSupport{UsageUnit: unit},
	}

	client := node.WithClients(gen.WrapMeterApi(meter.NewModelServer(model)), gen.WrapMeterInfo(modelInfoServer))
	root.Announce(name, node.HasTrait(meter.TraitName, client))

	store := memstore.New()

	for _, event := range events {
		rec, err := proto.Marshal(&gen.MeterReading{
			Usage: event,
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
