package soundsensorpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedSoundSensorApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterSoundSensorApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetSoundLevel(_ context.Context, request *gen.GetSoundLevelRequest) (*gen.SoundLevel, error) {
	return m.model.GetSoundLevel(resource.WithReadMask(request.ReadMask))
}

func (m *ModelServer) PullSoundLevel(request *gen.PullSoundLevelRequest, server grpc.ServerStreamingServer[gen.PullSoundLevelResponse]) error {
	for change := range m.model.PullSoundLevel(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullSoundLevelResponse{Changes: []*gen.PullSoundLevelResponse_Change{{
			Name:       request.Name,
			ChangeTime: timestamppb.New(change.ChangeTime),
			SoundLevel: change.Value,
		}}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
