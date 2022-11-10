package testlight

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"go.etcd.io/bbolt"
)

type MostRecentScanKey struct {
	Name string
	Kind rpc.Test
}

type TestResultRecord struct {
	UUID             string `boltholdKey:"UUID"`
	Name             string `boltholdIndex:"Name"`
	Kind             rpc.Test
	After            time.Time // the earliest time that the test might have completed
	Before           time.Time // for passes: the time when the test was recorded, for failures the most recent time when failure was reported
	TestPass         bool
	AchievedDuration time.Duration // only used for successful duration tests
}

// updateScanTime will update the MostRecentScanRecord for the device of the given name in the database.
// If a record was previously present in the database, then it will be returned in old with existed=true.
// Otherwise, existed=false and the zero time are returned.
func updateScanTime(db *bolthold.Store, name string, kind rpc.Test, scanTime time.Time) (old time.Time, existed bool, err error) {
	key := MostRecentScanKey{
		Name: name,
		Kind: kind,
	}
	err = db.Bolt().Update(func(tx *bbolt.Tx) error {
		// get the existing record for the device, if it exists
		err := db.TxGet(tx, key, &old)
		if errors.Is(err, bolthold.ErrNotFound) {
			existed = false
			old = time.Time{}
		} else if err != nil {
			return err
		} else {
			existed = true
		}

		// store the updated record
		err = db.TxUpsert(tx, key, scanTime)
		return err
	})
	return
}

func saveTestResult(db *bolthold.Store, record TestResultRecord) error {
	if record.TestPass {
		return db.Insert(genUUIDKey(), record)
	} else {
		return db.Bolt().Update(func(tx *bbolt.Tx) error {
			return mergeFailedTestResult(db, tx, record)
		})
	}
}

func genUUIDKey() string {
	return uuid.New().String()
}

// must only be called with record.Success == false
// If most recent stored test record with matching Name and Kind is also a failure, then merges the provided record into
// the stored record by extending the Before time. This is necessary because failed test results can't be manually cleared,
// which would result in lots of duplicated failure records.
// If no such matching record is found, inserts the record into the database.
func mergeFailedTestResult(db *bolthold.Store, tx *bbolt.Tx, record TestResultRecord) error {
	if record.TestPass {
		panic("merging not applicable to test passes")
	}

	var latestMatching TestResultRecord
	err := db.TxFindOne(tx, &latestMatching,
		bolthold.Where("Name").Eq(record.Name).
			And("Kind").Eq(record.Kind).
			SortBy("Before").Reverse().
			Limit(1),
	)

	if errors.Is(err, bolthold.ErrNotFound) {
		return db.TxInsert(tx, genUUIDKey(), record)
	} else if err != nil {
		return err
	}

	// update the stored record with the new time
	if record.Before.After(latestMatching.Before) {
		latestMatching.Before = record.Before
	}
	err = db.TxUpdate(tx, latestMatching.UUID, latestMatching)
	return err
}
