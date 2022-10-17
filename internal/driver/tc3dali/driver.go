package tc3dali

import (
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
	"go.uber.org/multierr"
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

//goland:noinspection GoMixedReceiverTypes
func (n *NetID) UnmarshalJSON(buf []byte) error {
	var str string
	err := json.Unmarshal(buf, &str)
	if err != nil {
		return err
	}

	parsed, err := ParseNetID(str)
	if err != nil {
		return err
	}
	*n = parsed
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (n NetID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d.%d.%d.%d.%d.%d", n[0], n[1], n[2], n[3], n[4], n[5])), nil
}

func ParseNetID(raw string) (n NetID, err error) {
	_, err = fmt.Sscanf(raw, "%d.%d.%d.%d.%d.%d", &n[0], &n[1], &n[2], &n[3], &n[4], &n[5])
	return
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
	InstanceTypeOccupancySensor InstanceType = "occupancySensor"
)

func (c *ControlDeviceConfig) hasInstance(want InstanceType) bool {
	for _, have := range c.InstanceTypes {
		if have == want {
			return true
		}
	}
	return false
}

const DriverName = "tc3dali"

func Factory(ctx context.Context, services driver.Services, rawConfig json.RawMessage) (driver.Driver, error) {
	var config Config
	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// TODO: validate and normalise the config
	port, err := adsdll.Connect()
	if err != nil {
		return nil, fmt.Errorf("adsdll.Connect: %w", err)
	}

	var (
		spawned int
		errs    error
	)
	for _, bus := range config.Buses {
		// create a new device.Device for each bus, because Device isn't safe for concurrent access
		dev, err := device.Open(port, ads.Addr{
			NetId: ads.NetId(config.ADS.NetID),
			Port:  config.ADS.Port,
		})
		if err != nil {
			errs = multierr.Append(errs, err)
			services.Logger.Error("failed to connect to ADS PLC Device", zap.Error(err),
				zap.Uint8s("netID", config.ADS.NetID[:]), zap.Uint16("port", config.ADS.Port))
			continue
		}
		services.Tasks.Spawn(ctx, bus.Name, BusTask(bus, dev, services))
		spawned++
	}

	// if we wanted to run some DALI buses, but all failed, return failure for the driver
	if len(config.Buses) > 0 && spawned == 0 {
		return nil, errs
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
			return
		}

		err = InitBus(ctx, config, busBridge, services)
		return
	}
}

// InitBus exposes the DALI devices on a single DALI bus over Smart Core, by registering them with services.Node.
// It exposes the following devices:
//   - The bus itself,  implementing the Light trait with DALI broadcast commands
//   - Each control gear, implementing Light
//   - Each declared control gear group, implementing Light
//   - Each occupancy control device, implementing OccupancySensor
func InitBus(ctx context.Context, config BusConfig, busBridge bridge.Dali, services driver.Services) error {
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
			logger:           services.Logger.Named("control-device").With(zap.String("sc-device-name", deviceName)),
			enableEventsOnce: newOnce(),
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
