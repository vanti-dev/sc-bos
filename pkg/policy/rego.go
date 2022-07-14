package policy

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/open-policy-agent/opa/ast"
)

//go:embed rego
var regoSources embed.FS

var RegoCompiler *ast.Compiler

func init() {
	sources := make(map[string]string)
	err := fs.WalkDir(regoSources, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".rego") {
			return nil
		}

		contents, err := regoSources.ReadFile(path)
		if err != nil {
			return err
		}

		sources[path] = string(contents)
		return nil
	})
	if err != nil {
		panic(err)
	}

	RegoCompiler = ast.MustCompileModules(sources)
}
