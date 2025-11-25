package meter

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedMeterApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterMeterApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetMeterReading(_ context.Context, request *gen.GetMeterReadingRequest) (*gen.MeterReading, error) {
	return m.model.GetMeterReading(resource.WithReadMask(request.ReadMask))
}

func (m *ModelServer) PullMeterReadings(request *gen.PullMeterReadingsRequest, server gen.MeterApi_PullMeterReadingsServer) error {
	for change := range m.model.PullMeterReadings(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullMeterReadingsResponse{Changes: []*gen.PullMeterReadingsResponse_Change{{
			Name:         request.Name,
			ChangeTime:   timestamppb.New(change.ChangeTime),
			MeterReading: change.Value,
		}}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
