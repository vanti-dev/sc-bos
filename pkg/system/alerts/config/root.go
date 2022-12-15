package config

import (
	"github.com/vanti-dev/sc-bos/pkg/util/pgxutil"
)

type Root struct {
	Storage *Storage `json:"storage,omitempty"`
}

type Storage struct {
	Type string `json:"type,omitempty"`
	pgxutil.ConnectConfig
}
