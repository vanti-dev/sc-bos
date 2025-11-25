package temperaturepb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedTemperatureApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterTemperatureApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetTemperature(_ context.Context, request *gen.GetTemperatureRequest) (*gen.Temperature, error) {
	return m.model.GetTemperature(resource.WithReadMask(request.ReadMask))
}

func (m *ModelServer) UpdateTemperature(_ context.Context, request *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
	return m.model.UpdateTemperature(request.Temperature, resource.WithUpdateMask(request.UpdateMask))
}

func (m *ModelServer) PullTemperature(request *gen.PullTemperatureRequest, server grpc.ServerStreamingServer[gen.PullTemperatureResponse]) error {
	for change := range m.model.PullTemperature(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullTemperatureResponse{Changes: []*gen.PullTemperatureResponse_Change{{
			Name:        request.Name,
			ChangeTime:  timestamppb.New(change.ChangeTime),
			Temperature: change.Value,
		}}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
