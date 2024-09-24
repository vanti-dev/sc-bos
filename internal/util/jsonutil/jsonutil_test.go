package jsonutil

import (
	"testing"
)

func TestMarshalObjects(t *testing.T) {
	type testStruct struct {
		A string `json:"a"`
		B string `json:"b,omitempty"`
	}

	type testCase struct {
		inputs []any
		expect string
	}

	cases := map[string]testCase{
		"empty": {
			inputs: []any{},
			expect: "{}",
		},
		"single": {
			inputs: []any{
				map[string]any{"a": "a"},
			},
			expect: `{"a":"a"}`,
		},
		"add": {
			inputs: []any{
				map[string]any{"a": "a"},
				map[string]any{"b": "b"},
			},
			expect: `{"a":"a","b":"b"}`,
		},
		"replace": {
			inputs: []any{
				map[string]any{"a": "a"},
				map[string]any{"a": "b"},
			},
			expect: `{"a":"b"}`,
		},
		"struct_1": {
			inputs: []any{
				testStruct{A: "a"},
				testStruct{},
			},
			expect: `{"a":""}`,
		},
		"struct_2": {
			inputs: []any{
				testStruct{B: "b"},
				testStruct{B: ""},
			},
			expect: `{"a":"","b":"b"}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			data, err := MarshalObjects(tc.inputs...)
			if err != nil {
				t.Fatalf("MarshalObjects failed: %v", err)
			}
			if string(data) != tc.expect {
				t.Errorf("MarshalObjects returned %q, expected %q", data, tc.expect)
			}
		})
	}
}
