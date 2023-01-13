package xovis

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type occupancyServer struct {
	traits.UnimplementedOccupancySensorApiServer
	client      *Client
	multiSensor bool
	logicID     int
}

func (o *occupancyServer) GetOccupancy(
	ctx context.Context, request *traits.GetOccupancyRequest,
) (*traits.Occupancy, error) {
	res, err := GetLiveLogic(o.client, o.multiSensor, o.logicID)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	occupancyCount, ok := decodeOccupancyCounts(res.Logic.Counts)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition,
			"more than one count received from sensor; are you sure this is an occupancy zone?")
	}

	occupancy := &traits.Occupancy{
		PeopleCount: int32(occupancyCount),
	}
	if occupancyCount > 0 {
		occupancy.State = traits.Occupancy_OCCUPIED
	} else {
		occupancy.State = traits.Occupancy_UNOCCUPIED
	}
	return occupancy, nil
}

func (o *occupancyServer) PullOccupancy(
	request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer,
) error {
	return status.Error(codes.Unimplemented, "PullOccupancy not implemented for this device")
}

func decodeOccupancyCounts(counts []Count) (occupancy int, ok bool) {
	if len(counts) != 1 || counts[0].Name != "balance" {
		ok = false
		return
	}
	occupancy = counts[0].Value
	ok = true
	return
}
