package main

import (
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/meter"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/historypb"
	"github.com/vanti-dev/sc-bos/pkg/history/memstore"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

// announceMeter with events in order
func announceMeter(root node.Announcer, name, unit string, events []float32) error {
	model := meter.NewModel()

	modelInfoServer := &meter.InfoServer{
		UnimplementedMeterInfoServer: traits.UnimplementedMeterInfoServer{},
		MeterReading:                 &traits.MeterReadingSupport{Unit: unit},
	}

	client := node.WithClients(meter.WrapApi(meter.NewModelServer(model)), meter.WrapInfo(modelInfoServer))
	root.Announce(name, node.HasTrait(trait.Meter, client))

	store := memstore.New()

	for _, event := range events {
		rec, err := proto.Marshal(&traits.MeterReading{
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
