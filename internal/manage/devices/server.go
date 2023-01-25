// Package devices provides mechanisms for querying devices on a node.
// If you want to find all the lighting devices on floor 2, this is the package you need.
package devices

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Server struct {
	gen.UnimplementedDevicesApiServer

	parentName string
	// ChildPageSize overrides the default page size used when querying the parent trait for children
	ChildPageSize int32

	node *node.Node
}

func NewServer(n *node.Node) *Server {
	return &Server{
		parentName: n.Name(),
		node:       n,
	}
}

func (s *Server) Register(server *grpc.Server) {
	gen.RegisterDevicesApiServer(server, s)
}

func (s *Server) ListDevices(_ context.Context, request *gen.ListDevicesRequest) (*gen.ListDevicesResponse, error) {
	var pageToken PageToken
	if request.PageToken != "" {
		var err error
		pageToken, err = decodePageToken(request.PageToken)
		if err != nil {
			return nil, err
		}
	}

	pageSize := 50
	if request.PageSize > 0 {
		pageSize = int(request.PageSize)
		if pageSize > 1000 {
			pageSize = 1000
		}
	}

	// note: allMetadata is already sorted by name
	allMetadata := s.node.ListAllMetadata(
		resource.WithReadMask(subMask(request.ReadMask, "metadata")),
		resource.WithInclude(func(id string, item proto.Message) bool {
			device := &gen.Device{
				Name:     id,
				Metadata: item.(*traits.Metadata),
			}
			return deviceMatchesQuery(request.Query, device)
		}),
	)
	nextIndex := 0
	if pageToken.LastName != "" {
		nextIndex = sort.Search(len(allMetadata), func(i int) bool {
			return allMetadata[i].Name >= pageToken.LastName
		})
		if nextIndex < len(allMetadata) && allMetadata[nextIndex].Name == pageToken.LastName {
			nextIndex++
		}
		pageToken.LastName = ""
	}

	var devices []*gen.Device
	for _, md := range allMetadata[nextIndex:] {
		device := &gen.Device{
			Name:     md.Name,
			Metadata: md,
		}
		if len(devices) == pageSize {
			// we found another device but we don't want to include it in the response,
			// we'll use the info to know whether to populate the next page token
			pageToken.LastName = devices[len(devices)-1].Name
			break
		}
		devices = append(devices, device)
	}

	res := &gen.ListDevicesResponse{
		Devices:   devices,
		TotalSize: int32(len(allMetadata)),
	}
	if pageToken.LastName != "" {
		ptStr, err := pageToken.encode()
		if err != nil {
			return nil, err
		}
		res.NextPageToken = ptStr
	}
	return res, nil
}

func (s *Server) PullDevices(request *gen.PullDevicesRequest, server gen.DevicesApi_PullDevicesServer) error {
	changes := s.node.PullAllMetadata(server.Context(),
		resource.WithUpdatesOnly(request.UpdatesOnly),
		resource.WithReadMask(subMask(request.ReadMask, "metadata")),
		resource.WithInclude(func(id string, item proto.Message) bool {
			device := &gen.Device{
				Name:     id,
				Metadata: item.(*traits.Metadata),
			}
			return deviceMatchesQuery(request.Query, device)
		}),
	)
	for change := range changes {
		resChange := &gen.PullDevicesResponse_Change{
			Name:       change.Name,
			ChangeTime: timestamppb.New(change.ChangeTime),
			Type:       change.Type,
		}
		if change.OldValue != nil {
			resChange.OldValue = &gen.Device{Name: change.Name, Metadata: change.OldValue}
		}
		if change.NewValue != nil {
			resChange.NewValue = &gen.Device{Name: change.Name, Metadata: change.NewValue}
		}
		res := &gen.PullDevicesResponse{Changes: []*gen.PullDevicesResponse_Change{resChange}}
		if err := server.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func subMask(mask *fieldmaskpb.FieldMask, prefix string) *fieldmaskpb.FieldMask {
	pd := fmt.Sprintf("%s.", prefix)
	if mask != nil {
		// the request read mask is prefixed with "metadata.", we need to remove that
		newMask := &fieldmaskpb.FieldMask{}
		for _, path := range mask.Paths {
			if path == prefix {
				// need all fields
				return nil
			}
			if strings.HasPrefix(path, pd) {
				newMask.Paths = append(newMask.Paths, path[len("metadata."):])
			}
		}
		return newMask
	}
	return nil
}

func deviceMatchesQuery(query *gen.Device_Query, device *gen.Device) bool {
	if query == nil {
		return true
	}
	for _, condition := range query.Conditions {
		if !conditionMatches(condition, device) {
			return false
		}
	}

	// this means a query with no conditions always returns true
	return true
}

func conditionMatches(cond *gen.Device_Query_Condition, device *gen.Device) bool {
	switch c := cond.Value.(type) {
	case *gen.Device_Query_Condition_StringEqual:
		return isMessageValueEqualString(cond.Field, c.StringEqual, device)
	default:
		return false
	}
}
