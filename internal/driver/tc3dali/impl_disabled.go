//go:build notc3dali

package tc3dali

import (
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"go.uber.org/zap"
)

func newBusBuilder(_ ADSConfig) (busBuilder, error) {
	return &mockBusBuilder{}, nil
}

type mockBusBuilder struct{}

func (bb *mockBusBuilder) buildBus(config BusConfig, logger *zap.Logger) (dali.Dali, error) {
	return dali.NewMock(logger), nil
}
