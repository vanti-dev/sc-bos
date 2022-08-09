package policy

import (
	"context"
	"embed"
	"io/fs"
	"strings"
	"sync"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
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

// LoadRegoCached loads a rego.PartialResult for the given query using the global RegoCompiler.
// Results are cached; the partial evaluation is performed once the first time and re-used for subsequent calls.
// If LoadRegoCached returns an error, then future calls with the same query will always return the same error.
func LoadRegoCached(query string) (rego.PartialResult, error) {
	regoCacheM.Lock()
	entry, ok := regoCache[query]
	if !ok {
		entry = &regoCacheEntry{}
		regoCache[query] = entry
	}
	regoCacheM.Unlock()

	entry.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r := rego.New(
			rego.Compiler(RegoCompiler),
			rego.Query(query),
		)
		entry.partialResult, entry.err = r.PartialResult(ctx)
	})
	return entry.partialResult, entry.err
}

type regoCacheEntry struct {
	once          sync.Once
	partialResult rego.PartialResult
	err           error
}

var (
	regoCache  = make(map[string]*regoCacheEntry)
	regoCacheM sync.Mutex
)
