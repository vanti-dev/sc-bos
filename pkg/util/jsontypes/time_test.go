package jsontypes

import (
	"fmt"
	"testing"
	"time"
)

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       time.Duration
		want    string
		wantErr bool
	}{
		{"zero", 0, "0s", false},
		{"1s", time.Second, "1s", false},
		{"10s", 10 * time.Second, "10s", false},
		{"1m", time.Minute, "1m", false},
		{"1m1s", time.Minute + time.Second, "1m1s", false},
		{"1m10s", time.Minute + 10*time.Second, "1m10s", false},
		{"1.1s", time.Second + 100*time.Millisecond, "1.1s", false},
		{"1h", time.Hour, "1h", false}, // fails
		{"1h1s", time.Hour + time.Second, "1h1s", false},
		{"1ms", time.Millisecond, "1ms", false},
		{"1ns", time.Nanosecond, "1ns", false},
		{"1.001ms", time.Millisecond + time.Microsecond, "1.001ms", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{Duration: tt.d}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != fmt.Sprintf(`"%s"`, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want \"%v\"", got, tt.want)
			}
		})
	}
}
