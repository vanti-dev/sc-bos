package accesspb

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/util/resources"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type Model struct {
	accessAttempt *resource.Value // of *gen.AccessAttempt
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.AccessAttempt{})}
	opts = append(defaultOpts, opts...)
	return &Model{
		accessAttempt: resource.NewValue(opts...),
	}
}

func (m *Model) GetLastAccessAttempt(opts ...resource.ReadOption) (*gen.AccessAttempt, error) {
	v := m.accessAttempt.Get(opts...)
	return v.(*gen.AccessAttempt), nil
}

func (m *Model) UpdateLastAccessAttempt(accessAttempt *gen.AccessAttempt, opts ...resource.WriteOption) (*gen.AccessAttempt, error) {
	v, err := m.accessAttempt.Set(accessAttempt, opts...)
	if err != nil {
		return nil, err
	}
	return v.(*gen.AccessAttempt), nil
}

func (m *Model) PullAccessAttempts(ctx context.Context, opts ...resource.ReadOption) <-chan PullAccessAttemptsChange {
	return resources.PullValue[*gen.AccessAttempt](ctx, m.accessAttempt.Pull(ctx, opts...))
}

type PullAccessAttemptsChange = resources.ValueChange[*gen.AccessAttempt]
