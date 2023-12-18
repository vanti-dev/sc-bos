package pestsense

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type PestSensor struct {
	traits.UnimplementedOccupancySensorApiServer

	Id        string
	Occupancy *resource.Value
}

func NewPestSensor(id string) *PestSensor {
	return &PestSensor{
		Id:        id,
		Occupancy: resource.NewValue(resource.WithInitialValue(&traits.Occupancy{}), resource.WithNoDuplicates()),
	}
}

func (s *PestSensor) GetOccupancy(ctx context.Context, request *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	value := s.Occupancy.Get()
	occupancy := value.(*traits.Occupancy)
	return occupancy, nil
}

func (o *PestSensor) PullOccupancy(request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	// TODO: refresh initial occupancy value

	changes := o.Occupancy.Pull(ctx, resource.WithBackpressure(false))
	for change := range changes {
		occupancy := change.Value.(*traits.Occupancy)
		resChange := &traits.PullOccupancyResponse_Change{
			Occupancy:  occupancy,
			ChangeTime: timestamppb.New(change.ChangeTime),
		}
		res := &traits.PullOccupancyResponse{
			Changes: []*traits.PullOccupancyResponse_Change{resChange},
		}

		err := server.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}
