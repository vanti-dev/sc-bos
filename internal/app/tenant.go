package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/vanti-dev/sc-bos/internal/auth/tenant"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"go.uber.org/multierr"
	"os"
	"path/filepath"
)

const tenantsFilename = "tenants.json"
const usersFilename = "users.json"

type credentialsFile []identity

type identity struct {
	Title   string         `json:"title,omitempty"`
	ID      string         `json:"id,omitempty"`
	Secrets []secretConfig `json:"secrets,omitempty"`
	Zones   []string       `json:"zones,omitempty"`
}

type secretConfig struct {
	Note string `json:"note,omitempty"`
	Hash string `json:"hash,omitempty"`
}

func clientVerifier(config SystemConfig, manager node.Remote) (tenant.Verifier, error) {
	localTenants, err := loadVerifier(config, tenantsFilename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			// if the file exists, but we can't read it, we should let someone know
			return nil, err
		}
		// reading the local tenant data failed, we return this error each time as part of the secret verification
		localTenants = tenant.NeverVerify(err)
	}

	// remoteTenants verifies tenant access using a remote service defined via TenantApiClient and managerConn
	remoteTenants := tenant.VerifierFunc(func(ctx context.Context, id, secret string) (data tenant.SecretData, err error) {
		conn, err := manager.Connect(ctx)
		if err != nil {
			return data, err
		}
		if conn == nil {
			return data, errors.New("no remote clientVerifier")
		}
		return tenant.RemoteVerify(ctx, id, secret, gen.NewTenantApiClient(conn))
	})
	tenantVerifier := tenant.FirstSuccessfulVerifier([]tenant.Verifier{
		localTenants,
		remoteTenants,
	})
	return tenantVerifier, nil
}

func passwordVerifier(config SystemConfig) (tenant.Verifier, error) {
	return loadVerifier(config, usersFilename)
}

func loadVerifier(config SystemConfig, filename string) (tenant.Verifier, error) {
	file, err := os.ReadFile(filepath.Join(config.DataDir, filename))
	if err != nil {
		return nil, fmt.Errorf("reading file %w", err)
	}

	var configFile credentialsFile
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
