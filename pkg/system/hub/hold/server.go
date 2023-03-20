package hold

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var errNotEnabled = status.Error(codes.FailedPrecondition, "not enabled")

// Server forwards all HubApiServer calls to impl if present or returns errNotEnabled.
type Server struct {
	gen.UnimplementedHubApiServer

	mu   sync.Mutex
	impl gen.HubApiClient
}

func (s *Server) Register(server *grpc.Server) {
	gen.RegisterHubApiServer(server, s)
}

func (s *Server) Fill(impl gen.HubApiClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.impl = impl
}

func (s *Server) Empty() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.impl = nil
}

func (s *Server) GetHubNode(ctx context.Context, request *gen.GetHubNodeRequest) (*gen.HubNode, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.GetHubNode(ctx, request)
}

func (s *Server) EnrollHubNode(ctx context.Context, request *gen.EnrollHubNodeRequest) (*gen.HubNode, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.EnrollHubNode(ctx, request)
}

func (s *Server) ListHubNodes(ctx context.Context, request *gen.ListHubNodesRequest) (*gen.ListHubNodesResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.ListHubNodes(ctx, request)
}

func (s *Server) TestHubNode(ctx context.Context, request *gen.TestHubNodeRequest) (*gen.TestHubNodeResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.TestHubNode(ctx, request)
}

func (s *Server) client() (gen.HubApiClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.impl == nil {
		return nil, errNotEnabled
	}
	return s.impl, nil
}
