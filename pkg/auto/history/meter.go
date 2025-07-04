package history

import (
	"context"
	"math"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func (a *automation) collectMeterReadingChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	client := gen.NewMeterApiClient(a.clients.ClientConn())

	last := newDeduper[*gen.MeterReading](cmp.Equal(cmp.FloatValueApprox(0, 0.0001), cmp.TimeValueWithin(time.Second)))

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullMeterReadings(ctx, &gen.PullMeterReadingsRequest{Name: source.Name, UpdatesOnly: true, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				if change.GetMeterReading() == nil || math.IsNaN(float64(change.GetMeterReading().Usage)) {
					continue
				}
				if !last.Changed(change.GetMeterReading()) {
					continue
				}

				payload, err := proto.Marshal(change.GetMeterReading())
				if err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case changes <- payload:
				}
			}
		}
	}
	pollFn := func(ctx context.Context, changes chan<- []byte) error {
		resp, err := client.GetMeterReading(ctx, &gen.GetMeterReadingRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		if resp == nil || math.IsNaN(float64(resp.Usage)) {
			return nil
		}
		if !last.Changed(resp) {
			return nil
		}

		payload, err := proto.Marshal(resp)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case changes <- payload:
		}
		return nil
	}

	if err := collectChanges(ctx, source, pullFn, pollFn, payloads, a.logger); err != nil {
		a.logger.Warn("collection aborted", zap.Error(err))
	}
}
