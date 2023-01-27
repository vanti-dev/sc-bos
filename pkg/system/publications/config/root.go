package config

import (
	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/system"
)

type Root struct {
	system.Config
	Storage *Storage `json:"storage,omitempty"`
}

type Storage struct {
	Type string `json:"type,omitempty"`
	pgxutil.ConnectConfig
}
