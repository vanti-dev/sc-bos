package history

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-bos/pkg/auto/history/config"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/cmp"
)

func (a *automation) collectSoundSensorChanges(ctx context.Context, source config.Source, payloads chan<- []byte) {
	client := gen.NewSoundSensorApiClient(a.clients.ClientConn())

	last := newDeduper[*gen.SoundLevel](cmp.Equal(cmp.FloatValueApprox(0, 0.0001)))

	pullFn := func(ctx context.Context, changes chan<- []byte) error {
		stream, err := client.PullSoundLevel(ctx, &gen.PullSoundLevelRequest{Name: source.Name, UpdatesOnly: true, ReadMask: source.ReadMask.PB()})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}
			for _, change := range msg.Changes {
				if !last.Changed(change.GetSoundLevel()) {
					continue
				}

				payload, err := proto.Marshal(change.GetSoundLevel())
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
		resp, err := client.GetSoundLevel(ctx, &gen.GetSoundLevelRequest{Name: source.Name, ReadMask: source.ReadMask.PB()})
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
