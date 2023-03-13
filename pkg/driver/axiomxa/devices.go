package axiomxa

import (
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/mps"
)

// devices allows lookup and conversion between axiom controllers and devices and Smart Core device names.
type devices struct {
	byNetDevice map[netDevice]config.Device // values are smart core names
}

func devicesFromConfig(ds []config.Device) *devices {
	out := &devices{byNetDevice: make(map[netDevice]config.Device)}
	for _, d := range ds {
		out.byNetDevice[netDevice{network: d.NetworkDesc, device: d.DeviceDesc}] = d
	}
	return out
}

func (d *devices) Find(fields mps.Fields) (config.Device, bool) {
	nd := netDevice{network: fields.NetworkDesc, device: fields.DeviceDesc}
	dv, ok := d.byNetDevice[nd]
	return dv, ok
}

func (d *devices) SmartCoreName(fields mps.Fields) (string, bool) {
	dv, ok := d.Find(fields)
	return dv.Name, ok
}

func (d *devices) UDMITopicPrefix(fields mps.Fields) (string, bool) {
	dv, ok := d.Find(fields)
	if !ok {
		return "", false
	}
	if dv.UDMITopicPrefix == "" {
		return dv.Name, true
	}
	return dv.UDMITopicPrefix, true
}

type netDevice struct {
	network, device string
}
