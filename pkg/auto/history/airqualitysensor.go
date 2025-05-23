package history

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
)

func (a *automation) collectAirQualityChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	var client traits.AirQualitySensorApiClient
	if err := a.clients.Client(&client); err != nil {
		a.logger.Error("collection aborted", zap.Error(err))
		return
	}

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullAirQuality(ctx, &traits.PullAirQualityRequest{Name: source.Name, UpdatesOnly: true, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				payload, err := proto.Marshal(change.AirQuality)
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
		resp, err := client.GetAirQuality(ctx, &traits.GetAirQualityRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
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
