package hold

import (
	"context"
	"sync"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotEnabled = status.Error(codes.FailedPrecondition, "not enabled")

// Server forwards all NodeApiServer calls to impl if present or returns errNotEnabled.
type Server struct {
	gen.UnimplementedNodeApiServer

	mu   sync.Mutex
	impl gen.NodeApiClient
}

func (s *Server) Register(server *grpc.Server) {
	gen.RegisterNodeApiServer(server, s)
}

func (s *Server) Fill(impl gen.NodeApiClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.impl = impl
}

func (s *Server) Empty() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.impl = nil
}

func (s *Server) GetNodeRegistration(ctx context.Context, request *gen.GetNodeRegistrationRequest) (*gen.NodeRegistration, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.GetNodeRegistration(ctx, request)
}

func (s *Server) CreateNodeRegistration(ctx context.Context, request *gen.CreateNodeRegistrationRequest) (*gen.NodeRegistration, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.CreateNodeRegistration(ctx, request)
}

func (s *Server) ListNodeRegistrations(ctx context.Context, request *gen.ListNodeRegistrationsRequest) (*gen.ListNodeRegistrationsResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.ListNodeRegistrations(ctx, request)
}

func (s *Server) TestNodeCommunication(ctx context.Context, request *gen.TestNodeCommunicationRequest) (*gen.TestNodeCommunicationResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.TestNodeCommunication(ctx, request)
}

func (s *Server) client() (gen.NodeApiClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.impl == nil {
		return nil, errNotEnabled
	}
	return s.impl, nil
}
