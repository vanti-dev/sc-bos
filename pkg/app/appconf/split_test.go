package appconf

import (
	"os"
	"testing"
)

// the idea is that every driver that wants to be configurable defines a .split.json file (aka a split file)
// i.e. metadata.json would have a companion metadata.split.json file, which defines how the config can be split up
// into atomic parts which can be independently edited without affecting the rest of the file
// then if something is not defined in the .split.json file, it is not editable
// we are testing that the way we are defining how to split the config makes sense
// we are defining the structure for the *.split.json files which themselves define how the *.json file can be split

// TestLoadSplitFile sanity check that we can load the split structure from the metadata.split.json file
func TestLoadSplitFile(t *testing.T) {
	_, err := readSplits("testdata/metadata.split.json")

	if err != nil {
		t.Errorf("error reading split file: %s", err)
	}
}

// now we have read the split file, we know how to split the file into parts we want to edit
// we can now create the db & the ext cache file structures based on what we have read in
// TestCreateSplitPathsStructure tests we can create the directory structure for db & ext cache
func TestCreateSplitPathsStructure(t *testing.T) {

	t.Run("TestCreateSplitPathsStructure", func(t *testing.T) {
		splits, err := readSplits("testdata/metadata.split.json")

		if err != nil {
			t.Errorf("error reading split file: %s", err)
		}

		dbRootPath := "testdata/db"
		err = writeSplitStructure(dbRootPath, splits)

		if err != nil {
			t.Errorf("error writing split file structure: %s", err)
		}

		expectedPaths := []string{
			"testdata/db/metadata/location/floor",
			"testdata/db/metadata/product/manufacturer",
			"testdata/db/metadata/product/model",
		}

		for _, p := range expectedPaths {
			_, err := os.ReadFile(p)
			if err != nil {
				t.Errorf("error writing split file structure: %s", err)
			}
			// not sure yet about the file contents // todo need to check
		}
	})
}
