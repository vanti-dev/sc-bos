package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"go.uber.org/multierr"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet"
)

// BacnetConfigForFloor is basically just a wrapper around the root config with support for automations and traits
type BacnetConfigForFloor struct {
	Root
	Devices     []Device
	Automations []map[string]any // map looks like config that embeds auto.Config
	Traits      []map[string]any // map looks like config that embeds config.Trait
}

// WriteBacnetConfig writes bacnet config for each floor to the given directories for each floor
// configPerFloor is a map of floor name to the bacnet config for that floor
// dirForFloor is a map of floor name to the directory for that floor. They keys are the same as configPerFloor and
// there must be a directory value defined for each floor. This allows you to put 2 floor configs in the same directory.
// configRoot is the root directory for the config, defaults to "config" if empty
// scPrefix is the prefix for the driver's Name. The bacnet driver name will become <sc-prefix>/floor-xx/drivers/bms
func WriteBacnetConfig(configPerFloor map[string]*BacnetConfigForFloor, dirForFloor map[string]string, configRoot string,
	scPrefix string) error {
	if configRoot == "" {
		configRoot = "config"
	}
	for floor, cfg := range configPerFloor {
		sortDevices(cfg.Devices)
		sortRawConfig(cfg.Traits)
		sortRawConfig(cfg.Automations)

		configDir, ok := dirForFloor[floor]
		if !ok {
			return fmt.Errorf("no config dir for floor %s", floor)
		}
		configDir = path.Join(configRoot, configDir)

		err := writeToDir(configDir, floor, scPrefix, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeToDir(dir string, floor string, scPrefix string, bacnetCfg *BacnetConfigForFloor) error {
	var errs error
	outFile := filepath.Join(dir, fmt.Sprintf("%s.bms.part.json", floor))
	cfg := appconf.Config{}
	driverCfg := Defaults()
	// copy over any customisations we have made
	driverCfg = bacnetCfg.Root
	driverCfg.BaseConfig = driver.BaseConfig{
		Name: path.Join(scPrefix, fmt.Sprintf("floor-%s", floor), "drivers", "bms"),
		Type: bacnet.DriverName,
	}
	driverCfg.Devices = bacnetCfg.Devices

	for _, t := range bacnetCfg.Traits {
		rawTrait, err := marshalTrait(t)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		driverCfg.Traits = append(driverCfg.Traits, rawTrait)
	}

	rawConfig, err := marshalDriver(driverCfg)
	if errs != nil {
		errs = multierr.Append(errs, err)
	} else {
		cfg.Drivers = append(cfg.Drivers, rawConfig)
	}

	for _, automation := range bacnetCfg.Automations {
		rawConfig, err := marshalAuto(automation)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		cfg.Automation = append(cfg.Automation, rawConfig)
	}

	if errs != nil {
		return errs
	}

	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func marshalTrait(data map[string]any) (RawTrait, error) {
	t := Trait{
		Name: data["name"].(string),
		Kind: data["kind"].(trait.Name),
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return RawTrait{}, err
	}
	return RawTrait{
		Trait: t,
		Raw:   raw,
	}, nil
}

func marshalAuto(data map[string]any) (auto.RawConfig, error) {
	c := auto.Config{
		Name: data["name"].(string),
		Type: data["type"].(string),
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return auto.RawConfig{}, err
	}
	return auto.RawConfig{
		Config: c,
		Raw:    raw,
	}, nil
}

func marshalDriver(data Root) (driver.RawConfig, error) {
	d := data.BaseConfig
	raw, err := json.Marshal(data)
	if err != nil {
		return driver.RawConfig{}, err
	}
	return driver.RawConfig{
		BaseConfig: d,
		Raw:        raw,
	}, nil
}

func sortDevices(devices []Device) {
	for _, device := range devices {
		sort.Slice(device.Objects, func(i, j int) bool {
			ida := device.Objects[i].ID
			idb := device.Objects[j].ID
			if ida.Type < idb.Type {
				return true
			}
			if ida.Type > idb.Type {
				return false
			}
			return ida.Instance < idb.Instance
		})
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Name < devices[j].Name
	})
}

func sortRawConfig(rawConfig []map[string]any) {
	sort.Slice(rawConfig, func(i, j int) bool {
		n1, n2 := rawConfig[i]["name"].(string), rawConfig[j]["name"].(string)
		if n1 == n2 {
			if _, ok := rawConfig[i]["kind"]; ok {
				if _, ok2 := rawConfig[j]["kind"]; ok2 {
					return rawConfig[i]["kind"].(trait.Name) < rawConfig[j]["kind"].(trait.Name)
				}
			} else if _, ok := rawConfig[i]["type"]; ok {
				if _, ok2 := rawConfig[j]["type"]; ok2 {
					return rawConfig[i]["type"].(trait.Name) < rawConfig[j]["type"].(trait.Name)
				}
			}
		}
		return n1 < n2
	})
}
