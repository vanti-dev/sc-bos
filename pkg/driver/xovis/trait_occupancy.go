package xovis

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

type occupancyServer struct {
	traits.UnimplementedOccupancySensorApiServer
	bus         *minibus.Bus[PushData]
	client      *Client
	multiSensor bool
	logicID     int
}

var errDataFormat = status.Error(codes.FailedPrecondition, "data received from sensor did not match expected format")

func (o *occupancyServer) GetOccupancy(ctx context.Context, request *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	res, err := GetLiveLogic(o.client, o.multiSensor, o.logicID)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	occupancy := decodeOccupancyCounts(res.Logic.Counts)
	if occupancy == nil {
		return nil, errDataFormat
	}

	return occupancy, nil
}

func (o *occupancyServer) PullOccupancy(request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	// fetch the initial occupancy state
	res, err := GetLiveLogic(o.client, o.multiSensor, o.logicID)
	if err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}
	occupancy := decodeOccupancyCounts(res.Logic.Counts)
	if occupancy == nil {
		return errDataFormat
	}
	if !request.UpdatesOnly {
		err = server.Send(&traits.PullOccupancyResponse{
			Changes: []*traits.PullOccupancyResponse_Change{
				{
					Name:       request.Name,
					ChangeTime: timestamppb.New(res.Time),
					Occupancy:  occupancy,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	for data := range o.bus.Listen(ctx) {
		if data.LogicsData == nil {
			continue
		}
		records, ok := findLogicRecords(data.LogicsData, o.logicID)
		if !ok {
			continue
		}

		var changes []*traits.PullOccupancyResponse_Change
		for _, record := range records {
			occupancy := decodeOccupancyCounts(record.Counts)
			if occupancy == nil {
				continue
			}

			changes = append(changes, &traits.PullOccupancyResponse_Change{
				Name:       request.Name,
				ChangeTime: timestamppb.New(record.To),
				Occupancy:  occupancy,
			})
		}

		err = server.Send(&traits.PullOccupancyResponse{Changes: changes})
		if err != nil {
			return err
		}
	}

	return nil
}

// returns nil if the counts don't match the expected format for an occupancy logic
func decodeOccupancyCounts(counts []Count) (occupancy *traits.Occupancy) {
	if len(counts) != 1 || counts[0].Name != "balance" {
		return nil
	}
	var state traits.Occupancy_State
	if counts[0].Value > 0 {
		state = traits.Occupancy_OCCUPIED
	} else {
		state = traits.Occupancy_UNOCCUPIED
	}
	return &traits.Occupancy{
		State:       state,
		PeopleCount: int32(counts[0].Value),
	}
}
