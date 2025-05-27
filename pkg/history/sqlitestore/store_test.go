package sqlitestore

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

func TestDatabase_Insert(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	// Using fixed origin time for test determinism
	originTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	source := "test-source"

	// Insert multiple records sequentially
	records := []struct {
		payload []byte
		at      time.Time
	}{
		{[]byte("test-payload-1"), originTime.Add(-2 * time.Hour)},
		{[]byte("test-payload-2"), originTime.Add(-1 * time.Hour)},
		{[]byte("test-payload-3"), originTime},
	}

	// Use zero-sized slice with capacity
	insertedRecords := make([]history.Record, 0, len(records))

	// Insert all records and collect the results
	for _, r := range records {
		record, err := db.Insert(ctx, source, r.at, r.payload)
		if err != nil {
			t.Fatalf("unexpected error inserting record: %v", err)
		}

		if record.ID == "" {
			t.Errorf("expected valid record ID, got record.ID=%q", record.ID)
		}

		if record.Payload == nil || string(record.Payload) != string(r.payload) {
			t.Errorf("expected payload to match %q, got %q", r.payload, record.Payload)
		}

		insertedRecords = append(insertedRecords, record)
	}

	// Verify all records were stored correctly
	verifyRecords(t, db, ctx, source, records)
}

// verifyRecords retrieves all records for a source from the database and verifies they match expected records.
func verifyRecords(t *testing.T, db *Database, ctx context.Context, source string, expected []struct {
	payload []byte
	at      time.Time
}) {
	t.Helper()

	// Create a slice with capacity to hold all expected records
	retrievedRecords := make([]history.Record, len(expected))

	// Fetch all records using zero-value Records to indicate full range
	count, err := db.Read(ctx, source, history.Record{}, history.Record{}, retrievedRecords, false)
	if err != nil {
		t.Fatalf("unexpected error reading records: %v", err)
	}

	if count != len(expected) {
		t.Errorf("expected %d records, got %d", len(expected), count)
	}

	// Verify each record has correct data
	for i := 0; i < count; i++ {
		if retrievedRecords[i].ID == "" {
			t.Errorf("record %d: missing ID", i)
		}

		// Verify the payload matches what we inserted
		expectedPayload := string(expected[i].payload)
		actualPayload := string(retrievedRecords[i].Payload)
		if actualPayload != expectedPayload {
			t.Errorf("record %d: expected payload %q, got %q", i, expectedPayload, actualPayload)
		}

		// Verify timestamps are close (within 1 second)
		timeDiff := retrievedRecords[i].CreateTime.Sub(expected[i].at).Abs()
		if timeDiff > time.Second {
			t.Errorf("record %d: timestamp difference too large: %v", i, timeDiff)
		}
	}
}

func TestDatabase_InsertBulk(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	records := []InsertRecord{
		{Source: "source-1", CreateTime: time.Now(), Payload: []byte("payload-1")},
		{Source: "source-2", CreateTime: time.Now(), Payload: []byte("payload-2")},
	}

	ids, err := db.InsertBulk(ctx, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ids) != len(records) {
		t.Errorf("expected %d IDs, got %d", len(records), len(ids))
	}

	for i, id := range ids {
		if id == "" {
			t.Errorf("expected valid ID for record %d, got empty string", i)
		}
	}
}

func TestDatabase_InsertBulk_DuplicateSources(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	records := []InsertRecord{
		{Source: "duplicate-source", CreateTime: time.Now(), Payload: []byte("payload-1")},
		{Source: "duplicate-source", CreateTime: time.Now(), Payload: []byte("payload-2")},
	}

	ids, err := db.InsertBulk(ctx, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ids) != len(records) {
		t.Errorf("expected %d IDs, got %d", len(records), len(ids))
	}
}

func TestDatabase_Read(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	source := "test-source"
	payload := []byte("test-payload")
	at := time.Now()

	_, err := db.Insert(ctx, source, at, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Using zero-value Records to read the entire dataset
	from := history.Record{} // zero value = beginning of dataset
	to := history.Record{}   // zero value = end of dataset
	into := make([]history.Record, 1)

	count, err := db.Read(ctx, source, from, to, into, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
}

func TestDatabase_Count(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	source := "test-source"
	payload := []byte("test-payload")
	at := time.Now()

	_, err := db.Insert(ctx, source, at, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Using zero-value Records to count the entire dataset
	count, err := db.Count(ctx, source, history.Record{}, history.Record{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
}

func TestDatabase_Size(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	size, err := db.Size(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size <= 0 {
		t.Errorf("expected database size to be greater than 0, got %d", size)
	}
}

func newTestDB(t *testing.T) *Database {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	db, err := Open(ctx, dbPath, WithLogger(logger))
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close test database: %v", err)
		}
		stat, err := os.Stat(dbPath)
		if err != nil {
			t.Logf("failed to stat test database file: %v", err)
		} else {
			t.Logf("database file size: %d bytes", stat.Size())
		}
	})
	return db
}
