package hpd3

import (
	"context"
	"math"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type occupancyServer struct {
	traits.UnimplementedOccupancySensorApiServer
	traits.UnimplementedOccupancySensorInfoServer

	client Client
	logger *zap.Logger
}

func (s *occupancyServer) GetOccupancy(ctx context.Context, _ *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	points, err := FetchPoints(ctx, s.client, PointPresence1, PointNumberOfPeopleTotal)
	if err != nil {
		s.logger.Error("failed to fetch occupancy points from device", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "failed to fetch occupancy points from device")
	}

	var state traits.Occupancy_State
	if points.Presence1 {
		state = traits.Occupancy_OCCUPIED
	} else {
		state = traits.Occupancy_UNOCCUPIED
	}
	return &traits.Occupancy{
		State:       state,
		PeopleCount: int32(points.NumberOfPeopleTotal),
	}, nil
}

func (s *occupancyServer) DescribeOccupancy(context.Context, *traits.DescribeOccupancyRequest) (*traits.OccupancySupport, error) {
	return &traits.OccupancySupport{
		ResourceSupport: &types.ResourceSupport{
			Readable:   true,
			Writable:   false,
			Observable: false,
		},
		MaxPeople: math.MaxInt32,
	}, nil
}
