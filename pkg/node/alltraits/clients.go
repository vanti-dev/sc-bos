package alltraits

import (
	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/emergencylight"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/mqttpb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/button"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/color"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
)

type ClientFactory func(conn grpc.ClientConnInterface) any

var apiClientFactories map[trait.Name]ClientFactory

// RegisterAPIClientFactory registers a {trait}ApiClient factory for the named trait.
// This factory will take president over the default generated factory.
// Should be called before any call to APIClient, typically in init().
func RegisterAPIClientFactory(t trait.Name, f ClientFactory) {
	if apiClientFactories == nil {
		apiClientFactories = make(map[trait.Name]ClientFactory)
	}
	apiClientFactories[t] = f
}

// APIClient returns the {trait}ApiClient implementation for the named trait.
// For example passing trait.OnOff would return traits.NewOnOffApiClient.
// Returns nil if the trait is not known.
func APIClient(conn grpc.ClientConnInterface, t trait.Name) any {
	// todo: I feel this should really live in sc-golang somewhere
	if d, ok := apiClientFactories[t]; ok {
		return d(conn)
	}

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
	case color.TraitName:
		return gen.NewColorApiClient(conn)
	case emergencylight.TraitName:
		return gen.NewDaliApiClient(conn)
	case meter.TraitName:
		return gen.NewMeterApiClient(conn)
	case mqttpb.TraitName:
		return gen.NewMqttServiceClient(conn)
	case statuspb.TraitName:
		return gen.NewStatusApiClient(conn)
	case udmipb.TraitName:
		return gen.NewUdmiServiceClient(conn)
	}
	return nil
}

var historyClientFactories map[trait.Name]ClientFactory

// RegisterHistoryClientFactory registers a {trait}HistoryClient factory for the named trait.
// This factory will take president over the default generated factory.
// Should be called before any call to HistoryClient, typically in init().
func RegisterHistoryClientFactory(t trait.Name, f ClientFactory) {
	if historyClientFactories == nil {
		historyClientFactories = make(map[trait.Name]ClientFactory)
	}
	historyClientFactories[t] = f
}

// HistoryClient returns the {trait}HistoryClient implementation for the named trait.
// For example passing trait.Meter would return traits.NewMeterHistoryClient.
// Returns nil if the trait is not known.
func HistoryClient(conn grpc.ClientConnInterface, t trait.Name) any {
	if d, ok := historyClientFactories[t]; ok {
		return d(conn)
	}

	switch t {
	// Smart Core traits
	case trait.Electric:
		return gen.NewElectricHistoryClient(conn)
	case trait.OccupancySensor:
		return gen.NewOccupancySensorHistoryClient(conn)

		// (not yet) Smart Core traits
	case meter.TraitName:
		return gen.NewMeterHistoryClient(conn)
	case statuspb.TraitName:
		return gen.NewStatusHistoryClient(conn)
	default:
		return nil
	}
}

var infoClientFactories map[trait.Name]ClientFactory

// RegisterInfoClientFactory registers a {trait}InfoClient factory for the named trait.
// This factory will take president over the default generated factory.
// Should be called before any call to InfoClient, typically in init().
func RegisterInfoClientFactory(t trait.Name, f ClientFactory) {
	if infoClientFactories == nil {
		infoClientFactories = make(map[trait.Name]ClientFactory)
	}
	infoClientFactories[t] = f
}

// InfoClient returns the {trait}InfoClient implementation for the named trait.
// For example passing trait.Meter would return traits.NewMeterInfoClient.
// Returns nil if the trait is not known or has no info aspect.
func InfoClient(conn grpc.ClientConnInterface, t trait.Name) any {
	// todo: I feel this should really live in sc-golang somewhere
	if d, ok := infoClientFactories[t]; ok {
		return d(conn)
	}

	switch t {
	case trait.AirQualitySensor:
		return traits.NewAirQualitySensorInfoClient(conn)
	case trait.AirTemperature:
		return traits.NewAirTemperatureInfoClient(conn)
	case trait.Booking:
		return traits.NewBookingInfoClient(conn)
	case trait.BrightnessSensor:
		return traits.NewBrightnessSensorInfoClient(conn)
	case trait.Channel:
		return traits.NewChannelInfoClient(conn)
	case trait.Count:
		return traits.NewCountInfoClient(conn)
	case trait.Electric:
		return traits.NewElectricInfoClient(conn)
	case trait.Emergency:
		return traits.NewEmergencyInfoClient(conn)
	case trait.EnergyStorage:
		return traits.NewEnergyStorageInfoClient(conn)
	case trait.EnterLeaveSensor:
		return traits.NewEnterLeaveSensorInfoClient(conn)
	case trait.ExtendRetract:
		return traits.NewExtendRetractInfoClient(conn)
	case trait.FanSpeed:
		return traits.NewFanSpeedInfoClient(conn)
	case trait.Hail:
		return traits.NewHailInfoClient(conn)
	case trait.InputSelect:
		return traits.NewInputSelectInfoClient(conn)
	case trait.Light:
		return traits.NewLightInfoClient(conn)
	case trait.LockUnlock:
		return traits.NewLockUnlockInfoClient(conn)
	case trait.Metadata:
		return traits.NewMetadataInfoClient(conn)
	case trait.Microphone:
		return traits.NewMicrophoneInfoClient(conn)
	case trait.Mode:
		return traits.NewModeInfoClient(conn)
	case trait.MotionSensor:
		// return traits.NewMotionSensorInfoClient(conn)
		return nil
	case trait.OccupancySensor:
		return traits.NewOccupancySensorInfoClient(conn)
	case trait.OnOff:
		return traits.NewOnOffInfoClient(conn)
	case trait.OpenClose:
		return traits.NewOpenCloseInfoClient(conn)
	case trait.Parent:
		return traits.NewParentInfoClient(conn)
	case trait.Publication:
		// return traits.NewPublicationInfoClient(conn)
		return nil
	case trait.Ptz:
		return traits.NewPtzInfoClient(conn)
	case trait.Speaker:
		return traits.NewSpeakerInfoClient(conn)
	case trait.Vending:
		return traits.NewVendingInfoClient(conn)

		// sc-bos private traits
	case button.TraitName:
		// return gen.NewButtonInfoClient(conn)
		return nil
	case color.TraitName:
		return gen.NewColorInfoClient(conn)
	case emergencylight.TraitName:
		// return gen.NewDaliInfoClient(conn)
		return nil
	case meter.TraitName:
		// return gen.NewMeterInfoClient(conn)
		return nil
	case mqttpb.TraitName:
		// return gen.NewMqttInfoClient(conn)
		return nil
	case statuspb.TraitName:
		// return gen.NewStatusInfoClient(conn)
		return nil
	case udmipb.TraitName:
		// return gen.NewUdmiInfoClient(conn)
		return nil
	}
	return nil
}
