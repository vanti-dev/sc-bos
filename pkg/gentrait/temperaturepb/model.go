package temperaturepb

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Model struct {
	temperature *resource.Value // of *gen.Temperature
}

func NewModel(opts ...resource.Option) *Model {
	args := calcModelArgs(opts...)
	return &Model{
		temperature: resource.NewValue(args.temperatureOpts...),
	}
}

func (m *Model) Temperature(opts ...resource.ReadOption) (*gen.Temperature, error) {
	got := m.temperature.Get(opts...)
	return got.(*gen.Temperature), nil
}

func (m *Model) SetTemperature(v *gen.Temperature, opts ...resource.WriteOption) (*gen.Temperature, error) {
	got, err := m.temperature.Set(v, opts...)
	if err != nil {
		return nil, err
	}
	return got.(*gen.Temperature), nil
}

type TemperatureChange struct {
	Value         *gen.Temperature
	ChangeTime    time.Time
	LastSeedValue bool
}

func (m *Model) PullTemperature(ctx context.Context, opts ...resource.ReadOption) <-chan TemperatureChange {
	send := make(chan TemperatureChange)
	go func() {
		defer close(send)
		for change := range m.temperature.Pull(ctx, opts...) {
			send <- TemperatureChange{
				Value:         change.Value.(*gen.Temperature),
				ChangeTime:    change.ChangeTime,
				LastSeedValue: change.LastSeedValue,
			}
		}
	}()
	return send
}
