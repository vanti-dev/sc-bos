package history

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/util/pull"
)

func (a *automation) collectOccupancyChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	var client traits.OccupancySensorApiClient
	if err := a.clients.Client(&client); err != nil {
		a.logger.Error("collection aborted", zap.Error(err))
		return
	}

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullOccupancy(ctx, &traits.PullOccupancyRequest{Name: source.Name, UpdatesOnly: true, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				payload, err := proto.Marshal(change.Occupancy)
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
		demand, err := client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})
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

func (a *automation) sampleOccupancyChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	var client traits.OccupancySensorApiClient
	if err := a.clients.Client(&client); err != nil {
		a.logger.Error("sampling aborted", zap.Error(err))
		return
	}

	t := time.Now()

	var prev *traits.Occupancy

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(source.Sample.Schedule.Next(t))):
			resp, err := client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})

			t = time.Now()
			if err != nil {
				a.logger.Warn("sample aborted", zap.Error(err))
				continue
			}

			if proto.Equal(prev, resp) {
				continue
			}

			prev = resp

			payload, err := proto.Marshal(resp)
			if err != nil {
				a.logger.Warn("sample aborted", zap.Error(err))
				continue
			}

			select {
			case <-ctx.Done():
				return
			case payloads <- payload:
			}
		}
	}

}
