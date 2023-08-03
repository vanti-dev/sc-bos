// Package pgxstore provides an implementation of history.Store backed by a Postgres database.
// The historical records for all stores are stored in a single table disambiguated via the source parameter.
package pgxstore

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

//go:embed schema.sql
var schemaSql string

func SetupDB(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, schemaSql)
		return err
	})
}

func New(ctx context.Context, source, connStr string) (history.Store, error) {
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("connect %w", err)
	}

	return SetupStoreFromPool(ctx, source, pool)
}

func SetupStoreFromPool(ctx context.Context, source string, pool *pgxpool.Pool) (history.Store, error) {
	err := SetupDB(ctx, pool)
	if err != nil {
		return nil, fmt.Errorf("setup %w", err)
	}

	return NewStoreFromPool(source, pool), nil
}

func NewStoreFromPool(source string, pool *pgxpool.Pool) history.Store {
	return &Store{
		slice: slice{pool: pool, source: source},
		now:   time.Now,
	}
}

type Store struct {
	slice

	now func() time.Time
}

func (s *Store) Append(ctx context.Context, payload []byte) (history.Record, error) {
	now := s.now()
	r := history.Record{
		CreateTime: now,
		Payload:    payload,
	}

	row := s.pool.QueryRow(ctx, "INSERT INTO history (source, create_time, payload) VALUES ($1, $2, $3) RETURNING id",
		s.source, now, payload)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return history.Record{}, err
	}

	r.ID = strconv.FormatInt(id, 10)
	return r, nil
}

type slice struct {
	pool     *pgxpool.Pool
	source   string // distinguishes between this store and other stores that use the same table
	from, to history.Record
}

func (s slice) Slice(from, to history.Record) history.Slice {
	return slice{
		pool:   s.pool,
		source: s.source,
		from:   from,
		to:     to,
	}
}

func (s slice) Read(ctx context.Context, into []history.Record) (int, error) {
	var where []string
	var args []any
	where, args = s.sourceClause(where, args)
	where, args, err := s.rangeClause(where, args)
	if err != nil {
		return 0, err
	}

	sql := fmt.Sprintf("SELECT id, create_time, payload FROM history WHERE %s ORDER BY id ASC LIMIT %v", strings.Join(where, " AND "), len(into))
	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var id int64
		err := rows.Scan(&id, &into[i].CreateTime, &into[i].Payload)
		if err != nil {
			return 0, err
		}
		into[i].ID = idFromSql(id)
		i++
	}
	return i, nil
}

func (s slice) Len(ctx context.Context) (int, error) {
	var where []string
	var args []any
	where, args = s.sourceClause(where, args)
	where, args, err := s.rangeClause(where, args)
	if err != nil {
		return 0, err
	}

	sql := fmt.Sprintf("SELECT COUNT(*) FROM history WHERE %s", strings.Join(where, " AND "))
	row := s.pool.QueryRow(ctx, sql, args...)
	var count int
	err = row.Scan(&count)
	return count, err
}

func (s *slice) sourceClause(clauses []string, args []any) ([]string, []any) {
	return append(clauses, fmt.Sprintf("source = $%d", len(args)+1)), append(args, s.source)
}

func (s *slice) rangeClause(clauses []string, args []any) ([]string, []any, error) {
	switch {
	case s.from.ID != "":
		id, err := idToSql(s.from.ID)
		if err != nil {
			return nil, nil, err
		}
		clauses = append(clauses, fmt.Sprintf("id >= $%d", len(args)+1))
		args = append(args, id)
	case !s.from.CreateTime.IsZero():
		clauses = append(clauses, fmt.Sprintf("create_time >= $%d", len(args)+1))
		args = append(args, s.from.CreateTime)
	}
	switch {
	case s.to.ID != "":
		id, err := idToSql(s.to.ID)
		if err != nil {
			return nil, nil, err
		}
		clauses = append(clauses, fmt.Sprintf("id < $%d", len(args)+1))
		args = append(args, id)
	case !s.to.CreateTime.IsZero():
		clauses = append(clauses, fmt.Sprintf("create_time < $%d", len(args)+1))
		args = append(args, s.to.CreateTime)
	}
	return clauses, args, nil
}

func idToSql(id string) (int64, error) {
	return strconv.ParseInt(id, 16, 64)
}

func idFromSql(id int64) string {
	return strconv.FormatInt(id, 16)
}
