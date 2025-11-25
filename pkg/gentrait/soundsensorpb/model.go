package soundsensorpb

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/util/resources"
	"github.com/smart-core-os/sc-golang/pkg/resource"
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
	return resources.PullValue[*gen.SoundLevel](ctx, m.soundLevel.Pull(ctx, opts...))
}

func (m *Model) UpdateSoundLevel(soundLevel *gen.SoundLevel, opts ...resource.WriteOption) (*gen.SoundLevel, error) {
	res, err := m.soundLevel.Set(soundLevel, opts...)
	if err != nil {
		return nil, err
	}
	return res.(*gen.SoundLevel), nil
}

type PullSoundLevelChange = resources.ValueChange[*gen.SoundLevel]
