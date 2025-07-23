package slices

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewSorted(t *testing.T) {
	t.Run("presorted", func(t *testing.T) {
		s := NewSorted("a", "b", "c")
		assertAllItems(t, s, "a", "b", "c")
	})
	t.Run("unsorted", func(t *testing.T) {
		s := NewSorted("c", "a", "b")
		assertAllItems(t, s, "a", "b", "c")
	})
}

func TestSorted_Set(t *testing.T) {
	tests := []struct {
		name  string
		start []string
		set   string
		want  []string
	}{
		{"set after", []string{"5", "6", "7"}, "8", []string{"5", "6", "7", "8"}},
		{"set before", []string{"5", "6", "7"}, "4", []string{"4", "5", "6", "7"}},
		{"set middle", []string{"5", "6", "7"}, "6", []string{"5", "6", "7"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSorted(tt.start...)
			_, old, replaced := s.Set(tt.set)
			wantReplace := slices.Contains(tt.start, tt.set)
			if replaced != wantReplace {
				t.Errorf("expected replaced to be %v, got %v", wantReplace, replaced)
			}
			if wantReplace {
				if old != tt.set {
					t.Errorf("expected old value to be the same as set value, got %s", old)
				}
			} else {
				if old != "" {
					t.Errorf("expected old value to be empty, got %s", old)
				}
			}
			assertAllItems(t, s, tt.want...)
		})
	}
}

func cmpReverse(a, b string) int {
	return -1 * strings.Compare(a, b)
}

func TestNewSortedFunc(t *testing.T) {
	t.Run("presorted", func(t *testing.T) {
		s := NewSortedFunc(cmpReverse, "c", "b", "a")
		assertAllItems(t, s, "c", "b", "a")
	})
	t.Run("unsorted", func(t *testing.T) {
		s := NewSortedFunc(cmpReverse, "a", "c", "b")
		assertAllItems(t, s, "c", "b", "a")
	})
}

func TestSorted_Set_func(t *testing.T) {
	tests := []struct {
		name  string
		start []string
		set   string
		want  []string
	}{
		{"set after", []string{"7", "6", "5"}, "4", []string{"7", "6", "5", "4"}},
		{"set before", []string{"7", "6", "5"}, "8", []string{"8", "7", "6", "5"}},
		{"set middle", []string{"7", "6", "5"}, "6", []string{"7", "6", "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSortedFunc(cmpReverse, tt.start...)
			_, old, replaced := s.Set(tt.set)
			wantReplace := slices.Contains(tt.start, tt.set)
			if replaced != wantReplace {
				t.Errorf("expected replaced to be %v, got %v", wantReplace, replaced)
			}
			if wantReplace {
				if old != tt.set {
					t.Errorf("expected old value to be the same as set value, got %s", old)
				}
			} else {
				if old != "" {
					t.Errorf("expected old value to be empty, got %s", old)
				}
			}
			assertAllItems(t, s, tt.want...)
		})
	}
}

func assertAllItems(t *testing.T, s *Sorted[string], expected ...string) {
	t.Helper()
	var all []string
	for _, item := range s.All {
		all = append(all, item)
	}
	if wantLen, gotLen := len(expected), len(all); wantLen != gotLen {
		t.Errorf("expected %d items, got %d", wantLen, gotLen)
	}
	if diff := cmp.Diff(expected, all); len(diff) != 0 {
		t.Errorf("expected items to be sorted, (-want,+got)\n%s", diff)
	}
}
