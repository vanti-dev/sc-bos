package pressurepb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
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
	resourceOpts := []resource.WriteOption{resource.WithUpdateMask(request.UpdateMask)}

	if request.GetDelta() {
		resourceOpts = append(resourceOpts, resource.InterceptBefore(func(old, new proto.Message) {
			oldPressure, ok := old.(*gen.Pressure)
			if !ok {
				return
			}
			newPressure, ok := new.(*gen.Pressure)
			if !ok {
				return
			}
			*newPressure.TargetPressure += *oldPressure.TargetPressure
		}))
	}
	return m.model.UpdatePressure(request.Pressure, resourceOpts...)
}
