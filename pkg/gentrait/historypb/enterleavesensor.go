package historypb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

type EnterLeaveSensorServer struct {
	gen.UnimplementedEnterLeaveHistoryServer
	store history.Store // payloads of *traits.EnterLeaveEvent
}

func NewEnterLeaveSensorServer(store history.Store) *EnterLeaveSensorServer {
	return &EnterLeaveSensorServer{store: store}
}

func (e *EnterLeaveSensorServer) Register(server *grpc.Server) {
	gen.RegisterEnterLeaveHistoryServer(server, e)
}

func (e *EnterLeaveSensorServer) Unwrap() any {
	return e.store
}

var enterLeaveEventPager = NewPageReader(func(r history.Record) (*gen.EnterLeaveEventRecord, error) {
	v := &traits.EnterLeaveEvent{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.EnterLeaveEventRecord{
		RecordTime:      timestamppb.New(r.CreateTime),
		EnterLeaveEvent: v,
	}, nil
})

func (e *EnterLeaveSensorServer) ListEnterLeaveEventHistory(ctx context.Context, request *gen.ListEnterLeaveHistoryRequest) (*gen.ListEnterLeaveHistoryResponse, error) {
	page, size, nextToken, err := enterLeaveEventPager.ListRecords(ctx, e.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}

	return &gen.ListEnterLeaveHistoryResponse{
		TotalSize:         int32(size),
		NextPageToken:     nextToken,
		EnterLeaveRecords: page,
	}, nil
}
