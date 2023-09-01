package hpd3

import (
	"context"
	"math"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
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
	points, err := FetchPoints(ctx, s.client, pointPresence1, pointNumberOfPeopleTotal)
	if err != nil {
		s.logger.Error("failed to fetch occupancy points from device", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "failed to fetch occupancy points from device")
	}

	presence, presenceOK := points[pointPresence1].(bool)
	number, numberOK := points[pointNumberOfPeopleTotal].(float64)
	if !presenceOK || !numberOK {
		s.logger.Error("some occupancy points missing or invalid", zap.Strings("keys", maps.Keys(points)))
		return nil, status.Error(codes.Internal, "some occupancy points missing or invalid")
	}

	var state traits.Occupancy_State
	if presence {
		state = traits.Occupancy_OCCUPIED
	} else {
		state = traits.Occupancy_UNOCCUPIED
	}
	return &traits.Occupancy{
		State:       state,
		PeopleCount: int32(number),
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
