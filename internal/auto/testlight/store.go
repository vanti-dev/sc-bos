package testlight

import (
	"errors"
	"strconv"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type LatestStatusRecord struct {
	Name       string `boltholdKey:"Name"`
	LastUpdate time.Time

	Faults []gen.EmergencyLightFault
}

type EventRecord struct {
	ID        uint64    `boltholdKey:"ID"`
	Name      string    `boltholdIndex:"Name"`
	Timestamp time.Time `boltholdIndex:"Timestamp"`

	// use a separate Kind field so it's indexable
	Kind             EventKind `boltholdIndex:"Kind"`
	StatusReport     *gen.EmergencyLightingEvent_StatusReport
	DurationTestPass *gen.EmergencyLightingEvent_DurationTestPass
	// no structure require for function test passes; that carries no data (and therefore can't be serialised)
}

type EventKind int

const (
	StatusReportEvent EventKind = iota + 1
	FunctionTestPassEvent
	DurationTestPassEvent
)

func saveFunctionTestPass(db *bolthold.Store, tx *bbolt.Tx, name string, timestamp time.Time) (err error) {
	record := EventRecord{
		Name:      name,
		Timestamp: timestamp,
		Kind:      FunctionTestPassEvent,
	}
	err = db.TxInsert(tx, bolthold.NextSequence(), record)
	return
}

func saveDurationTestPass(db *bolthold.Store, tx *bbolt.Tx, name string, timestamp time.Time, result time.Duration) (err error) {
	record := EventRecord{
		Name:      name,
		Timestamp: timestamp,
		Kind:      DurationTestPassEvent,
		DurationTestPass: &gen.EmergencyLightingEvent_DurationTestPass{
			AchievedDuration: durationpb.New(result),
		},
	}
	err = db.TxInsert(tx, bolthold.NextSequence(), record)
	return
}

func saveStatusReport(db *bolthold.Store, tx *bbolt.Tx, name string, timestamp time.Time, faults []gen.EmergencyLightFault) (err error) {
	record := EventRecord{
		Name:      name,
		Timestamp: timestamp,
		Kind:      StatusReportEvent,
		StatusReport: &gen.EmergencyLightingEvent_StatusReport{
			Faults: sortDeduplicateFaults(faults),
		},
	}
	err = db.TxInsert(tx, bolthold.NextSequence(), record)
	return
}

// updates the LatestStatusRecord for an emergency light
// creates the record if it does not exist
// Returns changed=true if a new record was created, or if the new faults list is meaningfully different from the old
// list.
func updateLatestStatus(db *bolthold.Store, tx *bbolt.Tx, name string, now time.Time, faults []gen.EmergencyLightFault) (changed bool, err error) {
	var existing LatestStatusRecord
	err = db.TxGet(tx, name, &existing)
	if errors.Is(err, bolthold.ErrNotFound) {
		changed = true
	} else if err != nil {
		return
	} else {
		changed = !faultsEquivalent(existing.Faults, faults)
	}

	record := LatestStatusRecord{
		Name:       name,
		LastUpdate: now,
		Faults:     sortDeduplicateFaults(faults), // normalise to produce a consistent order without duplicates
	}

	err = db.TxUpsert(tx, name, record)
	return
}

func findLatestStatusPaged(db *bolthold.Store, pageToken string, pageSize int) (page []LatestStatusRecord, nextToken string, err error) {
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	err = db.Find(&page,
		// this still works when request.PageToken=="", because all other strings are greater than the empty string
		bolthold.Where(bolthold.Key).Gt(pageToken).
			Limit(pageSize+1), // include one more, so we can tell if another page is required
	)
	if err != nil {
		return
	}

	more := len(page) > pageSize
	if more {
		page = page[:pageSize]
		nextToken = page[len(page)-1].Name
	}
	return
}

func findEventsPaged(db *bolthold.Store, pageToken string, pageSize int) (page []EventRecord, nextToken string, err error) {
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	query := &bolthold.Query{}
	if pageToken != "" {
		var boundaryId uint64
		boundaryId, err = strconv.ParseUint(pageToken, 10, 64)
		if err != nil {
			err = status.Errorf(codes.InvalidArgument, "invalid page_token: %s", err.Error())
			return
		}
		query = query.And(bolthold.Key).Gt(boundaryId)
	}

	err = db.Find(&page,
		query.Limit(pageSize+1), // include one more, so we can tell if another page is required
	)
	if err != nil {
		return
	}

	more := len(page) > pageSize
	if more {
		page = page[:pageSize]
		nextToken = strconv.FormatUint(page[len(page)-1].ID, 10)
	}
	return
}
