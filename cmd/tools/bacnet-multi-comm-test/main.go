package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vanti-dev/gobacnet"
	bactypes "github.com/vanti-dev/gobacnet/types"

	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
)

var (
	dir         = flag.String("dir", ".", "directory to scan for config")
	configFile  = flag.String("config-file", "area-controller.local.json", "file name to look for and load")
	resultsFile = flag.String("results-file", "results.csv", "file name to save results to")
)

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		log.Printf("uhoh: %s", err)
	}
}

// run will recursively scan dir for any files named configFile, and load any bacnet driver config from them.
// each bacnet device will be sent a request, and if successful will be marked "found"
// results are saved to a CSV
func run() error {
	bacnetConfigs, err := loadConfigs()
	if err != nil {
		return err
	}
	results := make(map[string][]*result, len(bacnetConfigs))
	for key, cfg := range bacnetConfigs {
		log.Printf("config: %s, devices: %d", key, len(cfg.Devices))

		client, err := gobacnet.NewClient(cfg.LocalInterface, int(cfg.LocalPort))
		if err != nil {
			return err
		}
		client.Log.SetLevel(logrus.InfoLevel)
		address, err := client.LocalUDPAddress()
		if err == nil {
			log.Printf("bacnet client configured: local=%s, localInterface=%s, localPort:%d", address, cfg.LocalInterface, cfg.LocalPort)
		}

		for _, device := range cfg.Devices {
			_, err := findDevice(client, device)
			res := &result{
				name:     device.Name,
				deviceId: fmt.Sprintf("%d", device.ID),
			}
			if err == nil {
				res.found = true
			}
			if device.Comm != nil {
				res.address = device.Comm.IP.String()
			}
			results[key] = append(results[key], res)
		}
		client.Close()
	}
	var rows [][]string
	rows = append(rows, []string{"Name", "Address", "Device ID", "Found"})
	for _, res := range results {
		for _, re := range res {
			rows = append(rows, re.toRow())
		}
	}
	f, err := os.OpenFile(*resultsFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.WriteAll(rows)
	if err != nil {
		return err
	}
	w.Flush()
	log.Printf("results written to: %s", f.Name())
	return nil
}

func loadConfigs() (map[string]config.Root, error) {
	var configPaths []string
	fileSystem := os.DirFS(*dir)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, *configFile) {
			configPaths = append(configPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	log.Printf("configPaths: %s", strings.Join(configPaths, ";"))

	bacnetConfigs := make(map[string]config.Root, len(configPaths))

	for _, configPath := range configPaths {
		var localConfig appconf.Config
		rawLocalConfig, err := os.ReadFile(filepath.Join(*dir, configPath))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rawLocalConfig, &localConfig)
		if err != nil {
			return nil, fmt.Errorf("%s JSON unmarshal: %w", configPath, err)
		}
		for _, driver := range localConfig.Drivers {
			if driver.Type == bacnet.DriverName {
				bacnetConfig, err := config.ReadBytes(driver.Raw)
				if err != nil {
					return nil, err
				}
				bacnetConfigs[configPath+":"+driver.Name] = bacnetConfig
			}
		}
	}
	return bacnetConfigs, nil
}

func findDevice(client *gobacnet.Client, device config.Device) (bactypes.Device, error) {
	fail := func(err error) (bactypes.Device, error) {
		return bactypes.Device{}, err
	}

	if device.Comm == nil {
		id := device.ID
		is, err := client.WhoIs(int(id), int(id))
		if err != nil {
			return fail(err)
		}
		if len(is) == 0 {
			return fail(fmt.Errorf("no devices found (via WhoIs) with id %d", id))
		}
		return is[0], nil
	}

	addr, err := device.Comm.ToAddress()
	if err != nil {
		return fail(err)
	}
	bacDevices, err := client.RemoteDevices(addr, device.ID)
	if err != nil {
		return fail(err)
	}
	return bacDevices[0], nil
}

type result struct {
	name     string
	address  string
	deviceId string
	found    bool
}

func (r *result) toRow() []string {
	return []string{r.name, r.address, r.deviceId, fmt.Sprintf("%t", r.found)}
}
