package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
)

func TestDB_records(t *testing.T) {
	db := newDBTester(t)
	want := []Record{
		db.InsertRecord("dev1", "toner", "main1", "aux1"),
		db.InsertRecord("dev1", "toner", "main2", "aux2"),
		db.InsertRecord("dev2", "paper", "main1", "aux1"),
	}
	db.AssertRecords(want...)
}

func TestDB_Read(t *testing.T) {
	db := newDBTester(t)
	all := []Record{
		db.InsertRecord("dev1", "toner", "main1", "aux1"),
		db.InsertRecord("dev2", "paper", "main1", "aux1"),
		db.InsertRecord("dev1", "toner", "main2", "aux2"),
		db.InsertRecord("dev2", "paper", "main2", "aux2"),
		db.InsertRecord("dev1", "paper", "main3", "aux3"),
	}

	db.Run("all", func(t *dbTester) {
		t.AssertRecords(all...)
	})
	db.Run("page asc", func(t *dbTester) {
		page := make([]Record, 2)
		n, err := db.Read(t.Context(), CheckID{}, 0, 0, false, page)
		t.AssertReadResponse(n, err, all[:2], page)
	})
	db.Run("page desc", func(t *dbTester) {
		page := make([]Record, 2)
		n, err := db.Read(t.Context(), CheckID{}, 0, 0, true, page)
		want := []Record{all[4], all[3]}
		t.AssertReadResponse(n, err, want, page)
	})
	db.Run("filter name+id", func(t *dbTester) {
		got := make([]Record, len(all))
		want := []Record{all[0], all[2]}
		n, err := db.Read(t.Context(), CheckID{"dev1", "toner"}, 0, 0, false, got)
		t.AssertReadResponse(n, err, want, got)
	})
	db.Run("filter name", func(t *dbTester) {
		got := make([]Record, len(all))
		want := []Record{all[0], all[2], all[4]}
		n, err := db.Read(t.Context(), CheckID{"dev1", ""}, 0, 0, false, got)
		t.AssertReadResponse(n, err, want, got)
	})
	db.Run("filter from+to", func(t *dbTester) {
		got := make([]Record, len(all))
		want := []Record{all[1], all[2], all[3]}
		from := all[1].ID
		to := all[4].ID
		n, err := db.Read(t.Context(), CheckID{}, from, to, false, got)
		t.AssertReadResponse(n, err, want, got)
	})
	db.Run("filter from+to desc", func(t *dbTester) {
		got := make([]Record, len(all))
		want := []Record{all[3], all[2], all[1]}
		from := all[1].ID
		to := all[4].ID
		n, err := db.Read(t.Context(), CheckID{}, from, to, true, got)
		t.AssertReadResponse(n, err, want, got)
	})
	db.Run("filter from", func(t *dbTester) {
		got := make([]Record, len(all))
		want := []Record{all[2], all[3], all[4]}
		from := all[2].ID
		n, err := db.Read(t.Context(), CheckID{}, from, 0, false, got)
		t.AssertReadResponse(n, err, want, got)
	})
	db.Run("filter to", func(t *dbTester) {
		got := make([]Record, len(all))
		want := []Record{all[0], all[1], all[2]}
		to := all[3].ID
		n, err := db.Read(t.Context(), CheckID{}, 0, to, false, got)
		t.AssertReadResponse(n, err, want, got)
	})
}

type dbTester struct {
	*testing.T
	*DB
	epoch   time.Time
	timeInc time.Duration
}

func newDBTester(t *testing.T) *dbTester {
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
	return &dbTester{T: t, DB: db, epoch: time.Date(2021, 9, 4, 12, 0, 0, 0, time.UTC)}
}

func (t *dbTester) Run(name string, f func(t *dbTester)) {
	t.Helper()
	parent := t
	t.T.Run(name, func(t *testing.T) {
		f(&dbTester{T: t, DB: parent.DB, epoch: parent.epoch, timeInc: parent.timeInc})
	})
}

func (t *dbTester) InsertRecord(name, id, main, aux string) Record {
	t.Helper()
	rec := Record{
		Name:       name,
		CheckID:    id,
		CreateTime: t.tick(),
		Main:       []byte(main),
		Aux:        []byte(aux),
	}
	got, err := t.Insert(t.Context(), rec)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if err := t.testReadRecord(rec, got); err != nil {
		t.Fatalf("Insert returned invalid record: %v", err)
	}
	return got
}

func (t *dbTester) AssertRecords(want ...Record) {
	t.Helper()

	// space for more than we expect
	buf := make([]Record, len(want)+1)

	n, err := t.Read(t.Context(), CheckID{}, 0, 0, false, buf)
	t.AssertReadResponse(n, err, want, buf)
}

func (t *dbTester) AssertReadResponse(n int, err error, want, got []Record) {
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != len(want) {
		t.Fatalf("Read returned %d records, want %d", n, len(want))
	}
	var i int
	for i = 0; i < n; i++ {
		if err := t.testReadRecord(want[i], got[i]); err != nil {
			t.Errorf("Read record[%v]: %v", i, err)
		}
	}
	// remaining got records should be zero
	for ; i < len(got); i++ {
		if !got[i].IsZero() {
			t.Errorf("Read record[%v]: got %+v, want zero", i, got[i])
		}
	}
}

func (t *dbTester) AssertCount(want int) {
	t.Helper()
	n, err := t.Count(t.Context(), CheckID{}, 0, 0)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if n != want {
		t.Fatalf("Count returned %d, want %d", n, want)
	}
}

func (t *dbTester) testReadRecord(want, got Record) error {
	t.Helper()
	var err error
	if got.ID == 0 {
		err = errors.Join(err, fmt.Errorf("missing ID"))
	}
	if want.ID == 0 {
		got.ID = 0 // ignore ID in comparison
	}
	if diff := got.CreateTime.Sub(want.CreateTime).Abs(); diff > time.Millisecond {
		err = errors.Join(err, fmt.Errorf("create time diff %v > 1ms", diff))
	}
	got.CreateTime = want.CreateTime // avoid time precision loss during comparison
	if diff := cmp.Diff(want, got); diff != "" {
		err = errors.Join(err, fmt.Errorf("record mismatch (-want +got):\n%s", diff))
	}
	return err
}

func (t *dbTester) tick() time.Time {
	now := t.epoch.Add(t.timeInc)
	t.timeInc += 10 * time.Minute
	return now
}
