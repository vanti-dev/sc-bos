package pubcache

import (
	"context"
	"os"
	"testing"

	"go.etcd.io/bbolt"
)

func TestBoltStorage(t *testing.T) {
	ctx := context.Background()
	db, closeAndDelete := createTempDB()
	defer closeAndDelete()

	storage := NewBoltStorage(db, []byte("pubcache"))

	t.Run("publicationRoundTrip", func(t *testing.T) {
		_ = publicationRoundTrip(ctx, t, storage)
	})
	t.Run("rejectsNilPublication", func(t *testing.T) {
		_ = rejectsNilPublication(ctx, t, storage)
	})
	t.Run("rejectsEmptyPublicationID", func(t *testing.T) {
		_ = rejectsEmptyPublicationID(ctx, t, storage)
	})
}

func createTempDB() (db *bbolt.DB, closeAndDelete func()) {
	f, err := os.CreateTemp(os.TempDir(), "*.bolt")
	if err != nil {
		panic(err)
	}
	path := f.Name()
	// we don't need the file handle, close it now
	if err := f.Close(); err != nil {
		_ = os.Remove(path)
		panic(err)
	}

	// open a DB
	db, err = bbolt.Open(path, 0600, nil)
	if err != nil {
		_ = os.Remove(path)
		panic(err)
	}

	closeAndDelete = func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}

		err = os.Remove(path)
		if err != nil {
			panic(err)
		}
	}
	return
}
