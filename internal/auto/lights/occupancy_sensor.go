package lights

import (
	"context"
	"github.com/smart-core-os/sc-api/go/traits"
	"go.uber.org/zap"
)

// OccupancySensorPatches contributes patches for changing the state based on occupancy sensor readings.
type OccupancySensorPatches struct {
	name   string
	client traits.OccupancySensorApiClient
	logger *zap.Logger
}

func (o *OccupancySensorPatches) Subscribe(ctx context.Context, changes chan<- Patcher) error {
	// remove our signal when we shouldn't be contributing anymore
	defer func() {
		changes <- clearOccupancyTransition(o.name)
	}()
	return subscribe(ctx, o, changes, withLogger(o.logger.Named("occupancy")))
}

func (o *OccupancySensorPatches) pull(ctx context.Context, changes chan<- Patcher) error {
	stream, err := o.client.PullOccupancy(ctx, &traits.PullOccupancyRequest{Name: o.name})
	if err != nil {
		return err
	}

	for {
		change, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case changes <- (*pullOccupancyTransition)(change):
		}
	}
}

func (o *OccupancySensorPatches) poll(ctx context.Context, changes chan<- Patcher) error {
	res, err := o.client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: o.name})
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case changes <- getOccupancyPatcher{o.name, res}:
		return nil
	}
}

type pullOccupancyTransition traits.PullOccupancyResponse

func (o *pullOccupancyTransition) Patch(s *ReadState) {
	r := (*traits.PullOccupancyResponse)(o)

	for _, change := range r.Changes {
		s.Occupancy[change.Name] = change.Occupancy
	}
}

type getOccupancyPatcher struct {
	name string
	res  *traits.Occupancy
}

func (g getOccupancyPatcher) Patch(s *ReadState) {
	s.Occupancy[g.name] = g.res
}

type clearOccupancyTransition string

func (c clearOccupancyTransition) Patch(s *ReadState) {
	delete(s.Occupancy, string(c))
}
