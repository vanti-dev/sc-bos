package alldrivers

import (
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/airthings"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/gallagher"
	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/hikcentral"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock"
	"github.com/vanti-dev/sc-bos/pkg/driver/opcua"
	"github.com/vanti-dev/sc-bos/pkg/driver/pestsense"
	"github.com/vanti-dev/sc-bos/pkg/driver/proxy"
	seWiserKnx "github.com/vanti-dev/sc-bos/pkg/driver/se/wiser-knx"
	shellyTrv "github.com/vanti-dev/sc-bos/pkg/driver/shelly/trv"
	steinelHpd "github.com/vanti-dev/sc-bos/pkg/driver/steinel/hpd"
	"github.com/vanti-dev/sc-bos/pkg/driver/xovis"
)

// Factories returns a new map containing all known driver factories.
func Factories() map[string]driver.Factory {
	return map[string]driver.Factory{
		airthings.DriverName:  airthings.Factory,
		bacnet.DriverName:     bacnet.Factory,
		gallagher.DriverName:  gallagher.Factory,
		helvarnet.DriverName:  helvarnet.Factory,
		hikcentral.DriverName: hikcentral.Factory,
		mock.DriverName:       mock.Factory,
		opcua.DriverName:      opcua.Factory,
		pestsense.DriverName:  pestsense.Factory,
		proxy.DriverName:      proxy.Factory,
		seWiserKnx.DriverName: seWiserKnx.Factory,
		shellyTrv.DriverName:  shellyTrv.Factory,
		steinelHpd.DriverName: steinelHpd.Factory,
		xovis.DriverName:      xovis.Factory,
	}
}
