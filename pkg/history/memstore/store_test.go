package memstore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

func TestStore_Append_gc(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2023-11-16T10:00:00Z")
	now = now.Round(time.Second) // get rid of millis, etc
	// All contains n records, each 1 hour apart in ascending order, with the last record being now.
	const n = 10
	all := make(slice, n)
	for i := range all {
		t := now.Add(time.Duration(i-n+1) * time.Hour)
		all[i] = history.Record{ID: createTimeToID(t), CreateTime: t, Payload: []byte(fmt.Sprint(i + 1))}
	}

	tests := []struct {
		maxAge   time.Duration
		maxCount int64
		want     slice
	}{
		{0 * time.Hour, 0, all},     // no gc needed, no cap
		{n * time.Hour, 0, all},     // no gc needed, age cap is too big
		{0 * time.Hour, n, all},     // no gc needed, count cap is to big
		{n * time.Hour, n, all},     // no gc needed, caps are to big
		{2 * time.Hour, 0, all[7:]}, // There are three records between now and -2h: 10, 9, and 8
		{0 * time.Hour, 3, all[7:]},
		{3 * time.Hour, 3, all[7:]}, // count beats age
		{2 * time.Hour, 4, all[7:]}, // age beats count
	}
	for _, tt := range tests {
		name := fmt.Sprintf("maxAge=%v,maxCount=%d", tt.maxAge, tt.maxCount)
		t.Run(name, func(t *testing.T) {
			s := &Store{
				slice:    all[: n-1 : n-1], // make sure slice has no spare capacity to avoid Append assigning to the underlying array of all
				maxAge:   tt.maxAge,
				maxCount: tt.maxCount,
				now:      func() time.Time { return now },
			}
			_, err := s.Append(context.Background(), all[n-1].Payload)
			if err != nil {
				t.Errorf("Append() error = %v", err)
				return
			}
			if diff := cmp.Diff(tt.want, s.slice); diff != "" {
				t.Errorf("Append() slice diff (-want +got):\n%s", diff)
			}
		})
	}
}
