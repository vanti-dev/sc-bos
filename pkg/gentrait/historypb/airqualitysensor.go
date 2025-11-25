package historypb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/history"
)

type AirQualitySensorServer struct {
	gen.UnimplementedAirQualitySensorHistoryServer
	store history.Store // payloads of *traits.AirQuality
}

func NewAirQualitySensorServer(store history.Store) *AirQualitySensorServer {
	return &AirQualitySensorServer{store: store}
}

func (m *AirQualitySensorServer) Register(server *grpc.Server) {
	gen.RegisterAirQualitySensorHistoryServer(server, m)
}

func (m *AirQualitySensorServer) Unwrap() any {
	return m.store
}

var airQualityPager = NewPageReader(func(r history.Record) (*gen.AirQualityRecord, error) {
	v := &traits.AirQuality{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.AirQualityRecord{
		RecordTime: timestamppb.New(r.CreateTime),
		AirQuality: v,
	}, nil
})

func (m *AirQualitySensorServer) ListAirQualityHistory(ctx context.Context, request *gen.ListAirQualityHistoryRequest) (*gen.ListAirQualityHistoryResponse, error) {
	page, size, nextToken, err := airQualityPager.ListRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}

	return &gen.ListAirQualityHistoryResponse{
		TotalSize:         int32(size),
		NextPageToken:     nextToken,
		AirQualityRecords: page,
	}, nil
}
