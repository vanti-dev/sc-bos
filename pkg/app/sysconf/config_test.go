package sysconf

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
)

func TestLoad(t *testing.T) {
	t.Run("args override json", func(t *testing.T) {
		oldArgs := os.Args
		t.Cleanup(func() {
			os.Args = oldArgs
		})
		os.Args = []string{
			"app",
			"--sysconf", "testdata/conf.json",
			"--listen-grpc", ":2222",
		}
		dst := Default()
		if err := Load(&dst); err != nil {
			t.Fatal(err)
		}

		want := Default()
		want.ConfigDirs = []string{"testdata"}
		want.ConfigFiles = []string{"conf.json"}
		want.ListenGRPC = ":2222"
		want.Normalize()

		if diff := cmp.Diff(want, dst, cmpopts.IgnoreTypes(zap.Config{}, Certs{})); diff != "" {
			t.Errorf("Load() mismatch (-want +got):\n%s", diff)
		}
	})
}
