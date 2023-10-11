package hd2

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type Occupancy struct {
	traits.UnimplementedOccupancySensorApiServer

	logger       *zap.Logger
	pollInterval time.Duration

	client *Client

	occupancy *resource.Value
}

func NewOccupancySensor(client *Client, logger *zap.Logger, pollInterval time.Duration) Occupancy {
	if pollInterval <= 0 {
		pollInterval = time.Second * 60
	}

	occupancy := Occupancy{
		client:       client,
		logger:       logger,
		pollInterval: pollInterval,
		occupancy:    resource.NewValue(resource.WithInitialValue(&traits.Occupancy{}), resource.WithNoDuplicates()),
	}

	occupancy.GetUpdate()

	go occupancy.startPoll(context.Background())

	return occupancy
}

func (a *Occupancy) startPoll(ctx context.Context) error {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			err := a.GetUpdate()
			if err != nil {
				a.logger.Error("error refreshing thermostat data", zap.Error(err))
			}
		}
	}
}

func (a *Occupancy) GetOccupancy(ctx context.Context, req *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	err := a.GetUpdate()
	if err != nil {
		return nil, err
	}
	return a.occupancy.Get().(*traits.Occupancy), nil
}

func (a *Occupancy) PullOccupancy(request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	changes := a.occupancy.Pull(ctx)

	for change := range changes {
		v := change.Value.(*traits.Occupancy)

		err := server.Send(&traits.PullOccupancyResponse{
			Changes: []*traits.PullOccupancyResponse_Change{
				{Name: request.GetName(), ChangeTime: timestamppb.New(change.ChangeTime), Occupancy: v},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Occupancy) GetUpdate() error {
	response := SensorResponse{}
	err := doGetRequest(a.client, &response, "sensor")
	if err != nil {
		return err
	}

	peopleCount := 0

	state := traits.Occupancy_STATE_UNSPECIFIED
	if response.TruePresence1 {
		state = traits.Occupancy_OCCUPIED
	} else {
		state = traits.Occupancy_UNOCCUPIED
	}

	if response.ZonePeople0 > 0 {
		state = traits.Occupancy_OCCUPIED
		peopleCount = response.ZonePeople0
	}

	a.occupancy.Set(&traits.Occupancy{
		State:       state,
		PeopleCount: int32(peopleCount),
	})

	return nil
}
