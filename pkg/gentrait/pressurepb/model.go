package pressurepb

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/util/resources"
)

type Model struct {
	pressure *resource.Value // of *gen.Pressure
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.Pressure{})}
	opts = append(defaultOpts, opts...)
	return &Model{
		pressure: resource.NewValue(opts...),
	}
}

func (m *Model) GetPressure() (*gen.Pressure, error) {
	return m.pressure.Get().(*gen.Pressure), nil
}

func (m *Model) UpdatePressure(pressure *gen.Pressure, opts ...resource.WriteOption) (*gen.Pressure, error) {
	res, err := m.pressure.Set(pressure, opts...)
	if err != nil {
		return nil, err
	}
	return res.(*gen.Pressure), nil
}

func (m *Model) PullPressure(ctx context.Context, opts ...resource.ReadOption) <-chan PullPressureChange {
	return resources.PullValue[*gen.Pressure](ctx, m.pressure.Pull(ctx, opts...))
}

type PullPressureChange = resources.ValueChange[*gen.Pressure]
