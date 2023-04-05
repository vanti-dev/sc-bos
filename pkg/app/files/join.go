package files

import (
	"path/filepath"
	"strings"
)

// Path adds dataDir as a prefix to path only if path is not empty and path does not start with "/".
func Path(dataDir, path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return joinIfPresent(dataDir, path)
}

func joinIfPresent(dir, path string) string {
	if path == "" {
		return ""
	}
	return filepath.Join(dir, path)
}
