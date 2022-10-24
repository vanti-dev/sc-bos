package config

import (
	"embed"
	"encoding/json"
	"io/fs"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func TestJSON(t *testing.T) {
	// We're mostly testing that json marshal/unmarshal don't error here

	fileBytes, err := fs.ReadFile(testdata, "testdata/sample.json")
	if err != nil {
		t.Fatal(err)
	}

	var root Root
	err = json.Unmarshal(fileBytes, &root)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Got root %+v", root)
	t.Logf("root.cov %+v", root.COV)
	t.Logf("root.discovert %+v", root.Discovery)

	outBytes, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Written: %s", outBytes)
}
