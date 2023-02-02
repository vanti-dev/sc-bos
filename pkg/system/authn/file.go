package authn

import (
	"errors"
	"fmt"

	"github.com/vanti-dev/sc-bos/internal/auth/tenant"
	"github.com/vanti-dev/sc-bos/pkg/system/authn/config"
	"go.uber.org/multierr"
)

// loadFileVerifier returns a tenant.Verifier that checks credentials against those found in a json file.
func loadFileVerifier(idConfig *config.Identities, dataDir, defaultFilename string) (tenant.Verifier, error) {
	ids, err := idConfig.Load(dataDir, defaultFilename)
	if err != nil {
		return nil, fmt.Errorf("local accounts: %w", err)
	}
	if len(ids) == 0 {
		return tenant.NeverVerify(errors.New("no local accounts")), nil
	}

	verifier := &tenant.MemoryVerifier{}
	var allErrs error
	for _, t := range ids {
		err := verifier.AddRecord(tenant.SecretData{TenantID: t.ID, Zones: t.Zones})
		if err != nil {
			allErrs = multierr.Append(allErrs, err)
			continue
		}
		for _, secret := range t.Secrets {
			_, err := verifier.AddSecretHash(t.ID, []byte(secret.Hash))
			allErrs = multierr.Append(allErrs, err)
		}
	}
	if allErrs != nil {
		return nil, fmt.Errorf("local accounts: %w", allErrs)
	}
	return verifier, nil
}
