package temperaturepb

import (
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// DefaultModelOptions holds the default options for the model.
var DefaultModelOptions = []resource.Option{
	WithInitialTemperature(&gen.Temperature{
		SetPoint: &types.Temperature{ValueCelsius: 21.5},
		Measured: &types.Temperature{ValueCelsius: 20.12},
	}),
}

// ModelOption defined the base type for all options that apply to this traits model.
type ModelOption interface {
	resource.Option
	applyModel(args *modelArgs)
}

// WithTemperatureOption configures the temperature resource of the model.
func WithTemperatureOption(opts ...resource.Option) resource.Option {
	return modelOptionFunc(func(args *modelArgs) {
		args.temperatureOpts = append(args.temperatureOpts, opts...)
	})
}

// WithInitialTemperature returns an option that configures the model to initialise with the given temperature.
func WithInitialTemperature(temperature *gen.Temperature) resource.Option {
	return WithTemperatureOption(resource.WithInitialValue(temperature))
}

func calcModelArgs(opts ...resource.Option) modelArgs {
	args := new(modelArgs)
	args.apply(DefaultModelOptions...)
	args.apply(opts...)
	return *args
}

type modelArgs struct {
	temperatureOpts []resource.Option
}

func (a *modelArgs) apply(opts ...resource.Option) {
	for _, opt := range opts {
		if v, ok := opt.(ModelOption); ok {
			v.applyModel(a)
			continue
		}
		a.temperatureOpts = append(a.temperatureOpts, opt)
	}
}

func modelOptionFunc(fn func(args *modelArgs)) ModelOption {
	return modelOption{resource.EmptyOption{}, fn}
}

type modelOption struct {
	resource.Option
	fn func(args *modelArgs)
}

func (m modelOption) applyModel(args *modelArgs) {
	m.fn(args)
}
