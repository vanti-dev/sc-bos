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

type Server struct {
	gen.UnimplementedTenantApiServer

	mu   sync.Mutex
	impl gen.TenantApiClient
}

func (s *Server) Register(server *grpc.Server) {
	gen.RegisterTenantApiServer(server, s)
}

func (s *Server) Fill(impl gen.TenantApiClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.impl = impl
}

func (s *Server) Empty() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.impl = nil
}

func (s *Server) ListTenants(ctx context.Context, request *gen.ListTenantsRequest) (*gen.ListTenantsResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.ListTenants(ctx, request)
}

func (s *Server) PullTenants(request *gen.PullTenantsRequest, server gen.TenantApi_PullTenantsServer) error {
	c, err := s.client()
	if err != nil {
		return err
	}
	stream, err := c.PullTenants(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (s *Server) CreateTenant(ctx context.Context, request *gen.CreateTenantRequest) (*gen.Tenant, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.CreateTenant(ctx, request)
}

func (s *Server) GetTenant(ctx context.Context, request *gen.GetTenantRequest) (*gen.Tenant, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.GetTenant(ctx, request)
}

func (s *Server) UpdateTenant(ctx context.Context, request *gen.UpdateTenantRequest) (*gen.Tenant, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.UpdateTenant(ctx, request)
}

func (s *Server) DeleteTenant(ctx context.Context, request *gen.DeleteTenantRequest) (*gen.DeleteTenantResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.DeleteTenant(ctx, request)
}

func (s *Server) PullTenant(request *gen.PullTenantRequest, server gen.TenantApi_PullTenantServer) error {
	c, err := s.client()
	if err != nil {
		return err
	}
	stream, err := c.PullTenant(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (s *Server) AddTenantZones(ctx context.Context, request *gen.AddTenantZonesRequest) (*gen.Tenant, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.AddTenantZones(ctx, request)
}

func (s *Server) RemoveTenantZones(ctx context.Context, request *gen.RemoveTenantZonesRequest) (*gen.Tenant, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.RemoveTenantZones(ctx, request)
}

func (s *Server) ListSecrets(ctx context.Context, request *gen.ListSecretsRequest) (*gen.ListSecretsResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.ListSecrets(ctx, request)
}

func (s *Server) PullSecrets(request *gen.PullSecretsRequest, server gen.TenantApi_PullSecretsServer) error {
	c, err := s.client()
	if err != nil {
		return err
	}
	stream, err := c.PullSecrets(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (s *Server) CreateSecret(ctx context.Context, request *gen.CreateSecretRequest) (*gen.Secret, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.CreateSecret(ctx, request)
}

func (s *Server) VerifySecret(ctx context.Context, request *gen.VerifySecretRequest) (*gen.Secret, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.VerifySecret(ctx, request)
}

func (s *Server) GetSecret(ctx context.Context, request *gen.GetSecretRequest) (*gen.Secret, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.GetSecret(ctx, request)
}

func (s *Server) UpdateSecret(ctx context.Context, request *gen.UpdateSecretRequest) (*gen.Secret, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.UpdateSecret(ctx, request)
}

func (s *Server) DeleteSecret(ctx context.Context, request *gen.DeleteSecretRequest) (*gen.DeleteSecretResponse, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.DeleteSecret(ctx, request)
}

func (s *Server) PullSecret(request *gen.PullSecretRequest, server gen.TenantApi_PullSecretServer) error {
	c, err := s.client()
	if err != nil {
		return err
	}
	stream, err := c.PullSecret(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (s *Server) RegenerateSecret(ctx context.Context, request *gen.RegenerateSecretRequest) (*gen.Secret, error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}
	return c.RegenerateSecret(ctx, request)
}

func (s *Server) client() (gen.TenantApiClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.impl == nil {
		return nil, errNotEnabled
	}
	return s.impl, nil
}
