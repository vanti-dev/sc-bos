package transport

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
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
	send := make(chan PullTransportChange)

	go func() {
		defer close(send)
		for change := range m.transport.Pull(ctx, opts...) {
			val := change.Value.(*gen.Transport)
			select {
			case <-ctx.Done():
				return
			case send <- PullTransportChange{Value: val, ChangeTime: change.ChangeTime}:
			}
		}
	}()

	return send
}

type PullTransportChange struct {
	Value      *gen.Transport
	ChangeTime time.Time
}
