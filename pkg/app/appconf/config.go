package appconf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/util/slices"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

var readFile = os.ReadFile // for testing

type Config struct {
	Name string `json:"name,omitempty"`
	// an array of other config files to read and merge.
	// If any included config files have further includes, they will also be loaded.
	// name will be ignored from any included files, and all other values will be merged.
	Includes   []string           `json:"includes,omitempty"`
	Drivers    []driver.RawConfig `json:"drivers,omitempty"`
	Automation []auto.RawConfig   `json:"automation,omitempty"`
	Zones      []zone.RawConfig   `json:"zones,omitempty"`
}

func (c *Config) mergeWith(other *Config) {
	c.Drivers = append(c.Drivers, other.Drivers...)
	c.Automation = append(c.Automation, other.Automation...)
	c.Zones = append(c.Zones, other.Zones...)
	// special case for includes - de-duplicate
	for i := 0; i < len(other.Includes); i++ {
		inc := other.Includes[i]
		if slices.Contains(inc, c.Includes) {
			continue
		}
		c.Includes = append(c.Includes, inc)
	}
}

// LoadLocalConfig will load Config from a local file, as well as any included files
func LoadLocalConfig(dir, file string) (*Config, error) {
	path := filepath.Join(dir, file)
	conf, err := configFromFile(path)
	if err != nil {
		return nil, err
	}
	// if we successfully loaded config, also load included files
	_, err = loadIncludes(dir, conf, conf.Includes, nil)
	return conf, err // return the config we have, and any errors
}

// loadIncludes will go through each include, load the configs, merge the configs, then load any further includes
func loadIncludes(dir string, dst *Config, includes, seen []string) ([]string, error) {
	var errs error
	var configs []*Config
	// load first layer of includes
	for i := 0; i < len(includes); i++ {
		include := includes[i]
		path := filepath.Join(dir, include)
		if slices.Contains(path, seen) {
			continue
		}
		seen = append(seen, path) // track files we've seen, to avoid getting in a loop
		extraConf, err := configFromFile(path)
		if err != nil {
			errs = multierr.Append(errs, err)
		} else {
			configs = append(configs, extraConf)
		}
	}
	// merge
	for i := 0; i < len(configs); i++ {
		base.mergeWith(configs[i]) // don't append to include in the first loop
	}
	// load all deeper includes
	for i := 0; i < len(configs); i++ {
		alsoSeen, err := loadIncludes(dir, dst, configs[i].Includes, seen)
		seen = append(seen, alsoSeen...)
		errs = multierr.Append(errs, err)
	}
	return seen, errs
}

// configFromFile will load Config for a local file
func configFromFile(path string) (*Config, error) {
	var conf Config
	raw, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file %s: %w", path, err)
	}
	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return nil, fmt.Errorf("config JSON unmarshal %s: %w", path, err)
	}
	return &conf, nil
}
