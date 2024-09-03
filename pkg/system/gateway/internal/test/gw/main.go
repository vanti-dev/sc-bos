package gw

import (
	"context"
	"log"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/app"
	"github.com/vanti-dev/sc-bos/pkg/app/sysconf"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/lighttest"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/gateway"
	"github.com/vanti-dev/sc-bos/pkg/system/hub"
	"github.com/vanti-dev/sc-bos/pkg/system/tenants"
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

	alltraits.AddSupportFor(controller.Node, trait.Parent, trait.Metadata, trait.OnOff)

	err = controller.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func loadSystemConfig() (sysconf.Config, error) {
	systemConfig := sysconf.Default()
	gwFactory := gateway.Factory(&lighttest.Holder{})
	systemConfig.SystemFactories = map[string]system.Factory{
		gateway.Name:       gwFactory,
		gateway.LegacyName: gwFactory,
		// todo: remove these services
		"hub":     hub.Factory(),
		"tenants": tenants.Factory,
	}

	err := sysconf.Load(&systemConfig)
	return systemConfig, err
}
