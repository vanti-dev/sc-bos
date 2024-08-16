package split

import (
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/vanti-dev/sc-bos/pkg/util/maps"
)

func Diff(a, b any, schema []Split) []Patch {
	if equal(a, b) {
		return nil
	}

	switch a := a.(type) {
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
			return []Patch{{Value: b}}
		}
		return diffMap(a, b, schema)

	case []any:
		b, ok := b.([]any)
		if !ok {
			return []Patch{{Value: b}}
		}
		return diffArray(a, b, schema)

	default:
		if equal(a, b) {
			return nil
		} else {
			return []Patch{{Value: b}}
		}
	}
}

func diffMap(a, b map[string]any, schema []Split) []Patch {
	tree := buildFieldTree(schema)
	a, aPages := splitMap(a, tree)
	b, bPages := splitMap(b, tree)

	patches := diffPages(aPages, bPages)
	if !equal(a, b) {
		patches = append(patches, Patch{Value: b})
	}

	return patches
}

func diffPages(a, b []page) []Patch {
	pageLess := func(a, b page) bool {
		if len(a.Path) < len(b.Path) {
			return true
		} else if len(a.Path) > len(b.Path) {
			return false
		}

		for i := range a.Path {
			if a.Path[i] < b.Path[i] {
				return true
			} else if a.Path[i] > b.Path[i] {
				return false
			}
		}
		return false
	}
	// by sorting the pages, we can step through them in order (linear time)
	// to find pages with matching paths
	slices.SortFunc(a, pageLess)
	slices.SortFunc(b, pageLess)

	type change struct {
		Path    []string
		Deleted bool
		A, B    page
	}
	var changes []change

	for len(a) > 0 && len(b) > 0 {
		if pageLess(a[0], b[0]) {
			changes = append(changes, change{Path: a[0].Path, A: a[0], Deleted: true})
			a = a[1:]
		} else if pageLess(b[0], a[0]) {
			changes = append(changes, change{Path: b[0].Path, B: b[0]})
			b = b[1:]
		} else {
			if !equal(a[0].Value, b[0].Value) {
				changes = append(changes, change{Path: a[0].Path, A: a[0], B: b[0]})
			}
			a = a[1:]
			b = b[1:]
		}
	}
	for _, p := range a {
		changes = append(changes, change{Path: p.Path, A: p, Deleted: true})
	}
	for _, p := range b {
		changes = append(changes, change{Path: p.Path, B: p})
	}

	var patches []Patch
	for _, c := range changes {
		if c.Deleted {
			patches = append(patches, Patch{Path: fieldPathSegments(c.Path), Deleted: true})
			continue
		}

		subpatches := Diff(c.A.Value, c.B.Value, c.B.Split.Splits)
		for i := range subpatches {
			// prepend the current field path to the subpatch path so that paths are absolute
			subpatches[i].Path = append(fieldPathSegments(c.Path), subpatches[i].Path...)
		}
		patches = append(patches, subpatches...)
	}

	return patches
}

func fieldPathSegments(path []string) []PathSegment {
	var segs []PathSegment
	for _, p := range path {
		segs = append(segs, PathSegment{Field: p})
	}
	return segs
}

type fieldTree map[string]fieldTreeEntry

// fieldTreeEntry contains either a Split or a fieldTree, not both
type fieldTreeEntry struct {
	Split  Split
	Fields fieldTree
}

func (e fieldTreeEntry) IsLeaf() bool {
	return len(e.Fields) == 0
}

func buildFieldTree(schema []Split) fieldTree {
	tree := make(fieldTree)
	for _, split := range schema {
		node := tree
		path := split.Path
		for len(path) > 0 {
			key := path[0]
			if len(path) == 1 {
				node[key] = fieldTreeEntry{Split: split}
			} else {
				if _, ok := node[key]; !ok {
					node[key] = fieldTreeEntry{Fields: make(fieldTree)}
				}
				node = node[key].Fields
			}
			path = path[1:]
		}
	}
	return tree
}

func splitMap(m map[string]any, fields fieldTree) (map[string]any, []page) {
	if len(fields) == 0 {
		return m, nil
	}
	m = maps.Clone(m)

	var pages []page
	for k, v := range m {
		subfields, ok := fields[k]
		if !ok {
			// don't need to delete this field or any of its children
			continue
		}
		if subfields.IsLeaf() {
			// leaf of the fieldTree, need to delete this key
			delete(m, k)
			pages = append(pages, page{Path: []string{k}, Value: v, Split: subfields.Split})
		} else if submap, ok := v.(map[string]any); ok {
			submap, subpages := splitMap(submap, subfields.Fields)
			m[k] = submap
			for _, p := range subpages {
				p.Path = append([]string{k}, p.Path...)
				pages = append(pages, p)
			}
		} else {
			// not a map, don't need to delete this field or any of its children
		}
	}
	return m, pages
}

type page struct {
	Path  []string
	Split Split
	Value any
}

func diffArray(a, b []any, schema []Split) []Patch {
	panic("not implemented")
}

func equal(a, b any) bool {
	switch a := a.(type) {
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
			return false
		}
		if !equalKeys(a, b) {
			return false
		}
		for k, v := range a {
			if !equal(v, b[k]) {
				return false
			}
		}
		return true

	case []any:
		b, ok := b.([]any)
		if !ok {
			return false
		}
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if !equal(a[i], b[i]) {
				return false
			}
		}
		return true

	default:
		return a == b
	}
}

func equalKeys(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			return false
		}
	}
	return true
}

func ApplyPatch(data any, diff Patch) (any, error) {
	if len(diff.Path) == 0 {
		return patchValue(data, diff.Value)
	}

	segment := diff.Path[0]
	if segment.IsField() {
		m, ok := data.(map[string]any)
		if !ok {
			return data, fmt.Errorf("cannot patch %T with field", data)
		}
		if diff.Deleted {
			delete(m, segment.Field)
			return m, nil
		}

		fieldValue, ok := m[segment.Field]
		if !ok {
			fieldValue = emptyValue(diff.Path[1:])
		}
		patched, err := ApplyPatch(fieldValue, Patch{
			Path:  diff.Path[1:],
			Value: diff.Value,
		})
		if err != nil {
			return data, err
		}
		m[segment.Field] = patched
		return m, nil
	} else if segment.IsArrayElem() {
		a, ok := data.([]any)
		if !ok {
			return data, fmt.Errorf("cannot patch %T with array element", data)
		}

		i := slices.IndexFunc(a, func(a any) bool {
			m, ok := a.(map[string]any)
			if !ok {
				return false
			}
			v, ok := m[segment.ArrayKey]
			if !ok {
				return false
			}
			return v == segment.ArrayElem
		})
		var existing any
		if i >= 0 {
			existing = a[i]
		}

		patched, err := ApplyPatch(existing, Patch{
			Path:  diff.Path[1:],
			Value: diff.Value,
		})
		if err != nil {
			return data, err
		}
		if i >= 0 {
			a[i] = patched
		} else {
			a = append(a, patched)
		}
		return a, nil

	} else {
		return data, ErrInvalidPathSegment
	}
}

func patchValue(dst any, value any) (any, error) {
	if dst == nil {
		return value, nil
	}
	if _, ok := value.(Point); ok {
		// points represent the limit of the patch
		// nothing to apply
		return dst, nil
	}

	// replace the top-level value
	switch dst := dst.(type) {
	case map[string]any:
		if patch, ok := value.(map[string]any); ok {
			return patchMap(dst, patch)
		} else {
			return dst, fmt.Errorf("cannot patch map with %T", value)
		}

	case []any:
		if patch, ok := value.([]any); ok {
			dst = dst[:0]
			dst = append(dst, patch...)
			return dst, nil
		} else {
			return dst, fmt.Errorf("cannot patch array with %T", value)
		}

	default:
		return value, nil
	}
}

func patchMap(m map[string]any, patch map[string]any) (map[string]any, error) {
	patched := make(map[string]any)
	for k, v := range patch {
		if _, ok := v.(Point); ok {
			// we need to preserve the original value, if it exists
			if original, ok := m[k]; ok {
				patched[k] = original
			}
		} else {
			patchedField, err := patchValue(m[k], v)
			if err != nil {
				return m, err
			}
			patched[k] = patchedField
		}
	}

	return patched, nil
}

// creates a new empty value that can be indexed by the given path
func emptyValue(segs []PathSegment) any {
	if len(segs) == 0 {
		return nil
	}

	seg := segs[0]
	if seg.IsField() {
		return map[string]any{}
	} else if seg.IsArrayElem() {
		return []any{}
	} else {
		panic("invalid path segment")
	}
}

type Patch struct {
	Path    []PathSegment `json:"path"`
	Value   any           `json:"value"`
	Deleted bool          `json:"deleted"`
}

type Point struct {
	Kind PointKind
}

func (p *Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"$split": string(p.Kind)})
}

type PointKind string

const (
	ObjectPoint PointKind = "object"
	ArrayPoint  PointKind = "array"
)

type PathSegment struct {
	Field     string
	ArrayKey  string
	ArrayElem string
}

func (ps *PathSegment) IsField() bool {
	return ps.Field != "" && ps.ArrayKey == ""
}

func (ps *PathSegment) IsArrayElem() bool {
	return ps.Field == "" && ps.ArrayKey != ""
}

func (ps *PathSegment) MarshalJSON() ([]byte, error) {
	if ps.IsField() {
		return json.Marshal(ps.Field)
	} else if ps.IsArrayElem() {
		return json.Marshal(map[string]string{ps.ArrayKey: ps.ArrayElem})
	} else {
		return nil, ErrInvalidPathSegment
	}
}

func (ps *PathSegment) UnmarshalJSON(data []byte) error {
	var m map[string]string
	if err := json.Unmarshal(data, &m); err == nil {
		if len(m) != 1 {
			return ErrInvalidPathSegment
		}
		for k, v := range m {
			ps.ArrayKey = k
			ps.ArrayElem = v
		}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		ps.Field = s
		return nil
	}

	return ErrInvalidPathSegment
}

type Split struct {
	SplitKey    string             `json:"splitKey"`
	Key         string             `json:"key"`
	Path        []string           `json:"path"`
	Splits      []Split            `json:"splits"`
	SplitsByKey map[string][]Split `json:"splitsByKey"`
}

var ErrInvalidPathSegment = errors.New("invalid path segment")
