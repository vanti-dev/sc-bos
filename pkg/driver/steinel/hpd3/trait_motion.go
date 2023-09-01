package hpd3

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type motionServer struct {
	traits.UnimplementedMotionSensorApiServer
	traits.UnimplementedMotionSensorSensorInfoServer

	client Client
	logger *zap.Logger
}

func (s *motionServer) DescribeMotionDetection(context.Context, *traits.DescribeMotionDetectionRequest) (*traits.MotionDetectionSupport, error) {
	return &traits.MotionDetectionSupport{
		ResourceSupport: &types.ResourceSupport{
			Readable:   true,
			Writable:   false,
			Observable: false,
		},
	}, nil
}

func (s *motionServer) GetMotionDetection(ctx context.Context, _ *traits.GetMotionDetectionRequest) (*traits.MotionDetection, error) {
	value, err := s.client.FetchSensorData(ctx, pointMotion1)
	if err != nil {
		s.logger.Error("unable to fetch motion points from device", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "unable to fetch motion points from device")
	}
	motion, ok := value.(bool)
	if !ok {
		s.logger.Error("data point is not a boolean", zap.String("point", pointMotion1), zap.Any("value", value))
		return nil, status.Error(codes.Internal, "device returned an unexpected data type")
	}

	var state traits.MotionDetection_State
	if motion {
		state = traits.MotionDetection_DETECTED
	} else {
		state = traits.MotionDetection_NOT_DETECTED
	}
	return &traits.MotionDetection{
		State: state,
	}, nil
}
