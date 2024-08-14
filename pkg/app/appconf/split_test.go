package appconf

import (
	"encoding/json"
	"io/fs"
	"net/netip"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/traits"
	lightingconfig "github.com/vanti-dev/sc-bos/pkg/auto/lights/config"
	bacnetconfig "github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
	areaconfig "github.com/vanti-dev/sc-bos/pkg/zone/area/config"
)

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

func unwrapValue(value any) any {

	result := value
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		result = reflect.ValueOf(value).Elem().Interface()
		switch result.(type) {
		case *netip.AddrPort:
			result = result.(*netip.AddrPort).String()
		case jsontypes.Duration:
			result = result.(jsontypes.Duration).String()
			result = strings.TrimSuffix(result.(string), "0s")
		}
	}
	return result
}

type alternateKey struct {
	Path string `json:"path,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Test that we can update the metadata section of the app config
func TestMetadataConfigPatch(t *testing.T) {

	// first set up the mock filesystem, read the test ext config
	var mockFs = MockFs{fs: afero.NewMemMapFs()}
	readFile = mockFs.mockReadFile
	writeFile = mockFs.mockWriteFile
	mkdirAll = mockFs.mockMkdirAll
	readDir = mockFs.mockReadDir
	isDir = mockFs.mockIsDir
	mockFsConfigFileName := "fstest.metadata.json"

	file, err := os.ReadFile("testdata/metadata.json")
	if err != nil {
		t.Errorf("error reading config file: %s", err)
	}
	err = writeFile(mockFsConfigFileName, file, 0664)
	if err != nil {
		t.Fatal(err)
	}

	rootPath := filepath.Join("testdata", "db")

	appConfig, err := LoadLocalConfig("", mockFsConfigFileName)

	if err != nil {
		t.Errorf("failed to LoadLocalConfig: %s", err)
	}

	tests := []struct {
		name      string      // name of the test
		which     interface{} // the field to be updated
		preExpect any         // the value of the field before the update
		patchFile string      // the file containing the new value
		change    any         // the change we are applying
		// at the end of each test the value of `which` should be equal to `change`
	}{
		{
			name:      "Floor",
			which:     &appConfig.Metadata.Location.Floor,
			preExpect: "Floor 1",
			patchFile: filepath.Join("testdata", "db", "metadata", "Location", "Floor"),
			change:    "New Floor",
		},
		{
			name:      "Manufacturer",
			which:     &appConfig.Metadata.Product.Manufacturer,
			preExpect: "Vanti",
			patchFile: filepath.Join("testdata", "db", "metadata", "Product", "Manufacturer"),
			change:    "New Manufacturer",
		},
		{
			name:      "Model",
			which:     &appConfig.Metadata.Product.Model,
			preExpect: "Smart Core BOS",
			patchFile: filepath.Join("testdata", "db", "metadata", "Product", "Model"),
			change:    "New Model",
		},
		{
			name:      "Membership",
			which:     appConfig.Metadata.Membership,
			preExpect: traits.Metadata_Membership{Subsystem: "smart"},
			patchFile: filepath.Join("testdata", "db", "metadata", "Membership"),
			change:    traits.Metadata_Membership{Subsystem: "New Subsystem"},
		},
		{
			name:      "Traits",
			which:     appConfig.Metadata.Traits,
			preExpect: []*traits.TraitMetadata{{Name: "oldTrait"}},
			patchFile: filepath.Join("testdata", "db", "metadata", "Traits"),
			change:    []*traits.TraitMetadata{{Name: "newTrait"}},
		},
		{
			name:  "MoreMap",
			which: &appConfig.Metadata.More,
			preExpect: map[string]string{
				"type":     "sensor",
				"function": "temperature",
			},
			patchFile: filepath.Join("testdata", "db", "metadata", "More"),
			change: map[string]string{
				"type":     "newType",
				"function": "newFunction",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// check the value at the beginning is the original from the ext
			if diff := cmp.Diff(tt.preExpect, unwrapValue(tt.which), protocmp.Transform()); diff != "" {
				t.Errorf("unexpected value at the beginning (-want +got):\n%s", diff)
			}

			// write the file containing the new value to the correct place
			err := writePageFile(tt.patchFile, tt.change)
			if err != nil {
				t.Errorf("failed to write patch file: %s", err)
			}

			// merge changes from the files in rootPath directory into appConfig
			err = mergeDbWithExtConfig(appConfig, rootPath)
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}

			// the value in appConfig should have the new value from the patchFile
			if diff := cmp.Diff(tt.change, unwrapValue(tt.which), protocmp.Transform()); diff != "" {
				t.Errorf("unexpected value at the end (-want +got):\n%s", diff)
			}
		})
	}
}

// test we can update the drivers section of the config, using the bacnet driver as the test
func TestBacnetDriverConfigPatch(t *testing.T) {

	// first set up the mock filesystem, read & add the metadata & split files
	// is more readable to do it this way
	var mockFs = MockFs{fs: afero.NewMemMapFs()}
	readFile = mockFs.mockReadFile
	writeFile = mockFs.mockWriteFile
	mkdirAll = mockFs.mockMkdirAll
	readDir = mockFs.mockReadDir
	isDir = mockFs.mockIsDir
	mockFsConfigFileName := "fstest.bms.json"

	file, err := os.ReadFile("testdata/bms.json")
	if err != nil {
		t.Errorf("error reading config file: %s", err)
	}

	err = writeFile(mockFsConfigFileName, file, 0664)
	if err != nil {
		t.Fatal(err)
	}

	rootPath := filepath.Join("testdata", "db")
	appConfig, err := LoadLocalConfig("", mockFsConfigFileName)
	if err != nil {
		t.Errorf("failed to LoadLocalConfig: %s", err)
	}

	var bacnetConfig bacnetconfig.Root
	err = json.Unmarshal(appConfig.Drivers[0].Raw, &bacnetConfig)
	if err != nil {
		t.Errorf("failed to unmarshall bacnet config: %s", err)
	}

	tests := []struct {
		name         string        // name of the test
		which        interface{}   // the field to be updated
		preExpect    any           // the value of the field before the update
		patchFile    string        // the file containing the new value
		change       any           // the change we are applying
		alternateKey *alternateKey // an alternate key to use instead of the default `name`
	}{
		{
			name:      "localInterface",
			which:     &bacnetConfig.LocalInterface,
			preExpect: "eth0",
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "localInterface"),
			change:    "New Interface",
		},
		{
			name:      "localPort",
			which:     &bacnetConfig.LocalPort,
			preExpect: uint16(47808),
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName(normaliseDeviceName("floor-01/bms")), "localPort"),
			change:    uint16(12345),
		},
		{
			name:      "device1IP",
			which:     &bacnetConfig.Devices[0].Comm.IP,
			preExpect: "172.16.8.115:47808",
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "comm", "ip"),
			change:    "188.88.8.71:8888",
		},
		{
			name:      "device2IP",
			which:     &bacnetConfig.Devices[1].Comm.IP,
			preExpect: "172.16.8.117:47808",
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE2"), "comm", "ip"),
			change:    "22.22.2.71:2222",
		},
		{
			name:      "metadata_title",
			which:     &bacnetConfig.Devices[0].Metadata.Appearance.Title,
			preExpect: "Floor 1 Controller North",
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "metadata", "appearance", "title"),
			change:    "New Title",
		},
		{
			name:      "metadata_location",
			which:     &bacnetConfig.Devices[0].Metadata.Location,
			preExpect: &traits.Metadata_Location{Floor: "Floor 1", Zone: "North"},
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "metadata", "location"),
			change:    &traits.Metadata_Location{Floor: "New Floor", Zone: "New Zone"},
		},
		{
			name:      "object_title",
			which:     &bacnetConfig.Devices[0].Objects[0].Title,
			preExpect: "CPU Board Temperature",
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "objects", normaliseDeviceName("CPUBoardTemp"), "title"),
			change:    "New Object Title",
		},
		{
			name:      "object_name_by_id",
			which:     &bacnetConfig.Devices[0].Objects[1].Name,
			preExpect: "SpVAVFeedback",
			patchFile: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "objects", normaliseDeviceName("AnalogInput:1101"), "name"),
			change:    "New Object Name",
			alternateKey: &alternateKey{
				Path: filepath.Join("testdata", "db", "drivers", normaliseDeviceName("floor-01/bms"), "devices", normaliseDeviceName("uk-ocw/floors/01/devices/CE1"), "objects"),
				Key:  "id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// check the value at the beginning is the original from the ext
			if diff := cmp.Diff(tt.preExpect, unwrapValue(tt.which), protocmp.Transform()); diff != "" {
				t.Errorf("unexpected value at the beginning (-want +got):\n%s", diff)
			}

			// write the file containing the new value to the correct place
			err := writePageFile(tt.patchFile, tt.change)
			if err != nil {
				t.Errorf("failed to write patch file: %s", err)
			}

			if tt.alternateKey != nil {
				// tell the system to use a key other than 'name' to specify the value in slice/map we want
				err = writeAlternateKey(tt.alternateKey.Path, tt.alternateKey.Key)
				if err != nil {
					t.Errorf("failed to write alternate key file: %s", err)
				}
			}

			// merge changes from the files in rootPath directory into appConfig
			err = mergeDbWithExtConfig(appConfig, rootPath)
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}

			err = json.Unmarshal(appConfig.Drivers[0].Raw, &bacnetConfig)
			if err != nil {
				t.Errorf("failed to unmarshall bacnet config: %s", err)
			}

			// the value in appConfig should have the new value from the patchFile
			if diff := cmp.Diff(tt.change, unwrapValue(tt.which), protocmp.Transform()); diff != "" {
				t.Errorf("unexpected value at the end (-want +got):\n%s", diff)
			}
		})
	}
}

// test we can update the automations section of the config using lighting automation as an example
func TestAutomation(t *testing.T) {

	var mockFs = MockFs{fs: afero.NewMemMapFs()}
	readFile = mockFs.mockReadFile
	writeFile = mockFs.mockWriteFile
	mkdirAll = mockFs.mockMkdirAll
	readDir = mockFs.mockReadDir
	isDir = mockFs.mockIsDir
	mockFsConfigFileName := "fstest.automation.json"

	file, err := os.ReadFile("testdata/automation.json")
	if err != nil {
		t.Errorf("error reading config file: %s", err)
	}

	err = writeFile(mockFsConfigFileName, file, 0664)
	if err != nil {
		t.Fatal(err)
	}

	rootPath := filepath.Join("testdata", "db")
	appConfig, err := LoadLocalConfig("", mockFsConfigFileName)
	if err != nil {
		t.Errorf("failed to LoadLocalConfig: %s", err)
	}

	var lightsConfig lightingconfig.Root
	err = json.Unmarshal(appConfig.Automation[0].Raw, &lightsConfig)
	if err != nil {
		t.Errorf("failed to unmarshall auto config: %s", err)
	}

	tests := []struct {
		name         string        // name of the test
		which        interface{}   // the field to be updated
		preExpect    any           // the value of the field before the update
		patchFile    string        // the file containing the new value
		change       any           // the change we are applying
		alternateKey *alternateKey // an alternate key to use instead of the default `name`
	}{
		{
			name:      "unoccupiedOffDelay",
			which:     &lightsConfig.UnoccupiedOffDelay,
			preExpect: "20m",
			patchFile: filepath.Join("testdata", "db", "automation", normaliseDeviceName("Lights: Basecamp 2"), "unoccupiedOffDelay"),
			change:    "45m",
		},
		{
			name:      "unoccupiedOffDelay_Daytime",
			which:     &lightsConfig.Modes[0].UnoccupiedOffDelay,
			preExpect: "30m",
			patchFile: filepath.Join("testdata", "db", "automation", normaliseDeviceName("Lights: Basecamp 2"), "modes",
				normaliseDeviceName("Daytime work area"), "unoccupiedOffDelay"),
			change: "20m",
		},
		{
			name:      "unoccupiedOffDelay_Night",
			which:     &lightsConfig.Modes[1].UnoccupiedOffDelay,
			preExpect: "15m",
			patchFile: filepath.Join("testdata", "db", "automation", normaliseDeviceName("Lights: Basecamp 2"), "modes",
				normaliseDeviceName("Night work area"), "unoccupiedOffDelay"),
			change: "20m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// check the value at the beginning is the original from the ext
			if diff := cmp.Diff(tt.preExpect, unwrapValue(tt.which)); diff != "" {
				t.Errorf("unexpected value at the beginning (-want +got):\n%s", diff)
			}

			// write the file containing the new value to the correct place
			err := writePageFile(tt.patchFile, tt.change)
			if err != nil {
				t.Errorf("failed to write patch file: %s", err)
			}

			// merge changes from the files in rootPath directory into appConfig
			err = mergeDbWithExtConfig(appConfig, rootPath)
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}
			err = json.Unmarshal(appConfig.Automation[0].Raw, &lightsConfig)
			if err != nil {
				t.Errorf("failed to unmarshall lights config: %s", err)
			}

			// the value in appConfig should have the new value from the patchFile
			if diff := cmp.Diff(tt.change, unwrapValue(tt.which)); diff != "" {
				t.Errorf("unexpected value at the end (-want +got):\n%s", diff)
			}
		})
	}
}

// test we can update the automations section of the config using area zone type as example
func TestZones(t *testing.T) {

	var mockFs = MockFs{fs: afero.NewMemMapFs()}
	readFile = mockFs.mockReadFile
	writeFile = mockFs.mockWriteFile
	mkdirAll = mockFs.mockMkdirAll
	readDir = mockFs.mockReadDir
	isDir = mockFs.mockIsDir
	mockFsConfigFileName := "fstest.zones.json"

	file, err := os.ReadFile("testdata/zones.json")
	if err != nil {
		t.Errorf("error reading config file: %s", err)
	}

	err = writeFile(mockFsConfigFileName, file, 0664)
	if err != nil {
		t.Fatal(err)
	}

	rootPath := filepath.Join("testdata", "db")
	appConfig, err := LoadLocalConfig("", mockFsConfigFileName)
	if err != nil {
		t.Errorf("failed to LoadLocalConfig: %s", err)
	}

	var areaConfig areaconfig.Root
	err = json.Unmarshal(appConfig.Zones[0].Raw, &areaConfig)
	if err != nil {
		t.Errorf("failed to unmarshall zone config: %s", err)
	}

	tests := []struct {
		name      string      // name of the test
		which     interface{} // the field to be updated
		preExpect any         // the value of the field before the update
		patchFile string      // the file containing the new value
		change    any         // the change we are applying
	}{
		{
			name:      "metadataAppearanceTitle",
			which:     &areaConfig.Metadata.Appearance.Title,
			preExpect: "Audit Office",
			patchFile: filepath.Join("testdata", "db", "zones", normaliseDeviceName("informa/uk/5hp/zone/audit-office"), "metadata", "appearance", "title"),
			change:    "CEO Golf Club Storage",
		},
		{
			name:      "disabled",
			which:     &areaConfig.Disabled,
			preExpect: false,
			patchFile: filepath.Join("testdata", "db", "zones", normaliseDeviceName("informa/uk/5hp/zone/audit-office"), "disabled"),
			change:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// check the value at the beginning is the original from the ext
			if diff := cmp.Diff(tt.preExpect, unwrapValue(tt.which)); diff != "" {
				t.Errorf("unexpected value at the beginning (-want +got):\n%s", diff)
			}

			// write the file containing the new value to the correct place
			err := writePageFile(tt.patchFile, tt.change)
			if err != nil {
				t.Errorf("failed to write patch file: %s", err)
			}

			// merge changes from the files in rootPath directory into appConfig
			err = mergeDbWithExtConfig(appConfig, rootPath)
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}
			if err != nil {
				t.Errorf("failed to join app config & db: %s", err)
			}
			err = json.Unmarshal(appConfig.Zones[0].Raw, &areaConfig)
			if err != nil {
				t.Errorf("failed to unmarshall area config: %s", err)
			}

			// the value in appConfig should have the new value from the patchFile
			if diff := cmp.Diff(tt.change, unwrapValue(tt.which)); diff != "" {
				t.Errorf("unexpected value at the end (-want +got):\n%s", diff)
			}
		})
	}
}
