package jsonapi

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestKeepUnknown(t *testing.T) {
	t.Run("ptr", func(t *testing.T) {
		unknown := KeepUnknown[*keepUnknownType]{}
		in := `{
	"knownInt": 1,
	"knownString": "foo",
	"unknownInt": 4,
	"unknownFloat": 1.2,
	"unknownString": "bar",
	"unknownObj": {"foo": "bar"},
	"unknownArr": ["foo", "bar"]
}`

		err := json.Unmarshal([]byte(in), &unknown)
		if err != nil {
			t.Fatal(err)
		}

		want := KeepUnknown[*keepUnknownType]{
			Known: &keepUnknownType{
				KnownInt:    1,
				KnownString: "foo",
			},
			unknown: map[string]json.RawMessage{
				"unknownInt":    json.RawMessage(`4`),
				"unknownFloat":  json.RawMessage(`1.2`),
				"unknownString": json.RawMessage(`"bar"`),
				"unknownObj":    json.RawMessage(`{"foo": "bar"}`),
				"unknownArr":    json.RawMessage(`["foo", "bar"]`),
			},
		}
		if diff := cmp.Diff(want.Known, unknown.Known); diff != "" {
			t.Errorf("Unmarshal known (-want,+got)\n%s", diff)
		}
		if diff := cmp.Diff(want.unknown, unknown.unknown); diff != "" {
			t.Errorf("Unmarshal unknown (-want,+got)\n%s", diff)
		}

		gotJson, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		wantMap, gotMap := make(map[string]any), make(map[string]any)
		if err := json.Unmarshal([]byte(in), &wantMap); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(gotJson, &gotMap); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(wantMap, gotMap); diff != "" {
			t.Errorf("Marshal (-want,+got)\n%s", diff)
		}
	})
	t.Run("value", func(t *testing.T) {
		unknown := KeepUnknown[keepUnknownType]{}
		in := `{
	"knownInt": 1,
	"knownString": "foo",
	"unknownInt": 4,
	"unknownFloat": 1.2,
	"unknownString": "bar",
	"unknownObj": {"foo": "bar"},
	"unknownArr": ["foo", "bar"]
}`

		err := json.Unmarshal([]byte(in), &unknown)
		if err != nil {
			t.Fatal(err)
		}

		want := KeepUnknown[keepUnknownType]{
			Known: keepUnknownType{
				KnownInt:    1,
				KnownString: "foo",
			},
			unknown: map[string]json.RawMessage{
				"unknownInt":    json.RawMessage(`4`),
				"unknownFloat":  json.RawMessage(`1.2`),
				"unknownString": json.RawMessage(`"bar"`),
				"unknownObj":    json.RawMessage(`{"foo": "bar"}`),
				"unknownArr":    json.RawMessage(`["foo", "bar"]`),
			},
		}
		if diff := cmp.Diff(want.Known, unknown.Known); diff != "" {
			t.Errorf("Unmarshal known (-want,+got)\n%s", diff)
		}
		if diff := cmp.Diff(want.unknown, unknown.unknown); diff != "" {
			t.Errorf("Unmarshal unknown (-want,+got)\n%s", diff)
		}

		gotJson, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		wantMap, gotMap := make(map[string]any), make(map[string]any)
		if err := json.Unmarshal([]byte(in), &wantMap); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(gotJson, &gotMap); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(wantMap, gotMap); diff != "" {
			t.Errorf("Marshal (-want,+got)\n%s", diff)
		}
	})
}

type keepUnknownType struct {
	KnownInt    int    `json:"knownInt,omitempty"`
	KnownString string `json:"knownString,omitempty"`
}
