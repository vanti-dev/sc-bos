package database

import (
	"context"
	"database/sql"
	"errors"
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

	expectID := 123
	err = db.WriteTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
CREATE TABLE test ( id INTEGER PRIMARY KEY );
INSERT INTO test (id) VALUES (?);
`, expectID)
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
		if id != expectID {
			return errors.New("unexpected id")
		}
		return nil
	})
	if err != nil {
		t.Errorf("ReadTx: %v", err)
	}
}

func TestOpenMemory(t *testing.T) {
	db := OpenMemory()
	defer db.Close()

}
