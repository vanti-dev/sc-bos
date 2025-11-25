package historypb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/history"
)

type StatusServer struct {
	gen.UnimplementedStatusHistoryServer
	store history.Store // payloads of *gen.StatusLog
}

func NewStatusServer(store history.Store) *StatusServer {
	return &StatusServer{store: store}
}

func (m *StatusServer) Register(server *grpc.Server) {
	gen.RegisterStatusHistoryServer(server, m)
}

func (m *StatusServer) Unwrap() any {
	return m.store
}

var currentStatusPager = NewPageReader(func(r history.Record) (*gen.StatusLogRecord, error) {
	v := &gen.StatusLog{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.StatusLogRecord{
		RecordTime:    timestamppb.New(r.CreateTime),
		CurrentStatus: v,
	}, nil
})

func (m *StatusServer) ListCurrentStatusHistory(ctx context.Context, request *gen.ListCurrentStatusHistoryRequest) (*gen.ListCurrentStatusHistoryResponse, error) {
	page, size, nextToken, err := currentStatusPager.ListRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}
	return &gen.ListCurrentStatusHistoryResponse{
		TotalSize:            int32(size),
		NextPageToken:        nextToken,
		CurrentStatusRecords: page,
	}, nil
}
