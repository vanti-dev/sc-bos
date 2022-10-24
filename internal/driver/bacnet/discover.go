package bacnet

import (
	"fmt"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/config"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"net"
)

func (d *Driver) findDevice(device config.Device) (bactypes.Device, error) {
	fail := func(err error) (bactypes.Device, error) {
		return bactypes.Device{}, err
	}

	if device.Comm == nil {
		id := device.ID
		is, err := d.client.WhoIs(int(id), int(id))
		if err != nil {
			return fail(err)
		}
		if len(is) == 0 {
			return fail(fmt.Errorf("no devices found (via WhoIs) with id %d", id))
		}
		return is[0], nil
	}

	udpAddr := net.UDPAddrFromAddrPort(*device.Comm.IP)
	addr := bactypes.UDPToAddress(udpAddr)
	bacDevices, err := d.client.RemoteDevices(addr, device.ID)
	if err != nil {
		return fail(err)
	}
	return bacDevices[0], nil
}
