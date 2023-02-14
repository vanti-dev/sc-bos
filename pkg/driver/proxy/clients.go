package proxy

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock/button"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"google.golang.org/grpc"
)

// newApiClientForTrait returns the *ApiClient implementation for the named trait.
// For example passing trait.OnOff would return traits.NewOnOffApiClient.
// Returns nil if the trait is not known.
func newApiClientForTrait(conn *grpc.ClientConn, t trait.Name) any {
	// todo: I feel this should really live in sc-golang somewhere

	switch t {
	case trait.AirQualitySensor:
		return traits.NewAirQualitySensorApiClient(conn)
	case trait.AirTemperature:
		return traits.NewAirTemperatureApiClient(conn)
	case trait.Booking:
		return traits.NewBookingApiClient(conn)
	case trait.BrightnessSensor:
		return traits.NewBrightnessSensorApiClient(conn)
	case trait.Channel:
		return traits.NewChannelApiClient(conn)
	case trait.Count:
		return traits.NewCountApiClient(conn)
	case trait.Electric:
		return traits.NewElectricApiClient(conn)
	case trait.Emergency:
		return traits.NewEmergencyApiClient(conn)
	case trait.EnergyStorage:
		return traits.NewEnergyStorageApiClient(conn)
	case trait.EnterLeaveSensor:
		return traits.NewEnterLeaveSensorApiClient(conn)
	case trait.ExtendRetract:
		return traits.NewExtendRetractApiClient(conn)
	case trait.FanSpeed:
		return traits.NewFanSpeedApiClient(conn)
	case trait.Hail:
		return traits.NewHailApiClient(conn)
	case trait.InputSelect:
		return traits.NewInputSelectApiClient(conn)
	case trait.Light:
		return traits.NewLightApiClient(conn)
	case trait.LockUnlock:
		return traits.NewLockUnlockApiClient(conn)
	case trait.Metadata:
		return traits.NewMetadataApiClient(conn)
	case trait.Microphone:
		return traits.NewMicrophoneApiClient(conn)
	case trait.Mode:
		return traits.NewModeApiClient(conn)
	case trait.MotionSensor:
		return traits.NewMotionSensorApiClient(conn)
	case trait.OccupancySensor:
		return traits.NewOccupancySensorApiClient(conn)
	case trait.OnOff:
		return traits.NewOnOffApiClient(conn)
	case trait.OpenClose:
		return traits.NewOpenCloseApiClient(conn)
	case trait.Parent:
		return traits.NewParentApiClient(conn)
	case trait.Publication:
		return traits.NewPublicationApiClient(conn)
	case trait.Ptz:
		return traits.NewPtzApiClient(conn)
	case trait.Speaker:
		return traits.NewSpeakerApiClient(conn)
	case trait.Vending:
		return traits.NewVendingApiClient(conn)

		// sc-bos private traits
	case button.TraitName:
		return gen.NewButtonApiClient(conn)
	case meter.TraitName:
		return gen.NewMeterApiClient(conn)
	}
	return nil
}
