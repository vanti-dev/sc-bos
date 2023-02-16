// Package memstore provides an implementation of history.Store with records stored in memory.
package memstore

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

type Store struct {
	slice // sorted by id, which is createTime+dedupe index
	now   func() time.Time
}

func New() *Store {
	return &Store{now: time.Now}
}

func (s *Store) Append(_ context.Context, payload []byte) (history.Record, error) {
	now := s.now()
	l := len(s.slice)
	if l > 0 && s.slice[l-1].CreateTime.After(now) {
		return history.Record{}, errors.New("time is running backwards")
	}

	r := history.Record{Payload: payload, CreateTime: now, ID: createTimeToID(now)}
	if l > 0 && s.slice[l-1].CreateTime == now {
		// todo: cap this memory growth. Should be ok for now as this isn't likely in production
		// 000 sorts after 00 after 0, so adding a 0 to the time maintains the sort order
		r.ID = s.slice[l-1].ID + "0"
	}
	s.slice = append(s.slice, r)
	return r, nil
}

type slice []history.Record // sorted by id, which is createTime+dedupe

func (rs slice) indexOf(r history.Record) (int, bool) {
	id := computeId(r)
	return sort.Find(len(rs), func(i int) int {
		return strings.Compare(id, rs[i].ID)
	})
}

func (rs slice) indexFrom(from history.Record) (int, bool) {
	if from.IsZero() {
		return 0, true
	}
	return rs.indexOf(from)
}

func (rs slice) indexTo(to history.Record) (int, bool) {
	if to.IsZero() {
		return len(rs), true
	}
	return rs.indexOf(to)
}

func (rs slice) Slice(from, to history.Record) history.Slice {
	fromIndex, _ := rs.indexFrom(from)
	toIndex, _ := rs.indexTo(to)
	return rs[fromIndex:toIndex]
}

func (rs slice) Read(_ context.Context, into []history.Record) (int, error) {
	return copy(into, rs), nil
}

func (rs slice) Len(_ context.Context) (int, error) {
	return len(rs), nil
}

func computeId(r history.Record) string {
	if r.ID != "" {
		return r.ID
	}
	return createTimeToID(r.CreateTime)
}

func createTimeToID(t time.Time) string {
	return t.In(time.UTC).Format(time.RFC3339Nano)
}
