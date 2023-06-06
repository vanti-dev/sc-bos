package node

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

var MetadataTraitNotSupported = errors.New("metadata is not supported")

type MetadataChange struct {
	Name          string
	ChangeTime    time.Time
	Type          types.ChangeType
	OldValue      *traits.Metadata
	NewValue      *traits.Metadata
	LastSeedValue bool
}

// ListAllMetadata returns a slice containing all metadata set via Announce.
func (n *Node) ListAllMetadata(opts ...resource.ReadOption) []*traits.Metadata {
	msgs := n.allMetadata.List(opts...)
	mds := make([]*traits.Metadata, len(msgs))
	for i, msg := range msgs {
		mds[i] = msg.(*traits.Metadata)
	}
	return mds
}

// PullAllMetadata returns a chan that emits MetadataChange whenever Announce or that announcement is undone.
func (n *Node) PullAllMetadata(ctx context.Context, opts ...resource.ReadOption) <-chan MetadataChange {
	mdC := make(chan MetadataChange)
	go func() {
		defer close(mdC)
		for change := range n.allMetadata.Pull(ctx, opts...) {
			emit := MetadataChange{
				Name:          change.Id,
				Type:          change.ChangeType,
				ChangeTime:    change.ChangeTime,
				LastSeedValue: change.LastSeedValue,
			}
			if change.OldValue != nil {
				emit.OldValue = change.OldValue.(*traits.Metadata)
			}
			if change.NewValue != nil {
				emit.NewValue = change.NewValue.(*traits.Metadata)
			}
			select {
			case <-ctx.Done():
				return
			case mdC <- emit:
			}
		}
	}()
	return mdC
}

func (n *Node) mergeMetadata(name string, md *traits.Metadata) (Undo, error) {
	metadataModel, undo, err := n.initMetadataModel(name)
	if err != nil {
		return undo, err
	}

	var newMd *traits.Metadata
	for i := 0; i < 5; i++ {
		var err error
		newMd, err = metadataModel.MergeMetadata(md)
		if isConcurrentUpdateDetectedError(err) && i < 4 {
			n.Logger.Debug("writing metadata, will try again", zap.Int("attempt", i), zap.String("name", name))
			continue
		}
		if err != nil {
			undo()
			return NilUndo, err
		}
		break // no err
	}

	for i := 0; i < 5; i++ {
		_, err := n.allMetadata.Update(name, newMd, resource.WithCreateIfAbsent())
		if isConcurrentUpdateDetectedError(err) && i < 4 {
			n.Logger.Debug("updating all metadata, will try again", zap.Int("attempt", i), zap.String("name", name))
			continue
		}
		if err != nil {
			undo()
			return NilUndo, err
		}
		break // no err
	}

	// todo: undo applying the metadata to the device
	return undo, nil
}

// initMetadataModel gets or creates a named metadata model, storing it in the nodes routers as needed.
func (n *Node) initMetadataModel(name string) (*metadata.Model, Undo, error) {
	undo := NilUndo
	metadataModel, err := n.metadataModel(name)
	if err != nil {
		if !n.isNotFound(err) {
			return nil, nil, err
		}
		metadataModel, undo = n.announceMetadata(name)

		// send that the metadata was removed if the merge is undone.
		undo = UndoAll(undo, func() {
			_, _ = n.allMetadata.Delete(name, resource.WithAllowMissing(true))
		})
	}
	return metadataModel, undo, nil
}

func (n *Node) metadataApiRouter() router.Router {
	for _, r := range n.routers {
		if _, ok := r.(*metadata.ApiRouter); ok {
			return r
		}
	}
	return nil
}

func (n *Node) metadataModel(name string) (*metadata.Model, error) {
	metadataApiRouter := n.metadataApiRouter()
	if metadataApiRouter == nil {
		return nil, MetadataTraitNotSupported
	}
	client, err := metadataApiRouter.Get(name)
	if err != nil {
		return nil, err
	}
	metadataModel, ok := wrap.UnwrapFully(client).(*metadata.Model)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "%v cannot store node metadata", name)
	}
	return metadataModel, nil
}

func (n *Node) isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func (n *Node) announceMetadata(name string) (*metadata.Model, Undo) {
	// auto add metadata support for devices that are asking to add metadata to that device
	md := &traits.Metadata{Name: name, Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}}}
	model := metadata.NewModel(resource.WithInitialValue(md))
	undo := n.announceLocked(name, HasTrait(trait.Metadata, WithClients(metadata.WrapApi(metadata.NewModelServer(model)))),
		// avoid cycles when announcing metadata for the first time
		HasNoAutoMetadata())
	return model, undo
}
