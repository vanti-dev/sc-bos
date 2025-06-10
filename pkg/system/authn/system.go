package authn

import (
	"context"
	"errors"
	"fmt"

	"github.com/vanti-dev/sc-bos/internal/auth/tenant"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/system/authn/config"
)

func (s *System) systemTenantVerifier(cfg config.Root) (tenant.Verifier, error) {
	// System account credentials are verified and exchanged to an access token using the
	// OAuth2 Client Credentials flow.
	// There are multiple actors who can verify credentials, these are setup here.
	// Only access tokens that were produced by this verification will be validated.

	// verify system accounts using file provided identities
	fileVerifier, err := loadFileVerifier(cfg.System.FileAccounts, s.configDirs, "tenants.json")
	if err != nil {
		// note, NotFound errors are already handled by loadFileVerifier
		return nil, fmt.Errorf("system %w", err)
	}

	// verify system accounts using the local systems/tenants package, via TenantApi
	nodeVerifier := tenant.NeverVerify(errors.New("tenant system verification not enabled"))
	if cfg.System.TenantAccounts {
		client := gen.NewTenantApiClient(s.clienter.ClientConn())
		nodeVerifier = tenant.VerifierFunc(func(ctx context.Context, id, secret string) (tenant.SecretData, error) {
			return tenant.RemoteVerify(ctx, id, secret, client)
		})
	}

	// verify system accounts using the cohort manager, via TenantApi
	cohortVerifier := tenant.NeverVerify(errors.New("cohort verification not enabled"))
	if cfg.System.CohortAccounts {
		cohortVerifier = tenant.VerifierFunc(func(ctx context.Context, id, secret string) (data tenant.SecretData, err error) {
			conn, err := s.cohortManager.Connect(ctx)
			if err != nil {
				return data, err
			}
			if conn == nil {
				return data, errors.New("cohort manager not available")
			}
			return tenant.RemoteVerify(ctx, id, secret, gen.NewTenantApiClient(conn))
		})
	}

	verifier := tenant.FirstSuccessfulVerifier{
		fileVerifier,
		nodeVerifier,
		cohortVerifier,
	}
	return verifier, nil
}
