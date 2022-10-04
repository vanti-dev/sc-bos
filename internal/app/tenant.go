package app

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/vanti-dev/bsp-ew/internal/auth/tenant"
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

func LocalTenantVerifier(config SystemConfig) (tenant.SecretSource, error) {
	file, err := os.ReadFile(filepath.Join(config.DataDir, tenantFile))
	if err != nil {
		return nil, fmt.Errorf("reading file %w", err)
	}

	var configFile tenantConfigFile
	if err := json.Unmarshal(file, &configFile); err != nil {
		return nil, fmt.Errorf("unmarshal json %w", err)
	}

	byHash := make(map[[32]byte]tenantConfig)
	for _, t := range configFile {
		for _, secret := range t.Secrets {
			var hash [32]byte
			copy(hash[:], secret.Hash)
			byHash[hash] = t
		}
	}

	return tenant.SecretSourceFunc(func(ctx context.Context, secret string) (data tenant.SecretData, err error) {
		hash := sha256.Sum256([]byte(secret))
		t, ok := byHash[hash]
		if !ok {
			return data, errors.New("unknown tenant")
		}
		data.TenantID = t.ID
		data.Zones = t.Zones
		return data, nil
	}), nil
}
