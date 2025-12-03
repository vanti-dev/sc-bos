package goproto

import (
	"bytes"
	"path/filepath"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestTransformGoFile(t *testing.T) {
	tests := []struct {
		name  string
		txtar string
		opts  TransformOption
	}{
		// RemoveGenImports tests
		{
			name:  "RemoveGenImports removes gen import from multi-line import block",
			txtar: "remove_gen_imports_multi_line.txtar",
			opts:  RemoveGenImports,
		},
		{
			name:  "RemoveGenImports handles single line import",
			txtar: "remove_gen_imports_single_line.txtar",
			opts:  RemoveGenImports,
		},
		{
			name:  "RemoveGenImports preserves other imports",
			txtar: "remove_gen_imports_preserves_others.txtar",
			opts:  RemoveGenImports,
		},
		// RemoveGenQualifiers tests
		{
			name:  "RemoveGenQualifiers removes gen. prefix from types",
			txtar: "remove_gen_qualifiers_types.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers removes multiple gen. prefixes",
			txtar: "remove_gen_qualifiers_multiple.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers preserves non-gen identifiers",
			txtar: "remove_gen_qualifiers_preserves_non_gen.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers handles gen at start of statement",
			txtar: "remove_gen_qualifiers_start_of_statement.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers handles gen at end of statement",
			txtar: "remove_gen_qualifiers_end_of_statement.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers doesn't remove gen without dot",
			txtar: "remove_gen_qualifiers_no_dot.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers doesn't remove gen at end of word",
			txtar: "remove_gen_qualifiers_end_of_word.txtar",
			opts:  RemoveGenQualifiers,
		},
		// RenamePackageToGen tests
		{
			name:  "RenamePackageToGen renames package genpb to gen",
			txtar: "rename_package_to_gen.txtar",
			opts:  RenamePackageToGen,
		},
		{
			name:  "RenamePackageToGen handles package with trailing comment",
			txtar: "rename_package_trailing_comment.txtar",
			opts:  RenamePackageToGen,
		},
		{
			name:  "RenamePackageToGen doesn't affect other packages",
			txtar: "rename_package_no_change.txtar",
			opts:  RenamePackageToGen,
		},
		// Combined options tests
		{
			name:  "combines RemoveGenImports and RemoveGenQualifiers",
			txtar: "combined_imports_and_qualifiers.txtar",
			opts:  RemoveGenImports | RemoveGenQualifiers,
		},
		{
			name:  "combines all three options",
			txtar: "combined_all_three.txtar",
			opts:  RemoveGenImports | RemoveGenQualifiers | RenamePackageToGen,
		},
		{
			name:  "no options leaves input unchanged",
			txtar: "no_options.txtar",
			opts:  0,
		},
		{
			name:  "RemoveGenQualifiers removes gen. from comments",
			txtar: "remove_gen_qualifiers_comments.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers handles gen. in multiline comments",
			txtar: "remove_gen_qualifiers_multiline_comments.txtar",
			opts:  RemoveGenQualifiers,
		},
		{
			name:  "RemoveGenQualifiers removes gen. from string literals",
			txtar: "remove_gen_qualifiers_strings.txtar",
			opts:  RemoveGenQualifiers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load test data from txtar file
			txtarPath := filepath.Join("testdata", tt.txtar)
			archive, err := txtar.ParseFile(txtarPath)
			if err != nil {
				t.Fatalf("failed to parse txtar file %q: %v", txtarPath, err)
			}

			// Extract input and expected from archive
			var input, want []byte
			for _, file := range archive.Files {
				switch file.Name {
				case "input.go":
					input = file.Data
				case "want.go":
					want = file.Data
				}
			}

			if input == nil {
				t.Fatalf("txtar file %q missing input.go", txtarPath)
			}
			if want == nil {
				t.Fatalf("txtar file %q missing want.go", txtarPath)
			}

			// Run the transformation
			result := TransformGoFile(input, tt.opts)

			// Compare results (trim trailing whitespace since go/format adds a final newline)
			if !bytes.Equal(bytes.TrimSpace(result), bytes.TrimSpace(want)) {
				t.Errorf("TransformGoFile() =\n%s\n\nwant:\n%s", result, want)
			}
		})
	}
}
