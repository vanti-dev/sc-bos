package lights

import (
	"context"
	"github.com/smart-core-os/sc-api/go/traits"
)

// OccupancySensorPatches contributes patches for changing the state based on occupancy sensor readings.
type OccupancySensorPatches struct {
	name   string
	client traits.OccupancySensorApiClient
}

func (o *OccupancySensorPatches) Subscribe(ctx context.Context, changes chan<- Patcher) error {
	stream, err := o.client.PullOccupancy(ctx, &traits.PullOccupancyRequest{Name: o.name})
	if err != nil {
		return err
	}

	// remove our signal when we shouldn't be contributing anymore
	defer func() {
		changes <- clearOccupancyTransition(o.name)
	}()

	for {
		change, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case changes <- (*occupancyTransition)(change):
		}
	}
}

type occupancyTransition traits.PullOccupancyResponse

func (o *occupancyTransition) Patch(s *ReadState) {
	r := (*traits.PullOccupancyResponse)(o)

	for _, change := range r.Changes {
		s.Occupancy[change.Name] = change.Occupancy
	}
}

type clearOccupancyTransition string

func (c clearOccupancyTransition) Patch(s *ReadState) {
	delete(s.Occupancy, string(c))
}
