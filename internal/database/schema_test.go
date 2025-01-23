package database

import (
	_ "embed"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/tools/txtar"
)

var (
	//go:embed testdata/no_migrations.txtar
	testDataNoMigrations []byte
	//go:embed testdata/not_contiguous.txtar
	testDataNotContiguous []byte
	//go:embed testdata/start_from_2.txtar
	testDataStartFrom2 []byte
	//go:embed testdata/valid.txtar
	testDataValid []byte
	//go:embed testdata/version_0.txtar
	testDataVersion0 []byte
)

func TestLoadVersionedSchema(t *testing.T) {
	type testCase struct {
		txtar     []byte
		expect    []Migration
		expectErr error
	}
	cases := map[string]testCase{
		"no migrations": {
			txtar:     testDataNoMigrations,
			expectErr: ErrNoMigrations,
		},
		"not contiguous": {
			txtar:     testDataNotContiguous,
			expectErr: ErrMissingVersion,
		},
		"start from 2": {
			txtar:     testDataStartFrom2,
			expectErr: ErrMissingVersion,
		},
		"valid": {
			txtar: testDataValid,
			expect: []Migration{
				{1, "-- dummy sql"},
				{2, "-- dummy sql"},
				{3, "-- dummy sql"},
			},
		},
		"version 0": {
			txtar:     testDataVersion0,
			expectErr: ErrVersionRange,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			source, err := txtar.FS(txtar.Parse(tc.txtar))
			if err != nil {
				t.Fatalf("txtar.FS: %v", err)
			}

			got, err := LoadVersionedSchema(source)
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
			diff := cmp.Diff(tc.expect, got.migrations,
				cmpopts.EquateEmpty(),
				cmp.Comparer(func(a, b string) bool {
					return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
				}),
			)
			if diff != "" {
				t.Errorf("unexpected schema (-want +got):\n%s", diff)
			}
		})
	}
}
