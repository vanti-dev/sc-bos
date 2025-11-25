package historypb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/history"
)

type TransportServer struct {
	gen.UnimplementedTransportHistoryServer
	store history.Store // payloads of *gen.Transport
}

func NewTransportServer(store history.Store) *TransportServer {
	return &TransportServer{store: store}
}

func (m *TransportServer) Register(server *grpc.Server) {
	gen.RegisterTransportHistoryServer(server, m)
}

func (m *TransportServer) Unwrap() any {
	return m.store
}

var transportPager = NewPageReader(func(r history.Record) (*gen.TransportRecord, error) {
	v := &gen.Transport{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.TransportRecord{
		RecordTime: timestamppb.New(r.CreateTime),
		Transport:  v,
	}, nil
})

func (m *TransportServer) ListTransportHistory(ctx context.Context, request *gen.ListTransportHistoryRequest) (*gen.ListTransportHistoryResponse, error) {
	page, size, nextToken, err := transportPager.ListRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}

	return &gen.ListTransportHistoryResponse{
		TotalSize:        int32(size),
		NextPageToken:    nextToken,
		TransportRecords: page,
	}, nil
}
