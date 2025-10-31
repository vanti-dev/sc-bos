package temperaturepb

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/util/resources"
)

type Model struct {
	temperature *resource.Value // of *gen.Temperature
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.Temperature{})}
	opts = append(defaultOpts, opts...)

	return &Model{
		temperature: resource.NewValue(opts...),
	}
}

func (m *Model) GetTemperature(opts ...resource.ReadOption) (*gen.Temperature, error) {
	return m.temperature.Get(opts...).(*gen.Temperature), nil
}

func (m *Model) UpdateTemperature(temperature *gen.Temperature, opts ...resource.WriteOption) (*gen.Temperature, error) {
	res, err := m.temperature.Set(temperature, opts...)
	if err != nil {
		return nil, err
	}
	return res.(*gen.Temperature), nil
}

func (m *Model) PullTemperature(ctx context.Context, opts ...resource.ReadOption) <-chan PullTemperatureChange {
	return resources.PullValue[*gen.Temperature](ctx, m.temperature.Pull(ctx, opts...))
}

type PullTemperatureChange = resources.ValueChange[*gen.Temperature]
