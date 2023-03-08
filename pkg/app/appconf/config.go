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
	// if a duplicate name is found (e.g. duplicate driver), then the first one will be used
	Includes   []string           `json:"includes,omitempty"`
	Drivers    []driver.RawConfig `json:"drivers,omitempty"`
	Automation []auto.RawConfig   `json:"automation,omitempty"`
	Zones      []zone.RawConfig   `json:"zones,omitempty"`
}

func (c *Config) mergeWith(other *Config) {
	// if any driver/auto/zone has a duplicate name it is ignored in favour of the one already present

	driverNames := c.driverNamesMap()
	autoNames := c.autoNamesMap()
	zoneNames := c.zoneNamesMap()
	for i := 0; i < len(other.Drivers); i++ {
		d := other.Drivers[i]
		if _, found := driverNames[d.Name]; !found {
			c.Drivers = append(c.Drivers, d)
		}
	}
	for i := 0; i < len(other.Automation); i++ {
		d := other.Automation[i]
		if _, found := autoNames[d.Name]; !found {
			c.Automation = append(c.Automation, d)
		}
	}
	for i := 0; i < len(other.Zones); i++ {
		d := other.Zones[i]
		if _, found := zoneNames[d.Name]; !found {
			c.Zones = append(c.Zones, d)
		}
	}
	// special case for includes - de-duplicate
	for i := 0; i < len(other.Includes); i++ {
		inc := other.Includes[i]
		if slices.Contains(inc, c.Includes) {
			continue
		}
		c.Includes = append(c.Includes, inc)
	}
}

func (c *Config) driverNamesMap() map[string]bool {
	names := make(map[string]bool, len(c.Drivers))
	for _, d := range c.Drivers {
		names[d.Name] = true
	}
	return names
}

func (c *Config) autoNamesMap() map[string]bool {
	names := make(map[string]bool, len(c.Automation))
	for _, d := range c.Automation {
		names[d.Name] = true
	}
	return names
}

func (c *Config) zoneNamesMap() map[string]bool {
	names := make(map[string]bool, len(c.Zones))
	for _, d := range c.Zones {
		names[d.Name] = true
	}
	return names
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
			dst.mergeWith(extraConf)
		}
	}
	// load all deeper includes
	for i := 0; i < len(configs); i++ {
		alsoSeen, err := loadIncludes(dir, dst, configs[i].Includes, seen)
		if err != nil {
			seen = alsoSeen
		}
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
