package alert

import (
	"context"

	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedAlertApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterAlertApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) ListAlerts(_ context.Context, request *gen.ListAlertsRequest) (*gen.ListAlertsResponse, error) {
	alert := m.model.GetAllAlerts()
	return &gen.ListAlertsResponse{
		Alerts: alert,
	}, nil
}
