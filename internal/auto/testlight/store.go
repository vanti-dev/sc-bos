package testlight

import (
	"errors"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"go.etcd.io/bbolt"
)

type MostRecentScanKey struct {
	Name string
	Kind rpc.Test
}

type TestResultRecord struct {
	Name             string   `boltholdIndex:"Name"`
	Kind             rpc.Test `boltholdIndex:"Kind"`
	CompletedAfter   time.Time
	CompletedBefore  time.Time
	Success          bool
	AchievedDuration time.Duration // only used for Kind == DurationTest
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
