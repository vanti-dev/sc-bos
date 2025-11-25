package transport

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedTransportApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterTransportApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetTransport(_ context.Context, request *gen.GetTransportRequest) (*gen.Transport, error) {
	return m.model.GetTransport(resource.WithReadMask(request.ReadMask))
}

func (m *ModelServer) PullTransport(request *gen.PullTransportRequest, server grpc.ServerStreamingServer[gen.PullTransportResponse]) error {
	for change := range m.model.PullTransport(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullTransportResponse{Changes: []*gen.PullTransportResponse_Change{{
			Name:       request.Name,
			ChangeTime: timestamppb.New(change.ChangeTime),
			Transport:  change.Value,
		}}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
