// Package history provides a store for historical records.
// This package contains the interfaces and general types, see various foostore packages for actual store implementations.
package history

import (
	"context"
	"time"
)

// Store defines an append-only collection of ordered records.
// The store can be sliced, like a go slice, to query for a range of records which can then be read.
type Store interface {
	// Append adds the given payload to the store, returning the Record as recorded.
	// The context can be used to abort the append operation if needed.
	Append(ctx context.Context, payload []byte) (Record, error)
	Slice
}

// Slice describes a read-only ordered segment of a Store.
type Slice interface {
	// Slice creates a new slice including records >= from and < to.
	// If from is zero then it is treated as the first record available,
	// if to is zero then it is the record immediately after the last record,
	// just like with go slices.
	// The record should have either ID and/or CreateTime set.
	Slice(from, to Record) Slice
	// Read reads records from the slice starting at the first record available, ending when into is full (according to len).
	// Read returns the number of records actually read.
	// The oldest record read will be placed into index 0, the second into index 1, and so on.
	Read(ctx context.Context, into []Record) (int, error)
	// ReadDesc reads records from the slice starting at the last record available, ending when into is full (according to len).
	// ReadDesc returns the number of records actually read.
	// The newest record read will be placed into index 0, the second into index 1, and so on.
	ReadDesc(ctx context.Context, into []Record) (int, error)
	// Len returns the number of records this slice represents.
	Len(ctx context.Context) (int, error)
}

// Record is a payload as it was at a point in time.
type Record struct {
	ID         string
	CreateTime time.Time
	Payload    []byte
}

// IsZero returns whether r is equivalent to Record{}, the zero record.
func (r Record) IsZero() bool {
	return r.ID == "" && r.CreateTime.IsZero() && len(r.Payload) == 0
}
