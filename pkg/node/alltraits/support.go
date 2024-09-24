package alltraits

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/booking"
	"github.com/smart-core-os/sc-golang/pkg/trait/brightnesssensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/channel"
	"github.com/smart-core-os/sc-golang/pkg/trait/count"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/smart-core-os/sc-golang/pkg/trait/emergency"
	"github.com/smart-core-os/sc-golang/pkg/trait/energystorage"
	"github.com/smart-core-os/sc-golang/pkg/trait/enterleavesensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/extendretract"
	"github.com/smart-core-os/sc-golang/pkg/trait/fanspeed"
	"github.com/smart-core-os/sc-golang/pkg/trait/hail"
	"github.com/smart-core-os/sc-golang/pkg/trait/inputselect"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/lockunlock"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/trait/microphone"
	"github.com/smart-core-os/sc-golang/pkg/trait/mode"
	"github.com/smart-core-os/sc-golang/pkg/trait/motionsensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/trait/openclose"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/smart-core-os/sc-golang/pkg/trait/ptz"
	"github.com/smart-core-os/sc-golang/pkg/trait/publication"
	"github.com/smart-core-os/sc-golang/pkg/trait/speaker"
	"github.com/smart-core-os/sc-golang/pkg/trait/vending"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/dalipb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/emergencylight"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/mqttpb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/wastepb"

	"github.com/vanti-dev/sc-bos/pkg/gentrait/button"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/color"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

var traitSupport = map[trait.Name]func(s node.Supporter){
	trait.AirQualitySensor: func(s node.Supporter) {
		{
			r := airqualitysensor.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(airqualitysensor.WrapApi(r)))
		}
		{
			r := airqualitysensor.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(airqualitysensor.WrapInfo(r)))
		}
		{
			r := gen.NewAirQualitySensorHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapAirQualitySensorHistory(r)))
		}
	},
	trait.AirTemperature: func(s node.Supporter) {
		{
			r := airtemperature.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(airtemperature.WrapApi(r)))
		}
		{
			r := airtemperature.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(airtemperature.WrapInfo(r)))
		}
		{
			r := gen.NewAirTemperatureHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapAirTemperatureHistory(r)))
		}
	},
	trait.Booking: func(s node.Supporter) {
		{
			r := booking.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(booking.WrapApi(r)))
		}
		{
			r := booking.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(booking.WrapInfo(r)))
		}
	},
	trait.BrightnessSensor: func(s node.Supporter) {
		{
			r := brightnesssensor.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(brightnesssensor.WrapApi(r)))
		}
		{
			r := brightnesssensor.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(brightnesssensor.WrapInfo(r)))
		}
	},
	trait.Channel: func(s node.Supporter) {
		{
			r := channel.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(channel.WrapApi(r)))
		}
		{
			r := channel.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(channel.WrapInfo(r)))
		}
	},
	trait.Count: func(s node.Supporter) {
		{
			r := count.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(count.WrapApi(r)))
		}
		{
			r := count.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(count.WrapInfo(r)))
		}
	},
	trait.Electric: func(s node.Supporter) {
		{
			r := electric.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(electric.WrapApi(r)))
		}
		{
			r := electric.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(electric.WrapInfo(r)))
		}
	},
	trait.Emergency: func(s node.Supporter) {
		{
			r := emergency.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(emergency.WrapApi(r)))
		}
		{
			r := emergency.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(emergency.WrapInfo(r)))
		}
		{
			r := gen.NewElectricHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapElectricHistory(r)))
		}
	},
	trait.EnergyStorage: func(s node.Supporter) {
		{
			r := energystorage.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(energystorage.WrapApi(r)))
		}
		{
			r := energystorage.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(energystorage.WrapInfo(r)))
		}
	},
	trait.EnterLeaveSensor: func(s node.Supporter) {
		{
			r := enterleavesensor.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(enterleavesensor.WrapApi(r)))
		}
		{
			r := enterleavesensor.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(enterleavesensor.WrapInfo(r)))
		}
	},
	trait.ExtendRetract: func(s node.Supporter) {
		{
			r := extendretract.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(extendretract.WrapApi(r)))
		}
		{
			r := extendretract.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(extendretract.WrapInfo(r)))
		}
	},
	trait.FanSpeed: func(s node.Supporter) {
		{
			r := fanspeed.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(fanspeed.WrapApi(r)))
		}
		{
			r := fanspeed.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(fanspeed.WrapInfo(r)))
		}
	},
	trait.Hail: func(s node.Supporter) {
		{
			r := hail.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(hail.WrapApi(r)))
		}
		{
			r := hail.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(hail.WrapInfo(r)))
		}
	},
	trait.InputSelect: func(s node.Supporter) {
		{
			r := inputselect.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(inputselect.WrapApi(r)))
		}
		{
			r := inputselect.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(inputselect.WrapInfo(r)))
		}
	},
	trait.Light: func(s node.Supporter) {
		{
			r := light.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(light.WrapApi(r)))
		}
		{
			r := light.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(light.WrapInfo(r)))
		}
	},
	trait.LockUnlock: func(s node.Supporter) {
		{
			r := lockunlock.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(lockunlock.WrapApi(r)))
		}
		{
			r := lockunlock.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(lockunlock.WrapInfo(r)))
		}
	},
	trait.Metadata: func(s node.Supporter) {
		{
			r := metadata.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(metadata.WrapApi(r)))
		}
		{
			r := metadata.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(metadata.WrapInfo(r)))
		}
	},
	trait.Microphone: func(s node.Supporter) {
		{
			r := microphone.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(microphone.WrapApi(r)))
		}
		{
			r := microphone.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(microphone.WrapInfo(r)))
		}
	},
	trait.Mode: func(s node.Supporter) {
		{
			r := mode.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(mode.WrapApi(r)))
		}
		{
			r := mode.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(mode.WrapInfo(r)))
		}
	},
	trait.MotionSensor: func(s node.Supporter) {
		{
			r := motionsensor.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(motionsensor.WrapApi(r)))
		}
		{
			r := motionsensor.NewSensorInfoRouter()
			s.Support(node.Routing(r), node.Clients(motionsensor.WrapSensorInfo(r)))
		}
	},
	trait.OccupancySensor: func(s node.Supporter) {
		{
			r := occupancysensor.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(occupancysensor.WrapApi(r)))
		}
		{
			r := occupancysensor.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(occupancysensor.WrapInfo(r)))
		}
		{
			r := gen.NewOccupancySensorHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapOccupancySensorHistory(r)))
		}
	},
	trait.OnOff: func(s node.Supporter) {
		{
			r := onoff.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(onoff.WrapApi(r)))
		}
		{
			r := onoff.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(onoff.WrapInfo(r)))
		}
	},
	trait.OpenClose: func(s node.Supporter) {
		{
			r := openclose.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(openclose.WrapApi(r)))
		}
		{
			r := openclose.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(openclose.WrapInfo(r)))
		}
	},
	trait.Parent: func(s node.Supporter) {
		{
			r := parent.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(parent.WrapApi(r)))
		}
		{
			r := parent.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(parent.WrapInfo(r)))
		}
	},
	trait.Publication: func(s node.Supporter) {
		{
			r := publication.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(publication.WrapApi(r)))
		}
	},
	trait.Ptz: func(s node.Supporter) {
		{
			r := ptz.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(ptz.WrapApi(r)))
		}
		{
			r := ptz.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(ptz.WrapInfo(r)))
		}
	},
	trait.Speaker: func(s node.Supporter) {
		{
			r := speaker.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(speaker.WrapApi(r)))
		}
		{
			r := speaker.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(speaker.WrapInfo(r)))
		}
	},
	trait.Vending: func(s node.Supporter) {
		{
			r := vending.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(vending.WrapApi(r)))
		}
		{
			r := vending.NewInfoRouter()
			s.Support(node.Routing(r), node.Clients(vending.WrapInfo(r)))
		}
	},

	// sc-bos private traits
	accesspb.TraitName: func(s node.Supporter) {
		{
			r := gen.NewAccessApiRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapAccessApi(r)))
		}
	},
	button.TraitName: func(s node.Supporter) {
		r := gen.NewButtonApiRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapButtonApi(r)))
	},
	color.TraitName: func(s node.Supporter) {
		{
			r := gen.NewColorApiRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapColorApi(r)))
		}
		{
			r := gen.NewColorInfoRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapColorInfo(r)))
		}
	},
	dalipb.TraitName: func(s node.Supporter) {
		{
			r := gen.NewDaliApiRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapDaliApi(r)))
		}
	},
	emergencylight.TraitName: func(s node.Supporter) {
		// We don't do anything here, there is no trait that this supports exclusively.
		// Manually expose the DaliApi on the node if you need this functionality.
	},
	meter.TraitName: func(s node.Supporter) {
		{
			r := gen.NewMeterApiRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapMeterApi(r)))
		}
		{
			r := gen.NewMeterInfoRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapMeterInfo(r)))
		}
		{
			r := gen.NewMeterHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapMeterHistory(r)))
		}
	},
	mqttpb.TraitName: func(s node.Supporter) {
		r := gen.NewMqttServiceRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapMqttService(r)))
	},
	statuspb.TraitName: func(s node.Supporter) {
		{
			r := gen.NewStatusApiRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapStatusApi(r)))
		}
		{
			r := gen.NewStatusHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapStatusHistory(r)))
		}
	},
	udmipb.TraitName: func(s node.Supporter) {
		r := gen.NewUdmiServiceRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapUdmiService(r)))
	},
	wastepb.TraitName: func(s node.Supporter) {
		r := gen.NewWasteApiRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapWasteApi(r)))
	},
}

// AddSupport adds support to n for all known traits.
func AddSupport(n node.Supporter) {
	for _, f := range traitSupport {
		f(n)
	}
}

// AddSupportFor adds support to n for the given traits.
func AddSupportFor(n node.Supporter, traits ...trait.Name) {
	for _, t := range traits {
		if f, ok := traitSupport[t]; ok {
			f(n)
		}
	}
}
