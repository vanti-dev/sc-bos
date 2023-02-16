package history

import (
	"context"
	"errors"

	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (a *automation) collectMeterReadingChanges(ctx context.Context, name string, payloads chan<- []byte) {
	var client gen.MeterApiClient
	if err := a.clients.Client(&client); err != nil {
		a.logger.Error("collection aborted", zap.Error(err))
		return
	}

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullMeterReadings(ctx, &gen.PullMeterReadingsRequest{Name: name, UpdatesOnly: true})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				payload, err := proto.Marshal(change.MeterReading)
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
		demand, err := client.GetMeterReading(ctx, &gen.GetMeterReadingRequest{Name: name})
		if err != nil {
			return err
		}
		payload, err := proto.Marshal(demand)
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

	err := pull.Changes(ctx, pull.NewFetcher(pullFn, pollFn), payloads, pull.WithLogger(a.logger))
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}
	if err != nil {
		a.logger.Warn("collection aborted", zap.Error(err))
	}
}
