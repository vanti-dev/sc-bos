package hd2

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"

	"github.com/vanti-dev/sc-bos-drivers/pkg/driver/steinel/hd2/config"
)

const DriverName = "steinel-hd2"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	d := &Driver{
		announcer: services.Node,
	}
	d.Service = service.New(service.MonoApply(d.applyConfig))
	d.logger = services.Logger.Named(DriverName)
	return d
}

type Driver struct {
	*service.Service[config.Root]
	logger    *zap.Logger
	announcer node.Announcer

	client *Client
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {
	d.client = NewInsecureClient(cfg.IpAddress)
	response := AccessResponse{}
	err := doGetRequest(d.client, &response, "access")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Response from hd2 driver access request: " + response.Rights)
	}

	return nil
}
