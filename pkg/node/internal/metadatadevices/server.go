package metadatadevices

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/devicespb"
)

// Server implements the MetadataApi backed by a devicespb.Collection.
type Server struct {
	traits.UnimplementedMetadataApiServer
	devices Collection
}

// Collection contains devices keyed by their name.
type Collection interface {
	GetDevice(name string, opts ...resource.ReadOption) (*gen.Device, error)
	PullDevice(ctx context.Context, name string, opts ...resource.ReadOption) <-chan devicespb.DeviceChange
}

func NewServer(devices Collection) *Server {
	return &Server{
		devices: devices,
	}
}

func (s Server) GetMetadata(_ context.Context, request *traits.GetMetadataRequest) (*traits.Metadata, error) {
	device, err := s.devices.GetDevice(request.Name)
	if err != nil {
		return nil, err
	}
	filter := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))
	return filter.FilterClone(device.Metadata).(*traits.Metadata), nil
}

func (s Server) PullMetadata(request *traits.PullMetadataRequest, g grpc.ServerStreamingServer[traits.PullMetadataResponse]) error {
	filter := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))
	for change := range s.devices.PullDevice(g.Context(), request.Name, resource.WithUpdatesOnly(request.UpdatesOnly)) {
		mdChange := deviceChangeToProto(change, filter)
		err := g.Send(&traits.PullMetadataResponse{Changes: []*traits.PullMetadataResponse_Change{mdChange}})
		if err != nil {
			return err
		}
	}
	return nil
}

func deviceChangeToProto(c devicespb.DeviceChange, filter *masks.ResponseFilter) *traits.PullMetadataResponse_Change {
	res := &traits.PullMetadataResponse_Change{
		Name:       c.Name,
		ChangeTime: timestamppb.New(c.ChangeTime),
	}
	if c.Device.GetMetadata() == nil {
		res.Metadata = &traits.Metadata{}
	} else {
		res.Metadata = filter.FilterClone(c.Device.GetMetadata()).(*traits.Metadata)
	}
	return res
}
