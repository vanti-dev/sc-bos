package appconf

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
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

func (m MockFs) mockIsDir(path string) (bool, error) {
	return afero.IsDir(m.fs, path)
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
	isDir = mockFs.mockIsDir
	mockFsSplitFileName := "fstest.metadata.split.json"
	mockFsConfigFileName := "fstest.metadata.json"

	file, err := os.ReadFile("testdata/metadata.split.json")
	if err != nil {
		t.Errorf("error reading split file: %s", err)
	}
	err = writeFile(mockFsSplitFileName, file, 0664)

	file, err = os.ReadFile("testdata/metadata.json")
	if err != nil {
		t.Errorf("error reading config file: %s", err)
	}
	writeFile(mockFsConfigFileName, file, 0664)

	assert := assert.New(t)
	dbRootPath := filepath.Join("testdata", "db")

	t.Run("TestBasicMetadataSplit", func(t *testing.T) {

		splits, err := readSplits(mockFsSplitFileName)

		if err != nil {
			t.Errorf("error reading split file: %s", err)
		}

		err = writeSplitStructure(dbRootPath, splits)

		if err != nil {
			t.Errorf("error writing split file structure: %s", err)
		}

		expectedPaths := []string{
			filepath.Join("testdata", "db", "metadata", "Location", "Floor"),
			filepath.Join("testdata", "db", "metadata", "Product", "Manufacturer"),
			filepath.Join("testdata", "db", "metadata", "Product", "Model"),
			filepath.Join("testdata", "db", "metadata", "Membership"),
			filepath.Join("testdata", "db", "metadata", "Traits"),
			filepath.Join("testdata", "db", "metadata", "More"),
		}

		for _, p := range expectedPaths {
			_, err := readFile(p)
			if err != nil {
				t.Errorf("failed checking expected paths: %s %s", p, err)
			}
		}

		appConfig, err := LoadLocalConfig("", mockFsConfigFileName)

		if err != nil {
			t.Errorf("failed to LoadLocalConfig: %s", err)
		}

		assert.Equal("Floor 1", appConfig.Metadata.Location.Floor)
		assert.Equal("Vanti", appConfig.Metadata.Product.Manufacturer)
		assert.Equal("Smart Core BOS", appConfig.Metadata.Product.Model)
		assert.Equal("smart", appConfig.Metadata.Membership.Subsystem)
		assert.Equal("sensor", appConfig.Metadata.More["type"])
		assert.Equal("temperature", appConfig.Metadata.More["function"])
		assert.Equal(appConfig.Metadata.Traits[0].Name, "oldTrait")
		assert.Equal(1, len(appConfig.Metadata.Traits))
		assert.Equal(2, len(appConfig.Metadata.More))

		// now update the db files to simulate a user edit.
		// i.e. we are creating the page files that would contain the user edits to overlay onto our app config
		writePageFile(expectedPaths[0], nil, "New Floor")
		writePageFile(expectedPaths[1], nil, "New Manufacturer")
		writePageFile(expectedPaths[2], nil, "New Model")
		writePageFile(expectedPaths[3], nil, traits.Metadata_Membership{Subsystem: "New Subsystem"})
		// we want to replace the whole array here to test that we can do that
		writePageFile(expectedPaths[4], nil, []traits.TraitMetadata{{Name: "newTrait"}})
		writePageFile(expectedPaths[5], nil, map[string]string{"type": "newType", "function": "newFunction"})

		err = mergeDbWithExtConfig(appConfig, dbRootPath)
		if err != nil {
			t.Errorf("failed to join app config & db: %s", err)
		}

		// at this point our appConfig should have been updated with the db values from above
		assert.Equal("New Floor", appConfig.Metadata.Location.Floor)
		assert.Equal("New Manufacturer", appConfig.Metadata.Product.Manufacturer)
		assert.Equal("New Model", appConfig.Metadata.Product.Model)
		assert.Equal("New Subsystem", appConfig.Metadata.Membership.Subsystem)
		assert.Equal("newTrait", appConfig.Metadata.Traits[0].Name)
		assert.Equal(1, len(appConfig.Metadata.Traits))
		assert.Equal(2, len(appConfig.Metadata.More))
		assert.Equal("newType", appConfig.Metadata.More["type"])
		assert.Equal("newFunction", appConfig.Metadata.More["function"])
		//_, err = mockFs.mockIsDir("")
		//if err != nil {
		//	return
		//}
	})
}

// tests the ability of the config system to update the property of a specific device in the config
// the device is specified using the "key" attribute in the split file
// when the split-config encounters the key property being present in the split file, for a given split
// it will search through the array / map of objects in the config for the object with the matching key
// and then follow the same process
func TestDeviceSpecificBmsPage(t *testing.T) {

	// first set up the mock filesystem, read & add the metadata & split files
	// is more readable to do it this way
	var mockFs = MockFs{fs: afero.NewMemMapFs()}
	readFile = mockFs.mockReadFile
	writeFile = mockFs.mockWriteFile
	mkdirAll = mockFs.mockMkdirAll
	readDir = mockFs.mockReadDir
	isDir = mockFs.mockIsDir
	mockFsSplitFileName := "fstest.bms.split.json"
	mockFsConfigFileName := "fstest.bms.json"

	file, err := os.ReadFile("testdata/bms.split.json")
	if err != nil {
		t.Errorf("error reading split file: %s", err)
	}
	err = writeFile(mockFsSplitFileName, file, 0664)

	file, err = os.ReadFile("testdata/bms.json")
	if err != nil {
		t.Errorf("error reading config file: %s", err)
	}

	writeFile(mockFsConfigFileName, file, 0664)

	assert := assert.New(t)
	dbRootPath := filepath.Join("testdata", "db")

	t.Run("TestDeviceSpecificBmsPage", func(t *testing.T) {

		splits, err := readSplits(mockFsSplitFileName)

		if err != nil {
			t.Errorf("error reading split file: %s", err)
		}

		err = writeSplitStructure(dbRootPath, splits)

		if err != nil {
			t.Errorf("error writing split file structure: %s", err)
		}

		expectedPaths := []string{
			filepath.Join("testdata", "db", "drivers", "bacnet", "localInterface"),
			filepath.Join("testdata", "db", "drivers", "bacnet", "localPort"),
			filepath.Join("testdata", "db", "drivers", "bacnet", "devices", "title"),
			filepath.Join("testdata", "db", "drivers", "bacnet", "devices", "comm"),
			//filepath.Join("testdata", "db", "Drivers", "Devices", "metadata", "appearance", "title"),
		}

		for _, p := range expectedPaths {
			_, err := readFile(p)
			if err != nil {
				t.Errorf("failed checking expected paths: %s %s", p, err)
			}
		}

		appConfig, err := LoadLocalConfig("", mockFsConfigFileName)

		if err != nil {
			t.Errorf("failed to LoadLocalConfig: %s", err)
		}

		assert.Equal(1, len(appConfig.Drivers))
		assert.Equal("floor-01/bms", appConfig.Drivers[0].Name)
		var bacnetConfig config.Root
		err = json.Unmarshal(appConfig.Drivers[0].Raw, &bacnetConfig)
		if err != nil {
			t.Errorf("failed to unmarshall bacnet config: %s", err)
		}
		assert.Equal("eth0", bacnetConfig.LocalInterface)
		assert.Equal(uint16(47808), bacnetConfig.LocalPort)
		assert.Equal(2, len(bacnetConfig.Devices))
		assert.Equal("uk-ocw/floors/01/devices/CE1", bacnetConfig.Devices[0].Name)
		assert.Equal("172.16.8.115:47808", bacnetConfig.Devices[0].Comm.IP.String())
		assert.Equal("uk-ocw/floors/01/devices/CE2", bacnetConfig.Devices[1].Name)
		assert.Equal("172.16.8.117:47808", bacnetConfig.Devices[1].Comm.IP.String())

		writePageFile(expectedPaths[0], nil, "New Interface") //local interface
		writePageFile(expectedPaths[1], nil, 12345)           //local port
		err = mergeDbWithExtConfig(appConfig, dbRootPath)
		if err != nil {
			t.Errorf("failed to join app config & db: %s", err)
		}

		var newBacnetConfig config.Root
		err = json.Unmarshal(appConfig.Drivers[0].Raw, &newBacnetConfig)
		if err != nil {
			t.Errorf("failed to unmarshall bacnet config: %s", err)
		}
		assert.Equal("New Interface", newBacnetConfig.LocalInterface)
		assert.Equal(uint16(12345), newBacnetConfig.LocalPort)
		// ok now we want to update the IP address of the device with the key "uk-ocw/floors/01/devices/CE1"

	})
}
