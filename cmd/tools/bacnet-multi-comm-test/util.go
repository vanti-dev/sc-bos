package main

import (
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
)

func shouldDiscoverObjects(cfg config.Root, device config.Device) bool {
	if device.DiscoverObjects != nil {
		return *device.DiscoverObjects
	}
	return cfg.DiscoverObjects
}
