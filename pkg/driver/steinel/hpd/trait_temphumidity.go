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

	temperature *resource.Value
}

func NewTemperatureSensor(client *Client, logger *zap.Logger, pollInterval time.Duration) TemperatureSensor {
	if pollInterval <= 0 {
		pollInterval = time.Second * 60
	}

	temperatureSensor := TemperatureSensor{
		client:       client,
		logger:       logger,
		pollInterval: pollInterval,
		temperature:  resource.NewValue(resource.WithInitialValue(&traits.AirTemperature{}), resource.WithNoDuplicates()),
	}

	temperatureSensor.GetUpdate()

	go temperatureSensor.startPoll(context.Background())

	return temperatureSensor
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

func (a *TemperatureSensor) GetAirTemperature(ctx context.Context, req *traits.GetAirTemperatureRequest) (*traits.AirTemperature, error) {
	err := a.GetUpdate()
	if err != nil {
		return nil, err
	}
	return a.temperature.Get().(*traits.AirTemperature), nil
}

func (a *TemperatureSensor) PullAirTemperature(request *traits.PullAirTemperatureRequest, server traits.AirTemperatureApi_PullAirTemperatureServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	changes := a.temperature.Pull(ctx)

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

	a.temperature.Set(&traits.AirTemperature{
		Mode:               0,
		TemperatureGoal:    nil,
		AmbientTemperature: &types.Temperature{ValueCelsius: response.Temperature},
		AmbientHumidity:    &humidity,
		DewPoint:           nil,
	})

	return nil
}
