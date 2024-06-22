package appconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// the idea is that every driver that wants to be configurable defines a .split.json file (aka a split file)
// i.e. metadata.json would have a companion metadata.split.json file, which defines how the config can be split up
// into atomic parts which can be independently edited without affecting the rest of the file
// then if something is not defined in the .split.json file, it is not editable
// we are testing that the way we are defining how to split the config makes sense
// we are defining the structure for the *.split.json files which themselves define how the *.json file can be split

// read the split file, so we know how to split the file into parts we want to edit
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
			"testdata/db/metadata/more",
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

// TestJoinDbWithExt test joining the database with the ext (appconf.Config)
func TestJoinDbWithExt(t *testing.T) {

	assert := assert.New(t)
	dbRootPath := "testdata/db"
	t.Run("TestCreateSplitPathsStructure", func(t *testing.T) {

		appConfig, err := LoadLocalConfig("testdata", "metadata.json")

		if err != nil {
			t.Errorf("failed to LoadLocalConfig: %s", err)
		}

		assert.Equal("Floor 1", appConfig.Metadata.Location.Floor)
		assert.Equal("Vanti", appConfig.Metadata.Product.Manufacturer)
		assert.Equal("Smart Core BOS", appConfig.Metadata.Product.Model)
		assert.Equal("smart", appConfig.Metadata.Membership.Subsystem)
		assert.Equal("sensor", appConfig.Metadata.More["type"])
		assert.Equal("temperature", appConfig.Metadata.More["function"])

		err = mergeDbWithExtConfig(appConfig, dbRootPath)
		if err != nil {
			t.Errorf("failed to join app config & db: %s", err)
		}

		// at this point our appConfig should have been updated with the db values
		assert.Equal("New Floor", appConfig.Metadata.Location.Floor)
		assert.Equal("New Manufacturer", appConfig.Metadata.Product.Manufacturer)
		assert.Equal("New Model", appConfig.Metadata.Product.Model)
	})
}
