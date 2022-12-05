package pgxutil

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"strings"
)

type ConnectConfig struct {
	URI          string `json:"uri,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}

func Connect(ctx context.Context, sysConf ConnectConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(sysConf.URI)
	if err != nil {
		return nil, err
	}

	if sysConf.PasswordFile != "" {
		passwordFile, err := os.ReadFile(sysConf.PasswordFile)
		if err != nil {
			return nil, err
		}

		poolConfig.ConnConfig.Password = strings.TrimSpace(string(passwordFile))
	}

	return pgxpool.ConnectConfig(ctx, poolConfig)
}
