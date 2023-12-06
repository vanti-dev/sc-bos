package hpd3

import (
	"context"

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

	client Client
	logger *zap.Logger
}

var supportedPoints = []*gen.PointMetadata{
	{
		Name: PointPresence1,
	},
	{
		Name: PointMotion1,
	},
	{
		Name: PointTemperature,
	},
	{
		Name: PointHumidity,
	},
	{
		Name: PointNumberOfPeopleTotal,
	},
	{
		Name: PointCO2,
	},
	{
		Name: PointVOC,
	},
}

var supportedPointNames = getPointNames(supportedPoints)

func (s *pointServer) DescribePoints(ctx context.Context, request *gen.DescribePointsRequest) (*gen.PointsSupport, error) {
	return &gen.PointsSupport{Points: supportedPoints}, nil
}

func (s *pointServer) GetPoints(ctx context.Context, request *gen.GetPointsRequest) (*gen.Points, error) {
	data, err := FetchPoints(ctx, s.client, supportedPointNames...)
	if err != nil {
		s.logger.Error("failed to fetch some data points", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "failed to fetch some data points from device")
	}

	values, err := structpb.NewStruct(data.AsMap())
	if err != nil {
		s.logger.Error("can't convert data points to proto struct", zap.Error(err))
		return nil, status.Error(codes.Internal, "data conversion failure")
	}

	return &gen.Points{Values: values}, nil
}

func getPointNames(pointMeta []*gen.PointMetadata) []string {
	var result []string
	for _, p := range pointMeta {
		result = append(result, p.Name)
	}
	slices.Sort(result)
	result = slices.Compact(result)
	return result
}
