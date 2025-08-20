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

func (d *Database) Insert(ctx context.Context, record Record) (Record, error) {
	toInsert := []Record{record}
	err := d.InsertBulk(ctx, toInsert)
	if err != nil {
		return Record{}, err
	}
	return toInsert[0], nil
}

// InsertBulk inserts multiple records into the database in a single transaction.
// Mutates records to set the ID field for each inserted record.
func (d *Database) InsertBulk(ctx context.Context, records []Record) error {
	err := d.db.WriteTx(ctx, func(tx *sql.Tx) (err error) {
		sources, err := newSourceAllocator(ctx, tx)
		if err != nil {
			return err
		}
		idAlloc, err := newRecordIDAllocator(ctx, tx)
		if err != nil {
			return err
		}

		stmt, err := tx.PrepareContext(ctx, "INSERT INTO history (id, source_id, payload) VALUES (?, ?, ?);")
		if err != nil {
			return err
		}
		defer func() {
			err = errors.Join(err, stmt.Close())
		}()

		for i, record := range records {
			srcID, err := sources.getOrInsertSource(ctx, record.Source)
			if err != nil {
				return err
			}
			recordID, err := idAlloc.allocateRecordID(ctx, record.CreateTime)
			if err != nil {
				return err
			}

			_, err = stmt.ExecContext(ctx, recordID, srcID, record.Payload)
			if err != nil {
				return err
			}
			records[i].ID = recordID
			records[i].CreateTime = recordID.Timestamp() // truncated and without time zone
		}
		return stmt.Close()
	})
	return err
}

// Read reads records from the database into the provided slice.
// If source is non-empty, it filters records by source.
// If from and to are non-zero, it filters records by ID range, greater-or-equal-to from, and less-than to.
// Returns number of records read and any error encountered.
// Size of the slice limits the number of records to read.
func (d *Database) Read(ctx context.Context, source string, from, to RecordID, desc bool, into []Record) (n int, err error) {
	err = d.read(ctx, source, from, to, desc, func(record Record) {
		into[n] = record
		n++
	})
	return n, err
}

func (d *Database) read(ctx context.Context, source string, from, to RecordID, desc bool, cb func(Record)) error {
	err := d.db.ReadTx(ctx, func(tx *sql.Tx) error {
		filters, args := buildFilters(source, from, to)

		order := "ASC"
		if desc {
			order = "DESC"
		}

		query := fmt.Sprintf(`
			SELECT history.id, history_sources.source, payload
            FROM history
            INNER JOIN history_sources ON history.source_id = history_sources.id
			WHERE %s
            ORDER BY history.id %s;
            `, filters, order)

		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer func() {
			_ = rows.Close()
		}()

		for rows.Next() {
			var (
				id      RecordID
				src     string
				payload []byte
			)
			err = rows.Scan(&id, &src, &payload)
			if err != nil {
				return err
			}
			cb(Record{ID: id, CreateTime: id.Timestamp(), Source: src, Payload: payload})
		}
		return rows.Err()
	})
	return err
}

func (d *Database) Count(ctx context.Context, source string, from, to RecordID) (int, error) {
	var count int
	err := d.db.ReadTx(ctx, func(tx *sql.Tx) error {
		filters, args := buildFilters(source, from, to)

		query := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM history
			INNER JOIN history_sources ON history.source_id = history_sources.id
			WHERE %s;
		`, filters)

		err := tx.QueryRowContext(ctx, query, args...).Scan(&count)
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

type Record struct {
	ID         RecordID
	Source     string
	CreateTime time.Time
	Payload    []byte
}

type Store struct {
	database *Database
	source   string

	from, to history.Record
}

func (s *Store) Append(ctx context.Context, payload []byte) (history.Record, error) {
	now := time.Now()
	record, err := s.database.Insert(ctx, Record{
		Source:     s.source,
		CreateTime: now,
		Payload:    payload,
	})
	if err != nil {
		return history.Record{}, err
	}
	return history.Record{
		ID:         record.ID.String(),
		CreateTime: record.CreateTime,
		Payload:    record.Payload,
	}, nil
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
	fromBound, err := calcBound(s.from)
	if err != nil {
		return 0, err
	}
	toBound, err := calcBound(s.to)
	if err != nil {
		return 0, err
	}

	var n int
	err = s.database.read(ctx, s.source, fromBound, toBound, desc, func(record Record) {
		into[n] = history.Record{
			ID:         record.ID.String(),
			CreateTime: record.CreateTime,
			Payload:    record.Payload,
		}
		n++
	})
	return n, err
}

func (s *Store) Len(ctx context.Context) (int, error) {
	fromBound, err := calcBound(s.from)
	if err != nil {
		return 0, err
	}
	toBound, err := calcBound(s.to)
	if err != nil {
		return 0, err
	}

	return s.database.Count(ctx, s.source, fromBound, toBound)
}

func calcBound(limit history.Record) (RecordID, error) {
	if limit.ID != "" {
		id, err := ParseRecordID(limit.ID)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	if !limit.CreateTime.IsZero() {
		id := MakeRecordID(limit.CreateTime, 0)
		return id, nil
	}

	// zero value is sentinel meaning "no bound"
	return 0, nil
}

// builds an SQL term for filtering records based on source and ID range.
// Also returns a slice of parameters to be passed when executing the query.
// If no filtering is to be performed, returns a dummy condition to maintain valid SQL syntax.
func buildFilters(source string, from, to RecordID) (string, []any) {
	var filters []string
	var args []any
	if source != "" {
		filters = []string{"history_sources.source = ?"}
		args = []any{source}
	}

	if from != 0 {
		filters = append(filters, "history.id >= ?")
		args = append(args, from)
	}

	if to != 0 {
		filters = append(filters, "history.id < ?")
		args = append(args, to)
	}

	if len(filters) > 0 {
		return strings.Join(filters, " AND "), args
	} else {
		return "1 = 1", nil // no filtering, return a dummy condition
	}
}

type RecordID int64

func MakeRecordID(ts time.Time, serial int) RecordID {
	if serial < 0 || serial >= 1_000_000 {
		panic("serial must be in range [0, 999999]")
	}
	return RecordID(ts.UnixMilli()*1_000_000 + int64(serial))
}

func ParseRecordID(s string) (RecordID, error) {
	id, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, ErrInvalidRecordID
	}
	if id < 0 {
		return 0, ErrInvalidRecordID
	}
	return RecordID(id), nil
}

func (id RecordID) Timestamp() time.Time {
	return time.UnixMilli(int64(id) / 1_000_000)
}

func (id RecordID) Serial() int {
	return int(int64(id) % 1_000_000)
}

func (id RecordID) Next() (RecordID, bool) {
	if id.Serial() == 1_000_000-1 {
		// incrementing the serial number would overflow the millisecond timestamp
		return 0, false
	}
	return RecordID(int64(id) + 1), true
}

func (id RecordID) String() string {
	return fmt.Sprintf("%016X", int64(id))
}

type sourceAllocator struct {
	sources    map[string]int64 // source name to ID mapping
	insertStmt *sql.Stmt        // prepared statement for inserting new sources if they don't exist
	queryStmt  *sql.Stmt        // prepared statement for querying source IDs
}

func newSourceAllocator(ctx context.Context, tx *sql.Tx) (*sourceAllocator, error) {
	sc := &sourceAllocator{
		sources: make(map[string]int64),
	}

	var err error
	sc.insertStmt, err = tx.PrepareContext(ctx, "INSERT OR IGNORE INTO history_sources (source) VALUES (?)")
	if err != nil {
		return nil, err
	}

	sc.queryStmt, err = tx.PrepareContext(ctx, "SELECT id FROM history_sources WHERE source = ?")
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (sa *sourceAllocator) getOrInsertSource(ctx context.Context, source string) (int64, error) {
	if id, ok := sa.sources[source]; ok {
		return id, nil
	}

	// Insert the source if it doesn't exist
	_, err := sa.insertStmt.ExecContext(ctx, source)
	if err != nil {
		return 0, err
	}

	// Query the ID of the source
	var id int64
	err = sa.queryStmt.QueryRowContext(ctx, source).Scan(&id)
	if err != nil {
		return 0, err
	}

	sa.sources[source] = id
	return id, nil
}

type recordIDAllocator struct {
	maxRecordIDByTS map[int64]RecordID // map from millisecond timestamp to the max record ID within that ts
	stmt            *sql.Stmt          // prepared statement for querying the max record ID for a given millisecond timestamp
}

func newRecordIDAllocator(ctx context.Context, tx *sql.Tx) (*recordIDAllocator, error) {
	ra := &recordIDAllocator{
		maxRecordIDByTS: make(map[int64]RecordID),
	}

	var err error
	ra.stmt, err = tx.PrepareContext(ctx, "SELECT MAX(id) FROM history WHERE id >= (1000000 * ?1) AND id < (1000000 * (?1 + 1))")
	if err != nil {
		return nil, err
	}

	return ra, nil
}

func (ra *recordIDAllocator) allocateRecordID(ctx context.Context, ts time.Time) (RecordID, error) {
	var (
		anyExist bool
		highest  RecordID
	)

	highest, ok := ra.maxRecordIDByTS[ts.UnixMilli()]
	if ok {
		anyExist = true
	} else {
		// Query the highest record ID for this millisecond timestamp
		var id sql.NullInt64
		err := ra.stmt.QueryRowContext(ctx, ts.UnixMilli()).Scan(&id)
		if err != nil {
			return 0, err
		}
		if id.Valid {
			highest = RecordID(id.Int64)
			anyExist = true
		}
	}

	var next RecordID
	if anyExist {
		next, ok = highest.Next()
		if !ok {
			return 0, ErrTooManyRecords
		}
	} else {
		next = MakeRecordID(ts, 0)
	}
	ra.maxRecordIDByTS[ts.UnixMilli()] = next
	return next, nil
}

var (
	ErrTooManyRecords  = errors.New("too many records")
	ErrInvalidRecordID = errors.New("invalid record ID format")
)
