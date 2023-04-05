package serviceapi

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/natefinch/atomic"
)

type Store interface {
	Save(ctx context.Context, path string, data []byte) error
}

type StoreDir string

func (d StoreDir) Save(_ context.Context, path string, data []byte) error {
	if err := os.MkdirAll(filepath.Join(string(d), filepath.Dir(path)), 0755); err != nil {
		return err
	}
	return atomic.WriteFile(filepath.Join(string(d), path), bytes.NewReader(data))
}

type Marshaller interface {
	MarshalConfig([]byte) ([]byte, error)
}

type marshallerFunc func([]byte) ([]byte, error)

func (f marshallerFunc) MarshalConfig(b []byte) ([]byte, error) {
	return f(b)
}

func MarshalArrayConfig(prop string) Marshaller {
	return marshallerFunc(func(b []byte) ([]byte, error) {
		out := map[string]any{
			prop: []any{
				json.RawMessage(b),
			},
		}
		return json.Marshal(out)
	})
}

func MarshalMapConfig(prop, key string) Marshaller {
	return marshallerFunc(func(b []byte) ([]byte, error) {
		out := map[string]any{
			prop: map[string]any{
				key: json.RawMessage(b),
			},
		}
		return json.Marshal(out)
	})
}
