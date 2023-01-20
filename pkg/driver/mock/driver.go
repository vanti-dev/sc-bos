package mock

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/booking"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-golang/pkg/trait/enterleavesensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/fanspeed"
	"github.com/smart-core-os/sc-golang/pkg/trait/hail"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/trait/mode"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-golang/pkg/trait/publication"
	"github.com/smart-core-os/sc-golang/pkg/trait/vending"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock/config"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/util/maps"
)

const DriverName = "mock"

var Factory driver.Factory = factory{}

type factory struct{}

func (_ factory) New(services driver.Services) task.Starter {
	return NewDriver(services)
}

func NewDriver(services driver.Services) task.Starter {
	d := &Driver{
		announcer: services.Node,
		known:     make(map[deviceTrait]node.Undo),
	}
	d.Lifecycle = task.NewLifecycle(d.applyConfig)
	d.Logger = services.Logger.Named(DriverName)
	return d
}

type Driver struct {
	*task.Lifecycle[config.Root]

	announcer node.Announcer
	known     map[deviceTrait]node.Undo
}

type deviceTrait struct {
	name  string
	trait trait.Name
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	toUndo := maps.Clone(d.known)
	for _, device := range cfg.Devices {
		var undos []node.Undo
		dt := deviceTrait{name: device.Name}

		// the device is still in the config, don't delete it
		delete(toUndo, dt)

		if u, ok := d.known[dt]; ok {
			undos = append(undos, u)
		}
		undos = append(undos, d.announcer.Announce(dt.name, node.HasMetadata(device.Metadata)))

		for _, traitMd := range device.Traits {
			dt.trait = trait.Name(traitMd.Name)

			// the trait is still in the device config, don't delete it
			delete(toUndo, dt)

			traitOpts := []node.TraitOption{
				node.NoAddMetadata(),
			}
			if _, ok := d.known[dt]; !ok {
				client := newMockClient(dt.trait)
				if client == nil {
					d.Logger.Sugar().Warnf("Cannot create mock client %s::%s", dt.name, dt.trait)
				} else {
					traitOpts = append(traitOpts, node.WithClients(client))
				}
			}
			undo := d.announcer.Announce(dt.name, node.HasTrait(dt.trait, traitOpts...))
			if u, ok := d.known[dt]; ok {
				d.known[dt] = node.UndoAll(u, undo)
			} else {
				d.known[dt] = undo
			}
			undos = append(undos, undo)
		}

		dt.trait = ""
		d.known[dt] = node.UndoAll(undos...)
	}

	return nil
}

func newMockClient(traitName trait.Name) any {
	switch traitName {
	case trait.AirQualitySensor:
		return airqualitysensor.WrapApi(airqualitysensor.NewModelServer(airqualitysensor.NewModel(&traits.AirQuality{})))
	case trait.AirTemperature:
		return airtemperature.WrapApi(airtemperature.NewModelServer(airtemperature.NewModel(&traits.AirTemperature{})))
	case trait.Booking:
		return booking.WrapApi(booking.NewModelServer(booking.NewModel()))
	case trait.BrightnessSensor:
		// todo: return brightnesssensor.WrapApi(brightnesssensor.NewModelServer(brightnesssensor.NewModel()))
		return nil
	case trait.Channel:
		// todo: return channel.WrapApi(channel.NewModelServer(channel.NewModel()))
		return nil
	case trait.Count:
		// todo: return count.WrapApi(count.NewModelServer(count.NewModel()))
		return nil
	case trait.Electric:
		return electric.WrapApi(electric.NewModelServer(electric.NewModel(clock.Real())))
	case trait.Emergency:
		// todo: return emergency.WrapApi(emergency.NewModelServer(emergency.NewModel()))
		return nil
	case trait.EnergyStorage:
		return energystorage.WrapApi(energystorage.NewModelServer(energystorage.NewModel()))
	case trait.EnterLeaveSensor:
		return enterleavesensor.WrapApi(enterleavesensor.NewModelServer(enterleavesensor.NewModel()))
	case trait.ExtendRetract:
		// todo: return extendretract.WrapApi(extendretract.NewModelServer(extendretract.NewModel()))
		return nil
	case trait.FanSpeed:
		return fanspeed.WrapApi(fanspeed.NewModelServer(fanspeed.NewModel()))
	case trait.Hail:
		return hail.WrapApi(hail.NewModelServer(hail.NewModel()))
	case trait.InputSelect:
		// todo: return inputselect.WrapApi(inputselect.NewModelServer(inputselect.NewModel()))
		return nil
	case trait.Light:
		// todo: return light.WrapApi(light.NewModelServer(light.NewModel()))
		return light.WrapApi(light.NewMemoryDevice())
	case trait.LockUnlock:
		// todo: return lockunlock.WrapApi(lockunlock.NewModelServer(lockunlock.NewModel()))
		return nil
	case trait.Metadata:
		return metadata.WrapApi(metadata.NewModelServer(metadata.NewModel()))
	case trait.Microphone:
		// todo: return microphone.WrapApi(microphone.NewModelServer(microphone.NewModel()))
		return nil
	case trait.Mode:
		return mode.WrapApi(mode.NewModelServer(mode.NewModel()))
	case trait.MotionSensor:
		// todo: return motionsensor.WrapApi(motionsensor.NewModelServer(motionsensor.NewModel()))
		return nil
	case trait.OccupancySensor:
		return occupancysensor.WrapApi(occupancysensor.NewModelServer(occupancysensor.NewModel(&traits.Occupancy{})))
	case trait.OnOff:
		return onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(traits.OnOff_STATE_UNSPECIFIED)))
	case trait.OpenClose:
		// todo: return openclose.WrapApi(openclose.NewModelServer(openclose.NewModel()))
		return nil
	case trait.Parent:
		return parent.WrapApi(parent.NewModelServer(parent.NewModel()))
	case trait.Publication:
		return publication.WrapApi(publication.NewModelServer(publication.NewModel()))
	case trait.Ptz:
		// todo: return ptz.WrapApi(ptz.NewModelServer(ptz.NewModel()))
		return nil
	case trait.Speaker:
		// todo: return speaker.WrapApi(speaker.NewModelServer(speaker.NewModel()))
		return nil
	case trait.Vending:
		return vending.WrapApi(vending.NewModelServer(vending.NewModel()))
	}

	return nil
}
