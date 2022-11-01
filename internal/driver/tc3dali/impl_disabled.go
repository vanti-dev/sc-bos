//go:build notc3dali

package tc3dali

import (
	"context"
	"github.com/vanti-dev/bsp-ew/internal/driver"
)

func applyConfig(_ context.Context, services driver.Services, _ Config) error {
	services.Logger.Warn("tc3dali driver disabled")
	return nil
}
