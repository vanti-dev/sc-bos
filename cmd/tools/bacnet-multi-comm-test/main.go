package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/property"
	bactypes "github.com/vanti-dev/gobacnet/types"

	"github.com/vanti-dev/sc-bos/pkg/app/appconf"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
)

var (
	dir         = flag.String("dir", ".", "directory to scan for config")
	configFile  = flag.String("config-file", "area-controller.local.json", "file name to look for and load")
	resultsFile = flag.String("results-file", "results", "file name to save results to, timestamp will be appended")
)

func main() {
	flag.Parse()
	start := time.Now()
	err := run()
	if err != nil {
		log.Printf("uhoh: %s", err)
	} else {
		elapsed := time.Since(start)
		log.Printf("done in %s", elapsed)
	}
}

// run will recursively scan dir for any files named configFile, and load any bacnet driver config from them.
// each bacnet device will be sent a request, and if successful will be marked "responding", and all configured
// objects will then be read
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
		client.Log.Out = &bytes.Buffer{} // else client.Close() will close os.Stderr
		if err != nil {
			return err
		}
		client.Log.SetLevel(logrus.InfoLevel)
		address, err := client.LocalUDPAddress()
		if err == nil {
			log.Printf("bacnet client configured: local=%s, localInterface=%s, localPort:%d", address, cfg.LocalInterface, cfg.LocalPort)
		}

		for _, device := range cfg.Devices {
			log.Printf("check deviceId: %d", device.ID)
			dev, err := bacnet.FindDevice(client, device)
			res := &result{
				name:     device.Name,
				deviceId: fmt.Sprintf("%d", device.ID),
				value:    "",
			}
			if err == nil {
				res.responding = true
			}
			if device.Comm != nil {
				res.address = device.Comm.IP.String()
			}
			results[key] = append(results[key], res)

			if res.responding {
				for _, obj := range device.Objects {
					objRes := &result{
						name:     fmt.Sprintf("%s/%s", device.Name, obj.Name),
						deviceId: obj.ID.String(),
						value:    "",
					}
					value, err := readProp(client, dev, obj, property.PresentValue)
					if err == nil {
						objRes.responding = true
						objRes.value = value
					}
					results[key] = append(results[key], objRes)
				}
			}
		}
		client.Close()
	}
	fileName := *resultsFile + "_" + time.Now().Format("2006-01-02T15_04") + ".csv"
	return writeResults(fileName, results)
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

func readProp(client *gobacnet.Client, device bactypes.Device, obj config.Object, prop property.ID) (any, error) {
	req := bactypes.ReadPropertyData{
		Object: bactypes.Object{
			ID: bactypes.ObjectID(obj.ID),
			Properties: []bactypes.Property{
				{ID: prop, ArrayIndex: bactypes.ArrayAll},
			},
		},
	}
	res, err := client.ReadProperty(device, req)
	if err != nil {
		return nil, err
	}
	if len(res.Object.Properties) == 0 {
		// Shouldn't happen, but has on occasion. I guess it depends how the device responds to our request
		return nil, errors.New("zero length object properties")
	}
	value := res.Object.Properties[0].Data
	if strings.HasPrefix(obj.ID.String(), "Binary") {
		value = value == 1
	}
	return value, nil
}

func writeResults(fileName string, results map[string][]*result) error {
	var rows [][]string
	rows = append(rows, []string{"Name", "Address", "Device ID", "Responding", "Value"})
	for _, res := range results {
		for _, re := range res {
			rows = append(rows, re.toRow())
		}
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	// clear the file
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(rows)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

type result struct {
	name       string
	address    string
	deviceId   string
	responding bool
	value      any
}

func (r *result) toRow() []string {
	return []string{r.name, r.address, r.deviceId, boolYesNo(r.responding), anyToString(r.value)}
}

func anyToString(value any) string {
	switch v := value.(type) {
	case bool:
		return fmt.Sprintf("%t", v)
	case float32:
		return fmt.Sprintf("%.2f", v)
	case uint32:
		return fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%s", value)
}

func boolYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
