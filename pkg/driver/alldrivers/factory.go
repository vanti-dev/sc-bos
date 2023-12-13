package alldrivers

import (
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock"
	"github.com/vanti-dev/sc-bos/pkg/driver/proxy"
	"github.com/vanti-dev/sc-bos/pkg/driver/steinel/hpd3"
	"github.com/vanti-dev/sc-bos/pkg/driver/xovis"
)

// Factories returns a new map containing all known driver factories.
func Factories() map[string]driver.Factory {
	return map[string]driver.Factory{
		bacnet.DriverName: bacnet.Factory,
		hpd3.DriverName:   hpd3.Factory,
		mock.DriverName:   mock.Factory,
		proxy.DriverName:  proxy.Factory,
		xovis.DriverName:  xovis.Factory,
	}
}
