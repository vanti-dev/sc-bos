package devices

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_parsePath(t *testing.T) {
	n := func(name string) pathSegment {
		return pathSegment{Name: name}
	}
	i := func(index int) pathSegment {
		return pathSegment{Index: index, IsIndex: true}
	}

	tests := []struct {
		path    string
		want    []pathSegment
		wantErr bool
	}{
		{"", []pathSegment{}, true},
		{".", []pathSegment{}, true},
		{"a", []pathSegment{n("a")}, false},
		{"a.", []pathSegment{n("a")}, true},
		{".a", []pathSegment{}, true},
		{"abc", []pathSegment{n("abc")}, false},
		{"世界", []pathSegment{n("世界")}, false},
		{"a.b", []pathSegment{n("a"), n("b")}, false},
		{"a.b.", []pathSegment{n("a"), n("b")}, true},
		{"a..b", []pathSegment{n("a")}, true},
		{"a.世界", []pathSegment{n("a"), n("世界")}, false},
		{"世.界", []pathSegment{n("世"), n("界")}, false},
		{"世界.b", []pathSegment{n("世界"), n("b")}, false},
		{"a[0]", []pathSegment{n("a"), i(0)}, false},
		{"a[10]", []pathSegment{n("a"), i(10)}, false},
		{"a[+10]", []pathSegment{n("a"), i(10)}, false},
		{"a[-10]", []pathSegment{n("a"), i(-10)}, false},
		{"a[foo]", []pathSegment{n("a")}, true},
		{"世界[1]", []pathSegment{n("世界"), i(1)}, false},
		{"a[1].b", []pathSegment{n("a"), i(1), n("b")}, false},
		{"a.b[1]", []pathSegment{n("a"), n("b"), i(1)}, false},
		{"a[", []pathSegment{n("a")}, true},
		{"a[1", []pathSegment{n("a")}, true},
		{"[", []pathSegment{}, true},
		{"[1]", []pathSegment{i(1)}, false},
		{"[1.1]", []pathSegment{}, true},
		{"[1][2]", []pathSegment{i(1), i(2)}, false},
		{"a[1][2]", []pathSegment{n("a"), i(1), i(2)}, false},
		{"a[1][2].b", []pathSegment{n("a"), i(1), i(2), n("b")}, false},
		{"a[1].b[2].c", []pathSegment{n("a"), i(1), n("b"), i(2), n("c")}, false},
		{"[1]b", []pathSegment{i(1)}, true},
		{"a[1]b", []pathSegment{n("a"), i(1)}, true},
		{"a[1]b[2]", []pathSegment{n("a"), i(1)}, true},
		{"a.[1]", []pathSegment{n("a")}, true},
		{"a[1].[2]", []pathSegment{n("a"), i(1)}, true},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			ps, err := parsePath(test.path)
			if (err != nil) != test.wantErr {
				t.Errorf("parsePath(%q) error = %v, wantErr %v", test.path, err, test.wantErr)
			}
			if diff := cmp.Diff(test.want, ps, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("parsePath(%q): -want +got:\n%s", test.path, diff)
			}
		})
	}
}
