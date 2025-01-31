package securityevent

import (
	"context"

	"google.golang.org/grpc"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedSecurityEventApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterSecurityEventApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) ListSecurityEvents(ctx context.Context, req *gen.ListSecurityEventsRequest) (*gen.ListSecurityEventsResponse, error) {
	return m.model.ListSecurityEvents(req)
}

// PullSecurityEvents returns a channel of security events
// If updatesOnly is false, only the previous 50 events will be sent before any new events
// For historical events use ListSecurityEvents
func (m *ModelServer) PullSecurityEvents(request *gen.PullSecurityEventsRequest, server gen.SecurityEventApi_PullSecurityEventsServer) error {
	return m.model.PullSecurityEventsWrapper(request, server)
}
