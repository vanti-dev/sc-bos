package pubcache

import (
	"context"
	"os"
	"testing"
)

func Test_FileStorage(t *testing.T) {
	dir, err := os.MkdirTemp("", "bsp-ew_test_pubcache_TestFileStorage")
	if err != nil {
		t.Fatal(err)
	}

	storage := NewFileStorage(dir)
	ctx := context.Background()

	t.Run("rejectsNilPublication", func(t *testing.T) {
		_ = rejectsNilPublication(ctx, t, storage)
	})
	t.Run("rejectsEmptyPublicationID", func(t *testing.T) {
		_ = rejectsEmptyPublicationID(ctx, t, storage)
	})
	t.Run("publicationRoundTrip", func(t *testing.T) {
		_ = publicationRoundTrip(ctx, t, storage)
	})
}
