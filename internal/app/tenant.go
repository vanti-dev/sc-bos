package app

import (
	"fmt"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/vanti-dev/bsp-ew/internal/auth/tenant"
	"go.uber.org/multierr"
	"os"
	"path/filepath"
)

const tenantFile = "tenants.json"

type tenantConfigFile []tenantConfig

type tenantConfig struct {
	Title   string               `json:"title,omitempty"`
	ID      string               `json:"id,omitempty"`
	Secrets []tenantSecretConfig `json:"secrets,omitempty"`
	Zones   []string             `json:"zones,omitempty"`
}

type tenantSecretConfig struct {
	Note string `json:"note,omitempty"`
	Hash string `json:"hash,omitempty"`
}

func LocalTenantVerifier(config SystemConfig) (tenant.Verifier, error) {
	file, err := os.ReadFile(filepath.Join(config.DataDir, tenantFile))
	if err != nil {
		return nil, fmt.Errorf("reading file %w", err)
	}

	var configFile tenantConfigFile
	if err := json.Unmarshal(file, &configFile); err != nil {
		return nil, fmt.Errorf("unmarshal json %w", err)
	}

	verifier := &tenant.MemoryVerifier{}

	for _, t := range configFile {
		e := verifier.AddRecord(tenant.SecretData{TenantID: t.ID, Zones: t.Zones})
		err = multierr.Append(err, e)
		if e != nil {
			continue
		}
		for _, secret := range t.Secrets {
			_, e := verifier.AddSecretHash(t.ID, []byte(secret.Hash))
			err = multierr.Append(err, e)
		}
	}

	return verifier, err
}
