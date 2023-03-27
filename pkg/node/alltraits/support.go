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
	"github.com/vanti-dev/sc-bos/pkg/gentrait/emergencylight"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/mqttpb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"

	"github.com/vanti-dev/sc-bos/pkg/gentrait/button"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/color"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

var traitSupport = map[trait.Name]func(s node.Supporter){
	trait.AirQualitySensor: func(s node.Supporter) {
		r := airqualitysensor.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(airqualitysensor.WrapApi(r)))
	},
	trait.AirTemperature: func(s node.Supporter) {
		r := airtemperature.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(airtemperature.WrapApi(r)))
	},
	trait.Booking: func(s node.Supporter) {
		r := booking.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(booking.WrapApi(r)))
	},
	trait.BrightnessSensor: func(s node.Supporter) {
		r := brightnesssensor.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(brightnesssensor.WrapApi(r)))
	},
	trait.Channel: func(s node.Supporter) {
		r := channel.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(channel.WrapApi(r)))
	},
	trait.Count: func(s node.Supporter) {
		r := count.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(count.WrapApi(r)))
	},
	trait.Electric: func(s node.Supporter) {
		r := electric.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(electric.WrapApi(r)))
	},
	trait.Emergency: func(s node.Supporter) {
		{
			r := emergency.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(emergency.WrapApi(r)))
		}
		{
			r := gen.NewElectricHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapElectricHistory(r)))
		}
	},
	trait.EnergyStorage: func(s node.Supporter) {
		r := energystorage.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(energystorage.WrapApi(r)))
	},
	trait.EnterLeaveSensor: func(s node.Supporter) {
		r := enterleavesensor.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(enterleavesensor.WrapApi(r)))
	},
	trait.ExtendRetract: func(s node.Supporter) {
		r := extendretract.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(extendretract.WrapApi(r)))
	},
	trait.FanSpeed: func(s node.Supporter) {
		r := fanspeed.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(fanspeed.WrapApi(r)))
	},
	trait.Hail: func(s node.Supporter) {
		r := hail.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(hail.WrapApi(r)))
	},
	trait.InputSelect: func(s node.Supporter) {
		r := inputselect.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(inputselect.WrapApi(r)))
	},
	trait.Light: func(s node.Supporter) {
		r := light.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(light.WrapApi(r)))
	},
	trait.LockUnlock: func(s node.Supporter) {
		r := lockunlock.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(lockunlock.WrapApi(r)))
	},
	trait.Metadata: func(s node.Supporter) {
		r := metadata.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(metadata.WrapApi(r)))
	},
	trait.Microphone: func(s node.Supporter) {
		r := microphone.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(microphone.WrapApi(r)))
	},
	trait.Mode: func(s node.Supporter) {
		r := mode.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(mode.WrapApi(r)))
	},
	trait.MotionSensor: func(s node.Supporter) {
		r := motionsensor.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(motionsensor.WrapApi(r)))
	},
	trait.OccupancySensor: func(s node.Supporter) {
		{
			r := occupancysensor.NewApiRouter()
			s.Support(node.Routing(r), node.Clients(occupancysensor.WrapApi(r)))
		}
		{
			r := gen.NewOccupancySensorHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapOccupancySensorHistory(r)))
		}
	},
	trait.OnOff: func(s node.Supporter) {
		r := onoff.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(onoff.WrapApi(r)))
	},
	trait.OpenClose: func(s node.Supporter) {
		r := openclose.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(openclose.WrapApi(r)))
	},
	trait.Parent: func(s node.Supporter) {
		r := parent.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(parent.WrapApi(r)))
	},
	trait.Publication: func(s node.Supporter) {
		r := publication.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(publication.WrapApi(r)))
	},
	trait.Ptz: func(s node.Supporter) {
		r := ptz.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(ptz.WrapApi(r)))
	},
	trait.Speaker: func(s node.Supporter) {
		r := speaker.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(speaker.WrapApi(r)))
	},
	trait.Vending: func(s node.Supporter) {
		r := vending.NewApiRouter()
		s.Support(node.Routing(r), node.Clients(vending.WrapApi(r)))
	},

	// sc-bos private traits
	button.TraitName: func(s node.Supporter) {
		r := gen.NewButtonApiRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapButtonApi(r)))
	},
	color.TraitName: func(s node.Supporter) {
		r := gen.NewColorApiRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapColorApi(r)))
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
			r := gen.NewMeterHistoryRouter()
			s.Support(node.Routing(r), node.Clients(gen.WrapMeterHistory(r)))
		}
	},
	mqttpb.TraitName: func(s node.Supporter) {
		r := gen.NewMqttServiceRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapMqttService(r)))
	},
	udmipb.TraitName: func(s node.Supporter) {
		r := gen.NewUdmiServiceRouter()
		s.Support(node.Routing(r), node.Clients(gen.WrapUdmiService(r)))
	},
}

// AddSupport adds support to n for all known traits.
func AddSupport(n node.Supporter) {
	for _, f := range traitSupport {
		f(n)
	}
}
