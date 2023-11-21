package config

import (
	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	system.Config
	Storage *Storage `json:"storage,omitempty"`
}

type StorageType string

const (
	StorageTypePostgres StorageType = "postgres"
	StorageTypeBolt     StorageType = "bolt"
)

type Storage struct {
	Type StorageType `json:"type,omitempty"`
	pgxutil.ConnectConfig
	TTL *TTL `json:"ttl,omitempty"`
	// Retention is the minimum time records should be stored for. Zero-value (not-specified) means "forever".
	// Records can be deleted after this period, but may be kept longer depending on the cleanup cycle (e.g. if records
	// are only pruned once a day, a record could be kept for retention + 1day). Not all storage types might support this.
	Retention jsontypes.Duration `json:"retention,omitempty"`
}

type TTL struct {
	MaxAge   jsontypes.Duration `json:"maxAge,omitempty"`
	MaxCount int64              `json:"maxCount,omitempty"`
}
