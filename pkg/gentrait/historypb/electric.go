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

type ElectricServer struct {
	gen.UnimplementedElectricHistoryServer
	store history.Store // payloads of *traits.ElectricDemand
}

func NewElectricServer(store history.Store) *ElectricServer {
	return &ElectricServer{store: store}
}

func (m *ElectricServer) Register(server *grpc.Server) {
	gen.RegisterElectricHistoryServer(server, m)
}

func (m *ElectricServer) Unwrap() any {
	return m.store
}

var electricDemandPager = NewPageReader(func(r history.Record) (*gen.ElectricDemandRecord, error) {
	v := &traits.ElectricDemand{}
	err := proto.Unmarshal(r.Payload, v)
	if err != nil {
		return nil, err
	}
	return &gen.ElectricDemandRecord{
		RecordTime:     timestamppb.New(r.CreateTime),
		ElectricDemand: v,
	}, nil
})

func (m *ElectricServer) ListElectricDemandHistory(ctx context.Context, request *gen.ListElectricDemandHistoryRequest) (*gen.ListElectricDemandHistoryResponse, error) {
	page, size, nextToken, err := electricDemandPager.ListRecords(ctx, m.store, request.Period, int(request.PageSize), request.PageToken, request.OrderBy)
	if err != nil {
		return nil, err
	}

	return &gen.ListElectricDemandHistoryResponse{
		TotalSize:             int32(size),
		NextPageToken:         nextToken,
		ElectricDemandRecords: page,
	}, nil
}
