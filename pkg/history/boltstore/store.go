// Package boltstore provides an implementation of history.Store with records stored in a bolthold database.
package boltstore

import (
	"context"
	"strconv"
	"time"

	"github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

type Store struct {
	slice // sorted by id, which is createTime+dedupe index
	now   func() time.Time
}

func NewFromDb(db *bolthold.Store, source string) (history.Store, error) {
	var bucket *bolt.Bucket
	err := db.Bolt().Update(func(tx *bolt.Tx) error {
		var err error
		bucket, err = tx.CreateBucketIfNotExists([]byte(source))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Store{
		slice: slice{
			db:     db,
			bucket: bucket,
		},
	}, nil
}

func (s *Store) Append(ctx context.Context, payload []byte) (history.Record, error) {
	now := s.now()
	r := history.Record{
		ID:         createTimeToID(now),
		CreateTime: now,
		Payload:    payload,
	}

	err := s.db.InsertIntoBucket(s.bucket, r.ID, r)
	if err != nil {
		return history.Record{}, err
	}
	return r, nil
}

func createTimeToID(now time.Time) string {
	return strconv.FormatInt(now.UnixNano(), 10)
}

type slice struct {
	db *bolthold.Store

	bucket   *bolt.Bucket // distinguishes between this store and other stores that use the same table
	from, to history.Record
}

func (s slice) Slice(from, to history.Record) history.Slice {
	return slice{
		db:     s.db,
		bucket: s.bucket,
		from:   from,
		to:     to,
	}
}

func (s slice) Read(ctx context.Context, into []history.Record) (int, error) {
	var query *bolthold.Query
	if !s.from.IsZero() {
		if s.from.ID != "" {
			query = bolthold.Where("ID").Ge(s.from.ID)
		} else if !s.from.CreateTime.IsZero() {
			query = bolthold.Where("CreateTime").Ge(s.from.CreateTime)
		}
	}
	if !s.to.IsZero() {
		if s.to.ID != "" {
			if query == nil {
				query = bolthold.Where("ID").Lt(s.to.ID)
			} else {
				query = query.And("ID").Lt(s.to.ID)
			}
		} else if !s.to.CreateTime.IsZero() {
			if query == nil {
				query = bolthold.Where("CreateTime").Lt(s.to.CreateTime)
			} else {
				query = query.And("CreateTime").Lt(s.to.CreateTime)
			}
		}
	}

	err := s.db.FindInBucket(s.bucket, &into, query)
	if err != nil {
		return 0, err
	}
	return len(into), nil
}

func (s slice) Len(ctx context.Context) (int, error) {
	tmp := make([]history.Record, 0)
	return s.Read(ctx, tmp)
}
