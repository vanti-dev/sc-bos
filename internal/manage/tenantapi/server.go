package tenantapi

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/jackc/pgx/v4"
	"github.com/vanti-dev/bsp-ew/internal/db"
	"github.com/vanti-dev/bsp-ew/internal/util/rpcutil"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	gen.UnimplementedTenantApiServer
	dbConn *pgx.Conn
	logger *zap.Logger
}

func (s *Server) ListTenants(ctx context.Context, request *gen.ListTenantsRequest) (*gen.ListTenantsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) PullTenants(request *gen.PullTenantsRequest, server gen.TenantApi_PullTenantsServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) CreateTenant(ctx context.Context, request *gen.CreateTenantRequest) (*gen.Tenant, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) GetTenant(ctx context.Context, request *gen.GetTenantRequest) (*gen.Tenant, error) {
	logger := rpcutil.ServerLogger(ctx, s.logger)

	var tenant *gen.Tenant
	err := s.dbConn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		tenant, err = db.GetTenant(ctx, tx, request.GetId())
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "tenant not found")
	} else if err != nil {
		logger.Error("db transaction failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "database read failed")
	}
	return tenant, nil
}

func (s *Server) UpdateTenant(ctx context.Context, request *gen.UpdateTenantRequest) (*gen.Tenant, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) DeleteTenant(ctx context.Context, request *gen.DeleteTenantRequest) (*gen.DeleteTenantResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) PullTenant(request *gen.PullTenantRequest, server gen.TenantApi_PullTenantServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) AddTenantZones(ctx context.Context, request *gen.AddTenantZonesRequest) (*gen.Tenant, error) {
	logger := rpcutil.ServerLogger(ctx, s.logger)

	err := s.dbConn.BeginFunc(ctx, func(tx pgx.Tx) error {
		return db.AddTenantZones(ctx, tx, request.TenantId, request.AddZoneNames)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "tenant not found")
	} else if err != nil {
		logger.Error("db transaction failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "database write failed")
	}

	return s.GetTenant(ctx, &gen.GetTenantRequest{Id: request.TenantId})
}

func (s *Server) RemoveTenantZones(ctx context.Context, request *gen.RemoveTenantZonesRequest) (*gen.Tenant, error) {
	logger := rpcutil.ServerLogger(ctx, s.logger)

	err := s.dbConn.BeginFunc(ctx, func(tx pgx.Tx) error {
		return db.RemoveTenantZones(ctx, tx, request.TenantId, request.RemoveZoneNames)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "tenant not found")
	} else if err != nil {
		logger.Error("db transaction failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "database write failed")
	}

	return s.GetTenant(ctx, &gen.GetTenantRequest{Id: request.TenantId})
}

func (s *Server) ListSecrets(ctx context.Context, request *gen.ListSecretsRequest) (*gen.ListSecretsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) PullSecrets(request *gen.PullSecretsRequest, server gen.TenantApi_PullSecretsServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) CreateSecret(ctx context.Context, request *gen.CreateSecretRequest) (*gen.Secret, error) {
	logger := rpcutil.ServerLogger(ctx, s.logger)
	secret := request.Secret

	var err error
	secret.Secret, err = genSecret()
	if err != nil {
		logger.Error("secret generation failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "secret generation failed")
	}
	secret.SecretHash = hashSecret(secret.Secret)

	err = s.dbConn.BeginFunc(ctx, func(tx pgx.Tx) (err error) {
		secret, err = db.CreateTenantSecret(ctx, tx, secret)
		return
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "tenant not found")
	} else if err != nil {
		logger.Error("db transaction failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "database transaction failed")
	}
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) GetSecret(ctx context.Context, request *gen.GetSecretRequest) (*gen.Secret, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) UpdateSecret(ctx context.Context, request *gen.UpdateSecretRequest) (*gen.Secret, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) DeleteSecret(ctx context.Context, request *gen.DeleteSecretRequest) (*gen.DeleteSecretResponse, error) {
	logger := rpcutil.ServerLogger(ctx, s.logger)
	err := s.dbConn.BeginFunc(ctx, func(tx pgx.Tx) error {
		return db.DeleteTenantSecret(ctx, tx, request.Id)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "secret not found")
	} else if err != nil {
		logger.Error("db transaction failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "database transaction failed")
	}
	return &gen.DeleteSecretResponse{}, nil
}

func (s *Server) PullSecret(request *gen.PullSecretRequest, server gen.TenantApi_PullSecretServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) RegenerateSecret(ctx context.Context, request *gen.RegenerateSecretRequest) (*gen.Secret, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func genSecret() (string, error) {
	secretBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, secretBytes)
	if err != nil {
		return "", err
	}

	encoding := base64.URLEncoding
	encoded := make([]byte, encoding.EncodedLen(len(secretBytes)))
	encoding.Encode(encoded, secretBytes)

	return string(encoded), nil
}

func hashSecret(secret string) (hash []byte) {
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}
