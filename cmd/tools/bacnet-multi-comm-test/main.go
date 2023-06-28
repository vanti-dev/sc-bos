package main

import (
	"bytes"
	"context"
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
	devicesOnly = flag.Bool("devices-only", false, "only check devices, not objects")
	timeout     = flag.Duration("timeout", time.Second*2, "timeout for requests")
	localPort   = flag.Int("local-port", 0, "local port to use for bacnet requests, overrides any found in config")
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

		port := int(cfg.LocalPort)
		if *localPort > 0 {
			port = *localPort
		}
		client, err := gobacnet.NewClient(cfg.LocalInterface, port)
		if err != nil {
			return err
		}
		client.Log.Out = &bytes.Buffer{} // else client.Close() will close os.Stderr
		client.Log.SetLevel(logrus.InfoLevel)
		address, err := client.LocalUDPAddress()
		if err == nil {
			log.Printf("bacnet client configured: local=%s, localInterface=%s, localPort:%d", address, cfg.LocalInterface, cfg.LocalPort)
		}

		for _, device := range cfg.Devices {
			ctx, cancel := context.WithTimeout(context.Background(), *timeout)
			defer cancel()
			if device.Comm != nil {
				log.Printf("check device: %s (id:%d)", device.Comm.IP, device.ID)
			} else {
				log.Printf("check deviceId: %d", device.ID)
			}
			dev, err := bacnet.FindDevice(ctx, client, device)
			res := &result{
				name:     device.Name,
				deviceId: fmt.Sprintf("%d", device.ID),
				value:    "",
			}
			if err == nil {
				res.responding = true
				readId := fmt.Sprintf("%d", dev.ID.Instance)
				if readId != res.deviceId {
					log.Printf("configured device %v is really %v", res.deviceId, readId)
					res.deviceId = readId
				}
				if res.name == "" {
					res.name, err = readObjectName(client, dev, dev.ID)
					device.Name = res.name
					if err != nil {
						log.Printf("unable to discover device name %v %v", dev.ID, err)
					}
				}
			}
			if comm := device.Comm; comm != nil {
				res.address = comm.IP.String()
				if des := comm.Destination; des != nil {
					res.network = des.Network
					res.mac = string(des.Address)
				}
			}
			results[key] = append(results[key], res)

			if res.responding && !*devicesOnly {
				cfgObjects := device.Objects

				// Discover device objects if we're asked to
				if shouldDiscoverObjects(cfg, device) {
					ctx, cancel := context.WithTimeout(context.Background(), (*timeout)*10)
					defer cancel()
					objects, err := client.Objects(ctx, dev)
					if err != nil {
						log.Printf("Unable to discover objects %v : %v", device.ID, err)
					}

					uniqueObjIds := make(map[string]bool, len(cfgObjects))
					for _, object := range cfgObjects {
						uniqueObjIds[object.ID.String()] = true
					}
					for _, obj := range objects.Objects {
						for _, object := range obj {
							id := object.ID
							if !uniqueObjIds[id.String()] {
								uniqueObjIds[id.String()] = true
								// log.Printf("discovered device object %v:%v (name=%v)", id.Type, id.Instance, object.Name)
								cfgObjects = append(cfgObjects, config.Object{ID: config.ObjectID(id), Name: object.Name})
							}
						}
					}
				}

				// Discover all the names for objects that don't have them
				for i, obj := range cfgObjects {
					if obj.Name == "" {
						obj.Name, err = readObjectName(client, dev, bactypes.ObjectID(obj.ID))
						if err != nil {
							log.Printf("unable to get name for %v/%v", device.Name, obj.ID)
							continue
						}
						cfgObjects[i] = obj
					}
				}

				// Fetch the PresentValue for the objects and record them in the results
				for _, obj := range cfgObjects {
					objRes := &result{
						name:     fmt.Sprintf("%s/%s", device.Name, obj.Name),
						deviceId: obj.ID.String(),
						value:    "",
					}
					ctx, cancel := context.WithTimeout(context.Background(), *timeout)
					defer cancel()
					value, err := readProp(ctx, client, dev, obj, property.PresentValue)
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

func readObjectName(client *gobacnet.Client, dev bactypes.Device, objId bactypes.ObjectID) (string, error) {
	ctx, clean := context.WithTimeout(context.Background(), *timeout)
	propVal, err := client.ReadProperty(ctx, dev, bactypes.ReadPropertyData{Object: bactypes.Object{ID: objId, Properties: []bactypes.Property{
		{ID: property.ObjectName, ArrayIndex: bactypes.ArrayAll},
	}}})
	clean()
	if err != nil {
		return "", err
	}
	for _, p := range propVal.Object.Properties {
		if p.ID == property.ObjectName {
			return p.Data.(string), nil
		}
	}

	return "", errors.New("response did not include ObjectName property")
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

func readProp(ctx context.Context, client *gobacnet.Client, device bactypes.Device, obj config.Object, prop property.ID) (any, error) {
	req := bactypes.ReadPropertyData{
		Object: bactypes.Object{
			ID: bactypes.ObjectID(obj.ID),
			Properties: []bactypes.Property{
				{ID: prop, ArrayIndex: bactypes.ArrayAll},
			},
		},
	}
	res, err := client.ReadProperty(ctx, device, req)
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
	rows = append(rows, []string{"Name", "Address", "Network", "MAC", "BACnet ID", "Responding", "Value"})
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
	network    uint16
	mac        string
	deviceId   string
	responding bool
	value      any
}

func (r *result) toRow() []string {
	return []string{r.name, r.address, netString(r.network), r.mac, r.deviceId, boolYesNo(r.responding), anyToString(r.value)}
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

func netString(n uint16) string {
	if n == 0 {
		return ""
	}
	return fmt.Sprintf("%d", n)
}

func boolYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
