package merge

import (
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func IntoTrait(client *gobacnet.Client, ctx known.Context, traitConfig config.RawTrait, logger *zap.Logger) (node.SelfAnnouncer, error) {
	// todo: implement some traits that pull data from different bacnet devices.
	switch traitConfig.Kind {
	case trait.FanSpeed:
		return newFanSpeed(client, ctx, traitConfig, logger)
	case trait.AirTemperature:
		return newAirTemperature(client, ctx, traitConfig, logger)
	case UdmiMergeName:
		return newUdmiMerge(client, ctx, traitConfig, logger)
	}
	return nil, ErrTraitNotSupported
}
