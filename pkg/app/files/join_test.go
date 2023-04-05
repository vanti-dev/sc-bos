package files

import (
	"testing"
)

func TestPath(t *testing.T) {
	tests := []struct {
		name    string
		dataDir string
		path    string
		want    string
	}{
		{"empty", "", "", ""},
		{"empty dataDir", "", "foo", "foo"},
		{"empty path", "bar", "", ""},
		{"absolute path", "bar", "/foo", "/foo"},
		{"relative path", "bar", "foo", "bar/foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Path(tt.dataDir, tt.path); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}
