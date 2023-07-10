package status

import (
	"fmt"
	"testing"
)

func Test_nominalValue(t *testing.T) {
	tests := []struct {
		v    any
		nom  any
		want bool
	}{
		{true, true, true},
		{true, false, false},
		{false, true, false},
		{false, false, true},
		{int(1), true, true},
		{int(1), false, false},
		{int(0), true, false},
		{int(0), false, true},
		{float32(1), true, true},
		{float32(1), false, false},
		{float32(0), true, false},
		{float32(0), false, true},
		{int64(10), float64(10), true},
		{int64(10), 10.1, false},
		{"foo", "foo", true},
		{"foo", "bar", false},
		{"foo", 1, false},
		{"foo", false, false},
		{1, "foo", false},
	}
	for _, tt := range tests {
		cmp := "=="
		if !tt.want {
			cmp = "!="
		}
		t.Run(fmt.Sprintf("%[1]T[%[1]v] %[2]s %[3]T[%[3]v]", tt.v, cmp, tt.nom), func(t *testing.T) {
			if got := nominalValue(tt.v, tt.nom); got != tt.want {
				t.Errorf("nominalValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
