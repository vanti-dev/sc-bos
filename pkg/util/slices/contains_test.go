package slices

import (
	"fmt"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	type testCase struct {
		haystack string
		needle   string
		want     bool
	}
	tests := []testCase{
		{"", "", false},
		{"1", "", false},
		{"1,2", "", false},
		{"", "1", false},
		{"1,2", "2", true},
		{"1", "2", false},
		{"1,2", "1", true},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("[%v] contains '%v' == %v", tt.haystack, tt.needle, tt.want)
		t.Run(name, func(t *testing.T) {
			haystack := strings.Split(tt.haystack, ",")
			if tt.haystack == "" {
				haystack = nil
			}
			if got := Contains(tt.needle, haystack); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
