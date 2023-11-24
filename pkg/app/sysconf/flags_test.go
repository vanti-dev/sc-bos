package sysconf

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadFromArgs_sysConfArg(t *testing.T) {
	t.Run("no flags", func(t *testing.T) {
		dst := &Config{
			ConfigDirs:  []string{"foo"},
			ConfigFiles: []string{"test.json"},
		}

		_, err := LoadFromArgs(dst)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff([]string{"foo"}, dst.ConfigDirs); diff != "" {
			t.Errorf("ConfigDirs mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]string{"test.json"}, dst.ConfigFiles); diff != "" {
			t.Errorf("ConfigFiles mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("overwrite dirs", func(t *testing.T) {
		dst := &Config{
			ConfigDirs:  []string{"foo"},
			ConfigFiles: []string{"test.json"},
		}

		_, err := LoadFromArgs(dst, "--sysconf", "new/conf.json")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff([]string{"new"}, dst.ConfigDirs); diff != "" {
			t.Errorf("ConfigDirs mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]string{"conf.json"}, dst.ConfigFiles); diff != "" {
			t.Errorf("ConfigFiles mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("comma dirs", func(t *testing.T) {
		dst := &Config{
			ConfigDirs:  []string{"foo"},
			ConfigFiles: []string{"test.json"},
		}

		_, err := LoadFromArgs(dst, "--sysconf", "new/conf.json,other/conf.json")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff([]string{"new", "other"}, dst.ConfigDirs); diff != "" {
			t.Errorf("ConfigDirs mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]string{"conf.json"}, dst.ConfigFiles); diff != "" {
			t.Errorf("ConfigFiles mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("multiple dirs", func(t *testing.T) {
		dst := &Config{
			ConfigDirs:  []string{"foo"},
			ConfigFiles: []string{"test.json"},
		}

		_, err := LoadFromArgs(dst, "--sysconf", "new/conf.json", "--sysconf", "other/conf.json")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff([]string{"new", "other"}, dst.ConfigDirs); diff != "" {
			t.Errorf("ConfigDirs mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff([]string{"conf.json"}, dst.ConfigFiles); diff != "" {
			t.Errorf("ConfigFiles mismatch (-want +got):\n%s", diff)
		}
	})
}
