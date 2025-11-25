package bacnet

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/smart-core-os/gobacnet"
	bactypes "github.com/smart-core-os/gobacnet/types"
	"github.com/smart-core-os/sc-bos/pkg/app/appconf"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/adapt"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/config"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/known"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/status"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/statuspb"
	"github.com/smart-core-os/sc-bos/pkg/node"
)

func TestYabeRoom(t *testing.T) {
	if _, ok := os.LookupEnv("YABE"); !ok {
		t.Skip("Skipping test that relies on YABE room simulator")
	}
	client, err := gobacnet.NewClient("bridge100", 0)
	if err != nil {
		t.Fatalf("NewClient error %v", err)
	}
	t.Cleanup(client.Close)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	iAm, err := client.WhoIs(ctx, 800_000, 900_000)
	if err != nil {
		t.Fatalf("WhoIs error %v", err)
	}

	t.Logf("IAm %+v", iAm)

	netAddr, err := iAm[0].Addr.UDPAddr()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Connecting to %v", netAddr)
	bacAddr := bactypes.UDPToAddress(&netAddr)
	t.Logf("After conversion: iAm Addr:%v, udp addr:%v, bac addr:%v", iAm[0].Addr, netAddr, bacAddr)
	devices, err := client.RemoteDevices(ctx, bacAddr, iAm[0].ID.Instance)
	if err != nil {
		t.Fatalf("RemoteDevices error %v", err)
	}
	t.Logf("RemoteDevices %+v", devices)
}

func TestSiteFaults(t *testing.T) {
	t.SkipNow() // only run this test when connected to site
	ctx, cleanup := context.WithTimeout(context.Background(), 30*time.Second)
	defer cleanup()

	appConfig, err := appconf.LoadLocalConfig("/path/to/site/config/dir", "bms.part.json")
	if err != nil {
		t.Fatal(err)
	}

	var bacnetConfig config.Root
	for _, driver := range appConfig.Drivers {
		if driver.Type == DriverName {
			bacnetConfig, err = config.ReadBytes(driver.Raw)
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	if bacnetConfig.Type == "" {
		t.Fatal("No BACnet driver config found")
	}

	client, err := gobacnet.NewClient("bridge100", 0, gobacnet.WithLogLevel(logrus.InfoLevel))
	if err != nil {
		t.Fatalf("NewClient error %v", err)
	}
	t.Cleanup(client.Close)

	knownMap := known.NewMap()
	statusMap := statuspb.NewMap(node.AnnouncerFunc(func(name string, features ...node.Feature) node.Undo {
		return node.NilUndo
	}))

	monitor := status.NewMonitor(client, known.SyncContext(&sync.Mutex{}, knownMap), statusMap)
	monitor.Logger, _ = zap.NewDevelopment(zap.AddStacktrace(zap.FatalLevel))

	for _, device := range bacnetConfig.Devices {
		ctx, cleanup := context.WithTimeout(ctx, 2*time.Second)
		bacDevice, err := FindDevice(ctx, client, device)
		cleanup()
		if err != nil {
			log.Printf("ERR: FindDevice %v %v", device.Name, err)
			continue
		}
		knownMap.StoreDevice(adapt.DeviceName(device), bacDevice, 0)
		for _, co := range device.Objects {
			bo := bactypes.Object{
				ID: bactypes.ObjectID(co.ID),
			}
			_ = knownMap.StoreObject(bacDevice, adapt.ObjectName(co), bo)
		}

		monitor.AddDevice(device.Name, device)
	}

	err = monitor.Poll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
