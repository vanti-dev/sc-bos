package merge

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/gobacnet"

	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func IntoTrait(client *gobacnet.Client, ctx known.Context, traitConfig config.RawTrait) (node.SelfAnnouncer, error) {
	// todo: implement some traits that pull data from different bacnet devices.
	switch traitConfig.Kind {
	case trait.FanSpeed:
		return newFanSpeed(client, ctx, traitConfig)
	case trait.AirTemperature:
		return newAirTemperature(client, ctx, traitConfig)
	}
	return nil, ErrTraitNotSupported
}
