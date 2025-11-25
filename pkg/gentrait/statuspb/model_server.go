package statuspb

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedStatusApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{
		model: model,
	}
}

func (s *ModelServer) GetCurrentStatus(_ context.Context, req *gen.GetCurrentStatusRequest) (*gen.StatusLog, error) {
	return s.model.GetCurrentStatus(req.ReadMask)
}

func (s *ModelServer) PullCurrentStatus(request *gen.PullCurrentStatusRequest, server gen.StatusApi_PullCurrentStatusServer) error {
	ctx := server.Context()
	changes := s.model.PullCurrentStatus(ctx, request.ReadMask, request.UpdatesOnly)
	for change := range changes {
		err := server.Send(&gen.PullCurrentStatusResponse{Changes: []*gen.PullCurrentStatusResponse_Change{{
			Name:          request.Name,
			CurrentStatus: change.StatusLog,
			ChangeTime:    timestamppb.New(change.ChangeTime),
		}}})
		if err != nil {
			return err
		}
	}
	return nil
}
