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

type AirTemperatureServer struct {
	gen.UnimplementedAirTemperatureHistoryServer
	store history.Store // payloads of *traits.AirTemperature
}

func NewAirTemperatureServer(store history.Store) *AirTemperatureServer {
	return &AirTemperatureServer{store: store}
}

func (m *AirTemperatureServer) Register(server *grpc.Server) {
	gen.RegisterAirTemperatureHistoryServer(server, m)
}

func (m *AirTemperatureServer) Unwrap() any {
	return m.store
}

var airTemperatureReadingPager = NewPageReader(func(r history.Record) (*gen.AirTemperatureRecord, error) {
	v := &traits.AirTemperature{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.AirTemperatureRecord{
		RecordTime:     timestamppb.New(r.CreateTime),
		AirTemperature: v,
	}, nil
})

func (m *AirTemperatureServer) ListAirTemperatureHistory(ctx context.Context, request *gen.ListAirTemperatureHistoryRequest) (*gen.ListAirTemperatureHistoryResponse, error) {
	page, size, nextToken, err := airTemperatureReadingPager.ListRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}

	return &gen.ListAirTemperatureHistoryResponse{
		TotalSize:             int32(size),
		NextPageToken:         nextToken,
		AirTemperatureRecords: page,
	}, nil
}
