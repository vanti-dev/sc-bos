package serviceticketpb

import (
	"context"

	"google.golang.org/grpc"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedServiceTicketApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterServiceTicketApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) CreateTicket(_ context.Context, req *gen.CreateTicketRequest) (*gen.Ticket, error) {
	return m.model.addTicket(req.Ticket), nil
}

func (m *ModelServer) UpdateTicket(_ context.Context, req *gen.UpdateTicketRequest) (*gen.Ticket, error) {
	return m.model.updateTicket(req.Ticket)
}
