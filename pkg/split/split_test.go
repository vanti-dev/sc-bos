package split

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDiff(t *testing.T) {
	type testCase struct {
		a      any
		b      any
		schema []Split
		expect []Patch
	}

	cases := map[string]testCase{
		"Empty": {
			a:      nil,
			b:      nil,
			schema: nil,
			expect: nil,
		},
		"Primitive": {
			a:      "foo",
			b:      "bar",
			schema: nil,
			expect: []Patch{{Value: "bar"}},
		},
		"TopLevelFields": {
			a: map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
			b: map[string]any{
				"foo": "bar",
				"baz": "waldo",
			},
			schema: nil,
			expect: []Patch{{Value: map[string]any{"foo": "bar", "baz": "waldo"}}},
		},
		"NestedFields": {
			a: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			b: map[string]any{
				"foo": map[string]any{
					"bar": "qux",
				},
			},
			schema: nil,
			expect: []Patch{{Value: map[string]any{"foo": map[string]any{"bar": "qux"}}}},
		},
		"SplitField": {
			a: map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
			b: map[string]any{
				"foo": "barbar",
				"baz": "qux",
			},
			schema: []Split{{Path: []string{"foo"}}},
			expect: []Patch{
				{
					Path:  []PathSegment{{Field: "foo"}},
					Value: "barbar",
				},
			},
		},
		"SplitFieldAdded": {
			a: map[string]any{
				"baz": "qux",
			},
			b: map[string]any{
				"foo": "barbar",
				"baz": "qux",
			},
			schema: []Split{{Path: []string{"foo"}}},
			expect: []Patch{
				{
					Path:  []PathSegment{{Field: "foo"}},
					Value: "barbar",
				},
			},
		},
		"SplitFieldRemoved": {
			a: map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
			b: map[string]any{
				"baz": "qux",
			},
			schema: []Split{{Path: []string{"foo"}}},
			expect: []Patch{
				{
					Path:    []PathSegment{{Field: "foo"}},
					Deleted: true,
				},
			},
		},
		"SplitFieldNested": {
			a: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
					"qux": "waldo",
				},
			},
			b: map[string]any{
				"foo": map[string]any{
					"bar": "bazbaz",
					"qux": "quxqux",
				},
			},
			schema: []Split{
				{
					Path: []string{"foo"},
					Splits: []Split{
						{Path: []string{"bar"}},
					},
				},
			},
			expect: []Patch{
				{
					Path: []PathSegment{{Field: "foo"}},
					Value: map[string]any{
						"bar": Ignore{},
						"qux": "quxqux",
					},
				},
				{
					Path:  []PathSegment{{Field: "foo"}, {Field: "bar"}},
					Value: "bazbaz",
				},
			},
		},
		"Array": {
			// because of the JSON roundtrip, the numbers must be float64
			a: map[string]any{
				"foo": []any{
					map[string]any{"id": float64(1), "addr": "foo"},
					map[string]any{"id": float64(2), "addr": "bar"},
				},
			},
			b: map[string]any{
				"foo": []any{
					map[string]any{"id": float64(1), "addr": "foo"},
					map[string]any{"id": float64(2), "addr": "baz"},
				},
			},
			schema: []Split{
				{
					Path: []string{"foo"},
					Key:  "id",
				},
			},
			expect: []Patch{
				{
					Path:  []PathSegment{{Field: "foo"}, {ArrayKey: "id", ArrayElem: float64(2)}},
					Value: map[string]any{"id": float64(2), "addr": "baz"},
				},
			},
		},
		"ArraySplitByKey": {
			a: map[string]any{
				"drivers": []any{
					map[string]any{
						"type": "a",
						"name": "driver-1",
						"objects": []any{
							map[string]any{"id": float64(1), "addr": "foo"},
							map[string]any{"id": float64(2), "addr": "bar"},
						},
					},
					map[string]any{
						"type":     "b",
						"name":     "driver-2",
						"settings": map[string]any{"mode": "on"},
					},
				},
			},
			b: map[string]any{
				"drivers": []any{
					map[string]any{
						"type": "a",
						"name": "driver-1",
						"objects": []any{
							map[string]any{"id": float64(1), "addr": "foo2"},
							map[string]any{"id": float64(2), "addr": "bar"},
						},
					},
					map[string]any{
						"type":     "b",
						"name":     "driver-2",
						"settings": map[string]any{"mode": "off"},
					},
				},
			},
			schema: []Split{
				{
					Path:     []string{"drivers"},
					Key:      "name",
					SplitKey: "type",
					SplitsByKey: map[string][]Split{
						"a": {
							{
								Path: []string{"objects"},
								Key:  "id",
							},
						},
						"b": {
							{
								Path: []string{"settings"},
							},
						},
					},
				},
			},
			expect: []Patch{
				{
					Path: []PathSegment{
						{Field: "drivers"},
						{ArrayKey: "name", ArrayElem: "driver-1"},
						{Field: "objects"},
						{ArrayKey: "id", ArrayElem: float64(1)},
					},
					Value: map[string]any{"id": float64(1), "addr": "foo2"},
				},
				{
					Path: []PathSegment{
						{Field: "drivers"},
						{ArrayKey: "name", ArrayElem: "driver-2"},
						{Field: "settings"},
					},
					Value: map[string]any{"mode": "off"},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			patches, err := Diff(tc.a, tc.b, tc.schema)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expect, patches); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestApplyPatch(t *testing.T) {
	base := map[string]any{
		"foo": "bar",
		"baz": []any{1, 2, 3},
		"qux": map[string]any{
			"flub": map[string]any{
				"garply": "waldo",
			},
			"objects": []any{
				map[string]any{"name": "foo", "address": 123},
				map[string]any{"name": "bar", "address": 456},
			},
		},
	}

	type testCase struct {
		patch  Patch
		expect any
		err    error
	}

	cases := map[string]testCase{
		"ReplaceTopLevel": {
			patch: Patch{
				Path: nil,
				Value: map[string]any{
					"replaced": "value",
				},
			},
			expect: map[string]any{
				"replaced": "value",
			},
		},
		"ReplaceScalar": {
			patch: Patch{
				Path:  []PathSegment{{Field: "foo"}},
				Value: "replaced",
			},
			expect: map[string]any{
				"foo": "replaced",
				"baz": []any{1, 2, 3},
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 123},
						map[string]any{"name": "bar", "address": 456},
					},
				},
			},
		},
		"ReplaceTopLevelPoint": {
			patch: Patch{
				Path: nil,
				Value: map[string]any{
					"newfoo": "replaced",
					"qux":    Ignore{},
				},
			},
			expect: map[string]any{
				"newfoo": "replaced",
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 123},
						map[string]any{"name": "bar", "address": 456},
					},
				},
			},
		},
		"ReplaceArrayElem": {
			patch: Patch{
				Path:  []PathSegment{{Field: "qux"}, {Field: "objects"}, {ArrayKey: "name", ArrayElem: "foo"}},
				Value: map[string]any{"name": "foo", "address": 789},
			},
			expect: map[string]any{
				"foo": "bar",
				"baz": []any{1, 2, 3},
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 789},
						map[string]any{"name": "bar", "address": 456},
					},
				},
			},
		},
		"AddArrayElem": {
			patch: Patch{
				Path:  []PathSegment{{Field: "qux"}, {Field: "objects"}, {ArrayKey: "name", ArrayElem: "baz"}},
				Value: map[string]any{"name": "baz", "address": 789},
			},
			expect: map[string]any{
				"foo": "bar",
				"baz": []any{1, 2, 3},
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 123},
						map[string]any{"name": "bar", "address": 456},
						map[string]any{"name": "baz", "address": 789},
					},
				},
			},
		},
		"AddField": {
			patch: Patch{
				Path:  []PathSegment{{Field: "newfield"}},
				Value: "newfieldvalue",
			},
			expect: map[string]any{
				"foo": "bar",
				"baz": []any{1, 2, 3},
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 123},
						map[string]any{"name": "bar", "address": 456},
					},
				},
				"newfield": "newfieldvalue",
			},
		},
		"AddNestedField": {
			patch: Patch{
				Path:  []PathSegment{{Field: "newfield1"}, {Field: "newfield2"}},
				Value: "newfieldvalue",
			},
			expect: map[string]any{
				"foo": "bar",
				"baz": []any{1, 2, 3},
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 123},
						map[string]any{"name": "bar", "address": 456},
					},
				},
				"newfield1": map[string]any{
					"newfield2": "newfieldvalue",
				},
			},
		},
		"DeleteField": {
			patch: Patch{
				Path:    []PathSegment{{Field: "foo"}},
				Deleted: true,
			},
			expect: map[string]any{
				"baz": []any{1, 2, 3},
				"qux": map[string]any{
					"flub": map[string]any{
						"garply": "waldo",
					},
					"objects": []any{
						map[string]any{"name": "foo", "address": 123},
						map[string]any{"name": "bar", "address": 456},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			dst := clone(base)
			dst, err := ApplyPatches(dst, []Patch{tc.patch})
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}
			if diff := cmp.Diff(tc.expect, dst); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatch_UnmarshalJSON(t *testing.T) {
	original := `
{
	"path": ["a", {"id": "theid"}],
	"value": {
		"foo": "foo",
	   	"bar": {
	 		"barbar": {"$split": "ignore"}
		}
	}
}
`
	expect := Patch{
		Path: []PathSegment{
			{Field: "a"},
			{ArrayKey: "id", ArrayElem: "theid"},
		},
		Value: map[string]any{
			"foo": "foo",
			"bar": map[string]any{
				"barbar": Ignore{},
			},
		},
	}

	var decoded Patch
	err := json.Unmarshal([]byte(original), &decoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(expect, decoded); diff != "" {
		t.Errorf("unexpected decode (-want +got):\n%s", diff)
	}

}

// tests that applying all the diffs between a and b to a results in b
func TestApplyPatch_Consistency(t *testing.T) {
	type testCase struct {
		a, b   any
		splits [][]Split // allows us to try multiple splits and check they all work
	}
	cases := map[string]testCase{
		"Primitive": {
			a:      "foo",
			b:      "bar",
			splits: nil,
		},
		"Object": {
			a: map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
			b: map[string]any{
				"foofoo": "barbar",
				"baz":    "qux2",
			},
			splits: [][]Split{
				{},
				{{Path: []string{"baz"}}},
				{{Path: []string{"foo"}}},
				{{Path: []string{"foo", "baz"}}},
				{{Path: []string{"foofoo"}}},
			},
		},
		"Array": {
			a: map[string]any{
				"objects": []any{
					map[string]any{"type": "a", "name": "a-1", "addr": "foo"},
					map[string]any{"type": "b", "name": "b-1", "id": 1},
					map[string]any{"type": "b", "name": "b-2", "id": 2},
				},
			},
			b: map[string]any{
				"objects": []any{
					map[string]any{"type": "a", "name": "a-1", "addr": "oof"},
					map[string]any{"type": "a", "name": "a-2", "addr": "foo2"},
					map[string]any{"type": "b", "name": "b-2", "id": 3},
				},
			},
			splits: [][]Split{
				{},
				{{Path: []string{"objects"}}},
				{{Path: []string{"objects"}, Key: "name"}},
				{{Path: []string{"objects"}, Key: "name", SplitKey: "type", SplitsByKey: map[string][]Split{
					"a": {{Path: []string{"addr"}}},
					"b": {{Path: []string{"id"}}},
				}}},
			},
		},
	}

	// ignore the order of entries in maps
	mapLess := func(a, b any) bool {
		aMap, okA := a.(map[string]any)
		bMap, okB := b.(map[string]any)
		if !okA || !okB {
			return false
		}
		nameA, okA := aMap["name"].(string)
		nameB, okB := bMap["name"].(string)
		if okA && okB {
			return nameA < nameB
		}
		return false
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			for i, splits := range tc.splits {
				patches, err := Diff(tc.a, tc.b, splits)
				if err != nil {
					t.Errorf("splits %d unexpected diff error: %v", i, err)
				}
				dst, err := ApplyPatches(tc.a, patches)
				if err != nil {
					t.Errorf("splits %d unexpected patch error: %v", i, err)
				}
				if diff := cmp.Diff(tc.b, dst, cmpopts.SortSlices(mapLess)); diff != "" {
					t.Errorf("splits %d unexpected result (-want +got):\n%s", i, diff)
				}
			}
		})
	}
}

func clone(dst any) any {
	switch dst := dst.(type) {
	case map[string]any:
		cloned := make(map[string]any, len(dst))
		for k, v := range dst {
			cloned[k] = clone(v)
		}
		return cloned

	case []any:
		cloned := make([]any, len(dst))
		for i, v := range dst {
			cloned[i] = clone(v)
		}
		return cloned

	default:
		return dst
	}
}
