package historypb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/history"
)

type MeterServer struct {
	gen.UnimplementedMeterHistoryServer
	store history.Store // payloads of *gen.MeterReading
}

func NewMeterServer(store history.Store) *MeterServer {
	return &MeterServer{store: store}
}

func (m *MeterServer) Register(server *grpc.Server) {
	gen.RegisterMeterHistoryServer(server, m)
}

func (m *MeterServer) Unwrap() any {
	return m.store
}

var meterReadingPager = NewPageReader(func(r history.Record) (*gen.MeterReadingRecord, error) {
	v := &gen.MeterReading{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.MeterReadingRecord{
		RecordTime:   timestamppb.New(r.CreateTime),
		MeterReading: v,
	}, nil
})

func (m *MeterServer) ListMeterReadingHistory(ctx context.Context, request *gen.ListMeterReadingHistoryRequest) (*gen.ListMeterReadingHistoryResponse, error) {
	page, size, nextToken, err := meterReadingPager.ListRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}

	return &gen.ListMeterReadingHistoryResponse{
		TotalSize:           int32(size),
		NextPageToken:       nextToken,
		MeterReadingRecords: page,
	}, nil
}
