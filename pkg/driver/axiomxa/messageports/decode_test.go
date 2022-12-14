package messageports

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
	"time"
)

func ExampleUnmarshal() {
	var count int
	var name string
	data := []byte("23,apples")

	Unmarshal(data, &count, &name)

	fmt.Printf("%d %s", count, name)
	// Output: 23 apples
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []any
		wantErr bool
	}{
		{"empty", nil, []any{}, false},
		{"string", []byte("hello"), []any{"hello"}, false},
		{"true", []byte("true"), []any{true}, false},
		{"false", []byte("false"), []any{false}, false},
		{"int", []byte("42"), []any{42}, false},
		{"int negative", []byte("-42"), []any{-42}, false},
		{"int8", []byte("42"), []any{int8(42)}, false},
		{"int8 too big", []byte("500"), []any{int8(0)}, true},
		{"int16", []byte("42"), []any{int16(42)}, false},
		{"int32", []byte("42"), []any{int32(42)}, false},
		{"int64", []byte("42"), []any{int64(42)}, false},
		{"uint", []byte("42"), []any{uint(42)}, false},
		{"uint8", []byte("42"), []any{uint8(42)}, false},
		{"uint8 too big", []byte("500"), []any{uint8(0)}, true},
		{"uint8 negative", []byte("-42"), []any{uint8(0)}, true},
		{"uint16", []byte("42"), []any{uint16(42)}, false},
		{"uint32", []byte("42"), []any{uint32(42)}, false},
		{"uint64", []byte("42"), []any{uint64(42)}, false},
		{"float32", []byte("42"), []any{float32(42)}, false},
		{"float64", []byte("42"), []any{float64(42)}, false},
		{"default", []byte("hello"), []any{interface{}("hello")}, false},
		{"text", []byte("hello"), []any{textUnmarshaler("tm:hello")}, false},
		{"binary", []byte("hello"), []any{binaryUnmarshaler("bm:hello")}, false},
		{"too much data", []byte("hello,world"), []any{"hello"}, false},
		{"too many targets", []byte("hello"), []any{"hello", ""}, false},
		{"time", []byte("02/11/2022 12:52:30"), []any{Time{time.Date(2022, 11, 2, 12, 52, 30, 0, time.UTC)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := make([]any, len(tt.want))
			for i, v := range tt.want {
				dst[i] = reflect.New(reflect.TypeOf(v)).Interface()
			}
			if err := Unmarshal(tt.data, dst...); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			// un-pointer everything for comparison
			vals := make([]any, len(dst))
			for i, vp := range dst {
				vals[i] = reflect.ValueOf(vp).Elem().Interface()
			}
			if diff := cmp.Diff(tt.want, vals); diff != "" {
				t.Fatalf("Incorrect unmarshal (-want,+got)\n%s", diff)
			}
		})
	}
}

type textUnmarshaler string

func (t *textUnmarshaler) UnmarshalText(text []byte) error {
	*t = textUnmarshaler("tm:" + string(text))
	return nil
}

type binaryUnmarshaler string

func (t *binaryUnmarshaler) UnmarshalBinary(text []byte) error {
	*t = binaryUnmarshaler("bm:" + string(text))
	return nil
}
