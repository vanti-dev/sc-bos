package accesspb

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
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
	send := make(chan PullAccessAttemptsChange)

	recv := m.accessAttempt.Pull(ctx, opts...)
	go func() {
		defer close(send)
		for change := range recv {
			value := change.Value.(*gen.AccessAttempt)
			send <- PullAccessAttemptsChange{
				Value:      value,
				ChangeTime: change.ChangeTime,
			}
		}
	}()

	return send
}

type PullAccessAttemptsChange struct {
	Value      *gen.AccessAttempt
	ChangeTime time.Time
}
