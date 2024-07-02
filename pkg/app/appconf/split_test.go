package appconf

import (
	"encoding/json"
	"io/fs"
	"net/netip"
	"os"
	"path/filepath"
	"reflect"
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

func unwrapValue(value interface{}) interface{} {

	result := value
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		result = reflect.ValueOf(value).Elem().Interface()

		switch result.(type) {
		case *netip.AddrPort:
			result = result.(*netip.AddrPort).String()
		}

	} else if val.Kind() == reflect.Slice {

	}
	return result
}

// read the split file, so we know how to split the file into parts we want to edit
// we can now create the db file structure based on what we have read in
// tests we can create the directory structure for db
// test modifying the db & joining the db with the ext (appconf.Config)
// appconf.Config should then contain the edits in db & the original values
func TestMetadataConfigPatch(t *testing.T) {

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

	appConfig, err := LoadLocalConfig("", mockFsConfigFileName)

	if err != nil {
		t.Errorf("failed to LoadLocalConfig: %s", err)
	}

	tests := []struct {
		name       string
		fs         MockFs
		which      interface{}
		preExpect  any
		patchFile  string
		patchValue any
		post       func(interface{}, interface{}, ...interface{}) bool
	}{
		{
			name:       "Floor",
			fs:         mockFs,
			which:      &appConfig.Metadata.Location.Floor,
			preExpect:  "Floor 1",
			patchFile:  filepath.Join("testdata", "db", "metadata", "Location", "Floor"),
			patchValue: "New Floor",
		},
		{
			name:       "Manufacturer",
			fs:         mockFs,
			which:      &appConfig.Metadata.Product.Manufacturer,
			preExpect:  "Vanti",
			patchFile:  filepath.Join("testdata", "db", "metadata", "Product", "Manufacturer"),
			patchValue: "New Manufacturer",
		},
		{
			name:       "Model",
			fs:         mockFs,
			which:      &appConfig.Metadata.Product.Model,
			preExpect:  "Smart Core BOS",
			patchFile:  filepath.Join("testdata", "db", "metadata", "Product", "Model"),
			patchValue: "New Model",
		},
		{
			name:       "Membership",
			fs:         mockFs,
			which:      appConfig.Metadata.Membership,
			preExpect:  traits.Metadata_Membership{Subsystem: "smart"},
			patchFile:  filepath.Join("testdata", "db", "metadata", "Membership"),
			patchValue: traits.Metadata_Membership{Subsystem: "New Subsystem"},
		},
		{
			name:       "Traits",
			fs:         mockFs,
			which:      appConfig.Metadata.Traits,
			preExpect:  []*traits.TraitMetadata{{Name: "oldTrait"}},
			patchFile:  filepath.Join("testdata", "db", "metadata", "Traits"),
			patchValue: []*traits.TraitMetadata{{Name: "newTrait"}},
		},
		{
			name:  "MoreMap",
			fs:    mockFs,
			which: &appConfig.Metadata.More,
			preExpect: map[string]string{
				"type":     "sensor",
				"function": "temperature",
			},
			patchFile: filepath.Join("testdata", "db", "metadata", "More"),
			patchValue: map[string]string{
				"type":     "newType",
				"function": "newFunction",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assert.Equal(tt.preExpect, unwrapValue(tt.which))

			err := writePageFile(tt.patchFile, "", tt.patchValue)
			if err != nil {
				t.Errorf("failed to write patch file: %s", err)
			}
			err = mergeDbWithExtConfig(appConfig, dbRootPath)
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}

			assert.Equal(tt.patchValue, unwrapValue(tt.which))
		})
	}
}

// tests the ability of the config system to update the property of a specific device in the config
// the device is specified using the "key" attribute in the split file
// when the split-config encounters the key property being present in the split file, for a given split
// it will search through the array / map of objects in the config for the object with the matching key
// and then follow the same process
func TestBacnetDriverConfigPatch(t *testing.T) {

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
	appConfig, err := LoadLocalConfig("", mockFsConfigFileName)
	if err != nil {
		t.Errorf("failed to LoadLocalConfig: %s", err)
	}

	var bacnetConfig config.Root
	err = json.Unmarshal(appConfig.Drivers[0].Raw, &bacnetConfig)
	if err != nil {
		t.Errorf("failed to unmarshall bacnet config: %s", err)
	}

	tests := []struct {
		name       string
		which      interface{}
		preExpect  any
		patchFile  string
		patchValue any
		post       func(interface{}, interface{}, ...interface{}) bool
	}{
		{
			name:       "localInterface",
			which:      &bacnetConfig.LocalInterface,
			preExpect:  "eth0",
			patchFile:  filepath.Join("testdata", "db", "drivers", "bacnet", "localInterface"),
			patchValue: "New Interface",
		},
		{
			name:       "localPort",
			which:      &bacnetConfig.LocalPort,
			preExpect:  uint16(47808),
			patchFile:  filepath.Join("testdata", "db", "drivers", "bacnet", "localPort"),
			patchValue: uint16(12345),
		},
		{
			name:       "device1IP",
			which:      &bacnetConfig.Devices[0].Comm.IP,
			preExpect:  "172.16.8.115:47808",
			patchFile:  filepath.Join("testdata", "db", "drivers", "bacnet", "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "comm", "ip"),
			patchValue: "188.88.8.71:8888",
		},
		{
			name:       "device2IP",
			which:      &bacnetConfig.Devices[1].Comm.IP,
			preExpect:  "172.16.8.117:47808",
			patchFile:  filepath.Join("testdata", "db", "drivers", "bacnet", "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE2"), "comm", "ip"),
			patchValue: "22.22.2.71:2222",
		},
		{
			name:       "metadata_title",
			which:      &bacnetConfig.Devices[0].Metadata.Appearance.Title,
			preExpect:  "Floor 1 Controller North",
			patchFile:  filepath.Join("testdata", "db", "drivers", "bacnet", "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "metadata", "appearance", "title"),
			patchValue: "New Title",
		},
		{
			name:       "metadata_location",
			which:      &bacnetConfig.Devices[0].Metadata.Location,
			preExpect:  &traits.Metadata_Location{Floor: "Floor 1", Zone: "North"},
			patchFile:  filepath.Join("testdata", "db", "drivers", "bacnet", "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "metadata", "location"),
			patchValue: &traits.Metadata_Location{Floor: "New Floor", Zone: "New Zone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assert.Equal(tt.preExpect, unwrapValue(tt.which))

			err := writePageFile(tt.patchFile, "", tt.patchValue)
			if err != nil {
				t.Errorf("failed to write patch file: %s", err)
			}
			err = mergeDbWithExtConfig(appConfig, dbRootPath)
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}
			err = json.Unmarshal(appConfig.Drivers[0].Raw, &bacnetConfig)
			if err != nil {
				t.Errorf("failed to unmarshall bacnet config: %s", err)
			}

			assert.Equal(tt.patchValue, unwrapValue(tt.which))
		})
	}
}
