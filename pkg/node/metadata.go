package node

import (
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadatapb"
)

// mergeAllMetadata uses the metadatapb.Merge algorithm to merge multiple metadata objects into one.
// The input metadata objects will not be modified.
// Metadata at higher indexes in the mds slice will override metadata at lower indexes.
func (n *Node) mergeAllMetadata(name string, mds ...*traits.Metadata) *traits.Metadata {
	if len(mds) == 0 {
		return nil
	}
	if len(mds) == 1 {
		md := mds[0]
		md.Name = name
		// for consistency, sort the traits by name
		slices.SortFunc(md.Traits, func(a, b *traits.TraitMetadata) int {
			return strings.Compare(a.Name, b.Name)
		})
		return md
	}

	md := &traits.Metadata{
		Name: name,
	}
	for _, m := range mds {
		m = proto.Clone(m).(*traits.Metadata)
		metadatapb.Merge(md, m)
		md = m
	}
	return md
}
