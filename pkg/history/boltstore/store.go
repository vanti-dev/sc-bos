// Package boltstore provides an implementation of history.Store with records stored in a bolthold database.
package boltstore

import (
	"context"
	"strconv"
	"time"

	"github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

type Store struct {
	slice // sorted by id, which is createTime+dedupe index
	now   func() time.Time

	logger *zap.Logger
}

func NewFromDb(db *bolthold.Store, source string, logger *zap.Logger) (history.Store, error) {
	b := []byte(source)
	err := db.Bolt().Update(func(tx *bolt.Tx) error {
		var err error
		logger.Debug("Creating bucket", zap.String("bucket", string(b)))
		_, err = tx.CreateBucketIfNotExists(b)
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
			bucket: b,
		},
		now:    time.Now,
		logger: logger,
	}, nil
}

func (s *Store) Append(ctx context.Context, payload []byte) (history.Record, error) {
	now := s.now()
	r := history.Record{
		ID:         createTimeToID(now),
		CreateTime: now,
		Payload:    payload,
	}

	err := s.db.Bolt().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		return s.db.InsertIntoBucket(b, r.ID, r)
	})
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

	bucket   []byte // distinguishes between this store and other stores that use the same db
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

	records := make([]history.Record, 0)

	err := s.db.Bolt().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		return s.db.FindInBucket(b, &records, query)
	})
	if err != nil {
		return 0, err
	}

	copy(into, records)

	return len(into), nil
}

func (s slice) Len(ctx context.Context) (int, error) {
	tmp := make([]history.Record, 0)
	return s.Read(ctx, tmp)
}
