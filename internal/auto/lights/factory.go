package lights

import "github.com/vanti-dev/bsp-ew/internal/auto"

const AutoType = "lights"

var Factory = auto.FactoryFunc(func(services *auto.Services) auto.Starter {
	return PirsTurnLightsOn(services.Node, services.Logger.Named("lights"))
})
