package history

import (
	"context"
	"time"
)

type Store interface {
	Append(ctx context.Context, payload []byte) (Record, error)
	Slice
}

type Slice interface {
	Slice(from, to Record) Slice
	Read(ctx context.Context, into []Record) (int, error)
	Len(ctx context.Context) (int, error)
}

type Record struct {
	ID         string
	CreateTime time.Time
	Payload    []byte
}

func (r Record) IsZero() bool {
	return r.ID == "" && r.CreateTime.IsZero() && len(r.Payload) == 0
}
