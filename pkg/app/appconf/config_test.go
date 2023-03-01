package appconf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"

	"github.com/vanti-dev/sc-bos/pkg/driver"
)

func ExampleConfig() {
	type ExampleDriverConfig struct {
		driver.BaseConfig
		Property string `json:"property"`
	}

	// language=json
	buf := []byte(`{
	"drivers": [
		{
			"name": "foo",
			"type": "example",
			"property": "bar"
		}
    ]
}`)
	var config Config
	err := json.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	for _, d := range config.Drivers {
		switch d.Type {
		case "example":
			var exampleConfig ExampleDriverConfig
			err := json.Unmarshal(d.Raw, &exampleConfig)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ExampleDriver name=%q property=%q\n", d.Name, exampleConfig.Property)
		default:
			fmt.Printf("unknown driver type %q\n", d.Type)
		}
	}
	// Output: ExampleDriver name="foo" property="bar"
}

func TestLoadLocalConfig(t *testing.T) {
	driverNamed := func(name string) driver.RawConfig {
		return driver.RawConfig{
			BaseConfig: driver.BaseConfig{Name: name},
			Raw:        json.RawMessage(fmt.Sprintf("{\"name\": \"%s\"}", name)),
		}
	}
	tests := []struct {
		name   string
		fs     fstest.MapFS
		dir    string
		file   string
		config *Config
	}{
		{
			name:   "empty",
			fs:     nil,
			dir:    "data",
			file:   "base.json",
			config: nil,
		},
		{
			name: "single file",
			fs: fstest.MapFS{
				filepath.Join("data", "base.json"): {
					Data: []byte(`{"name": "my-config"}`),
				},
			},
			dir:    "data",
			file:   "base.json",
			config: &Config{Name: "my-config"},
		},
		{
			name: "with include",
			fs: fstest.MapFS{
				filepath.Join("data", "base.json"): {
					Data: []byte(`{
						"name": "my-config",
						"includes": ["part-1.json"],
						"drivers": [{"name": "driver-1"}]
					}`),
				},
				filepath.Join("data", "part-1.json"): {
					Data: []byte(`{
						"name": "ignored",
						"drivers": [{"name": "driver-part-1"}]
					}`),
				},
			},
			dir:  "data",
			file: "base.json",
			config: &Config{
				Name:     "my-config",
				Includes: []string{"part-1.json"},
				Drivers: []driver.RawConfig{
					driverNamed("driver-1"),
					driverNamed("driver-part-1"),
				},
			},
		},
		{
			name: "with multiple nested includes",
			fs: fstest.MapFS{
				filepath.Join("data", "base.json"): {
					Data: []byte(`{
						"name": "my-config",
						"includes": ["part-1.json", "part-2.json", "part-3.json"],
						"drivers": [{"name": "driver-1"}]
					}`),
				},
				filepath.Join("data", "part-1.json"): {
					Data: []byte(`{
						"name": "ignored",
						"includes": ["part-1a.json", "part-1b.json"],
						"drivers": [{"name": "driver-part-1"}]
					}`),
				},
				filepath.Join("data", "part-1a.json"): {
					Data: []byte(`{
						"drivers": [{"name": "driver-part-1a"}]
					}`),
				},
				filepath.Join("data", "part-1b.json"): {
					Data: []byte(`{
						"drivers": [{"name": "driver-part-1b"}]
					}`),
				},
				filepath.Join("data", "part-2.json"): {
					Data: []byte(`{
						"name": "ignored",
						"includes": ["part-2a.json", "part-2b.json"],
						"drivers": [{"name": "driver-part-2"}]
					}`),
				},
				filepath.Join("data", "part-2a.json"): {
					Data: []byte(`{
						"drivers": [{"name": "driver-part-2a"}]
					}`),
				},
				filepath.Join("data", "part-2b.json"): {
					Data: []byte(`{
						"drivers": [{"name": "driver-part-2b"}]
					}`),
				},
				filepath.Join("data", "part-3.json"): {
					Data: []byte(`{
						"name": "ignored",
						"includes": ["part-3a.json", "part-3b.json"],
						"drivers": [{"name": "driver-part-3"}]
					}`),
				},
				filepath.Join("data", "part-3a.json"): {
					Data: []byte(`{
						"drivers": [{"name": "driver-part-3a"}]
					}`),
				},
				// part-3b is missing
			},
			dir:  "data",
			file: "base.json",
			config: &Config{
				Name: "my-config",
				Includes: []string{
					"part-1.json",
					"part-2.json",
					"part-3.json",
					"part-1a.json",
					"part-1b.json",
					"part-2a.json",
					"part-2b.json",
					"part-3a.json",
					"part-3b.json",
				},
				Drivers: []driver.RawConfig{
					driverNamed("driver-1"),
					driverNamed("driver-part-1"),
					driverNamed("driver-part-2"),
					driverNamed("driver-part-3"),
					driverNamed("driver-part-1a"),
					driverNamed("driver-part-1b"),
					driverNamed("driver-part-2a"),
					driverNamed("driver-part-2b"),
					driverNamed("driver-part-3a"),
				},
			},
		},
		{
			name: "avoids includes loop",
			fs: fstest.MapFS{
				filepath.Join("data", "base.json"): {
					Data: []byte(`{
						"name": "my-config",
						"includes": ["part-1.json"],
						"drivers": [{"name": "driver-1"}]
					}`),
				},
				filepath.Join("data", "part-1.json"): {
					Data: []byte(`{
						"name": "ignored",
						"includes": ["part-1a.json"],
						"drivers": [{"name": "driver-part-1"}]
					}`),
				},
				filepath.Join("data", "part-1a.json"): {
					Data: []byte(`{
						"includes": ["part-1.json"],
						"drivers": [{"name": "driver-part-1a"}]
					}`),
				},
			},
			dir:  "data",
			file: "base.json",
			config: &Config{
				Name: "my-config",
				Includes: []string{
					"part-1.json",
					"part-1a.json",
				},
				Drivers: []driver.RawConfig{
					driverNamed("driver-1"),
					driverNamed("driver-part-1"),
					driverNamed("driver-part-1a"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readFile = tt.fs.ReadFile
			conf, err := LoadLocalConfig(tt.dir, tt.file)
			if tt.config != nil && conf == nil && err != nil {
				t.Errorf("expected config, got error: %s", err)
			}
			if diff := cmp.Diff(tt.config, conf); diff != "" {
				t.Errorf("wrong config: %s", diff)
			}
			readFile = os.ReadFile
		})
	}
}
