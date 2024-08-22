package block

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPath_String(t *testing.T) {
	type testCase struct {
		path   Path
		expect string
	}
	cases := map[string]testCase{
		"empty": {
			path:   nil,
			expect: "/",
		},
		"field": {
			path:   Path{{Field: "foo"}},
			expect: "/foo",
		},
		"field_with_whitespace": {
			path:   Path{{Field: " foo bar "}},
			expect: `/" foo bar "`,
		},
		"field_leading_digit": {
			path:   Path{{Field: "42foo"}},
			expect: `/"42foo"`,
		},
		"field_following_digit": {
			path:   Path{{Field: "foo42"}},
			expect: "/foo42",
		},
		"field_unicode": {
			path:   Path{{Field: "🦄"}},
			expect: `/"🦄"`,
		},
		"nested_field": {
			path:   Path{{Field: "foo"}, {Field: "bar"}},
			expect: "/foo/bar",
		},
		"array_elem_str": {
			path:   Path{{Field: "foo"}, {ArrayKey: "name", ArrayElem: "bar"}},
			expect: `/foo[name="bar"]`,
		},
		"array_elem_int": {
			path:   Path{{Field: "foo"}, {ArrayKey: "id", ArrayElem: 42}},
			expect: `/foo[id=42]`,
		},
		"array_key_with_whitespace": {
			path:   Path{{Field: "foo"}, {ArrayKey: " name ", ArrayElem: "bar"}},
			expect: `/foo[" name "="bar"]`,
		},
		"array_key_unicode": {
			path:   Path{{Field: "foo"}, {ArrayKey: "nåme", ArrayElem: "bar"}},
			expect: `/foo["nåme"="bar"]`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.path.String(); got != tc.expect {
				t.Errorf("expected %q, got %q", tc.expect, got)
			}
		})
	}
}

type parsePathTest struct {
	path      string
	expect    Path
	shouldErr bool
}

var parsePathTests = map[string]parsePathTest{
	"empty": {
		path:      "",
		shouldErr: true,
	},
	"root": {
		path:   "/",
		expect: Path{},
	},
	"empty_field": {
		path:      "/foo/",
		shouldErr: true,
	},
	"empty_field_2": {
		path:      "//",
		shouldErr: true,
	},
	"field": {
		path:   "/foo_99-123",
		expect: Path{{Field: "foo_99-123"}},
	},
	"field_quoted": {
		path:   `/" føø bår "`,
		expect: Path{{Field: " føø bår "}},
	},
	"field_quoted_escape": {
		path:   `/"foo\"bar"`,
		expect: Path{{Field: `foo"bar`}},
	},
	"field_not_quoted": {
		path:      `/foo bar`,
		shouldErr: true,
	},
	"nested_field": {
		path:   "/foo/bar",
		expect: Path{{Field: "foo"}, {Field: "bar"}},
	},
	"array_elem_str": {
		path:   `/foo[name="bar"]`,
		expect: Path{{Field: "foo"}, {ArrayKey: "name", ArrayElem: "bar"}},
	},
	"array_elem_int": {
		path:   `/foo[id=42]`,
		expect: Path{{Field: "foo"}, {ArrayKey: "id", ArrayElem: 42.0}},
	},
	"double_subscript": {
		path:      `/foo[name="bar"][id=42]`,
		shouldErr: true,
	},
}

func TestParsePath(t *testing.T) {
	for name, tc := range parsePathTests {
		t.Run(name, func(t *testing.T) {
			got, err := ParsePath(tc.path)
			if tc.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				diff := cmp.Diff(got, tc.expect, cmpopts.EquateEmpty())
				if diff != "" {
					t.Errorf("unexpected result (-got +want):\n%s", diff)
				}
			}
		})
	}
}

func FuzzParsePath(f *testing.F) {
	for _, c := range parsePathTests {
		f.Add(c.path)
	}
	f.Fuzz(func(t *testing.T, path string) {
		// We don't care about the result, just that it doesn't panic.
		parsed, err := ParsePath(path)
		var parseErr *PathParseError
		if err != nil && parsed != nil {
			t.Error("expected nil result with error")
		}
		if errors.As(err, &parseErr) {
			if parseErr.Where > len(path) {
				t.Errorf("out-of-range error location: %d > %d", parseErr.Where, len(path))
			}
		}
	})
}

func TestPath_UnmarshalJSON(t *testing.T) {
	// should be able to parse in the string format (the same as accepted by ParsePath)
	for name, tc := range parsePathTests {
		t.Run("strfmt_"+name, func(t *testing.T) {
			data := []byte(strconv.Quote(tc.path))
			var p Path
			err := json.Unmarshal(data, &p)
			if tc.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				diff := cmp.Diff(p, tc.expect, cmpopts.EquateEmpty())
				if diff != "" {
					t.Errorf("unexpected result (-got +want):\n%s", diff)
				}
			}
		})
	}

	// should also accept the JSON array format
	type testCase struct {
		input     string
		path      Path
		shouldErr bool
	}
	cases := map[string]testCase{
		"empty": {
			input: "[]",
			path:  Path{},
		},
		"field": {
			input: `["foo"]`,
			path:  Path{{Field: "foo"}},
		},
		"nested_field": {
			input: `["foo", "bar"]`,
			path:  Path{{Field: "foo"}, {Field: "bar"}},
		},
		"array_elem_str": {
			input: `["foo", {"name": "bar"}]`,
			path:  Path{{Field: "foo"}, {ArrayKey: "name", ArrayElem: "bar"}},
		},
		"array_elem_num": {
			input: `["foo", {"id": 42}]`,
			path:  Path{{Field: "foo"}, {ArrayKey: "id", ArrayElem: 42.0}},
		},
		"invalid_array_elem": {
			input:     `["foo", {"id": "bar", "extra": "field"}]`,
			shouldErr: true,
		},
		"invalid_field": {
			input:     `["foo", 42]`,
			shouldErr: true,
		},
	}
	for name, tc := range cases {
		t.Run("arrfmt_"+name, func(t *testing.T) {
			data := []byte(tc.input)
			var p Path
			err := json.Unmarshal(data, &p)
			if tc.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				diff := cmp.Diff(p, tc.path, cmpopts.EquateEmpty())
				if diff != "" {
					t.Errorf("unexpected result (-got +want):\n%s", diff)
				}
			}
		})
	}
}
