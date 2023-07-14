package config

import (
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	auto.Config
	// Name of the device that stores the alerts.
	// Must implement AlertAdminApi.
	Destination string `json:"destination,omitempty"`
	// If true, all devices on the current node that implement Status will be monitored.
	// Additional sources may be defined via Sources.
	DiscoverSources bool `json:"discoverSources,omitempty"`
	// Name of the devices that implement Status trait and that are monitored.
	Sources []Source `json:"sources,omitempty"`
	// Delay querying the status of devices by this much, to allow them to boot up.
	DelayStart *jsontypes.Duration `json:"delayStart,omitempty"`
	// Device name prefixes to ignore.
	// Only used if DiscoverSources is true.
	IgnorePrefixes []string `json:"ignorePrefixes,omitempty"`
}

type Source struct {
	Name  string `json:"name,omitempty"`
	Floor string `json:"floor,omitempty"`
	Zone  string `json:"zone,omitempty"`
}
