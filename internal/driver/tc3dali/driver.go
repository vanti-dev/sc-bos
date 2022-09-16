package tc3dali

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/bridge"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
	"github.com/vanti-dev/twincat3-ads-go/pkg/adsdll"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
	"go.uber.org/zap"
)

type Config struct {
	driver.BaseConfig
	ADS   ADSConfig   `json:"ads"`
	Buses []BusConfig `json:"buses"`
}

type ADSConfig struct {
	NetID NetID  `json:"netID"`
	Port  uint16 `json:"port"`
}

type NetID ads.NetId

func (n *NetID) UnmarshalJSON(buf []byte) error {
	_, err := fmt.Fscanf(bytes.NewBuffer(buf), "%d.%d.%d.%d.%d.%d", &n[0], &n[1], &n[2], &n[3], &n[4], &n[5])
	return err
}

type BusConfig struct {
	Name           string                `json:"name"`
	ControlGear    []ControlGearConfig   `json:"controlGear"`
	ControlDevices []ControlDeviceConfig `json:"controlDevices"`
	BridgePrefix   string                `json:"bridgePrefix"`
}

type ControlGearConfig struct {
	Name         string  `json:"name"`
	ShortAddress uint8   `json:"shortAddress"`
	Groups       []uint8 `json:"groups"`
}

type ControlDeviceConfig struct {
	Name          string         `json:"name"`
	ShortAddress  uint8          `json:"shortAddress"`
	InstanceTypes []InstanceType `json:"instanceTypes"`
}

type InstanceType string

const (
	InstanceTypeOccupancySensor = "occupancySensor"
)

func (c *ControlDeviceConfig) hasInstance(want InstanceType) bool {
	for _, have := range c.InstanceTypes {
		if have == want {
			return true
		}
	}
	return false
}

func Factory(ctx context.Context, services driver.Services, rawConfig json.RawMessage) (driver.Driver, error) {
	var config Config
	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// TODO: validate and normalise the config
	port, err := adsdll.Connect()
	if err != nil {
		return nil, err
	}
	dev, err := device.Open(port, ads.Addr{
		NetId: ads.NetId(config.ADS.NetID),
		Port:  config.ADS.Port,
	})

	for _, bus := range config.Buses {
		services.Tasks.Spawn(ctx, bus.Name, BusTask(bus, dev, services))
	}

	return &driverImpl{
		config: config,
	}, nil
}

var _ driver.Factory = Factory

type driverImpl struct {
	config Config
}

func (d *driverImpl) Name() string {
	return d.config.Name
}

func BusTask(config BusConfig, dev device.Device, services driver.Services) task.Task {
	return func(ctx context.Context) (next task.Next, err error) {
		bridgeConfig := &bridge.Config{
			Device:                  dev,
			Logger:                  services.Logger.Named("bridge"),
			BridgeFBName:            config.BridgePrefix + "_bridge",
			ResponseMailboxName:     config.BridgePrefix + "_response",
			NotificationMailboxName: config.BridgePrefix + "_notification",
		}
		busBridge, err := bridgeConfig.Connect()
		if err != nil {
			services.Logger.Error("DALI bus bridge initialisation failure", zap.Error(err),
				zap.String("busName", config.Name),
				zap.String("prefix", config.BridgePrefix))
		}

		err = RunBus(ctx, config, busBridge, services)
		return
	}
}

func RunBus(ctx context.Context, config BusConfig, busBridge bridge.Dali, services driver.Services) error {
	busServer := &controlGearServer{
		bus:      busBridge,
		addrType: bridge.Broadcast,
	}
	services.Node.Announce(config.Name,
		node.HasTrait(trait.Light, node.WithClients(light.WrapApi(busServer))),
	)

	knownGroups := make(map[uint8]struct{})
	for _, gear := range config.ControlGear {
		gearName := gear.Name
		if gearName == "" {
			gearName = fmt.Sprintf("%s/control-gear/%d", config.Name, gear.ShortAddress)
		}
		for _, g := range gear.Groups {
			knownGroups[g] = struct{}{}
		}

		gearServer := &controlGearServer{
			bus:      busBridge,
			addrType: bridge.Short,
			addr:     gear.ShortAddress,
		}
		services.Node.Announce(gearName,
			node.HasTrait(trait.Light, node.WithClients(light.WrapApi(gearServer))),
		)
	}

	for _, dev := range config.ControlDevices {
		// for now we only support occupancy sensors
		if !dev.hasInstance(InstanceTypeOccupancySensor) {
			continue
		}

		deviceName := dev.Name
		if deviceName == "" {
			deviceName = fmt.Sprintf("%s/control-device/%d", config.Name, dev.ShortAddress)
		}

		devServer := &controlDeviceServer{
			bus:       busBridge,
			shortAddr: dev.ShortAddress,
			occupancy: resource.NewValue(resource.WithInitialValue(&traits.Occupancy{
				State: traits.Occupancy_STATE_UNSPECIFIED,
			})),
			logger: services.Logger.Named("control-device").With(zap.String("sc-device-name", deviceName)),
		}
		services.Node.Announce(deviceName,
			node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensor.WrapApi(devServer))),
		)
	}

	for group := range knownGroups {
		groupName := fmt.Sprintf("%s/groups/%d", config.Name, group)
		groupServer := &controlGearServer{
			bus:      busBridge,
			addr:     group,
			addrType: bridge.Group,
		}
		services.Node.Announce(groupName,
			node.HasTrait(trait.Light, node.WithClients(light.WrapApi(groupServer))),
		)
	}

	return nil
}
