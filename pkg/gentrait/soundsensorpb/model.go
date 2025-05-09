package soundsensorpb

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Model struct {
	soundLevel *resource.Value // of *gen.SoundLevel
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.SoundLevel{})}
	opts = append(defaultOpts, opts...)

	return &Model{
		soundLevel: resource.NewValue(opts...),
	}
}

func (m *Model) GetSoundLevel(opts ...resource.ReadOption) (*gen.SoundLevel, error) {
	return m.soundLevel.Get(opts...).(*gen.SoundLevel), nil
}

func (m *Model) PullSoundLevel(ctx context.Context, opts ...resource.ReadOption) <-chan PullSoundLevelChange {
	send := make(chan PullSoundLevelChange)

	go func() {
		defer close(send)
		for change := range m.soundLevel.Pull(ctx, opts...) {
			val := change.Value.(*gen.SoundLevel)
			select {
			case <-ctx.Done():
				return
			case send <- PullSoundLevelChange{Value: val, ChangeTime: change.ChangeTime}:
			}
		}
	}()

	return send
}

type PullSoundLevelChange struct {
	Value      *gen.SoundLevel
	ChangeTime time.Time
}
