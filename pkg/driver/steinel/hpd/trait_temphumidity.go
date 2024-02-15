package hpd

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type TemperatureSensor struct {
	traits.UnimplementedAirTemperatureApiServer

	logger       *zap.Logger
	pollInterval time.Duration

	client *Client

	TemperatureValue *resource.Value
}

func NewTemperatureSensor(client *Client, logger *zap.Logger, pollInterval time.Duration) *TemperatureSensor {
	if pollInterval <= 0 {
		pollInterval = time.Second * 60
	}

	return &TemperatureSensor{
		client:           client,
		logger:           logger,
		pollInterval:     pollInterval,
		TemperatureValue: resource.NewValue(resource.WithInitialValue(&traits.AirTemperature{}), resource.WithNoDuplicates()),
	}
}

// StartPollingForData starts a loop which fetches data from the sensor at a set interval
func (a *TemperatureSensor) StartPollingForData() {
	go func() {
		_ = a.startPoll(context.Background())
	}()
}

func (a *TemperatureSensor) startPoll(ctx context.Context) error {
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

func (a *TemperatureSensor) GetAirTemperature(_ context.Context, _ *traits.GetAirTemperatureRequest) (*traits.AirTemperature, error) {
	err := a.GetUpdate()
	if err != nil {
		return nil, err
	}
	return a.TemperatureValue.Get().(*traits.AirTemperature), nil
}

func (a *TemperatureSensor) PullAirTemperature(request *traits.PullAirTemperatureRequest, server traits.AirTemperatureApi_PullAirTemperatureServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	changes := a.TemperatureValue.Pull(ctx)

	for change := range changes {
		v := change.Value.(*traits.AirTemperature)

		err := server.Send(&traits.PullAirTemperatureResponse{
			Changes: []*traits.PullAirTemperatureResponse_Change{
				{Name: request.GetName(), ChangeTime: timestamppb.New(change.ChangeTime), AirTemperature: v},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *TemperatureSensor) GetUpdate() error {
	response := SensorResponse{}
	err := doGetRequest(a.client, &response, "sensor")
	if err != nil {
		return err
	}

	humidity := float32(response.Humidity) / 100

	a.TemperatureValue.Set(&traits.AirTemperature{
		Mode:               0,
		TemperatureGoal:    nil,
		AmbientTemperature: &types.Temperature{ValueCelsius: response.Temperature},
		AmbientHumidity:    &humidity,
		DewPoint:           nil,
	})

	return nil
}
