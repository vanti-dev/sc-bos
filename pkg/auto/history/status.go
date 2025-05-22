package history

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/util/pull"
)

func (a *automation) collectCurrentStatusChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	var client gen.StatusApiClient
	if err := a.clients.Client(&client); err != nil {
		a.logger.Error("collection aborted", zap.Error(err))
		return
	}

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullCurrentStatus(ctx, &gen.PullCurrentStatusRequest{Name: source.Name, UpdatesOnly: true, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				payload, err := proto.Marshal(change.CurrentStatus)
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
		demand, err := client.GetCurrentStatus(ctx, &gen.GetCurrentStatusRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})
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

func (a *automation) sampleCurrentStatusChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	var client gen.StatusApiClient
	if err := a.clients.Client(&client); err != nil {
		a.logger.Error("sampling aborted", zap.Error(err))
		return
	}

	t := time.Now()

	var prev *gen.StatusLog

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(source.Sample.Schedule.Next(t))):
			resp, err := client.GetCurrentStatus(ctx, &gen.GetCurrentStatusRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})

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
