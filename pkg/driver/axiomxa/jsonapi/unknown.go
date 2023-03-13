package jsonapi

import (
	"encoding/json"
	"reflect"
	"strings"
)

// KeepUnknown can be used to keep unknown fields in a json Unmarshal->Marshal flow.
// Known and unknown fields are marshalled in sorted order, as if both known and unknown keys were stored in the same map.
type KeepUnknown[T any] struct {
	Known   T
	unknown map[string]json.RawMessage
}

func (k *KeepUnknown[T]) UnmarshalJSON(bytes []byte) error {
	var known T
	var unknown map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &known); err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &unknown); err != nil {
		return err
	}

	// delete all known tags out of the unknown map
	rt := reflect.TypeOf(known)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if excludedJSONField(f) {
			continue
		}
		if tag, ok := f.Tag.Lookup("json"); ok {
			if tag == "-" {
				continue
			}
			name, _, _ := strings.Cut(tag, ",")
			delete(unknown, name)
		}
	}

	*k = KeepUnknown[T]{Known: known, unknown: unknown}
	return nil
}

func (k KeepUnknown[T]) MarshalJSON() ([]byte, error) {
	t := reflect.TypeOf(k.Known)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	v := reflect.ValueOf(k.Known)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	out := make(map[string]any, t.NumField()+len(k.unknown))
	for k, v := range k.unknown {
		out[k] = v
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if excludedJSONField(f) {
			continue
		}
		if tag, ok := f.Tag.Lookup("json"); ok {
			if tag == "-" {
				continue
			}
			name, _, _ := strings.Cut(tag, ",")
			// we're ignoring omitempty for now
			out[name] = v.Field(i).Interface()
		}
	}
	return json.Marshal(out)
}

func excludedJSONField(sf reflect.StructField) bool {
	// snippet copied from json/encode.go
	if sf.Anonymous {
		t := sf.Type
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		if !sf.IsExported() && t.Kind() != reflect.Struct {
			// Ignore embedded fields of unexported non-struct types.
			return true
		}
		// Do not ignore embedded fields of unexported struct types
		// since they may have exported fields.
	} else if !sf.IsExported() {
		// Ignore unexported non-embedded fields.
		return true
	}
	return false
}
