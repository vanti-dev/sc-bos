package node

import (
	"errors"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var AutoTraitMetadata = map[string]string{}
var MetadataTraitNotSupported = errors.New("metadata is not supported")

func (n *Node) mergeMetadata(name string, md *traits.Metadata) (Undo, error) {
	metadataModel, err := n.metadataModel(name)
	if err != nil {
		return NilUndo, err
	}
	_, err = metadataModel.MergeMetadata(md)
	return func() {
		// todo: undo applying the metadata to the device
	}, err
}

func (n *Node) addTraitMetadata(name string, traitName trait.Name, md map[string]string) (Undo, error) {
	metadataModel, err := n.metadataModel(name)
	if err != nil {
		return NilUndo, err
	}
	_, err = metadataModel.UpdateTraitMetadata(&traits.TraitMetadata{
		Name: string(traitName),
		More: md,
	})
	return func() {
		// todo: remove any trait metadata from metadataModel
	}, err
}

func (n *Node) metadataApiRouter() router.Router {
	metadataApiClient := metadata.WrapApi(traits.UnimplementedMetadataApiServer{})
	for _, r := range n.routers {
		if r.HoldsType(metadataApiClient) {
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
