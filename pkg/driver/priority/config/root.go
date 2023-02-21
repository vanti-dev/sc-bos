package config

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/prioritypb"
)

type Root struct {
	driver.BaseConfig

	Names
	Devices []Device `json:"devices,omitempty"`
}

type Names struct {
	Suffix      *string  `json:"suffix,omitempty"`
	Separator   *string  `json:"separator,omitempty"`
	Slots       []string `json:"slots,omitempty"`
	DefaultSlot *string  `json:"defaultSlot,omitempty"`
}

func (n Names) Options(opts ...prioritypb.Option) []prioritypb.Option {
	if n.Suffix != nil {
		opts = append(opts, prioritypb.WithSuffix(*n.Suffix))
	}
	if n.Separator != nil {
		opts = append(opts, prioritypb.WithSeparator(*n.Separator))
	}
	if len(n.Slots) > 0 {
		opts = append(opts, prioritypb.WithSlots(n.Slots...))
	}
	if n.DefaultSlot != nil {
		opts = append(opts, prioritypb.WithDefaultSlot(*n.DefaultSlot))
	}
	return opts
}

type Device struct {
	*traits.Metadata
}
