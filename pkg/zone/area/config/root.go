package config

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

type Root struct {
	zone.RawConfig // RawConfig as we're going to pass this to all the features

	Metadata *traits.Metadata `json:"metadata,omitempty"`
}
