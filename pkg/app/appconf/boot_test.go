package appconf

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/block"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

var sampleBlocks = Blocks(
	map[string][]block.Block{ // sample driver blocks
		"driverType1": {
			{
				Path: []string{"devices"},
				Key:  "name",
				Blocks: []block.Block{
					{
						Path: []string{"addr"},
					},
				},
			},
		},
	},
	map[string][]block.Block{ // sample automation blocks

	},
	map[string][]block.Block{ // sample zone blocks

	},
)

func TestBootConfig(t *testing.T) {
	store := &memStore{}

	// boot with an empty store - this should just return the local config
	conf, err := BootConfig(&Config{
		Name: "test",
		Drivers: []driver.RawConfig{
			{
				BaseConfig: driver.BaseConfig{Type: "driverType1", Name: "driver1"},
				Raw: json.RawMessage(`{
					"type": "driverType1",
					"name": "driver1",
					"devices": [
						{"name": "device1", "addr": "addr1"},
						{"name": "device2", "addr": "addr2"}	
					]
				}`),
			},
		},
	}, store, sampleBlocks, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	expect := &Config{
		Name: "test",
		Drivers: []driver.RawConfig{
			{
				BaseConfig: driver.BaseConfig{},
				Raw: json.RawMessage(`{
					"type": "driverType1",
					"name": "driver1",
					"devices": [
						{"name": "device1", "addr": "addr1"},
						{"name": "device2", "addr": "addr2"}	
					]
				}`),
			},
		},
	}
	diff := cmp.Diff(expect, conf,
		cmpopts.EquateEmpty(),
		protocmp.Transform(),
		cmp.Comparer(compareDriverConfig),
	)
	if diff != "" {
		t.Errorf("(1) unexpected config (-want +got):\n%s", diff)
	}

	// modify the active config directly - this simulates making live changes
	// change the addr of a single device
	store.active.Drivers[0].Raw = json.RawMessage(`{
		"type": "driverType1",
		"name": "driver1",
		"devices": [
			{"name": "device1", "addr": "addr1"},
			{"name": "device2", "addr": "addr2changed"}
		]
	}`)

	// modify the local config with a patch adds another device to driver1
	conf, err = BootConfig(&Config{
		Name: "test",
		Drivers: []driver.RawConfig{
			{
				BaseConfig: driver.BaseConfig{Type: "driverType1", Name: "driver1"},
				Raw: json.RawMessage(`{
					"type": "driverType1",
					"name": "driver1",
					"devices": [
						{"name": "device1", "addr": "addr1"},
						{"name": "device2", "addr": "addr2"},
						{"name": "device3", "addr": "addr3"}
					]
				}`),
			},
		},
	}, store, sampleBlocks, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	// check that both changes from local config were applied to the active config
	expect = &Config{
		Name: "test",
		Drivers: []driver.RawConfig{
			{
				BaseConfig: driver.BaseConfig{Name: "driver1", Type: "driverType1"},
				Raw: json.RawMessage(`{
					"type": "driverType1",
					"name": "driver1",
					"devices": [
						{"name": "device1", "addr": "addr1"},
						{"name": "device2", "addr": "addr2changed"},
						{"name": "device3", "addr": "addr3"}
					]
				}`),
			},
		},
	}

	diff = cmp.Diff(expect, conf,
		cmpopts.EquateEmpty(),
		protocmp.Transform(),
		cmp.Comparer(compareDriverConfig),
	)
	if diff != "" {
		t.Errorf("(2) unexpected config (-want +got):\n%s", diff)
	}
}

type memStore struct {
	local   *Config
	active  *Config
	patches [][]block.Patch
}

func (m *memStore) SwapLocalConfig(new *Config) (old *Config, err error) {
	old = cloneConfig(m.local)
	m.local = cloneConfig(new)
	return old, nil
}

func (m *memStore) GetActiveConfig() (*Config, error) {
	return cloneConfig(m.active), nil
}

func (m *memStore) SetActiveConfig(c *Config) error {
	m.active = cloneConfig(c)
	return nil
}

func (m *memStore) SavePatches(patches []block.Patch) (ref string, err error) {
	ref = strconv.Itoa(len(m.patches))
	m.patches = append(m.patches, slices.Clone(patches))
	return ref, nil
}

func cloneConfig(c *Config) *Config {
	if c == nil {
		return nil
	}
	metadata := c.Metadata
	if metadata != nil {
		metadata = proto.Clone(metadata).(*traits.Metadata)
	}
	return &Config{
		Name:       c.Name,
		Metadata:   metadata,
		Includes:   slices.Clone(c.Includes),
		Drivers:    slices.Clone(c.Drivers),
		Automation: slices.Clone(c.Automation),
		Zones:      slices.Clone(c.Zones),
		FilePath:   c.FilePath,
	}
}

func compareDriverConfig(a, b driver.RawConfig) bool {
	// reencode so equivalent values compare equal
	var aMap, bMap any
	if err := json.Unmarshal(a.Raw, &aMap); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b.Raw, &bMap); err != nil {
		panic(err)
	}

	aBuf, err := json.MarshalIndent(aMap, "", "  ")
	if err != nil {
		panic(err)
	}
	bBuf, err := json.MarshalIndent(bMap, "", "  ")
	if err != nil {
		panic(err)
	}
	return bytes.Equal(aBuf, bBuf)
}
