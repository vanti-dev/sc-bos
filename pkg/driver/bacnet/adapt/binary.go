package adapt

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoffpb"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/property"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

// BinaryObject adapts a binary bacnet object as smart core traits.
func BinaryObject(prefix string, client *gobacnet.Client, device bactypes.Device, object config.Object, statuses *statuspb.Map) (node.SelfAnnouncer, error) {
	switch object.Trait {
	case "":
		return nil, ErrNoDefault
	case trait.OnOff:
		model := onoffpb.NewModel()
		return &binaryOnOff{
			prefix:   prefix,
			client:   client,
			device:   device,
			object:   object,
			statuses: statuses,

			model:       model,
			ModelServer: onoffpb.NewModelServer(model),
		}, nil
	}

	return nil, ErrNoAdaptation
}

type binaryOnOff struct {
	prefix   string
	client   *gobacnet.Client
	device   bactypes.Device
	object   config.Object
	statuses *statuspb.Map

	model *onoffpb.Model
	*onoffpb.ModelServer
}

func (b *binaryOnOff) GetOnOff(ctx context.Context, request *traits.GetOnOffRequest) (*traits.OnOff, error) {
	read, err := b.client.ReadProperty(ctx, b.device, bactypes.ReadPropertyData{
		Object: bactypes.Object{
			ID: bactypes.ObjectID(b.object.ID),
			Properties: []bactypes.Property{
				{ID: property.PresentValue, ArrayIndex: bactypes.ArrayAll},
			},
		},
	})

	updateRequestErrorStatus(b.statuses, b.name(), "getOnOff", err)
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
		return nil, status.Errorf(codes.Internal, "expected bool or uint32 return type for binary value, got %v", v)
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
	return a.Announce(b.name(),
		node.HasTrait(trait.OnOff, node.WithClients(onoffpb.WrapApi(b))),
	)
}

func (b *binaryOnOff) name() string {
	return b.prefix + ObjectName(b.object)
}
