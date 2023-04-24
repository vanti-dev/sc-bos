package button

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/resource"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Model struct {
	buttonState *resource.Value // of *gen.ButtonState
}

func NewModel(initialPressState gen.ButtonState_Press) *Model {
	return &Model{
		buttonState: resource.NewValue(resource.WithInitialValue(&gen.ButtonState{
			State: initialPressState,
		})),
	}
}

func (m *Model) GetButtonState(options ...resource.ReadOption) *gen.ButtonState {
	return m.buttonState.Get(options...).(*gen.ButtonState)
}

func (m *Model) UpdateButtonState(value *gen.ButtonState, options ...resource.WriteOption) (*gen.ButtonState, error) {
	updated, err := m.buttonState.Set(value, options...)
	if err != nil {
		return nil, err
	}
	return updated.(*gen.ButtonState), nil
}

func (m *Model) PullButtonState(ctx context.Context, options ...resource.ReadOption) <-chan PullButtonStateChange {
	tx := make(chan PullButtonStateChange)

	rx := m.buttonState.Pull(ctx, options...)
	go func() {
		defer close(tx)
		for change := range rx {
			value := change.Value.(*gen.ButtonState)
			tx <- PullButtonStateChange{
				Value:         value,
				ChangeTime:    change.ChangeTime,
				LastSeedValue: change.LastSeedValue,
			}
		}
	}()
	return tx
}

type PullButtonStateChange struct {
	Value         *gen.ButtonState
	ChangeTime    time.Time
	LastSeedValue bool
}
