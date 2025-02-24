package node

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadatapb"
)

var MetadataTraitNotSupported = errors.New("metadata is not supported")

type MetadataChange = metadatapb.CollectionChange

// ListAllMetadata returns a slice containing all metadata set via Announce.
func (n *Node) ListAllMetadata(opts ...resource.ReadOption) []*traits.Metadata {
	return n.allMetadata.ListMetadata(opts...)
}

// PullAllMetadata returns a chan that emits MetadataChange whenever Announce or that announcement is undone.
func (n *Node) PullAllMetadata(ctx context.Context, opts ...resource.ReadOption) <-chan MetadataChange {
	return n.allMetadata.PullAllMetadata(ctx, opts...)
}

func (n *Node) mergeMetadata(name string, md *traits.Metadata) (Undo, error) {
	for i := 0; i < 5; i++ {
		_, err := n.allMetadata.MergeMetadata(name, md, resource.WithCreateIfAbsent())
		if isConcurrentUpdateDetectedError(err) && i < 4 {
			n.Logger.Debug("writing metadata, will try again", zap.Int("attempt", i), zap.String("name", name))
			continue
		}
		if err != nil {
			return NilUndo, err
		}
		break // no err
	}

	undo := Undo(func() {
		_, _ = n.allMetadata.DeleteMetadata(name, resource.WithAllowMissing(true))
	})

	return undo, nil
}

func (n *Node) isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
