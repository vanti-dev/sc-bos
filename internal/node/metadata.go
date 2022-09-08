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

const (
	MetadataRealism        = "scos.playground.realism"
	MetadataRealismVirtual = "virtual"
	MetadataRealismModel   = "model"

	MetadataDeviceType = "scos.playground.device-type"
)

var AutoTraitMetadata = map[string]string{
	MetadataRealism: MetadataRealismModel,
}
var MetadataTraitNotSupported = errors.New("metadata is not supported")

func (n *Local) addTraitMetadata(name string, traitName trait.Name, md map[string]string) error {
	metadataApiRouter := n.metadataApiRouter()
	if metadataApiRouter == nil {
		return MetadataTraitNotSupported
	}
	client, err := metadataApiRouter.Get(name)
	if err != nil {
		return err
	}
	metadataModel, ok := wrap.UnwrapFully(client).(*metadata.Model)
	if !ok {
		return status.Errorf(codes.FailedPrecondition, "%v cannot auto-create trait %v", name, traitName)
	}
	_, err = metadataModel.UpdateTraitMetadata(&traits.TraitMetadata{
		Name: string(traitName),
		More: md,
	})
	return err
}

func (n *Local) metadataApiRouter() router.Router {
	metadataApiClient := metadata.WrapApi(traits.UnimplementedMetadataApiServer{})
	for _, r := range n.routers {
		if r.HoldsType(metadataApiClient) {
			return r
		}
	}
	return nil
}
