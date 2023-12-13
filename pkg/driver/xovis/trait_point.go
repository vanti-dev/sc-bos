package xovis

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type pointServer struct {
	gen.UnimplementedPointApiServer
	gen.UnimplementedPointInfoServer

	client      *Client
	multiSensor bool
	logger      *zap.Logger
}

func (s *pointServer) DescribePoints(ctx context.Context, req *gen.DescribePointsRequest) (*gen.PointsSupport, error) {
	logics, err := GetLiveLogics(s.client, s.multiSensor)
	if err != nil {
		s.logger.Error("cannot communicate with device", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "cannot communicate with device")
	}

	pointValues := logicsToPoints(logics)

	var points []*gen.PointMetadata
	for pointName := range pointValues {
		points = append(points, &gen.PointMetadata{
			Name: pointName,
			Kind: gen.PointKind_COUNT,
		})
	}

	slices.SortFunc(points, func(a, b *gen.PointMetadata) bool {
		return a.Name < b.Name
	})

	return &gen.PointsSupport{Points: points}, nil
}

func (s *pointServer) GetPoints(ctx context.Context, req *gen.GetPointsRequest) (*gen.Points, error) {
	logics, err := GetLiveLogics(s.client, s.multiSensor)
	if err != nil {
		s.logger.Error("cannot communicate with device", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "cannot communicate with device")
	}

	values := logicsToPoints(logics)
	valuesAsStruct, err := structpb.NewStruct(castMapKeysToAny(values))
	if err != nil {
		s.logger.Error("cannot convert points into structpb.Struct", zap.Error(err))
		return nil, status.Error(codes.Internal, "point conversion failed")
	}

	return &gen.Points{
		Values: valuesAsStruct,
	}, nil
}

func logicsToPoints(logics LiveLogicsResponse) map[string]int {
	values := make(map[string]int)

	for _, logic := range logics.Logics {
		for _, count := range logic.Counts {
			pointName := fmt.Sprintf("logic-%d-count-%d", logic.ID, count.ID)
			values[pointName] = count.Value
		}
	}

	return values
}

func castMapKeysToAny[K comparable, V any](input map[K]V) map[K]any {
	output := make(map[K]any)
	for k, v := range input {
		output[k] = v
	}
	return output
}
