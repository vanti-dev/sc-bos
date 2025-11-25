package accesspb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedAccessApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterAccessApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetLastAccessAttempt(ctx context.Context, request *gen.GetLastAccessAttemptRequest) (*gen.AccessAttempt, error) {
	return m.model.GetLastAccessAttempt(resource.WithReadMask(request.GetReadMask()))
}

func (m *ModelServer) PullAccessAttempts(request *gen.PullAccessAttemptsRequest, server gen.AccessApi_PullAccessAttemptsServer) error {
	for change := range m.model.PullAccessAttempts(server.Context(), resource.WithReadMask(request.GetReadMask()), resource.WithUpdatesOnly(request.GetUpdatesOnly())) {
		msg := &gen.PullAccessAttemptsResponse{Changes: []*gen.PullAccessAttemptsResponse_Change{{
			Name:          request.Name,
			ChangeTime:    timestamppb.New(change.ChangeTime),
			AccessAttempt: change.Value,
		}}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
