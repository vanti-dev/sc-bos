package tc3dali

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vanti-dev/bsp-ew/internal/driver"
)

const DriverName = "tc3dali"

func Factory(ctx context.Context, services driver.Services, rawConfig json.RawMessage) (driver.Driver, error) {
	var config Config
	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	err = applyConfig(ctx, services, config)
	if err != nil {
		return nil, err
	}

	return &driverImpl{
		config: config,
	}, nil
}

var _ driver.Factory = Factory

type driverImpl struct {
	config Config
}

func (d *driverImpl) Name() string {
	return d.config.Name
}
