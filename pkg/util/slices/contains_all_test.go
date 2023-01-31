package slices

import (
	"fmt"
	"strings"
	"testing"
)

func TestContainsAll(t *testing.T) {
	type testCase struct {
		haystack string
		needle   string
		want     bool
	}
	tests := []testCase{
		{"", "", true},
		{"1", "", false},
		{"1,2", "", false},
		{"", "1", false},
		{"", "1,2", false},
		{"1,2", "1,2", true},
		{"1", "1,2", false},
		{"1,2", "1", true},
		{"1,2,3,4", "1,2", true},
		{"3,1,4,2", "1,2", true},
		{"3,1,4,2", "2,1", true},
		{"3,1,4,2", "2,1,a", false},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("[%v] contains [%v] == %v", tt.haystack, tt.needle, tt.want)
		t.Run(name, func(t *testing.T) {
			if got := ContainsAll(strings.Split(tt.needle, ","), strings.Split(tt.haystack, ",")); got != tt.want {
				t.Errorf("ContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
