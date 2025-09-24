// Package nodeopts provides common and private options for the node package.
package nodeopts

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/devicespb"
)

type Option interface {
	apply(*Struct)
}

type optionFunc func(*Struct)

func (f optionFunc) apply(o *Struct) {
	f(o)
}

// WithStore sets the Store used by the Node to Store its announced devices.
func WithStore(store Store) Option {
	return optionFunc(func(o *Struct) {
		o.Store = store
	})
}

// Join combines multiple options into a single struct.
func Join(opts ...Option) Struct {
	var o Struct
	for _, opt := range opts {
		opt.apply(&o)
	}
	return o
}

// Struct contains all options for a Node as a struct for easy access.
type Struct struct {
	Store Store
}

func (s Struct) apply(o *Struct) {
	if s.Store != nil {
		o.Store = s.Store
	}
}

// Store describes how a node stores its announced devices.
type Store interface {
	GetDevice(name string, opts ...resource.ReadOption) (*gen.Device, error)
	PullDevice(ctx context.Context, name string, opts ...resource.ReadOption) <-chan devicespb.DeviceChange
	ListDevices(opts ...resource.ReadOption) []*gen.Device
	PullDevices(ctx context.Context, opts ...resource.ReadOption) <-chan devicespb.DevicesChange
	Update(d *gen.Device, opts ...resource.WriteOption) (*gen.Device, error)
	Delete(name string, opts ...resource.WriteOption) (*gen.Device, error)
}
