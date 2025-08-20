package sqlitestore

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestDatabase_Insert(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	// Using fixed origin time for test determinism
	originTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	source := "test-source"

	// Insert multiple records sequentially
	records := []Record{
		{Source: source, Payload: []byte("test-payload-1"), CreateTime: originTime.Add(-2 * time.Hour)},
		{Source: source, Payload: []byte("test-payload-2"), CreateTime: originTime.Add(-1 * time.Hour)},
		{Source: source, Payload: []byte("test-payload-3"), CreateTime: originTime},
	}

	// Use zero-sized slice with capacity
	insertedRecords := make([]Record, 0, len(records))

	// Insert all records and collect the results
	for _, r := range records {
		record, err := db.Insert(ctx, r)
		if err != nil {
			t.Fatalf("unexpected error inserting record: %v", err)
		}

		if record.ID == 0 {
			t.Errorf("expected valid record ID, got record.ID=%q", record.ID)
		}

		if record.Payload == nil || string(record.Payload) != string(r.Payload) {
			t.Errorf("expected payload to match %q, got %q", r.Payload, record.Payload)
		}

		insertedRecords = append(insertedRecords, record)
	}

	// Verify all records were stored correctly
	verifyRecords(t, db, ctx, records)
}

// verifyRecords retrieves all records for a source from the database and verifies they match expected records.
func verifyRecords(t *testing.T, db *Database, ctx context.Context, expected []Record) {
	t.Helper()

	// Create a slice with capacity to hold all expected records
	retrievedRecords := make([]Record, len(expected)+1)

	// Fetch all records using zero-value Records to indicate full range
	count, err := db.Read(ctx, "", 0, 0, false, retrievedRecords)
	if err != nil {
		t.Fatalf("unexpected error reading records: %v", err)
	}

	if count != len(expected) {
		t.Errorf("expected %d records, got %d", len(expected), count)
	}

	// Verify each record has correct data
	for i := range count {
		if retrievedRecords[i].ID == 0 {
			t.Errorf("record %d: missing ID", i)
		}

		// Verify the payload matches what we inserted
		expectedPayload := string(expected[i].Payload)
		actualPayload := string(retrievedRecords[i].Payload)
		if actualPayload != expectedPayload {
			t.Errorf("record %d: expected payload %q, got %q", i, expectedPayload, actualPayload)
		}

		// Verify timestamps are close (within 1 second)
		timeDiff := retrievedRecords[i].CreateTime.Sub(expected[i].CreateTime).Abs()
		if timeDiff > time.Millisecond {
			t.Errorf("record %d: timestamp difference too large: %v", i, timeDiff)
		}
	}
}

func TestDatabase_InsertBulk(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	records := []Record{
		{Source: "source-1", CreateTime: time.Now(), Payload: []byte("payload-1")},
		{Source: "source-2", CreateTime: time.Now(), Payload: []byte("payload-2")},
	}

	err := db.InsertBulk(ctx, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, record := range records {
		if record.ID == 0 {
			t.Errorf("expected valid ID for record %d, got zero", i)
		}
	}
}

func TestDatabase_InsertBulk_DuplicateSources(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	records := []Record{
		{Source: "duplicate-source", CreateTime: time.Now(), Payload: []byte("payload-1")},
		{Source: "duplicate-source", CreateTime: time.Now(), Payload: []byte("payload-2")},
	}

	err := db.InsertBulk(ctx, records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

}

func TestDatabase_Read(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	source := "test-source"
	payload := []byte("test-payload")
	at := time.Now()

	_, err := db.Insert(ctx, Record{
		Source:     source,
		CreateTime: at,
		Payload:    payload,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Using zero-value Records to read the entire dataset
	from := RecordID(0) // zero value = beginning of dataset
	to := RecordID(0)   // zero value = end of dataset
	into := make([]Record, 1)

	count, err := db.Read(ctx, source, from, to, false, into)
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

	_, err := db.Insert(ctx, Record{
		Source:     source,
		CreateTime: at,
		Payload:    payload,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Using zero-value Records to count the entire dataset
	count, err := db.Count(ctx, source, 0, 0)
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

func TestDatabase_LargeScale(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large scale test in short mode")
	}

	db := newTestDB(t)
	ctx := context.Background()

	const (
		totalRecords = 1_000_000
		numSources   = 50
		batchSize    = 10_000 // insert in batches for efficiency
		payloadSize  = 100    // size of each payload in bytes
	)

	logger := zaptest.NewLogger(t)

	logger.Info("Starting large scale database test",
		zap.Int("total_records", totalRecords),
		zap.Int("num_sources", numSources),
		zap.Int("batch_size", batchSize))

	start := time.Now()
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// Insert records in batches
	var expected []Record
	for batch := 0; batch < totalRecords/batchSize; batch++ {
		records := make([]Record, batchSize)

		for i := 0; i < batchSize; i++ {
			recordNum := batch*batchSize + i
			sourceID := recordNum % numSources
			record := Record{
				Source:     fmt.Sprintf("source-%d", sourceID),
				CreateTime: baseTime.Add(time.Duration(recordNum) * time.Millisecond),
				Payload:    generateRandomPayload(200),
			}
			records[i] = record
			expected = append(expected, record)
		}

		err := db.InsertBulk(ctx, records)
		if err != nil {
			t.Fatalf("failed to insert batch %d: %v", batch, err)
		}

		// Log progress every 20 batches
		if batch%20 == 0 {
			elapsed := time.Since(start)
			recordsInserted := (batch + 1) * batchSize
			rate := float64(recordsInserted) / elapsed.Seconds()
			logger.Info("Insertion progress",
				zap.Int("records_inserted", recordsInserted),
				zap.Duration("elapsed", elapsed),
				zap.Float64("records_per_second", rate))
		}
	}

	insertDuration := time.Since(start)
	logger.Info("Insertion completed",
		zap.Duration("total_time", insertDuration),
		zap.Float64("avg_records_per_second", float64(totalRecords)/insertDuration.Seconds()))

	// Verify record count
	count, err := db.Count(ctx, "", 0, 0)
	if err != nil {
		t.Fatalf("failed to count records: %v", err)
	}

	if count != totalRecords {
		t.Errorf("expected %d records, got %d", totalRecords, count)
	}

	// Measure database size
	dbSize, err := db.Size(ctx)
	if err != nil {
		t.Fatalf("failed to get database size: %v", err)
	}

	avgSizePerRecord := float64(dbSize) / float64(totalRecords)
	overheadPerRecord := avgSizePerRecord - float64(payloadSize)

	logger.Info("Database size analysis",
		zap.Int64("total_size_bytes", dbSize),
		zap.Float64("total_size_mb", float64(dbSize)/(1024*1024)),
		zap.Float64("avg_bytes_per_record", avgSizePerRecord),
		zap.Float64("overhead_bytes_per_record", overheadPerRecord),
		zap.Int("total_records", count))

	// Verify we can still read records efficiently
	readStart := time.Now()
	into := make([]Record, 1000)
	readCount, err := db.Read(ctx, "source-0", 0, 0, false, into)
	if err != nil {
		t.Fatalf("failed to read records: %v", err)
	}
	readDuration := time.Since(readStart)

	logger.Info("Read performance test",
		zap.Int("records_read", readCount),
		zap.Duration("read_time", readDuration))

	verifyRecords(t, db, ctx, expected)
}

func generateRandomPayload(size int) []byte {
	payload := make([]byte, size)
	_, err := rand.Read(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random payload: %v", err))
	}
	return payload
}

func newTestDB(t testing.TB) *Database {
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
	t.Logf("created test database %s", dbPath)
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
