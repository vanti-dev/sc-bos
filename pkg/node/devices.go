package node

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/devicespb"
)

func (n *Node) GetDevice(name string, opts ...resource.ReadOption) (*gen.Device, error) {
	return n.devices.GetDevice(name, opts...)
}

func (n *Node) PullDevice(ctx context.Context, name string, opts ...resource.ReadOption) <-chan devicespb.DeviceChange {
	return n.devices.PullDevice(ctx, name, opts...)
}

func (n *Node) ListDevices(opts ...resource.ReadOption) []*gen.Device {
	return n.devices.ListDevices(opts...)
}

func (n *Node) PullDevices(ctx context.Context, opts ...resource.ReadOption) <-chan devicespb.DevicesChange {
	return n.devices.PullDevices(ctx, opts...)
}
