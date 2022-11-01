package tc3dali

import (
	"context"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/task"
)

const DriverName = "tc3dali"

var Factory driver.Factory = factory{}

type factory struct{}

func (_ factory) New(services driver.Services) task.Starter {
	return NewDriver(services)
}

func NewDriver(services driver.Services) *driver.Lifecycle[Config] {
	d := driver.NewLifecycle(func(ctx context.Context, cfg Config) error {
		return applyConfig(ctx, services, cfg)
	})
	d.Logger = services.Logger.Named("tc3dali")
	return d
}
