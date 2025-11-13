package merge

import (
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/securityevent"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/temperaturepb"
	transportpb "github.com/vanti-dev/sc-bos/pkg/gentrait/transport"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func IntoTrait(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, traitConfig config.RawTrait, logger *zap.Logger) (node.SelfAnnouncer, error) {
	// todo: implement some traits that pull data from different bacnet devices.
	switch traitConfig.Kind {
	case trait.AirQualitySensor:
		return newAirQualitySensor(client, devices, statuses, traitConfig, logger)
	case trait.AirTemperature:
		return newAirTemperature(client, devices, statuses, traitConfig, logger)
	case trait.Electric:
		return newElectric(client, devices, statuses, traitConfig, logger)
	case trait.Emergency:
		return newEmergency(client, devices, statuses, traitConfig, logger)
	case trait.EnergyStorage:
		return newEnergyStorage(client, devices, statuses, traitConfig, logger)
	case trait.FanSpeed:
		return newFanSpeed(client, devices, statuses, traitConfig, logger)
	case trait.Light:
		return newLight(client, devices, statuses, traitConfig, logger)
	case meter.TraitName:
		return newMeter(client, devices, statuses, traitConfig, logger)
	case trait.Mode:
		return newMode(client, devices, statuses, traitConfig, logger)
	case trait.OccupancySensor:
		return newOccupancy(client, devices, statuses, traitConfig, logger)
	case trait.OnOff:
		return newOnOff(client, devices, statuses, traitConfig, logger)
	case accesspb.TraitName:
		return newAccess(client, devices, statuses, traitConfig, logger)
	case securityevent.TraitName:
		return newSecurityEvent(client, devices, statuses, traitConfig, logger)
	case statuspb.TraitName:
		return newStatus(client, devices, statuses, traitConfig, logger)
	case temperaturepb.TraitName:
		return newTemperature(client, devices, statuses, traitConfig, logger)
	case transportpb.TraitName:
		return newTransport(client, devices, statuses, traitConfig, logger)
	case UdmiMergeName, udmipb.TraitName:
		return newUdmiMerge(client, devices, statuses, traitConfig, logger)
	}
	return nil, ErrTraitNotSupported
}
