package testlight

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestUpdateLatestStatus(t *testing.T) {
	db, cleanup := prepareEmptyTestDB()
	defer cleanup()

	deviceName := "ns/test/example"

	doUpdateLatestStatus := func(now time.Time, expectChanged bool, faults ...gen.EmergencyLightFault) {
		var changed bool
		err := db.Bolt().Update(func(tx *bbolt.Tx) (err error) {
			changed, err = updateLatestStatus(db, tx, deviceName, now, faults)
			return
		})
		if err != nil {
			t.Fatalf("failed to run transaction: %s", err.Error())
		}
		if changed != expectChanged {
			t.Errorf("expected changed=%v but got changed=%v", expectChanged, changed)
		}
		return
	}
	expectLatestStatus := func(expected LatestStatusRecord) {
		var actual LatestStatusRecord
		err := db.Get(deviceName, &actual)
		if err != nil {
			t.Fatalf("failed to get LatestStatusRecord entry: %s", err.Error())
		}
		diff := cmp.Diff(expected, actual,
			cmpopts.EquateEmpty(),
			protocmp.Transform(),
		)
		if diff != "" {
			t.Errorf("unexpected LatestStatusRecord (-want +got):\n%s", diff)
		}
		return
	}
	expectNotExists := func() {
		var actual LatestStatusRecord
		err := db.Get(deviceName, &actual)
		if !errors.Is(err, bolthold.ErrNotFound) {
			t.Errorf("expected ErrNotFound but got: %s", err.Error())
		}
	}

	expectNotExists()
	t1 := time.Date(2022, time.November, 23, 15, 7, 0, 0, time.UTC)

	// when the record is first inserted, it is always considered changed, even with no faults
	doUpdateLatestStatus(t1, true)
	expectLatestStatus(LatestStatusRecord{
		Name:       deviceName,
		LastUpdate: t1,
		Faults:     nil,
	})
	// if we update again with no faults, then the faults list has not changed
	// the timestamp should still be updated though
	t2 := t1.Add(time.Hour)
	doUpdateLatestStatus(t2, false)
	expectLatestStatus(LatestStatusRecord{
		Name:       deviceName,
		LastUpdate: t2,
		Faults:     nil,
	})
	// now update again, with a different faults list. This should produce a change.
	t3 := t2.Add(time.Hour)
	doUpdateLatestStatus(t3, true, gen.EmergencyLightFault_BATTERY_FAULT, gen.EmergencyLightFault_LAMP_FAULT)
	expectLatestStatus(LatestStatusRecord{
		Name:       deviceName,
		LastUpdate: t3,
		Faults: []gen.EmergencyLightFault{
			gen.EmergencyLightFault_BATTERY_FAULT,
			gen.EmergencyLightFault_LAMP_FAULT,
		},
	})
	// update again, but with a faults list which is equivalent (but not exactly equal, because of ordering)
	// this is still considered the same
	t4 := t3.Add(time.Hour)
	doUpdateLatestStatus(t4, false, gen.EmergencyLightFault_LAMP_FAULT, gen.EmergencyLightFault_BATTERY_FAULT)
	expectLatestStatus(LatestStatusRecord{
		Name:       deviceName,
		LastUpdate: t4,
		Faults: []gen.EmergencyLightFault{
			gen.EmergencyLightFault_BATTERY_FAULT,
			gen.EmergencyLightFault_LAMP_FAULT,
		},
	})
}

func TestFindEventsPaged(t *testing.T) {
	db, cleanup := prepareEmptyTestDB()
	defer cleanup()
	expectedEvents := testEvents()
	addEvents(t, db, expectedEvents)

	// get all the events
	page, token, err := findEventsPaged(db, "", 10)
	if err != nil {
		t.Errorf("get all events: %s", err.Error())
	} else {
		if token != "" {
			t.Errorf("got an unexpected next page token: %q", token)
		}

		diff := cmp.Diff(expectedEvents, page, protocmp.Transform(), cmpopts.EquateEmpty())
		if diff != "" {
			t.Errorf("get all events: mismatch (-want +got):\n%s", diff)
		}
	}

	// get the first event
	page, token, err = findEventsPaged(db, "", 1)
	if err != nil {
		t.Errorf("get first event: %s", err.Error())
	} else {
		if token == "" {
			t.Error("missing next page token")
		}

		diff := cmp.Diff(expectedEvents[:1], page, protocmp.Transform(), cmpopts.EquateEmpty())
		if diff != "" {
			t.Errorf("get all events: mismatch (-want +got):\n%s", diff)
		}
	}
	// use token to get the next two events
	page, token, err = findEventsPaged(db, token, 2)
	if err != nil {
		t.Errorf("get first event: %s", err.Error())
	} else {
		if token != "" {
			t.Errorf("got an unexpected next page token: %q", token)
		}

		diff := cmp.Diff(expectedEvents[1:], page, protocmp.Transform(), cmpopts.EquateEmpty())
		if diff != "" {
			t.Errorf("get all events: mismatch (-want +got):\n%s", diff)
		}
	}
}

func addEvents(t *testing.T, db *bolthold.Store, events []EventRecord) {
	err := db.Bolt().Update(func(tx *bbolt.Tx) error {
		for _, event := range events {
			err := db.TxInsert(tx, event.ID, event)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("can't add events to test db: %s", err.Error())
	}
}

func testEvents() []EventRecord {
	t := time.Date(2022, time.November, 24, 0, 15, 0, 0, time.UTC)
	nextT := func() (out time.Time) {
		out = t
		t = t.Add(time.Hour)
		return
	}

	var id uint64
	nextID := func() (out uint64) {
		out = id
		id += 1
		return
	}

	name1 := "ns/test/emergency-light-1"
	name2 := "ns/test/emergency-light-2"

	return []EventRecord{
		{
			ID:        nextID(),
			Name:      name1,
			Timestamp: t,
			Kind:      FunctionTestPassEvent,
		},
		{
			ID:        nextID(),
			Name:      name1,
			Timestamp: nextT(),
			Kind:      StatusReportEvent,
			StatusReport: &gen.EmergencyLightingEvent_StatusReport{
				Faults: []gen.EmergencyLightFault{},
			},
		},
		{
			ID:        nextID(),
			Name:      name2,
			Timestamp: nextT(),
			Kind:      StatusReportEvent,
			StatusReport: &gen.EmergencyLightingEvent_StatusReport{
				Faults: []gen.EmergencyLightFault{gen.EmergencyLightFault_COMMUNICATION_FAILURE},
			},
		},
	}
}

func prepareEmptyTestDB() (store *bolthold.Store, cleanup func()) {
	f, err := os.CreateTemp("", "*.bolt")
	if err != nil {
		panic(fmt.Errorf("failed to create temp file for test db: %w", err))
	}
	name := f.Name()
	// we don't need the file handle, will open by name
	_ = f.Close()

	store, err = bolthold.Open(name, 0700, nil)
	if err != nil {
		panic(fmt.Errorf("failed to init db in %q: %w", name, err))
	}

	cleanup = func() {
		err := store.Close()
		if err != nil {
			log.Printf("error closing temporary bolt db at %q: %s", name, err.Error())
		}

		err = os.Remove(name)
		if err != nil {
			log.Printf("error removing temporary db at %q: %s", name, err.Error())
		}
	}
	return
}
