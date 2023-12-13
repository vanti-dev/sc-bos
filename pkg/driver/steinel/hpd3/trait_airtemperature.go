package hpd3

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type airTemperatureServer struct {
	traits.UnimplementedAirTemperatureApiServer
	traits.UnimplementedAirTemperatureInfoServer

	client Client
	logger *zap.Logger
}

func (s *airTemperatureServer) DescribeAirTemperature(context.Context, *traits.DescribeAirTemperatureRequest) (*traits.AirTemperatureSupport, error) {
	return &traits.AirTemperatureSupport{
		ResourceSupport: &types.ResourceSupport{
			Readable:   true,
			Writable:   false,
			Observable: false,
		},
		NativeUnit: types.TemperatureUnit_CELSIUS,
	}, nil
}

func (s *airTemperatureServer) GetAirTemperature(ctx context.Context, _ *traits.GetAirTemperatureRequest) (*traits.AirTemperature, error) {
	points, err := FetchPoints(ctx, s.client, PointTemperature, PointHumidity)
	if err != nil {
		s.logger.Error("failed to fetch air temperature points", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "failed to fetch air temperature points")
	}

	humidity := float32(points.Humidity)
	return &traits.AirTemperature{
		AmbientTemperature: &types.Temperature{ValueCelsius: points.Temperature},
		AmbientHumidity:    &humidity,
	}, nil
}
