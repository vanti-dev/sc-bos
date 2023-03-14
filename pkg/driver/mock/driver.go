package mock

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
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
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/button"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/util/maps"
)

const DriverName = "mock"

var Factory driver.Factory = factory{}

type factory struct{}

func (_ factory) New(services driver.Services) service.Lifecycle {
	return NewDriver(services)
}

func NewDriver(services driver.Services) *Driver {
	d := &Driver{
		announcer: services.Node,
		known:     make(map[deviceTrait]node.Undo),
	}
	d.Service = service.New(d.applyConfig, service.WithOnStop[config.Root](d.Clean))
	d.logger = services.Logger.Named(DriverName)
	return d
}

type Driver struct {
	*service.Service[config.Root]

	logger    *zap.Logger
	announcer node.Announcer
	known     map[deviceTrait]node.Undo
}

type deviceTrait struct {
	name  string
	trait trait.Name
}

func (d *Driver) Clean() {
	for _, undo := range d.known {
		undo()
	}
	d.known = make(map[deviceTrait]node.Undo)
}

func (d *Driver) applyConfig(_ context.Context, cfg config.Root) error {
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

			var traitOpts []node.TraitOption
			var undo []node.Undo
			if u, ok := d.known[dt]; ok {
				undo = append(undo, u)
			}

			if _, ok := d.known[dt]; !ok {
				client, slc := newMockClient(dt.trait)
				if client == nil {
					d.logger.Sugar().Warnf("Cannot create mock client %s::%s", dt.name, dt.trait)
				} else {
					traitOpts = append(traitOpts, node.WithClients(client))

					// start any mock trait automations - e.g. updating occupancy sensors
					if slc != nil {
						_, err := slc.Start()
						if err != nil {
							d.logger.Sugar().Warnf("Unable to start mock trait automation %s::%s %v", dt.name, dt.trait, err)
						} else {
							undo = append(undo, func() {
								_, _ = slc.Stop()
							})
						}
					}
				}
			}
			undo = append(undo, d.announcer.Announce(dt.name, node.HasTrait(dt.trait, traitOpts...)))
			d.known[dt] = node.UndoAll(undo...)
			undos = append(undos, undo...)
		}

		dt.trait = ""
		d.known[dt] = node.UndoAll(undos...)
	}

	for k, undo := range toUndo {
		undo()
		delete(d.known, k)
	}

	return nil
}

func newMockClient(traitName trait.Name) (any, service.Lifecycle) {
	switch traitName {
	case trait.AirQualitySensor:
		return airqualitysensor.WrapApi(airqualitysensor.NewModelServer(airqualitysensor.NewModel(&traits.AirQuality{}))), nil
	case trait.AirTemperature:
		return airtemperature.WrapApi(airtemperature.NewModelServer(airtemperature.NewModel(&traits.AirTemperature{}))), nil
	case trait.Booking:
		return booking.WrapApi(booking.NewModelServer(booking.NewModel())), nil
	case trait.BrightnessSensor:
		// todo: return brightnesssensor.WrapApi(brightnesssensor.NewModelServer(brightnesssensor.NewModel())), nil
		return nil, nil
	case trait.Channel:
		// todo: return channel.WrapApi(channel.NewModelServer(channel.NewModel())), nil
		return nil, nil
	case trait.Count:
		// todo: return count.WrapApi(count.NewModelServer(count.NewModel())), nil
		return nil, nil
	case trait.Electric:
		return electric.WrapApi(electric.NewModelServer(electric.NewModel(clock.Real()))), nil
	case trait.Emergency:
		// todo: return emergency.WrapApi(emergency.NewModelServer(emergency.NewModel())), nil
		return nil, nil
	case trait.EnergyStorage:
		return energystorage.WrapApi(energystorage.NewModelServer(energystorage.NewModel())), nil
	case trait.EnterLeaveSensor:
		return enterleavesensor.WrapApi(enterleavesensor.NewModelServer(enterleavesensor.NewModel())), nil
	case trait.ExtendRetract:
		// todo: return extendretract.WrapApi(extendretract.NewModelServer(extendretract.NewModel())), nil
		return nil, nil
	case trait.FanSpeed:
		return fanspeed.WrapApi(fanspeed.NewModelServer(fanspeed.NewModel())), nil
	case trait.Hail:
		return hail.WrapApi(hail.NewModelServer(hail.NewModel())), nil
	case trait.InputSelect:
		// todo: return inputselect.WrapApi(inputselect.NewModelServer(inputselect.NewModel())), nil
		return nil, nil
	case trait.Light:
		// todo: return light.WrapApi(light.NewModelServer(light.NewModel())), nil
		return light.WrapApi(light.NewMemoryDevice()), nil
	case trait.LockUnlock:
		// todo: return lockunlock.WrapApi(lockunlock.NewModelServer(lockunlock.NewModel())), nil
		return nil, nil
	case trait.Metadata:
		return metadata.WrapApi(metadata.NewModelServer(metadata.NewModel())), nil
	case trait.Microphone:
		// todo: return microphone.WrapApi(microphone.NewModelServer(microphone.NewModel())), nil
		return nil, nil
	case trait.Mode:
		return mode.WrapApi(mode.NewModelServer(mode.NewModel())), nil
	case trait.MotionSensor:
		// todo: return motionsensor.WrapApi(motionsensor.NewModelServer(motionsensor.NewModel())), nil
		return nil, nil
	case trait.OccupancySensor:
		model := occupancysensor.NewModel(&traits.Occupancy{})
		return occupancysensor.WrapApi(occupancysensor.NewModelServer(model)), occupancySensorAuto(model)
	case trait.OnOff:
		return onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(traits.OnOff_STATE_UNSPECIFIED))), nil
	case trait.OpenClose:
		// todo: return openclose.WrapApi(openclose.NewModelServer(openclose.NewModel())), nil
		return nil, nil
	case trait.Parent:
		return parent.WrapApi(parent.NewModelServer(parent.NewModel())), nil
	case trait.Publication:
		return publication.WrapApi(publication.NewModelServer(publication.NewModel())), nil
	case trait.Ptz:
		// todo: return ptz.WrapApi(ptz.NewModelServer(ptz.NewModel())), nil
		return nil, nil
	case trait.Speaker:
		// todo: return speaker.WrapApi(speaker.NewModelServer(speaker.NewModel())), nil
		return nil, nil
	case trait.Vending:
		return vending.WrapApi(vending.NewModelServer(vending.NewModel())), nil

	case button.TraitName:
		return gen.WrapButtonApi(button.NewModelServer(button.NewModel(gen.ButtonState_UNPRESSED))), nil
	case meter.TraitName:
		model := meter.NewModel()
		return gen.WrapMeterApi(meter.NewModelServer(model)), meterAuto(model)
	}

	return nil, nil
}

func occupancySensorAuto(model *occupancysensor.Model) *service.Service[string] {
	slc := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					state := traits.Occupancy_State(rand.Intn(3) + 1)
					_, _ = model.SetOccupancy(&traits.Occupancy{State: state}, resource.WithUpdatePaths("state"))
				}
			}
		}()
		return nil
	}), service.WithParser(func(data []byte) (string, error) {
		return string(data), nil
	}))
	_, _ = slc.Configure([]byte{}) // call configure to ensure we load when start is called.
	return slc
}

func meterAuto(model *meter.Model) *service.Service[string] {
	slc := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			timer := time.NewTimer((30 * time.Second) + time.Duration(rand.Float32())*time.Minute)
			start := &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			}
			value := rand.Float32() * 100
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					value += rand.Float32() * 100
					state := gen.MeterReading{
						Usage:     value,
						StartTime: start,
						EndTime: &timestamppb.Timestamp{
							Seconds: time.Now().Unix(),
							Nanos:   0,
						},
					}
					_, _ = model.UpdateMeterReading(&state)
					timer = time.NewTimer((30 * time.Second) + time.Duration(rand.Float32())*time.Minute)
				}
			}
		}()
		return nil
	}), service.WithParser(func(data []byte) (string, error) {
		return string(data), nil
	}))
	_, _ = slc.Configure([]byte{}) // call configure to ensure we load when start is called.
	return slc
}
