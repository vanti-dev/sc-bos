package hd2

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type AirQualitySensor struct {
	traits.UnimplementedAirQualitySensorApiServer

	logger       *zap.Logger
	pollInterval time.Duration

	client *Client

	airQuality *resource.Value
}

func NewAirQualitySensor(client *Client, logger *zap.Logger, pollInterval time.Duration) AirQualitySensor {
	if pollInterval <= 0 {
		pollInterval = time.Second * 60
	}

	airQualitySensor := AirQualitySensor{
		client:       client,
		logger:       logger,
		pollInterval: pollInterval,
		airQuality:   resource.NewValue(resource.WithInitialValue(&traits.AirQuality{}), resource.WithNoDuplicates()),
	}

	airQualitySensor.GetUpdate()

	go airQualitySensor.startPoll(context.Background())

	return airQualitySensor
}

func (a *AirQualitySensor) startPoll(ctx context.Context) error {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			err := a.GetUpdate()
			if err != nil {
				a.logger.Error("error refreshing airQuality data", zap.Error(err))
			}
		}
	}
}

func (a *AirQualitySensor) GetAirQuality(ctx context.Context, req *traits.GetAirQualityRequest) (*traits.AirQuality, error) {
	err := a.GetUpdate()
	if err != nil {
		return nil, err
	}
	return a.airQuality.Get().(*traits.AirQuality), nil
}

func (a *AirQualitySensor) PullAirQuality(request *traits.PullAirQualityRequest, server traits.AirQualitySensorApi_PullAirQualityServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	changes := a.airQuality.Pull(ctx)

	for change := range changes {
		v := change.Value.(*traits.AirQuality)

		err := server.Send(&traits.PullAirQualityResponse{
			Changes: []*traits.PullAirQualityResponse_Change{
				{Name: request.GetName(), ChangeTime: timestamppb.New(change.ChangeTime), AirQuality: v},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AirQualitySensor) GetUpdate() error {
	response := SensorResponse{}
	err := doGetRequest(a.client, &response, "sensor")
	if err != nil {
		return err
	}

	co2 := float32(response.CO2)
	voc := float32(response.VOC)
	airPressure := float32(response.AirPressure)
	infectionRisk := float32(response.AerosolRiskOfInfection)

	a.airQuality.Set(&traits.AirQuality{
		CarbonDioxideLevel:       &co2,
		VolatileOrganicCompounds: &voc,
		AirPressure:              &airPressure,
		InfectionRisk:            &infectionRisk,
	})

	return nil
}
