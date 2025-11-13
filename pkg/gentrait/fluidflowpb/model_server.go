package fluidflowpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedFluidFlowApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{
		model: model,
	}
}
func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterFluidFlowApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetFluidFlow(_ context.Context, _ *gen.GetFluidFlowRequest) (*gen.FluidFlow, error) {
	return m.model.GetFluidFlow()
}

func (m *ModelServer) PullFluidFlow(request *gen.PullFluidFlowRequest, server gen.FluidFlowApi_PullFluidFlowServer) error {
	for change := range m.model.PullFluidFlow(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullFluidFlowResponse{Changes: []*gen.PullFluidFlowResponse_Change{
			{
				Name:       request.Name,
				ChangeTime: timestamppb.New(change.ChangeTime),
				Flow:       change.Value,
			},
		}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (m *ModelServer) UpdateFluidFlow(_ context.Context, request *gen.UpdateFluidFlowRequest) (*gen.FluidFlow, error) {
	return m.model.UpdateFluidFlow(request.Flow, resource.WithUpdateMask(request.UpdateMask))
}
