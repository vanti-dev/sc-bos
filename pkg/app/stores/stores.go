package stores

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"

	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
)

var (
	ErrStoreClosed        = errors.New("closed")
	ErrStoreNotConfigured = errors.New("not configured")
)

// Config configures shared storage (dbs) for systems on this node.
type Config struct {
	Postgres *PostgresConfig `json:"postgres,omitempty"`
}

type PostgresConfig struct {
	pgxutil.ConnectConfig
}

func New(cfg *Config) *Stores {
	s := &Stores{}
	if cfg != nil && cfg.Postgres != nil {
		s.postgresStore.cfg = cfg.Postgres
	}
	return s
}

// Stores provides access to shared storage connections/clients.
type Stores struct {
	postgresStore
}

// Close closes all stores.
func (s *Stores) Close() error {
	return multierr.Combine(
		s.postgresStore.close(),
	)
}

type postgresStore struct {
	cfg *PostgresConfig

	mu          sync.Mutex
	r, w, admin *pgxpool.Pool
	err         error
	closed      bool
}

// Postgres returns shared postgres connection pools.
// The pools can be used for read, write, and admin (alter table) operations.
func (s *postgresStore) Postgres() (r, w, admin *pgxpool.Pool, _ error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.r != nil || s.err != nil {
		return s.r, s.w, s.admin, s.err
	}

	fail := func(err error) (_, _, _ *pgxpool.Pool, _ error) {
		s.err = err
		return nil, nil, nil, err
	}
	if s.cfg == nil {
		return fail(fmt.Errorf("%w: postgres", ErrStoreNotConfigured))
	}

	// todo: support r, w, and admin pools
	pool, err := pgxutil.Connect(context.Background(), s.cfg.ConnectConfig)
	if err != nil {
		return fail(err)
	}
	s.r, s.w, s.admin = pool, pool, pool
	return s.r, s.w, s.admin, nil
}

func (s *postgresStore) close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = fmt.Errorf("%w: postgres", ErrStoreClosed)
	if s.r == nil {
		return nil
	}

	s.r.Close()
	s.w.Close()
	s.admin.Close()
	s.r = nil
	s.w = nil
	s.admin = nil

	return nil
}
