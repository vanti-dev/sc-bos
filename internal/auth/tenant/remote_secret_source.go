package tenant

import (
	"context"
	"crypto/sha256"

	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RemoteSecretSource struct {
	Logger *zap.Logger
	Client gen.TenantApiClient
}

func (r *RemoteSecretSource) VerifySecret(ctx context.Context, secretStr string) (data SecretData, err error) {
	return RemoteVerify(ctx, secretStr, r.Client, r.Logger)
}

var _ SecretSource = (*RemoteSecretSource)(nil)

func RemoteVerify(ctx context.Context, secretStr string, client gen.TenantApiClient, logger *zap.Logger) (data SecretData, err error) {
	hash := sha256.Sum256([]byte(secretStr))

	secretRecord, err := client.GetSecretByHash(ctx, &gen.GetSecretByHashRequest{
		SecretHash: hash[:],
	})
	if err != nil {
		if status.Code(err) != codes.NotFound {
			logger.Error("failed to retrieve secret", zap.Error(err))
		}
		return
	}

	tenant, err := client.GetTenant(ctx, &gen.GetTenantRequest{
		Id: secretRecord.Tenant.GetId(),
	})
	if err != nil {
		if status.Code(err) != codes.NotFound {
			logger.Error("failed to retrieve tenant", zap.Error(err),
				zap.String("tenant_id", secretRecord.Tenant.GetId()))
		}
		return
	}

	data = SecretData{
		TenantID: tenant.Id,
		Zones:    tenant.ZoneNames,
	}
	return
}
