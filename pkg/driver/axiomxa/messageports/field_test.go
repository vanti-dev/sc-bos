package messageports

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPattern_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		pattern []Field
		want    Fields
		wantErr bool
	}{
		{
			"example",
			"02/03/2023 15:44:23,65536,Access granted: reader,ACU-06 Comms intake,Reader door 1,7349,1280458934,Anthony Ellis",
			[]Field{Timestamp, EventID, EventDesc, NetworkDesc, DeviceDesc, CardID, CardNumber, CardholderDesc},
			Fields{
				Timestamp:      Time{time.Date(2023, 3, 2, 15, 44, 23, 0, time.UTC)},
				EventID:        uint(65536),
				EventDesc:      "Access granted: reader",
				NetworkDesc:    "ACU-06 Comms intake",
				DeviceDesc:     "Reader door 1",
				CardID:         uint(7349),
				CardNumber:     uint64(1280458934),
				CardholderDesc: "Anthony Ellis",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pattern{Fields: tt.pattern}
			dst := Fields{}
			if err := p.Unmarshal([]byte(tt.data), &dst); (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, dst); diff != "" {
				t.Fatalf("Diff (-want,+got)\n%s", diff)
			}
		})
	}
}
