package database

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"

	"golang.org/x/exp/slices"
)

// Schema represents a database schema with (optionally) multiple versions.
//
// A Schema version contains migrations that move the database from one version to the next.
// A database starts on version 0 when empty, and each Migration is applied in order to reach the latest version.
//
// A single-version Schema has one Migration, from version 0 to version 1.
type Schema struct {
	// migrations are:
	// - sorted by version number
	// - contain no duplicate version numbers
	// - contiguous, with no gaps between version numbers
	// - start at version 1
	migrations []Migration
}

// Migration represents a single database schema Migration step from Version-1 to Version.
type Migration struct {
	Version uint32
	SQL     string
}

// LoadVersionedSchema loads a multi-version schema from the filesystem.
//
// Searches for files in the root of the given filesystem, where each file is named "<version>.sql",
// or "<version>_<description>.sql", and contains the SQL to migrate from version <version>-1 to <version>.
//
// Subdirectories and other files are ignored.
func LoadVersionedSchema(source fs.FS) (Schema, error) {
	dirEnts, err := fs.ReadDir(source, ".")
	if err != nil {
		return Schema{}, err
	}

	var migrations []Migration
	for _, dirEnt := range dirEnts {
		if dirEnt.Type() != 0 {
			// not a regular file
			continue
		}

		version, ok := parseMigrationFilename(dirEnt.Name())
		if !ok {
			continue
		}

		contents, err := fs.ReadFile(source, dirEnt.Name())
		if err != nil {
			return Schema{}, fmt.Errorf("read migration %q: %w", dirEnt.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			SQL:     string(contents),
		})
	}

	return NewVersionedSchema(migrations)
}

// MustLoadVersionedSchema is like LoadVersionedSchema but panics on error.
func MustLoadVersionedSchema(source fs.FS) Schema {
	schema, err := LoadVersionedSchema(source)
	if err != nil {
		panic(err)
	}
	return schema
}

// NewSchema returns a single-Migration schema with the given SQL.
//
// The returned schema has one Migration, from version 0 to version 1, that executes the given SQL.
func NewSchema(initSQL string) Schema {
	return Schema{
		migrations: []Migration{{Version: 1, SQL: initSQL}},
	}
}

func NewVersionedSchema(migrations []Migration) (Schema, error) {
	migrations = slices.Clone(migrations)
	slices.SortFunc(migrations, func(a, b Migration) int {
		return cmp.Compare(a.Version, b.Version)
	})
	s := Schema{migrations: migrations}
	err := s.validate()
	if err != nil {
		return Schema{}, err
	}
	return s, nil
}

func (s *Schema) validate() error {
	if len(s.migrations) == 0 {
		return ErrNoMigrations
	}
	lastVersion := uint32(0)
	for _, m := range s.migrations {
		if m.Version == 0 {
			return errors.New("invalid migration version number 0")
		} else if m.Version != lastVersion+1 {
			return fmt.Errorf("no migration from version %d to %d", lastVersion, lastVersion+1)
		}
	}
	return nil
}

func (s *Schema) find(version uint32) (Migration, bool) {
	i, ok := slices.BinarySearchFunc(s.migrations, version, func(m Migration, v uint32) int {
		return cmp.Compare(m.Version, v)
	})
	if !ok {
		return Migration{}, false
	}
	return s.migrations[i], true
}

var ErrNoMigrations = errors.New("no migrations found")

var migrationFilenameRegexp = regexp.MustCompile(`^(\d+)(_.*)?\.sql$`)

func parseMigrationFilename(filename string) (version uint32, ok bool) {
	m := migrationFilenameRegexp.FindStringSubmatch(filename)
	if m == nil {
		return 0, false
	}
	v64, err := strconv.ParseUint(m[1], 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(v64), true
}

func getApplicationID(ctx context.Context, tx *sql.Tx) (ApplicationID, error) {
	row := tx.QueryRowContext(ctx, "PRAGMA application_id;")
	var id ApplicationID
	err := row.Scan(&id)
	return id, err
}

func setApplicationID(ctx context.Context, tx *sql.Tx, id ApplicationID) error {
	// PRAGMAs don't support parameterised queries, have to put the value directly in the query string
	_, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA application_id = %d;", id))
	return err
}

func getUserVersion(ctx context.Context, tx *sql.Tx) (uint32, error) {
	row := tx.QueryRowContext(ctx, "PRAGMA user_version;")
	var version uint32
	err := row.Scan(&version)
	return version, err
}

func setUserVersion(ctx context.Context, tx *sql.Tx, version uint32) error {
	// PRAGMAs don't support parameterised queries, have to put the value directly in the query string
	_, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d;", version))
	return err
}
