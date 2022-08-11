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
// If the provided context is cancelled before the result is ready, the process will continue in the background and
// the context error is returned.
// If LoadRegoCached returns a non-context error, then future calls with the same query will always return the same error.
func LoadRegoCached(ctx context.Context, query string) (rego.PartialResult, error) {
	regoCacheM.Lock()
	entry, ok := regoCache[query]
	if !ok {
		entry = &regoCacheEntry{done: make(chan struct{})}
		regoCache[query] = entry
	}
	regoCacheM.Unlock()

	// each cache entry only gets one change to compile - it's a deterministic process, so if it fails once there's no
	// point trying again later
	entry.once.Do(func() {
		// run asynchronously so the compilation can complete in the background if ctx is cancelled early
		go func() {
			defer close(entry.done)
			// if a policy file takes more than 5 seconds to compile, something is wrong with it
			bgctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			r := rego.New(
				rego.Compiler(RegoCompiler),
				rego.Query(query),
			)
			entry.partialResult, entry.err = r.PartialResult(bgctx)
		}()
	})
	select {
	case <-entry.done:
		return entry.partialResult, entry.err
	case <-ctx.Done():
		return rego.PartialResult{}, ctx.Err()
	}
}

type regoCacheEntry struct {
	once          sync.Once
	done          chan struct{}
	partialResult rego.PartialResult
	err           error
}

var (
	regoCache  = make(map[string]*regoCacheEntry)
	regoCacheM sync.Mutex
)
