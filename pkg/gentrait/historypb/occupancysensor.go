package historypb

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OccupancySensorServer struct {
	gen.UnimplementedOccupancySensorHistoryServer
	store history.Store // payloads of *traits.Occupancy
}

func NewOccupancySensorServer(store history.Store) *OccupancySensorServer {
	return &OccupancySensorServer{store: store}
}

func (m *OccupancySensorServer) Register(server *grpc.Server) {
	gen.RegisterOccupancySensorHistoryServer(server, m)
}

func (m *OccupancySensorServer) Unwrap() any {
	return m.store
}

var occupancyPager = newPageReader(func(r history.Record) (*gen.OccupancyRecord, error) {
	v := &traits.Occupancy{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.OccupancyRecord{
		RecordTime: timestamppb.New(r.CreateTime),
		Occupancy:  v,
	}, nil
})

func (m *OccupancySensorServer) ListOccupancyHistory(ctx context.Context, request *gen.ListOccupancyHistoryRequest) (*gen.ListOccupancyHistoryResponse, error) {
	page, size, nextToken, err := occupancyPager.listRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken)
	if err != nil {
		return nil, err
	}

	return &gen.ListOccupancyHistoryResponse{
		TotalSize:        int32(size),
		NextPageToken:    nextToken,
		OccupancyRecords: page,
	}, nil
}
