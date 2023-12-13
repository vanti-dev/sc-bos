package hpd3

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type airQualityServer struct {
	traits.UnimplementedAirQualitySensorApiServer
	traits.UnimplementedAirQualitySensorInfoServer

	client Client
	logger *zap.Logger
}

func (s *airQualityServer) DescribeAirQuality(context.Context, *traits.DescribeAirQualityRequest) (*traits.AirQualitySupport, error) {
	return &traits.AirQualitySupport{
		ResourceSupport: &types.ResourceSupport{
			Readable:   true,
			Writable:   false,
			Observable: false,
		},
		CarbonDioxideLevel: &types.FloatBounds{Min: ptr[float32](-2.0), Max: ptr[float32](5000.0)},
		// the device specified VOC ranges in PPB so we have to convert down
		VolatileOrganicCompounds: &types.FloatBounds{Min: ptr[float32](-2.0 / 1000), Max: ptr[float32](5000.0 / 1000)},
	}, nil
}

func (s *airQualityServer) GetAirQuality(ctx context.Context, _ *traits.GetAirQualityRequest) (*traits.AirQuality, error) {
	points, err := FetchPoints(ctx, s.client, PointCO2, PointVOC)
	if err != nil {
		s.logger.Error("failed to fetch air quality points", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "failed to fetch air quality points")
	}

	return &traits.AirQuality{
		CarbonDioxideLevel:       ptr(float32(points.CO2)),
		VolatileOrganicCompounds: ptr(float32(points.VOC / 1000)),
		InfectionRisk:            nil,
	}, nil
}

func ptr[T any](t T) *T {
	return &t
}
