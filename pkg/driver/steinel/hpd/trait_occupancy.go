package hpd

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type Occupancy struct {
	traits.UnimplementedOccupancySensorApiServer
	gen.UnimplementedUdmiServiceServer

	logger       *zap.Logger
	pollInterval time.Duration

	client *Client

	OccupancyValue *resource.Value
}

func NewOccupancySensor(client *Client, logger *zap.Logger, pollInterval time.Duration) *Occupancy {
	if pollInterval <= 0 {
		pollInterval = time.Second * 60
	}

	return &Occupancy{
		client:         client,
		logger:         logger,
		pollInterval:   pollInterval,
		OccupancyValue: resource.NewValue(resource.WithInitialValue(&traits.Occupancy{}), resource.WithNoDuplicates()),
	}
}

// StartPollingForData starts a loop which fetches data from the sensor at a set interval
func (o *Occupancy) StartPollingForData() {
	go func() {
		_ = o.startPoll(context.Background())
	}()
}

func (o *Occupancy) startPoll(ctx context.Context) error {
	ticker := time.NewTicker(o.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			err := o.GetUpdate()
			if err != nil {
				o.logger.Error("error refreshing Occupancy data", zap.Error(err))
			}
		}
	}
}

func (o *Occupancy) GetOccupancy(_ context.Context, _ *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	err := o.GetUpdate()
	if err != nil {
		return nil, err
	}
	return o.OccupancyValue.Get().(*traits.Occupancy), nil
}

func (o *Occupancy) PullOccupancy(request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	changes := o.OccupancyValue.Pull(ctx)

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

func (o *Occupancy) GetUpdate() error {
	response := SensorResponse{}
	err := doGetRequest(o.client, &response, "sensor")
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

	o.OccupancyValue.Set(&traits.Occupancy{
		State:       state,
		PeopleCount: int32(peopleCount),
	})

	return nil
}
