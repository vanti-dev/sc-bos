package fluidflowpb

import (
	"context"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/util/resources"
)

type Model struct {
	fluidFlow *resource.Value // of *gen.FluidFlow
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.FluidFlow{})}
	opts = append(defaultOpts, opts...)
	return &Model{
		fluidFlow: resource.NewValue(opts...),
	}
}

func (m *Model) GetFluidFlow() (*gen.FluidFlow, error) {
	return m.fluidFlow.Get().(*gen.FluidFlow), nil
}

func (m *Model) UpdateFluidFlow(flow *gen.FluidFlow, opts ...resource.WriteOption) (*gen.FluidFlow, error) {
	res, err := m.fluidFlow.Set(flow, opts...)
	if err != nil {
		return nil, err
	}
	return res.(*gen.FluidFlow), nil
}

func (m *Model) PullFluidFlow(ctx context.Context, opts ...resource.ReadOption) <-chan FlowChange {
	return resources.PullValue[*gen.FluidFlow](ctx, m.fluidFlow.Pull(ctx, opts...))
}

type FlowChange = resources.ValueChange[*gen.FluidFlow]
