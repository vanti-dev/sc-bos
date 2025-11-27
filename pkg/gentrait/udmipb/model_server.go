package udmipb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedUdmiServiceServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterUdmiServiceServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetExportMessage(_ context.Context, request *gen.GetExportMessageRequest) (*gen.MqttMessage, error) {
	msg, err := m.model.GetExportMessage()
	if err != nil {
		return nil, err
	}
	if msg.Topic == "" && msg.Payload == "" {
		return nil, status.Error(codes.Unavailable, "no message available")
	}
	return msg, nil
}

func (m *ModelServer) PullExportMessages(request *gen.PullExportMessagesRequest, server grpc.ServerStreamingServer[gen.PullExportMessagesResponse]) error {
	for change := range m.model.PullExportMessages(server.Context(), resource.WithUpdatesOnly(!request.IncludeLast)) {
		msg := &gen.PullExportMessagesResponse{
			Name:    request.Name,
			Message: change.Value,
		}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
