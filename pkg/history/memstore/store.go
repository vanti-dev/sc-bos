// Package memstore provides an implementation of history.Store with records stored in memory.
package memstore

import (
	"context"
	"errors"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

type Store struct {
	slice // sorted by id, which is createTime+dedupe index
	now   func() time.Time

	maxAge   time.Duration
	maxCount int64
}

func New(opts ...Option) *Store {
	s := &Store{now: time.Now}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func SetNow(s *Store, now func() time.Time) func() {
	old := s.now
	s.now = now
	return func() {
		s.now = old
	}
}

func (s *Store) Append(_ context.Context, payload []byte) (history.Record, error) {
	s.slice.mtx.Lock()
	defer s.slice.mtx.Unlock()
	now := s.now()
	l := len(s.slice.records)
	if l > 0 && s.slice.records[l-1].CreateTime.After(now) {
		return history.Record{}, errors.New("time is running backwards")
	}

	r := history.Record{Payload: payload, CreateTime: now, ID: createTimeToID(now)}
	if l > 0 && s.slice.records[l-1].CreateTime == now {
		// todo: cap this memory growth. Should be ok for now as this isn't likely in production
		// 000 sorts after 00 after 0, so adding a 0 to the time maintains the sort order
		r.ID = s.slice.records[l-1].ID + "0"
	}
	s.slice.records = append(s.slice.records, r)
	s.gc(now)
	return r, nil
}

func (s *Store) gc(now time.Time) {
	if s.maxAge == 0 && s.maxCount == 0 {
		return
	}

	if s.maxAge > 0 {
		// find the first record that is "not older" than maxAge and drop everything before it
		if i, ok := s.indexOf(history.Record{CreateTime: now.Add(-s.maxAge)}); ok {
			s.slice.records = s.slice.records[i:]
		}
	}
	if s.maxCount > 0 {
		if l := int64(len(s.slice.records)); l > s.maxCount {
			s.slice.records = s.slice.records[l-s.maxCount:]
		}
	}
}

type records []history.Record // sorted by id, which is createTime+dedupe

type slice struct {
	records

	mtx sync.Mutex
}

func (rs *slice) indexOf(r history.Record) (int, bool) {
	id := computeId(r)
	return sort.Find(len(rs.records), func(i int) int {
		return strings.Compare(id, rs.records[i].ID)
	})
}

func (rs *slice) indexFrom(from history.Record) (int, bool) {
	if from.IsZero() {
		return 0, true
	}
	return rs.indexOf(from)
}

func (rs *slice) indexTo(to history.Record) (int, bool) {
	if to.IsZero() {
		return len(rs.records), true
	}
	return rs.indexOf(to)
}

func (rs *slice) Slice(from, to history.Record) history.Slice {
	rs.mtx.Lock()
	defer rs.mtx.Unlock()
	fromIndex, _ := rs.indexFrom(from)
	toIndex, _ := rs.indexTo(to)
	return &slice{records: rs.records[fromIndex:toIndex]}
}

func (rs *slice) Read(_ context.Context, into []history.Record) (int, error) {
	rs.mtx.Lock()
	defer rs.mtx.Unlock()
	return copy(into, rs.records), nil
}

func (rs *slice) ReadDesc(_ context.Context, into []history.Record) (int, error) {
	i := copy(into, rs.records[max(len(rs.records)-len(into), 0):])
	slices.Reverse(into[:i])
	return i, nil
}

func (rs *slice) Len(_ context.Context) (int, error) {
	return len(rs.records), nil
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
