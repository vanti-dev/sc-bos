package testlight

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/types/known/durationpb"
)

type LatestStatusRecord struct {
	Name       string `boltholdKey:"Name"`
	LastUpdate time.Time

	Faults []gen.EmergencyLightFault
}

type EventRecord struct {
	ID        string    `boltholdKey:"ID"`
	Name      string    `boltholdIndex:"Name"`
	Timestamp time.Time `boltholdIndex:"Timestamp"`

	// use a separate Kind field so it's indexable
	Kind             EventKind `boltholdIndex:"Kind"`
	StatusReport     *gen.EmergencyLightingEvent_StatusReport
	FunctionTestPass *gen.EmergencyLightingEvent_FunctionTestPass
	DurationTestPass *gen.EmergencyLightingEvent_DurationTestPass
}

type EventKind int

const (
	StatusReportEvent EventKind = iota + 1
	FunctionTestPassEvent
	DurationTestPassEvent
)

func saveFunctionTestPass(db *bolthold.Store, tx *bbolt.Tx, name string, timestamp time.Time) (id string, err error) {
	id = genUUIDKey()
	record := EventRecord{
		ID:               id,
		Name:             name,
		Timestamp:        timestamp,
		Kind:             FunctionTestPassEvent,
		FunctionTestPass: &gen.EmergencyLightingEvent_FunctionTestPass{},
	}
	err = db.TxInsert(tx, id, record)
	return
}

func saveDurationTestPass(db *bolthold.Store, tx *bbolt.Tx, name string, timestamp time.Time, result time.Duration) (id string, err error) {
	id = genUUIDKey()
	record := EventRecord{
		ID:        id,
		Name:      name,
		Timestamp: timestamp,
		Kind:      DurationTestPassEvent,
		DurationTestPass: &gen.EmergencyLightingEvent_DurationTestPass{
			AchievedDuration: durationpb.New(result),
		},
	}
	err = db.TxInsert(tx, id, record)
	return
}

func saveStatusReport(db *bolthold.Store, tx *bbolt.Tx, name string, timestamp time.Time, faults []gen.EmergencyLightFault) (id string, err error) {
	id = genUUIDKey()
	record := EventRecord{
		ID:        id,
		Name:      name,
		Timestamp: timestamp,
		Kind:      StatusReportEvent,
		StatusReport: &gen.EmergencyLightingEvent_StatusReport{
			Faults: faults,
		},
	}
	err = db.TxInsert(tx, id, record)
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
		Faults:     faults,
	}

	err = db.TxUpsert(tx, name, record)
	return
}

func genUUIDKey() string {
	return uuid.New().String()
}
