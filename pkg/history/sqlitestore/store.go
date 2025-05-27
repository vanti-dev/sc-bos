package sqlitestore

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/internal/sqlite"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

const appID = 0x5C0502

//go:embed migrations/*.sql
var migrationFS embed.FS

var schema = sqlite.MustLoadVersionedSchema(migrationFS, "migrations")

type Database struct {
	db       *sqlite.Database
	maxCount int64
	maxAge   time.Duration
}

func Open(ctx context.Context, path string, options ...Option) (*Database, error) {
	o := &opts{}
	for _, option := range options {
		option(o)
	}
	if o.logger == nil {
		o.logger = zap.NewNop()
	}

	db, err := sqlite.Open(ctx, path,
		sqlite.WithApplicationID(appID),
		sqlite.WithLogger(o.logger),
	)
	if err != nil {
		return nil, err
	}

	err = db.Migrate(ctx, schema)
	if err != nil {
		return nil, errors.Join(err, db.Close())
	}

	return &Database{db: db}, nil
}

func (d *Database) OpenStore(source string) *Store {
	return &Store{
		database: d,
		source:   source,
	}
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Insert(ctx context.Context, source string, at time.Time, payload []byte) (history.Record, error) {
	var id int64
	err := d.db.WriteTx(ctx, func(tx *sql.Tx) error {
		epoch, err := readEpoch(ctx, tx)
		if err != nil {
			return err
		}

		srcID, err := sourceID(ctx, tx, source, true)
		if err != nil {
			return err
		}

		offset := at.UTC().Sub(epoch).Milliseconds()

		err = tx.QueryRowContext(ctx, "INSERT INTO history (source_id, epoch_offset_ms, payload) VALUES (?, ?, ?) RETURNING id",
			srcID, offset, payload).Scan(&id)
		return err
	})
	return history.Record{
		ID:         strconv.FormatInt(id, 10),
		CreateTime: at,
		Payload:    payload,
	}, err
}

type InsertRecord struct {
	Source     string
	CreateTime time.Time
	Payload    []byte
}

// InsertBulk inserts multiple records into the database in a single transaction.
// Returns a slice of inserted record IDs in the same order as the input records.
func (d *Database) InsertBulk(ctx context.Context, records []InsertRecord) (ids []string, err error) {
	err = d.db.WriteTx(ctx, func(tx *sql.Tx) error {
		epoch, err := readEpoch(ctx, tx)
		if err != nil {
			return err
		}

		srcIDCache := make(map[string]int64)
		getSourceID := func(source string) (int64, error) {
			if id, ok := srcIDCache[source]; ok {
				return id, nil
			}
			srcID, err := sourceID(ctx, tx, source, true)
			if err != nil {
				return 0, err
			}
			srcIDCache[source] = srcID
			return srcID, nil
		}

		stmt, err := tx.PrepareContext(ctx, "INSERT INTO history (source_id, epoch_offset_ms, payload) VALUES (?, ?, ?) RETURNING id")
		if err != nil {
			return err
		}
		for _, record := range records {
			var id int64
			offset := record.CreateTime.UTC().Sub(epoch).Milliseconds()
			srcID, err := getSourceID(record.Source)
			if err != nil {
				return err
			}
			err = stmt.QueryRowContext(ctx, srcID, offset, record.Payload).Scan(&id)
			if err != nil {
				return errors.Join(err, stmt.Close())
			}
			ids = append(ids, strconv.FormatInt(id, 10))
		}
		return stmt.Close()
	})
	return ids, err
}

func (d *Database) Read(ctx context.Context, source string, from, to history.Record, into []history.Record, desc bool) (int, error) {
	var count int
	err := d.db.ReadTx(ctx, func(tx *sql.Tx) error {
		epoch, err := readEpoch(ctx, tx)
		if err != nil {
			return err
		}

		filters, args := buildFilters(source, epoch, from, to)

		order := "ASC"
		if desc {
			order = "DESC"
		}

		query := fmt.Sprintf(`
			SELECT history.id, epoch_offset_ms, payload
            FROM history
            INNER JOIN history_sources ON history.source_id = history_sources.id
			WHERE %s
            ORDER BY history.id %s
			LIMIT ?;
        `, filters, order)
		args = append(args, len(into))

		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			var (
				id       int64
				offsetMS int64
				payload  []byte
			)
			err = rows.Scan(&id, &offsetMS, &payload)
			if err != nil {
				return err
			}
			into[count] = buildRecord(id, epoch, offsetMS, payload)
			count++
		}
		return rows.Err()
	})
	return count, err
}

func (d *Database) Count(ctx context.Context, source string, from, to history.Record) (int, error) {
	var count int
	err := d.db.ReadTx(ctx, func(tx *sql.Tx) error {
		epoch, err := readEpoch(ctx, tx)
		if err != nil {
			return err
		}

		filters, args := buildFilters(source, epoch, from, to)

		query := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM history
			INNER JOIN history_sources ON history.source_id = history_sources.id
			WHERE %s;
		`, filters)

		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
		return err
	})
	return count, err
}

func (d *Database) Size(ctx context.Context) (int64, error) {
	var size int64
	err := d.db.ReadTx(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, "SELECT page_count * page_size FROM pragma_page_count(), pragma_page_size()").Scan(&size)
	})
	return size, err
}

type Store struct {
	database *Database
	source   string

	from, to history.Record
}

func (s *Store) Append(ctx context.Context, payload []byte) (history.Record, error) {
	now := time.Now()
	record, err := s.database.Insert(ctx, s.source, now, payload)
	if err != nil {
		return history.Record{}, err
	}
	return record, nil
}

func (s *Store) Slice(from, to history.Record) history.Slice {
	return &Store{
		database: s.database,
		source:   s.source,
		from:     from,
		to:       to,
	}
}

func (s *Store) Read(ctx context.Context, into []history.Record) (int, error) {
	return s.read(ctx, into, false)
}

func (s *Store) ReadDesc(ctx context.Context, into []history.Record) (int, error) {
	return s.read(ctx, into, true)
}

func (s *Store) read(ctx context.Context, into []history.Record, desc bool) (int, error) {
	return s.database.Read(ctx, s.source, s.from, s.to, into, desc)
}

func (s *Store) Len(ctx context.Context) (int, error) {
	return s.database.Count(ctx, s.source, s.from, s.to)
}

func readEpoch(ctx context.Context, tx *sql.Tx) (time.Time, error) {
	var epochStr string
	err := tx.QueryRowContext(ctx, "SELECT value FROM history_meta WHERE key = 'epoch'").Scan(&epochStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(sqlite.DateTimeFormat, epochStr)
}

func sourceID(ctx context.Context, tx *sql.Tx, source string, create bool) (int64, error) {
	if create {
		_, err := tx.ExecContext(ctx, "INSERT OR IGNORE INTO history_sources (source) VALUES (?)", source)
		if err != nil {
			return 0, err
		}
	}

	var id int64
	err := tx.QueryRowContext(ctx, "SELECT id FROM history_sources WHERE source = ?", source).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func buildFilters(source string, epoch time.Time, from, to history.Record) (string, []any) {
	filters := []string{"history_sources.source = ?"}
	args := []any{source}

	if !from.IsZero() {
		if from.ID != "" {
			filters = append(filters, "history.id >= ?")
			args = append(args, from.ID)
		} else if !from.CreateTime.IsZero() {
			filters = append(filters, "history.epoch_offset_ms >= ?")
			args = append(args, toEpochOffset(epoch, from.CreateTime))
		}
	}

	if !to.IsZero() {
		if to.ID != "" {
			filters = append(filters, "history.id < ?")
			args = append(args, to.ID)
		} else if !to.CreateTime.IsZero() {
			filters = append(filters, "history.epoch_offset_ms < ?")
			args = append(args, toEpochOffset(epoch, to.CreateTime))
		}
	}

	return strings.Join(filters, " AND "), args
}

func buildRecord(id int64, epoch time.Time, offsetMS int64, payload []byte) history.Record {
	return history.Record{
		ID:         strconv.FormatInt(id, 10),
		CreateTime: epoch.Add(time.Duration(offsetMS) * time.Millisecond),
		Payload:    payload,
	}
}

func toEpochOffset(epoch time.Time, at time.Time) int64 {
	return at.UTC().Sub(epoch).Milliseconds()
}

func fromEpochOffset(epoch time.Time, offsetMS int64) time.Time {
	return epoch.Add(time.Duration(offsetMS) * time.Millisecond)
}
