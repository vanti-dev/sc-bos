// Package history provides a store for historical records.
// This package contains the interfaces and general types, see various foostore packages for actual store implementations.
package history

import (
	"context"
	"strings"
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
	// Slice returns a slice from the records in this slice where the records are also >= from and < to.
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

// Compare compares this record r with b, returning -1 if r is before b, 1 if r is after b, and 0 if they are equal.
// A zero CreateTime is considered before any non-zero CreateTime, an empty ID is considered before any non-empty ID.
func (r Record) Compare(b Record) int {
	// time.Time{} is before any other time.Time
	i := r.CreateTime.Compare(b.CreateTime)
	if i != 0 {
		return i
	}
	// "" already compares before other strings
	return strings.Compare(r.ID, b.ID)
}

// CompareZeroAfter compares this record r with b, returning -1 if r is before b, 1 if r is after b, and 0 if they are equal.
// A zero CreateTime is considered after any non-zero CreateTime, an empty ID is considered after any non-empty ID.
func (r Record) CompareZeroAfter(b Record) int {
	if b.CreateTime.IsZero() && !r.CreateTime.IsZero() {
		return -1
	}
	if r.CreateTime.IsZero() && !b.CreateTime.IsZero() {
		return 1
	}
	i := r.CreateTime.Compare(b.CreateTime)
	if i != 0 {
		return i
	}

	if b.ID == "" && r.ID != "" {
		return -1
	}
	if r.ID == "" && b.ID != "" {
		return 1
	}
	return strings.Compare(r.ID, b.ID)
}

// IntersectRecords returns the smallest [from, to) range that is common to both input ranges.
// If there is no such range, then from will compare after to.
// If the range is exactly one record, then from will be equal to to.
func IntersectRecords(from1, to1, from2, to2 Record) (from, to Record) {
	if from1.Compare(from2) > 0 {
		from = from1
	} else {
		from = from2
	}
	if to1.CompareZeroAfter(to2) < 0 {
		to = to1
	} else {
		to = to2
	}
	return from, to
}
