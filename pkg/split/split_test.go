package split

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
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
					Path:  []PathSegment{{Field: "foo"}, {Field: "bar"}},
					Value: "bazbaz",
				},
				{
					Path: []PathSegment{{Field: "foo"}},
					Value: map[string]any{
						"qux": "quxqux",
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(tc.expect, Diff(tc.a, tc.b, tc.schema)); diff != "" {
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
					"qux":    Point{Kind: ObjectPoint},
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
			dst, err := ApplyPatch(dst, tc.patch)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}
			if diff := cmp.Diff(tc.expect, dst); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
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
