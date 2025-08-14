package devicespb

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadatapb"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// Collection is a list of unique devices.
type Collection struct {
	names *resource.Collection // of *gen.Device, keyed by device name
}

func NewCollection(opts ...resource.Option) *Collection {
	return &Collection{
		names: resource.NewCollection(opts...),
	}
}

// ListDevices returns a slice containing all devices in the collection.
// The devices will be ordered by their name, ascending.
// Use opts for filter and project the returned devices.
func (c *Collection) ListDevices(opts ...resource.ReadOption) []*gen.Device {
	res := c.names.List(opts...)
	devices := make([]*gen.Device, 0, len(res))
	for _, r := range res {
		devices = append(devices, r.(*gen.Device))
	}
	return devices
}

// PullDevices returns a channel that will receive changes to the devices in the collection.
// Unless [resource.WithUpdatesOnly] is true, the channel will receive all current devices as ADDs.
func (c *Collection) PullDevices(ctx context.Context, opts ...resource.ReadOption) <-chan DevicesChange {
	send := make(chan DevicesChange)
	recv := c.names.Pull(ctx, opts...)
	go func() {
		defer close(send)
		for change := range recv {
			select {
			case <-ctx.Done():
				return
			case send <- devicesChangeFromResource(change):
			}
		}
	}()
	return send
}

// GetDevice returns the device with the given name from the collection.
// If the device does not exist, an error with codes.NotFound will be returned.
func (c *Collection) GetDevice(name string, opts ...resource.ReadOption) (*gen.Device, error) {
	res, ok := c.names.Get(name, opts...)
	if !ok {
		return nil, status.Error(codes.NotFound, name)
	}
	return res.(*gen.Device), nil
}

// PullDevice returns a channel that will receive changes to the device with the given name.
// If the device is deleted, the channel will close.
func (c *Collection) PullDevice(ctx context.Context, name string, opts ...resource.ReadOption) <-chan DeviceChange {
	send := make(chan DeviceChange)
	recv := c.names.PullID(ctx, name, opts...)
	go func() {
		defer close(send)
		for d := range recv {
			select {
			case <-ctx.Done():
				return
			case send <- deviceChangeFromResource(name, d):
			}
		}
	}()
	return send
}

// Merge merges d with the existing device in the collection.
// An error will be returned if the device does not exist, and resource.WithCreateIfAbsent is not in opts.
// nil fields in d, like Metadata, will be interpreted as absent during the merge.
// This is different to how Update and normal proto merging works,
// where nil fields are interpreted as a delete instruction.
func (c *Collection) Merge(d *gen.Device, opts ...resource.WriteOption) (*gen.Device, error) {
	opts = append([]resource.WriteOption{
		resource.InterceptBefore(mergeInterceptor),
	}, opts...)
	update, err := c.names.Update(d.Name, d, opts...)
	if err != nil {
		return nil, err
	}
	return update.(*gen.Device), nil
}

// Delete removes the device with the given name from the collection.
func (c *Collection) Delete(name string, opts ...resource.WriteOption) (*gen.Device, error) {
	old, err := c.names.Delete(name, opts...)
	var oldDevice *gen.Device
	if old != nil {
		oldDevice = old.(*gen.Device)
	}
	return oldDevice, err
}

type DevicesChange struct {
	Name          string
	ChangeTime    time.Time
	ChangeType    types.ChangeType
	OldValue      *gen.Device
	NewValue      *gen.Device
	LastSeedValue bool
}

func devicesChangeFromResource(change *resource.CollectionChange) DevicesChange {
	dc := DevicesChange{
		Name:          change.Id,
		ChangeTime:    change.ChangeTime,
		ChangeType:    change.ChangeType,
		LastSeedValue: change.LastSeedValue,
	}
	if change.OldValue != nil {
		dc.OldValue = change.OldValue.(*gen.Device)
	}
	if change.NewValue != nil {
		dc.NewValue = change.NewValue.(*gen.Device)
	}
	return dc
}

type DeviceChange struct {
	Name       string
	ChangeTime time.Time
	Device     *gen.Device
}

func deviceChangeFromResource(name string, change *resource.ValueChange) DeviceChange {
	dc := DeviceChange{
		Name:       name,
		ChangeTime: change.ChangeTime,
	}
	if change.Value != nil {
		dc.Device = change.Value.(*gen.Device)
	}
	return dc
}

func mergeInterceptor(old, new proto.Message) {
	oldVal := old.(*gen.Device)
	newVal := new.(*gen.Device)
	switch {
	case newVal.Metadata == nil:
		newVal.Metadata = oldVal.Metadata // no metadata in new, use old
	case oldVal.Metadata == nil:
		break // old has no metadata, keep new as is
	default:
		// both have metadata, merge them
		metadatapb.Merge(oldVal.Metadata, newVal.Metadata)
	}
}
