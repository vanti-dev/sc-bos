package prioritypb

import (
	"context"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModelServer struct {
	gen.UnimplementedPriorityApiServer
	model Model
}

func newModelServer(model Model) *ModelServer {
	return &ModelServer{
		model: model,
	}
}

func (m *ModelServer) ClearPriorityEntry(_ context.Context, request *gen.ClearPriorityValueRequest) (*gen.ClearPriorityValueResponse, error) {
	switch id := request.Id.(type) {
	case *gen.ClearPriorityValueRequest_EntryIndex:
		return &gen.ClearPriorityValueResponse{}, m.model.ClearIndex(id.EntryIndex)
	case *gen.ClearPriorityValueRequest_EntryName:
		return &gen.ClearPriorityValueResponse{}, m.model.ClearName(id.EntryName)
	}
	return nil, status.Errorf(codes.InvalidArgument, "unsupported id type %v", request.Id)
}
