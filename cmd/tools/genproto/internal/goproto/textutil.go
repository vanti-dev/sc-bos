package goproto

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/ast/astutil"
)

// TransformOption represents a bitset of transformations to apply.
type TransformOption uint8

const (
	// RemoveGenImports removes import lines for the gen package.
	RemoveGenImports TransformOption = 1 << iota
	// RemoveGenQualifiers removes "gen." qualifiers from identifiers.
	RemoveGenQualifiers
	// RenamePackageToGen changes package declaration from "genpb" to "gen".
	RenamePackageToGen
)

// genQualifierPattern matches "gen." followed by an identifier (uppercase letter followed by word characters).
// This pattern is used to remove gen. prefixes from comments and string literals.
var genQualifierPattern = regexp.MustCompile(`\bgen\.([A-Z]\w*)`)

// TransformGoFile parses the Go file once and applies all requested transformations.
func TransformGoFile(content []byte, opts TransformOption) []byte {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		// Fall back to returning original content if parsing fails
		return content
	}

	// Apply transformations
	if opts&RemoveGenImports != 0 {
		astutil.DeleteNamedImport(fset, f, "gen", "github.com/smart-core-os/sc-bos/pkg/gen")
	}

	if opts&RemoveGenQualifiers != 0 {
		// Process comments to remove gen. prefixes
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				c.Text = genQualifierPattern.ReplaceAllString(c.Text, "$1")
			}
		}

		// Process AST nodes
		astutil.Apply(f, func(c *astutil.Cursor) bool {
			n := c.Node()

			// Look for selector expressions like "gen.SomeType"
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "gen" {
					// Replace "gen.SomeType" with just "SomeType"
					c.Replace(&ast.Ident{
						NamePos: sel.Sel.NamePos,
						Name:    sel.Sel.Name,
						Obj:     sel.Sel.Obj,
					})
				}
			}

			// Look for string literals and remove gen. prefixes from their content
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				lit.Value = genQualifierPattern.ReplaceAllString(lit.Value, "$1")
			}

			return true
		}, nil)
	}

	if opts&RenamePackageToGen != 0 && f.Name.Name == "genpb" {
		f.Name.Name = "gen"
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return content
	}

	return buf.Bytes()
}
