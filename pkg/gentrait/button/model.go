package button

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/util/resources"
	"github.com/smart-core-os/sc-golang/pkg/resource"

	"github.com/smart-core-os/sc-bos/pkg/gen"
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
	return resources.PullValue[*gen.ButtonState](ctx, m.buttonState.Pull(ctx, options...))
}

type PullButtonStateChange = resources.ValueChange[*gen.ButtonState]
