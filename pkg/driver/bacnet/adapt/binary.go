package adapt

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/property"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

// BinaryObject adapts a binary bacnet object as smart core traits.
func BinaryObject(client *gobacnet.Client, device bactypes.Device, object config.Object, logger *zap.Logger) (node.SelfAnnouncer, error) {
	switch object.Trait {
	case "":
		return nil, ErrNoDefault
	case trait.OnOff:
		model := onoff.NewModel(traits.OnOff_STATE_UNSPECIFIED)
		return &binaryOnOff{
			client: client,
			device: device,
			object: object,
			logger: logger,

			model:       model,
			ModelServer: onoff.NewModelServer(model),
		}, nil
	}

	return nil, ErrNoAdaptation
}

type binaryOnOff struct {
	client *gobacnet.Client
	device bactypes.Device
	object config.Object
	logger *zap.Logger

	model *onoff.Model
	*onoff.ModelServer
}

func (b *binaryOnOff) GetOnOff(ctx context.Context, request *traits.GetOnOffRequest) (*traits.OnOff, error) {
	read, err := b.client.ReadProperty(b.device, bactypes.ReadPropertyData{
		Object: bactypes.Object{
			ID: bactypes.ObjectID(b.object.ID),
			Properties: []bactypes.Property{
				{ID: property.PresentValue, ArrayIndex: bactypes.ArrayAll},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var value bool
	switch v := read.Object.Properties[0].Data.(type) {
	case bool:
		value = v
	case uint32: // YABE room simulator uses this, not sure if that is "normal"
		value = v == 1
	default:
		b.logger.Warn("expected bool||uint32 return type", zap.Any("value", v))
	}

	var state traits.OnOff_State
	if value {
		state = traits.OnOff_ON
	} else {
		state = traits.OnOff_OFF
	}
	return b.model.UpdateOnOff(&traits.OnOff{State: state})
}

func (b *binaryOnOff) UpdateOnOff(ctx context.Context, request *traits.UpdateOnOffRequest) (*traits.OnOff, error) {
	return b.UnimplementedOnOffApiServer.UpdateOnOff(ctx, request)
}

func (b *binaryOnOff) PullOnOff(request *traits.PullOnOffRequest, server traits.OnOffApi_PullOnOffServer) error {
	return b.ModelServer.PullOnOff(request, server)
}

func (b *binaryOnOff) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(ObjectName(b.object),
		node.HasTrait(trait.OnOff, node.WithClients(onoff.WrapApi(b))),
	)
}
