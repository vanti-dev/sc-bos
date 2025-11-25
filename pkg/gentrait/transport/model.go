package transport

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/util/resources"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type Model struct {
	transport *resource.Value // of *gen.Transport
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.Transport{})}
	opts = append(defaultOpts, opts...)

	return &Model{
		transport: resource.NewValue(opts...),
	}
}

func (m *Model) GetTransport(opts ...resource.ReadOption) (*gen.Transport, error) {
	return m.transport.Get(opts...).(*gen.Transport), nil
}

func (m *Model) UpdateTransport(transport *gen.Transport, opts ...resource.WriteOption) (*gen.Transport, error) {
	res, err := m.transport.Set(transport, opts...)
	if err != nil {
		return nil, err
	}
	return res.(*gen.Transport), nil
}

func (m *Model) PullTransport(ctx context.Context, opts ...resource.ReadOption) <-chan PullTransportChange {
	return resources.PullValue[*gen.Transport](ctx, m.transport.Pull(ctx, opts...))
}

type PullTransportChange = resources.ValueChange[*gen.Transport]
