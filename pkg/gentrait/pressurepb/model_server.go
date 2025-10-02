package pressurepb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedPressureApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{
		model: model,
	}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterPressureApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetPressure(_ context.Context, request *gen.GetPressureRequest) (*gen.Pressure, error) {
	return m.model.GetPressure()
}

func (m *ModelServer) PullPressure(request *gen.PullPressureRequest, server gen.PressureApi_PullPressureServer) error {
	for change := range m.model.PullPressure(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullPressureResponse{Changes: []*gen.PullPressureResponse_Change{
			{
				Name:       request.Name,
				ChangeTime: timestamppb.New(change.ChangeTime),
				Pressure:   change.Value,
			},
		}}

		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (m *ModelServer) UpdatePressure(_ context.Context, request *gen.UpdatePressureRequest) (*gen.Pressure, error) {
	if request.GetDelta() {
		current, err := m.model.GetPressure()
		if err != nil {
			return nil, err
		}

		return m.model.UpdatePressure(&gen.Pressure{TargetPressure: ptr(*current.TargetPressure + *request.Pressure.TargetPressure)})
	}
	return m.model.UpdatePressure(request.Pressure, resource.WithUpdateMask(request.UpdateMask))
}

func ptr[T any](v T) *T {
	return &v
}
