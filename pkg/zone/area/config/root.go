package config

import (
	"encoding/json"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"go.uber.org/multierr"
)

type Root struct {
	zone.Config
	MetadataCfg
	Raw json.RawMessage `json:"-"`
}

type MetadataCfg struct {
	Metadata *traits.Metadata `json:"metadata,omitempty"`
}

func (r *Root) UnmarshalJSON(buf []byte) error {
	r.Raw = buf
	return multierr.Combine(
		json.Unmarshal(buf, &r.Config),
		json.Unmarshal(buf, &r.MetadataCfg),
	)
}
