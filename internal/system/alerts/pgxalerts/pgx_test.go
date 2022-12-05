package pgxalerts

import (
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"strings"
	"testing"
)

func Test_fieldMaskIncludesPath(t *testing.T) {
	tests := []struct {
		name string
		m    string
		p    string
		want bool
	}{
		{"nil", "", "", true},
		{"nil prop", "", "prop", true},
		{"has", "prop", "prop", true},
		{"includes", "bar,prop,foo", "prop", true},
		{"not includes", "bar,prop,foo", "baz", false},
		{"parent", "parent.child", "parent", true},
		{"parent.child", "parent.child", "parent.child", true},
		{"match parent", "parent", "parent.child", true},
		{"invert parent", "parent.child", "child", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m *fieldmaskpb.FieldMask
			if tt.m != "" {
				m = &fieldmaskpb.FieldMask{Paths: strings.Split(tt.m, ",")}
			}
			if got := fieldMaskIncludesPath(m, tt.p); got != tt.want {
				t.Errorf("fieldMaskIncludesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
