// Package split provides a way to compare and apply changes to nested data structures.
// It supports JSON-like data structures of map[string]any, []any and primitive types.
//
// The data structures are compared in logical sections defined by a schema of Split objects.
// Each Patch returned from Diff will replace the value of one such section.
// Array elements (which must be map[string]any) are identified by a key, rather than their index within the array.
package split

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/vanti-dev/sc-bos/pkg/util/maps"
)

// Diff finds the changes required to transform a into b, using the provided Splits to define logical sections of
// the data.
// Diffing will be performed on the JSON representation of the data, so the data must be JSON-serializable.
// The returned patches will be non-conflicting, meaning that they can be applied in any order and will produce
// the same result.
func Diff(a, b any, splits []Split) ([]Patch, error) {
	var err error
	a, err = convertToWorking(a)
	if err != nil {
		return nil, err
	}
	b, err = convertToWorking(b)
	if err != nil {
		return nil, err
	}

	patches := diff(a, b, Split{Splits: splits})
	// sort patches to create a deterministic output
	slices.SortStableFunc(patches, func(i, j Patch) int {
		return slices.CompareFunc(i.Path, j.Path, comparePathSegments)
	})
	return patches, nil
}

// converts an arbitrary Go value to a JSON-like tree of map[string]any, []any and primitive types
func convertToWorking(in any) (any, error) {
	serialised, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var out any
	err = json.Unmarshal(serialised, &out)
	return out, err
}

// converts from the working representation back to a Go JSON-deserializable value
func convertFromWorking[T any](working any) (T, error) {
	var out T
	serialised, err := json.Marshal(working)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(serialised, &out)
	return out, err
}

func comparePathSegments(a, b PathSegment) int {
	if a.IsField() && b.IsField() {
		return strings.Compare(a.Field, b.Field)
	} else if a.IsArrayElem() && b.IsArrayElem() {
		if c := strings.Compare(a.ArrayKey, b.ArrayKey); c != 0 {
			return c
		} else {
			return compareAny(a.ArrayElem, b.ArrayElem)
		}
	} else if a.IsField() {
		return -1
	} else if b.IsField() {
		return 1
	} else {
		return 0
	}
}

// compare strings before all other types, which are ordered according to their string representation
func compareAny(a, b any) int {
	aStr, okA := a.(string)
	bStr, okB := b.(string)
	if okA && okB {
		return strings.Compare(aStr, bStr)
	} else if okA && !okB {
		return -1
	} else if !okA && okB {
		return 1
	} else {
		// non-strings are compared by their string representation
		return strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
	}
}

func diff(a, b any, split Split) []Patch {
	if equal(a, b) {
		return nil
	}

	switch a := a.(type) {
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
			return []Patch{{Value: b}}
		}
		return diffMap(a, b, split.Splits)

	case []any:
		b, ok := b.([]any)
		if !ok {
			return []Patch{{Value: b}}
		}
		return diffArray(a, b, split)

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
	comparePages := func(a, b page) int {
		if len(a.Path) < len(b.Path) {
			return -1
		} else if len(a.Path) > len(b.Path) {
			return 1
		}

		for i := range a.Path {
			if a.Path[i] < b.Path[i] {
				return -1
			} else if a.Path[i] > b.Path[i] {
				return 1
			}
		}
		return 0
	}
	// by sorting the pages, we can step through them in order (linear time)
	// to find pages with matching paths
	slices.SortFunc(a, comparePages)
	slices.SortFunc(b, comparePages)

	type change struct {
		Path    []string
		Deleted bool
		A, B    page
	}
	var changes []change

	for len(a) > 0 && len(b) > 0 {
		c := comparePages(a[0], b[0])
		if c < 0 {
			changes = append(changes, change{Path: a[0].Path, A: a[0], Deleted: true})
			a = a[1:]
		} else if c > 0 {
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

		subpatches := diff(c.A.Value, c.B.Value, c.B.Split)
		prefixPatches(subpatches, fieldPathSegments(c.Path))
		patches = append(patches, subpatches...)
	}

	return patches
}

func prefixPatches(patches []Patch, prefix []PathSegment) {
	for i := range patches {
		var prefixed []PathSegment
		prefixed = append(prefixed, prefix...)
		prefixed = append(prefixed, patches[i].Path...)
		patches[i].Path = prefixed
	}
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

// splits a map into its own fields, and a list of pages that represent the fields that were split
// In the returned map, fields that were split are replaced with Ignore{}
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
	markIgnored(m, fields)
	return m, pages
}

// places an Ignore{} value at every position in m corresponding to a leaf in ft
func markIgnored(m map[string]any, ft fieldTree) {
	for k, entry := range ft {
		if entry.IsLeaf() {
			m[k] = Ignore{}
		} else if submap, ok := m[k].(map[string]any); ok {
			markIgnored(submap, entry.Fields)
		}
	}
}

// represents the position of a section within the structure, and its value
type page struct {
	Path  []string // hierarchy of object field names to reach this page
	Split Split    // the split that applies to Value
	Value any
}

func diffArray(a, b []any, split Split) []Patch {
	if split.Key == "" {
		// if the array doesn't have a key, we can only compare the whole thing
		if equal(a, b) {
			return nil
		} else {
			return []Patch{{Value: b}}
		}
	}

	type entry struct {
		A, B                 map[string]any
		ASplitKey, BSplitKey string
		Splits               []Split
	}
	entries := make(map[any]entry) // we allow any comparable type as array keys
	for _, v := range a {
		key, ok := extractArrayEntryKey(v, split.Key)
		if !ok {
			continue
		}
		// if there is no SplitKey, we'll use "" to indicate that
		splitKey, ok := extractArrayEntrySplitKey(v, split.SplitKey)
		if !ok {
			splitKey = ""
		}
		entries[key] = entry{
			A:         v.(map[string]any),
			ASplitKey: splitKey,
		}
	}
	for _, v := range b {
		key, ok := extractArrayEntryKey(v, split.Key)
		if !ok {
			continue
		}
		splitKey, ok := extractArrayEntrySplitKey(v, split.SplitKey)
		if !ok {
			splitKey = ""
		}
		e := entries[key]
		e.B = v.(map[string]any)
		e.BSplitKey = splitKey
		entries[key] = e
	}

	var patches []Patch
	for k, e := range entries {
		if e.B == nil {
			// array element was deleted
			patches = append(patches, Patch{
				Path:    []PathSegment{{ArrayKey: split.Key, ArrayElem: k}},
				Deleted: true,
			})
		} else if e.ASplitKey != e.BSplitKey {
			// array element moved over to different split set, not directly comparable, replace the entire element
			// (this also covers the case where a new element was added)
			patches = append(patches, Patch{
				Path:  []PathSegment{{ArrayKey: split.Key, ArrayElem: k}},
				Value: e.B,
			})
		} else {
			splits, ok := split.SplitsByKey[e.BSplitKey]
			if !ok {
				// fall back to the global splits
				// this case is used when SplitByKey isn't used
				splits = split.Splits
			}

			subpatches := diff(e.A, e.B, Split{Splits: splits})
			prefixPatches(subpatches, []PathSegment{{ArrayKey: split.Key, ArrayElem: k}})
			patches = append(patches, subpatches...)
		}
	}

	return patches
}

// extracts the logical key from an object which is used to identify it within an array
// if a key is returned, it is guaranteed to be comparable
func extractArrayEntryKey(v any, key string) (any, bool) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}
	keyValue, ok := m[key]
	if !ok || !reflect.ValueOf(keyValue).Comparable() {
		return nil, false
	}
	return keyValue, true
}

func extractArrayEntrySplitKey(v any, splitKey string) (string, bool) {
	m, ok := v.(map[string]any)
	if !ok {
		return "", false
	}
	splitKeyValue, ok := m[splitKey]
	if !ok {
		return "", false
	}
	splitKeyStr, ok := splitKeyValue.(string)
	return splitKeyStr, ok
}

// deep equality check for supporting the following nested types:
// - map[string]any
// - []any
// - comparable types
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

// returns true if a and b have the same key set
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

// ApplyPatches applies a slice of Patch to a data structure.
// The type must be both JSON-serializable and JSON-deserializable.
// A patched copy of the data structure will be returned. The original will not be modified.
func ApplyPatches[T any](data T, patches []Patch) (T, error) {
	workingCopy, err := convertToWorking(data)
	if err != nil {
		return data, err
	}
	for _, patch := range patches {
		var err error
		workingCopy, err = applyPatch(workingCopy, patch)
		if err != nil {
			return data, err
		}
	}
	return convertFromWorking[T](workingCopy)
}

// The patch will be performed in-place if possible, and the modified data structure will be returned.
// The data structure must be a JSON-like structure of map[string]any, []any and primitive types.
func applyPatch(data any, patch Patch) (any, error) {
	if len(patch.Path) == 0 {
		if patch.Deleted {
			return nil, errors.New("cannot delete root")
		} else {
			return patchValue(data, patch.Value)
		}
	}

	segment := patch.Path[0]
	if segment.IsField() {
		m, ok := data.(map[string]any)
		if !ok {
			return data, fmt.Errorf("cannot patch %T with field", data)
		}
		if len(patch.Path) == 1 && patch.Deleted {
			delete(m, segment.Field)
			return m, nil
		}

		fieldValue, ok := m[segment.Field]
		if !ok {
			fieldValue = emptyValue(patch.Path[1:])
		}
		patched, err := applyPatch(fieldValue, Patch{
			Path:    patch.Path[1:],
			Value:   patch.Value,
			Deleted: patch.Deleted,
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
			return reflect.ValueOf(v).Comparable() && v == segment.ArrayElem
		})
		if patch.Deleted && len(patch.Path) == 1 {
			if i >= 0 {
				a = slices.Delete(a, i, i+1)
			}
			return a, nil
		}
		var existing any
		if i >= 0 {
			existing = a[i]
		}

		patched, err := applyPatch(existing, Patch{
			Path:    patch.Path[1:],
			Value:   patch.Value,
			Deleted: patch.Deleted,
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
	if _, ok := value.(Ignore); ok {
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
		if _, ok := v.(Ignore); ok {
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
	// How to reach the section to be patched from the root of the data structure
	Path []PathSegment `json:"path"`
	// The new value to replace the section with
	// If an Ignore{} value is present, the data in this position will be left unchanged
	Value any `json:"value"`
	// If true, the section will be deleted
	// e.g. if the section is a field, the field will be removed
	//      if the section is an array element, the element will be removed
	Deleted bool `json:"deleted"`
}

func (p *Patch) UnmarshalJSON(data []byte) error {
	type patchAlias Patch
	err := json.Unmarshal(data, (*patchAlias)(p))
	if err != nil {
		return err
	}
	// make sure JSON maps of the form {"$split": "ignore"} are converted to Ignore{} values
	// to match the behavior of Ignore.MarshalJSON
	convertIgnores(p.Value)
	return nil
}

// Ignore is a special value that can be used in a Patch value
// to indicate that the data in this position should be left unchanged.
type Ignore struct{}

func (i *Ignore) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"$split": "ignore"})
}

// modifies data in-place to swap any values of the form map[string]any{"$split": "ignore"} with Ignore{}
func convertIgnores(data any) any {
	switch data := data.(type) {
	case map[string]any:
		if v, ok := data["$split"]; ok && len(data) == 1 && v == "ignore" {
			return Ignore{}
		}
		for k, v := range data {
			data[k] = convertIgnores(v)
		}

	case []any:
		for i, v := range data {
			data[i] = convertIgnores(v)
		}
	default:
	}
	return data
}

type PathSegment struct {
	Field     string
	ArrayKey  string
	ArrayElem any // must be comparable
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
		return json.Marshal(map[string]any{ps.ArrayKey: ps.ArrayElem})
	} else {
		return nil, ErrInvalidPathSegment
	}
}

func (ps *PathSegment) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err == nil {
		if len(m) != 1 {
			return ErrInvalidPathSegment
		}
		for k, v := range m {
			if !reflect.ValueOf(v).Comparable() {
				return ErrInvalidPathSegment
			}
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
