package statusalerts

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/util/chans"
	"github.com/smart-core-os/sc-bos/pkg/util/pull"
)

func pullFrom(ctx context.Context, name string, client gen.StatusApiClient, c chan<- *gen.StatusLog) error {
	puller := statusLogPuller{client: client, name: name}
	return pull.Changes[*gen.StatusLog](ctx, puller, c)
}

type statusLogPuller struct {
	client gen.StatusApiClient
	name   string
}

func (s statusLogPuller) Pull(ctx context.Context, changes chan<- *gen.StatusLog) error {
	stream, err := s.client.PullCurrentStatus(ctx, &gen.PullCurrentStatusRequest{Name: s.name})
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range res.Changes {
			if err := chans.SendContext(ctx, changes, change.CurrentStatus); err != nil {
				return err
			}
		}
	}
}

func (s statusLogPuller) Poll(ctx context.Context, changes chan<- *gen.StatusLog) error {
	status, err := s.client.GetCurrentStatus(ctx, &gen.GetCurrentStatusRequest{Name: s.name})
	if err != nil {
		return err
	}
	return chans.SendContext(ctx, changes, status)
}
