package prioritypb

import (
	"fmt"

	"github.com/smart-core-os/sc-api/go/traits"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type options struct {
	_suffix      string
	separator    string
	_defaultName string
	slotNames    []string
	logger       *zap.Logger
	_metadata    *traits.Metadata
}

func (o options) metadata(base string) *traits.Metadata {
	if o._metadata == nil {
		return nil
	}
	md := proto.Clone(o._metadata).(*traits.Metadata)
	md.Name = o.suffix(base)
	return md
}

func (o options) suffix(base string) string {
	return fmt.Sprintf("%s%s%s", base, o.separator, o._suffix)
}

func (o options) defaultName(base string) string {
	switch {
	case o._defaultName != "":
		return o.fqn(base, o._defaultName)
	case len(o.slotNames) == 0:
		return ""
	default:
		return o.fqn(base, o.slotNames[len(o.slotNames)/2])
	}
}

func (o options) fqn(base, leaf string) string {
	return fmt.Sprintf("%s%s%s%s%s", base, o.separator, o._suffix, o.separator, leaf)
}

func (o options) fqns(base string, leafs ...string) []string {
	out := make([]string, len(leafs))
	for i, leaf := range leafs {
		out[i] = o.fqn(base, leaf)
	}
	return out
}

type Option func(opts *options)

var defaultOptions = []Option{
	WithSuffix("priority"),
	WithSeparator("/"),
	WithSlots("1", "2", "3", "4", "5"),
	// default slot defaults to the middle slot
	WithLogger(zap.NewNop()),
}

func readOpts(opts ...Option) options {
	dst := &options{}
	for _, opt := range defaultOptions {
		opt(dst)
	}
	for _, opt := range opts {
		opt(dst)
	}
	return *dst
}

func WithSuffix(suffix string) Option {
	return func(opts *options) {
		opts._suffix = suffix
	}
}

func WithSeparator(separator string) Option {
	return func(opts *options) {
		opts.separator = separator
	}
}

func WithDefaultSlot(name string) Option {
	return func(opts *options) {
		opts._defaultName = name
	}
}

func WithSlots(names ...string) Option {
	return func(opts *options) {
		opts.slotNames = names
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithMetadata(md *traits.Metadata) Option {
	return func(opts *options) {
		opts._metadata = md
	}
}
