package appconf

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// the idea is that every driver that wants to be configurable defines a .split.json file (aka a split file)
// i.e. metadata.json would have a companion metadata.split.json file, which defines how the config can be split up
// into atomic parts which can be independently edited without affecting the rest of the file
// then if something is not defined in the .split.json file, it is not editable
// we are testing that the way we are defining how to split the config makes sense
// we are defining the structure for the *.split.json files which themselves define how the *.json file can be split

type MockFs struct {
	fs afero.Fs
}

func (m MockFs) mockWriteFile(name string, data []byte, perm fs.FileMode) error {
	return afero.WriteFile(m.fs, name, data, perm)
}

func (m MockFs) mockReadFile(name string) ([]byte, error) {
	return afero.ReadFile(m.fs, name)
}

func (m MockFs) mockReadDir(name string) ([]os.FileInfo, error) {
	return afero.ReadDir(m.fs, name)
}

func (m MockFs) mockMkdirAll(path string, perm os.FileMode) error {
	return m.fs.MkdirAll(path, perm)
}

// read the split file, so we know how to split the file into parts we want to edit
// we can now create the db file structure based on what we have read in
// tests we can create the directory structure for db
// test modifying the db & joining the db with the ext (appconf.Config)
// appconf.Config should then contain the edits in db & the original values
func TestBasicMetadataSplit(t *testing.T) {

	// first set up the mock filesystem, read & add the metadata & split files
	// is more readable to do it this way
	var mockFs = MockFs{fs: afero.NewMemMapFs()}
	readFile = mockFs.mockReadFile
	writeFile = mockFs.mockWriteFile
	mkdirAll = mockFs.mockMkdirAll
	readDir = mockFs.mockReadDir

	file, err := os.ReadFile("testdata/metadata.split.json")
	if err != nil {
		t.Errorf("error reading split file: %s", err)
	}
	err = writeFile("fstest.metadata.split.json", file, 0664)

	file, err = os.ReadFile("testdata/metadata.json")
	if err != nil {
		t.Errorf("error reading metadata file: %s", err)
	}
	writeFile("fstest.metadata.json", file, 0664)

	assert := assert.New(t)
	dbRootPath := filepath.Join("testdata", "db")

	t.Run("TestBasicMetadataSplit", func(t *testing.T) {

		splits, err := readSplits("fstest.metadata.split.json")

		if err != nil {
			t.Errorf("error reading split file: %s", err)
		}

		err = writeSplitStructure(dbRootPath, splits)

		if err != nil {
			t.Errorf("error writing split file structure: %s", err)
		}

		expectedPaths := []string{
			filepath.Join("testdata", "db", "metadata", "location", "floor"),
			filepath.Join("testdata", "db", "metadata", "product", "manufacturer"),
			filepath.Join("testdata", "db", "metadata", "product", "model"),
		}

		for _, p := range expectedPaths {
			_, err := readFile(p)
			if err != nil {
				t.Errorf("error writing split file structure: %s", err)
			}
			// not sure yet about the file contents // todo need to check
		}

		appConfig, err := LoadLocalConfig("", "fstest.metadata.json")

		if err != nil {
			t.Errorf("failed to LoadLocalConfig: %s", err)
		}

		assert.Equal("Floor 1", appConfig.Metadata.Location.Floor)
		assert.Equal("Vanti", appConfig.Metadata.Product.Manufacturer)
		assert.Equal("Smart Core BOS", appConfig.Metadata.Product.Model)
		assert.Equal("smart", appConfig.Metadata.Membership.Subsystem)
		assert.Equal("sensor", appConfig.Metadata.More["type"])
		assert.Equal("temperature", appConfig.Metadata.More["function"])

		// now update the db files to simulate a user edit
		mockFs.mockWriteFile(expectedPaths[0], []byte("New Floor"), 0664)
		mockFs.mockWriteFile(expectedPaths[1], []byte("New Manufacturer"), 0664)
		mockFs.mockWriteFile(expectedPaths[2], []byte("New Model"), 0664)

		err = mergeDbWithExtConfig(appConfig, dbRootPath)
		if err != nil {
			t.Errorf("failed to join app config & db: %s", err)
		}

		// at this point our appConfig should have been updated with the db values from above
		assert.Equal("New Floor", appConfig.Metadata.Location.Floor)
		assert.Equal("New Manufacturer", appConfig.Metadata.Product.Manufacturer)
		assert.Equal("New Model", appConfig.Metadata.Product.Model)
	})
}
