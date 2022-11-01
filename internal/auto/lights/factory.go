package lights

import (
	"github.com/vanti-dev/bsp-ew/internal/auto"
	"github.com/vanti-dev/bsp-ew/internal/task"
)

const AutoType = "lights"

var Factory = auto.FactoryFunc(func(services *auto.Services) task.Starter {
	return PirsTurnLightsOn(services.Node, services.Logger.Named("lights"))
})
