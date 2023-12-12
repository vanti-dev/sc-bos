package jsontypes

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// String is either a literal string or a string loaded from a file.
// If the contents of String are an absolute path or start with a dot, the contents are loaded from a file.
// Otherwise, the contents are used as-is.
type String string

// IsPath returns whether s looks like a filesystem path or not.
// Filesystem paths are either absolute paths or paths starting with a dot.
// Absolute paths are defined by [filepath.IsAbs].
func (s String) IsPath() bool {
	return strings.HasPrefix(string(s), ".") || filepath.IsAbs(string(s))
}

// Open returns a reader for the contents of s, the file or string.
func (s String) Open() (io.ReadCloser, error) {
	if s.IsPath() {
		return os.Open(string(s))
	}
	return io.NopCloser(strings.NewReader(string(s))), nil
}

// OpenBase returns a reader for the contents of s, the file or string.
// Paths are resolved relative to base.
func (s String) OpenBase(base string) (io.ReadCloser, error) {
	if s.IsPath() {
		var p = string(s)
		if !filepath.IsAbs(p) {
			p = filepath.Join(base, p)
		}
		return os.Open(p)
	}
	return io.NopCloser(strings.NewReader(string(s))), nil
}

// Read reads the contents of s, the file or string.
func (s String) Read() (string, error) {
	if !s.IsPath() {
		return string(s), nil
	}
	f, err := os.ReadFile(string(s))
	return string(f), err
}
