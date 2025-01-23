package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// tests that the database file is created in the correct location
// and that read and write transactions go to the same database file.
func TestOpen(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "db.sqlite3")
	db, err := Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	_, err = os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("database file not created")
	} else if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	testReadWrite(t, db)
}

func TestOpenMemory(t *testing.T) {
	db := OpenMemory()
	defer db.Close()

	testReadWrite(t, db)
}

func TestWithApplicationID(t *testing.T) {
	type testCase struct {
		initAppID   ApplicationID
		initVersion uint32
		withAppID   ApplicationID
		expectAppID ApplicationID
		expectErr   error
	}
	cases := map[string]testCase{
		"empty": {
			// for empty DB (with application_id = user_version = 0), we set the App ID
			withAppID:   123,
			expectAppID: 123,
		},
		"mismatch_version_0": {
			initAppID: 123,
			withAppID: 456,
			expectErr: ErrApplicationIDMismatch,
		},
		"mismatch_version_1": {
			initAppID:   123,
			initVersion: 1,
			withAppID:   456,
			expectErr:   ErrApplicationIDMismatch,
		},
		"zero_appid_nonzero_version": {
			initVersion: 1,
			withAppID:   456,
			expectErr:   ErrApplicationIDMismatch,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			init := fmt.Sprintf(`PRAGMA application_id = %d; PRAGMA user_version = %d;`, tc.initAppID, tc.initVersion)
			dbPath := freshDBWith(t, init)
			ctx := context.Background()

			db, err := Open(ctx, dbPath, WithApplicationID(tc.withAppID))
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("expected ErrApplicationIDMismatch, got %v", err)
			}
			if db == nil {
				return
			}
			defer func() {
				_ = db.Close()
			}()

			var appID ApplicationID
			err = db.ReadTx(ctx, func(tx *sql.Tx) error {
				row := tx.QueryRowContext(ctx, "PRAGMA application_id;")
				return row.Scan(&appID)
			})
			if err != nil {
				t.Fatalf("read app id: %v", err)
			}

			if appID != tc.expectAppID {
				t.Errorf("expected app id %d, got %d", tc.expectAppID, appID)
			}
		})
	}

}

func TestDatabase_Migrate(t *testing.T) {
	migrations := []Migration{
		{1, `CREATE TABLE test ( id INTEGER PRIMARY KEY );`},
		{2, `ALTER TABLE test ADD COLUMN name TEXT;`},
	}
	schema, err := NewVersionedSchema(migrations)
	if err != nil {
		t.Fatalf("NewVersionedSchema: %v", err)
	}

	ctx := context.Background()
	db, err := Open(ctx, filepath.Join(t.TempDir(), "db.sqlite3"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	err = db.Migrate(ctx, schema)
	if err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	err = db.WriteTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO test (id, name) VALUES (123, 'test');`)
		return err
	})
	if err != nil {
		t.Fatalf("WriteTx: %v", err)
	}

	// add some new migrations, to simulate loading the app again with an updated schema
	migrations = append(migrations, Migration{3, `ALTER TABLE test ADD COLUMN age INTEGER NOT NULL DEFAULT 30;`})
	schema, err = NewVersionedSchema(migrations)
	if err != nil {
		t.Fatalf("NewVersionedSchema: %v", err)
	}

	err = db.Migrate(ctx, schema)
	if err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var age int
	err = db.ReadTx(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `SELECT age FROM test WHERE id = 123;`)
		return row.Scan(&age)
	})
	if err != nil {
		t.Fatalf("ReadTx: %v", err)
	}
	if age != 30 {
		t.Errorf("expected age 30, got %d", age)
	}
}

func testReadWrite(t *testing.T, db *Database) {
	t.Helper()
	ctx := context.Background()
	err := db.WriteTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
CREATE TABLE test ( id INTEGER PRIMARY KEY );
INSERT INTO test (id) VALUES (123);
`)
		return err
	})
	if err != nil {
		t.Errorf("WriteTx: %v", err)
	}

	err = db.ReadTx(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, "SELECT id FROM test;")
		var id int
		err := row.Scan(&id)
		if err != nil {
			return err
		}
		if id != 123 {
			return errors.New("unexpected id")
		}
		return nil
	})
	if err != nil {
		t.Errorf("ReadTx: %v", err)
	}
}

func freshDBWith(t *testing.T, schema string) (path string) {
	t.Helper()
	ctx := context.Background()
	f, err := os.CreateTemp(t.TempDir(), "db-*.sqlite3")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	path = f.Name()
	_ = f.Close()

	db, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	err = db.WriteTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, schema)
		return err
	})
	if err != nil {
		t.Fatalf("WriteTx: %v", err)
	}

	return path
}
