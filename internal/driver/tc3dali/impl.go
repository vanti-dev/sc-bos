//go:build !notc3dali

package tc3dali

import (
	"context"
	"fmt"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/bridge"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
	"github.com/vanti-dev/twincat3-ads-go/pkg/adsdll"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func applyConfig(ctx context.Context, services driver.Services, config Config) error {
	// TODO: validate and normalise the config
	port, err := adsdll.Connect()
	if err != nil {
		return fmt.Errorf("adsdll.Connect: %w", err)
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
		return errs
	}
	return nil
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
func InitBus(ctx context.Context, config BusConfig, busBridge dali.Dali, services driver.Services) error {
	busServer := &controlGearServer{
		bus:      busBridge,
		addrType: dali.Broadcast,
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
			addrType: dali.Short,
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
			addrType: dali.Group,
		}
		services.Node.Announce(groupName,
			node.HasTrait(trait.Light, node.WithClients(light.WrapApi(groupServer))),
		)
	}

	return nil
}
