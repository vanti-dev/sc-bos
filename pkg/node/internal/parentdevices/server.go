package parentdevices

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/devicespb"
)

// Server implements the ParentApi where children are devices in a devicespb.Collection.
type Server struct {
	traits.UnimplementedParentApiServer
	devices    Collection
	parentName string
}

// Collection contains a list of devices.
type Collection interface {
	ListDevices(opts ...resource.ReadOption) []*gen.Device
	PullDevices(ctx context.Context, opts ...resource.ReadOption) <-chan devicespb.DevicesChange
}

func NewServer(parentName string, devices Collection) *Server {
	return &Server{
		parentName: parentName,
		devices:    devices,
	}
}

func (s *Server) ListChildren(_ context.Context, request *traits.ListChildrenRequest) (*traits.ListChildrenResponse, error) {
	if err := s.validateRequest(request); err != nil {
		return nil, err
	}
	devices, totalSize, nextPageToken, err := listPage(request, func() []*gen.Device {
		return s.devices.ListDevices(resource.WithInclude(s.deviceIncludeFunc))
	})
	if err != nil {
		return nil, err
	}

	res := &traits.ListChildrenResponse{
		Children:      make([]*traits.Child, len(devices)),
		TotalSize:     int32(totalSize),
		NextPageToken: nextPageToken,
	}

	mask := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))
	for i, device := range devices {
		c := deviceToChild(device)
		mask.Filter(c)
		res.Children[i] = c
	}
	return res, nil
}

func (s *Server) PullChildren(request *traits.PullChildrenRequest, stream grpc.ServerStreamingServer[traits.PullChildrenResponse]) error {
	if err := s.validateRequest(request); err != nil {
		return err
	}
	filter := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))
	for change := range s.devices.PullDevices(stream.Context(), resource.WithInclude(s.deviceIncludeFunc), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		changeProto := s.devicesChangeToProto(change, filter)
		err := stream.Send(&traits.PullChildrenResponse{Changes: []*traits.PullChildrenResponse_Change{changeProto}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) validateRequest(req interface{ GetName() string }) error {
	if name := req.GetName(); name != s.parentName {
		return status.Error(codes.NotFound, name)
	}
	return nil
}

func (s *Server) deviceIncludeFunc(id string, item proto.Message) bool {
	// skip the parent device, it can't be its own child
	if id == s.parentName {
		return false
	}
	deviceItem, ok := item.(*gen.Device)
	if !ok {
		return false
	}

	// todo: decide if this legacy behaviour is still valid
	// Node.Announce used to only add children when there were traits,
	// so maintain that logic
	if len(deviceItem.GetMetadata().GetTraits()) == 0 {
		return false
	}

	return true
}

func deviceToChild(d *gen.Device) *traits.Child {
	c := &traits.Child{Name: d.Name}
	for _, t := range d.GetMetadata().GetTraits() {
		c.Traits = append(c.Traits, &traits.Trait{Name: t.Name})
	}
	return c
}

func (s *Server) devicesChangeToProto(c devicespb.DevicesChange, filter *masks.ResponseFilter) *traits.PullChildrenResponse_Change {
	res := &traits.PullChildrenResponse_Change{
		Name:       s.parentName,
		Type:       c.ChangeType,
		ChangeTime: timestamppb.New(c.ChangeTime),
	}
	if c.OldValue != nil {
		res.OldValue = deviceToChild(c.OldValue)
		filter.Filter(res.OldValue)
	}
	if c.NewValue != nil {
		res.NewValue = deviceToChild(c.NewValue)
		filter.Filter(res.NewValue)
	}
	return res
}
