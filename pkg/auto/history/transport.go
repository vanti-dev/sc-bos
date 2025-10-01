package history

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func (a *automation) collectTransportChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	client := gen.NewTransportApiClient(a.clients.ClientConn())

	last := newDeduper[*gen.Transport](cmp.Equal(cmp.FloatValueApprox(0.1, 0.0001), cmp.TimeValueWithin(time.Second)))

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullTransport(ctx, &gen.PullTransportRequest{Name: source.Name, UpdatesOnly: true, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				if !last.Changed(change.GetTransport()) {
					continue
				}

				payload, err := proto.Marshal(change.GetTransport())
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
		resp, err := client.GetTransport(ctx, &gen.GetTransportRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
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
