package authn

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"go.uber.org/multierr"

	"github.com/vanti-dev/sc-bos/internal/auth/tenant"
	"github.com/vanti-dev/sc-bos/pkg/system/authn/config"
)

// loadFileVerifier returns a tenant.Verifier that checks credentials against those found in a json file.
func loadFileVerifier(idConfig *config.Identities, dataDirs []string, defaultFilename string) (tenant.Verifier, error) {
	ids := make([]config.Identity, 0)
	for _, dataDir := range dataDirs {
		_, err := os.Stat(path.Join(dataDir, defaultFilename))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				return nil, fmt.Errorf("local accounts: %w", err)
			}
		}
		i, err := idConfig.Load(dataDir, defaultFilename)
		if err != nil {
			return nil, fmt.Errorf("local accounts: %w", err)
		}
		ids = append(ids, i...)
	}
	if len(ids) == 0 {
		log.Printf("no local accounts found in %v", dataDirs)
		return tenant.NeverVerify(errors.New("no local accounts")), nil
	}

	verifier := &tenant.MemoryVerifier{}
	var allErrs error
	for _, t := range ids {
		err := verifier.AddRecord(tenant.SecretData{Title: t.Title, TenantID: t.ID, Zones: t.Zones, Roles: t.Roles})
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
