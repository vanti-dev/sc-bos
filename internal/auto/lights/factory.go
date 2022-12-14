package lights

import (
	"github.com/vanti-dev/sc-bos/internal/auto"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

const AutoType = "lights"

var Factory = auto.FactoryFunc(func(services auto.Services) task.Starter {
	return PirsTurnLightsOn(services.Node, services.Logger.Named("lights"))
})
