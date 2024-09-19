package ac

import (
	"context"
	"log"

	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/mock"
)

func Main() {
	ctx := context.Background()

	systemConfig, err := loadSystemConfig()
	if err != nil {
		log.Fatal(err)
	}

	controller, err := app.Bootstrap(ctx, systemConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = controller.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func loadSystemConfig() (sysconf.Config, error) {
	systemConfig := sysconf.Default()
	systemConfig.DriverFactories = map[string]driver.Factory{
		"mock": mock.Factory,
	}

	err := sysconf.Load(&systemConfig)
	return systemConfig, err
}
