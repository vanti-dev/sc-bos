package mock

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
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
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-golang/pkg/trait/publication"
	"github.com/smart-core-os/sc-golang/pkg/trait/vending"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/vanti-dev/sc-bos/pkg/block"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/button"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
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

func (_ factory) ConfigBlocks() []block.Block {
	return config.Blocks
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
				clients, slc := newMockClient(traitMd, device.Name, d.logger)
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

func newMockClient(traitMd *traits.TraitMetadata, deviceName string, logger *zap.Logger) ([]wrap.ServiceUnwrapper, service.Lifecycle) {
	switch trait.Name(traitMd.Name) {
	case trait.AirQualitySensor:
		model := airqualitysensor.NewModel(airqualitysensor.WithInitialAirQuality(auto.GetAirQualityState()))
		return []wrap.ServiceUnwrapper{airqualitysensor.WrapApi(airqualitysensor.NewModelServer(model))}, auto.AirQualitySensorAuto(model)
	case trait.AirTemperature:
		h := rand.Float32() * 100
		t := 15 + (rand.Float64() * 10)
		initial := traits.AirTemperature{
			Mode:               traits.AirTemperature_AUTO,
			AmbientTemperature: &types.Temperature{ValueCelsius: t},
			AmbientHumidity:    &h,
			TemperatureGoal: &traits.AirTemperature_TemperatureSetPoint{
				TemperatureSetPoint: &types.Temperature{ValueCelsius: t},
			},
		}
		model := airtemperature.NewModel(airtemperature.WithInitialAirTemperature(&initial))
		return []wrap.ServiceUnwrapper{airtemperature.WrapApi(airtemperature.NewModelServer(model))}, nil
	case trait.Booking:
		return []wrap.ServiceUnwrapper{booking.WrapApi(booking.NewModelServer(booking.NewModel()))}, nil
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
		model := electric.NewModel()
		return []wrap.ServiceUnwrapper{electric.WrapApi(electric.NewModelServer(model))}, auto.Electric(model)
	case trait.Emergency:
		// todo: return []any{emergency.WrapApi(emergency.NewModelServer(emergency.NewModel()))}, nil
		return nil, nil
	case trait.EnergyStorage:
		return []wrap.ServiceUnwrapper{energystorage.WrapApi(energystorage.NewModelServer(energystorage.NewModel()))}, nil
	case trait.EnterLeaveSensor:
		return []wrap.ServiceUnwrapper{enterleavesensor.WrapApi(enterleavesensor.NewModelServer(enterleavesensor.NewModel()))}, nil
	case trait.ExtendRetract:
		// todo: return []any{extendretract.WrapApi(extendretract.NewModelServer(extendretract.NewModel()))}, nil
		return nil, nil
	case trait.FanSpeed:
		presets := []fanspeed.Preset{
			{Name: "off", Percentage: 0},
			{Name: "low", Percentage: 15},
			{Name: "med", Percentage: 40},
			{Name: "high", Percentage: 75},
			{Name: "full", Percentage: 100},
		}
		model := fanspeed.NewModel(fanspeed.WithPresets(presets...))
		return []wrap.ServiceUnwrapper{fanspeed.WrapApi(fanspeed.NewModelServer(model))}, auto.FanSpeed(model, presets...)
	case trait.Hail:
		return []wrap.ServiceUnwrapper{hail.WrapApi(hail.NewModelServer(hail.NewModel()))}, nil
	case trait.InputSelect:
		// todo: return []any{inputselect.WrapApi(inputselect.NewModelServer(inputselect.NewModel()))}, nil
		return nil, nil
	case trait.Light:
		server := light.NewModelServer(light.NewModel(
			light.WithPreset(0, &traits.LightPreset{Name: "off", Title: "Off"}),
			light.WithPreset(40, &traits.LightPreset{Name: "low", Title: "Low"}),
			light.WithPreset(60, &traits.LightPreset{Name: "med", Title: "Normal"}),
			light.WithPreset(80, &traits.LightPreset{Name: "high", Title: "High"}),
			light.WithPreset(100, &traits.LightPreset{Name: "full", Title: "Full"}),
		))
		return []wrap.ServiceUnwrapper{light.WrapApi(server), light.WrapInfo(server)}, nil
	case trait.LockUnlock:
		// todo: return []any{lockunlock.WrapApi(lockunlock.NewModelServer(lockunlock.NewModel()))}, nil
		return nil, nil
	case trait.Metadata:
		return []wrap.ServiceUnwrapper{metadata.WrapApi(metadata.NewModelServer(metadata.NewModel()))}, nil
	case trait.Microphone:
		// todo: return []any{microphone.WrapApi(microphone.NewModelServer(microphone.NewModel()))}, nil
		return nil, nil
	case trait.Mode:
		return mockMode(traitMd, deviceName, logger)
	case trait.MotionSensor:
		// todo: return []any{motionsensor.WrapApi(motionsensor.NewModelServer(motionsensor.NewModel()))}, nil
		return nil, nil
	case trait.OccupancySensor:
		model := occupancysensor.NewModel()
		return []wrap.ServiceUnwrapper{occupancysensor.WrapApi(occupancysensor.NewModelServer(model))}, auto.OccupancySensorAuto(model)
	case trait.OnOff:
		return []wrap.ServiceUnwrapper{onoff.WrapApi(onoff.NewModelServer(onoff.NewModel()))}, nil
	case trait.OpenClose:
		return mockOpenClose(traitMd, deviceName, logger)
	case trait.Parent:
		return []wrap.ServiceUnwrapper{parent.WrapApi(parent.NewModelServer(parent.NewModel()))}, nil
	case trait.Publication:
		return []wrap.ServiceUnwrapper{publication.WrapApi(publication.NewModelServer(publication.NewModel()))}, nil
	case trait.Ptz:
		// todo: return []any{ptz.WrapApi(ptz.NewModelServer(ptz.NewModel()))}, nil
		return nil, nil
	case trait.Speaker:
		// todo: return []any{speaker.WrapApi(speaker.NewModelServer(speaker.NewModel())), nil
		return nil, nil
	case trait.Vending:
		return []wrap.ServiceUnwrapper{vending.WrapApi(vending.NewModelServer(vending.NewModel()))}, nil

	case accesspb.TraitName:
		model := accesspb.NewModel()
		return []wrap.ServiceUnwrapper{gen.WrapAccessApi(accesspb.NewModelServer(model))}, auto.Access(model)
	case button.TraitName:
		return []wrap.ServiceUnwrapper{gen.WrapButtonApi(button.NewModelServer(button.NewModel(gen.ButtonState_UNPRESSED)))}, nil
	case meter.TraitName:
		var (
			unit string
			ok   bool
		)
		if unit, ok = traitMd.GetMore()["unit"]; !ok {
			unit = "kWh"
		}

		model := meter.NewModel()
		info := &meter.InfoServer{MeterReading: &gen.MeterReadingSupport{
			ResourceSupport: &types.ResourceSupport{
				Readable:   true,
				Writable:   true,
				Observable: true,
			},
			Unit: unit,
		}}
		return []wrap.ServiceUnwrapper{gen.WrapMeterApi(meter.NewModelServer(model)), gen.WrapMeterInfo(info)}, auto.MeterAuto(model)
	case statuspb.TraitName:
		model := statuspb.NewModel()
		// set an initial value or Pull methods can hang
		_, _ = model.UpdateProblem(&gen.StatusLog_Problem{Name: deviceName, Level: gen.StatusLog_NOMINAL})
		return []wrap.ServiceUnwrapper{gen.WrapStatusApi(statuspb.NewModelServer(model))}, auto.Status(model, deviceName)
	case udmipb.TraitName:
		return []wrap.ServiceUnwrapper{gen.WrapUdmiService(auto.NewUdmiServer(logger, deviceName))}, nil
	}

	return nil, nil
}
