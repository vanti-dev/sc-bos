package mock

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
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
	"github.com/smart-core-os/sc-golang/pkg/trait/openclose"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-golang/pkg/trait/publication"
	"github.com/smart-core-os/sc-golang/pkg/trait/vending"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/button"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/modepb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"
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
				clients, slc := newMockClient(dt.trait, device.Name, d.logger)
				if len(clients) == 0 {
					d.logger.Sugar().Warnf("Cannot create mock client %s::%s", dt.name, dt.trait)
				} else {
					traitOpts = append(traitOpts, node.WithClients(clients...))

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

func newMockClient(traitName trait.Name, deviceName string, logger *zap.Logger) ([]any, service.Lifecycle) {
	switch traitName {
	case trait.AirQualitySensor:
		co2 := rand.Float32()
		voc := rand.Float32()
		ap := rand.Float32()
		ir := rand.Float32()
		model := traits.AirQuality{
			CarbonDioxideLevel:       &co2,
			VolatileOrganicCompounds: &voc,
			AirPressure:              &ap,
			Comfort:                  0,
			InfectionRisk:            &ir,
		}
		return []any{airqualitysensor.WrapApi(airqualitysensor.NewModelServer(airqualitysensor.NewModel(&model)))}, nil
	case trait.AirTemperature:
		h := rand.Float32()
		t := 15 + (rand.Float64() * 10)
		model := traits.AirTemperature{
			Mode:               traits.AirTemperature_AUTO,
			AmbientTemperature: &types.Temperature{ValueCelsius: t},
			AmbientHumidity:    &h,
			TemperatureGoal: &traits.AirTemperature_TemperatureSetPoint{
				TemperatureSetPoint: &types.Temperature{ValueCelsius: t},
			},
		}
		return []any{airtemperature.WrapApi(airtemperature.NewModelServer(airtemperature.NewModel(&model)))}, nil
	case trait.Booking:
		return []any{booking.WrapApi(booking.NewModelServer(booking.NewModel()))}, nil
	case trait.BrightnessSensor:
		// todo: return []any{brightnesssensor.WrapApi(brightnesssensor.NewModelServer(brightnesssensor.NewModel()))}, nil
		return nil, nil
	case trait.Channel:
		// todo: return []any{channel.WrapApi(channel.NewModelServer(channel.NewModel())), nil
		return nil, nil
	case trait.Count:
		// todo: return []any{count.WrapApi(count.NewModelServer(count.NewModel())), nil
		return nil, nil
	case trait.Electric:
		model := electric.NewModel(clock.Real())
		return []any{electric.WrapApi(electric.NewModelServer(model))}, auto.Electric(model)
	case trait.Emergency:
		// todo: return []any{emergency.WrapApi(emergency.NewModelServer(emergency.NewModel()))}, nil
		return nil, nil
	case trait.EnergyStorage:
		return []any{energystorage.WrapApi(energystorage.NewModelServer(energystorage.NewModel()))}, nil
	case trait.EnterLeaveSensor:
		return []any{enterleavesensor.WrapApi(enterleavesensor.NewModelServer(enterleavesensor.NewModel()))}, nil
	case trait.ExtendRetract:
		// todo: return []any{extendretract.WrapApi(extendretract.NewModelServer(extendretract.NewModel()))}, nil
		return nil, nil
	case trait.FanSpeed:
		return []any{fanspeed.WrapApi(fanspeed.NewModelServer(fanspeed.NewModel()))}, nil
	case trait.Hail:
		return []any{hail.WrapApi(hail.NewModelServer(hail.NewModel()))}, nil
	case trait.InputSelect:
		// todo: return []any{inputselect.WrapApi(inputselect.NewModelServer(inputselect.NewModel()))}, nil
		return nil, nil
	case trait.Light:
		// todo: return []any{light.WrapApi(light.NewModelServer(light.NewModel())), nil
		return []any{light.WrapApi(light.NewMemoryDevice())}, nil
	case trait.LockUnlock:
		// todo: return []any{lockunlock.WrapApi(lockunlock.NewModelServer(lockunlock.NewModel()))}, nil
		return nil, nil
	case trait.Metadata:
		return []any{metadata.WrapApi(metadata.NewModelServer(metadata.NewModel()))}, nil
	case trait.Microphone:
		// todo: return []any{microphone.WrapApi(microphone.NewModelServer(microphone.NewModel()))}, nil
		return nil, nil
	case trait.Mode:
		model := mode.NewModel()
		modes := model.Modes()
		infoServer := &modepb.InfoServer{Modes: &traits.ModesSupport{AvailableModes: modes}}
		return []any{mode.WrapApi(mode.NewModelServer(model)), mode.WrapInfo(infoServer)}, nil
	case trait.MotionSensor:
		// todo: return []any{motionsensor.WrapApi(motionsensor.NewModelServer(motionsensor.NewModel()))}, nil
		return nil, nil
	case trait.OccupancySensor:
		model := occupancysensor.NewModel(&traits.Occupancy{})
		return []any{occupancysensor.WrapApi(occupancysensor.NewModelServer(model))}, auto.OccupancySensorAuto(model)
	case trait.OnOff:
		return []any{onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(traits.OnOff_STATE_UNSPECIFIED)))}, nil
	case trait.OpenClose:
		model := openclose.NewModel()
		return []any{openclose.WrapApi(openclose.NewModelServer(model))}, auto.OpenClose(model)
	case trait.Parent:
		return []any{parent.WrapApi(parent.NewModelServer(parent.NewModel()))}, nil
	case trait.Publication:
		return []any{publication.WrapApi(publication.NewModelServer(publication.NewModel()))}, nil
	case trait.Ptz:
		// todo: return []any{ptz.WrapApi(ptz.NewModelServer(ptz.NewModel()))}, nil
		return nil, nil
	case trait.Speaker:
		// todo: return []any{speaker.WrapApi(speaker.NewModelServer(speaker.NewModel())), nil
		return nil, nil
	case trait.Vending:
		return []any{vending.WrapApi(vending.NewModelServer(vending.NewModel()))}, nil

	case accesspb.TraitName:
		model := accesspb.NewModel()
		return []any{gen.WrapAccessApi(accesspb.NewModelServer(model))}, auto.Access(model)
	case button.TraitName:
		return []any{gen.WrapButtonApi(button.NewModelServer(button.NewModel(gen.ButtonState_UNPRESSED)))}, nil
	case meter.TraitName:
		model := meter.NewModel()
		info := &meter.InfoServer{MeterReading: &gen.MeterReadingSupport{
			ResourceSupport: &types.ResourceSupport{
				Readable:   true,
				Writable:   true,
				Observable: true,
			},
			Unit: "kWh",
		}}
		return []any{gen.WrapMeterApi(meter.NewModelServer(model)), gen.WrapMeterInfo(info)}, auto.MeterAuto(model)
	case statuspb.TraitName:
		model := statuspb.NewModel()
		// set an initial value or Pull methods can hang
		model.UpdateProblem(&gen.StatusLog_Problem{Name: deviceName, Level: gen.StatusLog_NOMINAL})
		return []any{gen.WrapStatusApi(statuspb.NewModelServer(model))}, nil
	case udmipb.TraitName:
		return []any{gen.WrapUdmiService(auto.NewUdmiServer(logger, deviceName))}, nil
	}

	return nil, nil
}
