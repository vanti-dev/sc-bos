// Package devices provides mechanisms for querying devices on a node.
// If you want to find all the lighting devices on floor 2, this is the package you need.
package devices

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/grpc"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Server struct {
	gen.UnimplementedDevicesApiServer

	parentName string
	// ChildPageSize overrides the default page size used when querying the parent trait for children
	ChildPageSize int32

	node node.Clienter
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

func (s *Server) ListDevices(ctx context.Context, request *gen.ListDevicesRequest) (*gen.ListDevicesResponse, error) {
	var parentApi traits.ParentApiClient
	var mdApi traits.MetadataApiClient
	if err := s.node.Client(&parentApi); err != nil {
		return nil, err
	}
	if err := s.node.Client(&mdApi); err != nil {
		return nil, err
	}

	var devices []*gen.Device

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

	childRequest := &traits.ListChildrenRequest{Name: s.parentName}
	if s.ChildPageSize > 0 {
		childRequest.PageSize = s.ChildPageSize
	}
full:
	for {
		childRequest.PageToken = pageToken.ParentPageToken
		childrenPage, err := parentApi.ListChildren(ctx, childRequest)
		if err != nil {
			return nil, err
		}
		children := childrenPage.Children
		if pageToken.PageIndex > 0 {
			children = children[pageToken.PageIndex:]
		}
		for i, child := range children {
			md, err := mdApi.GetMetadata(ctx, &traits.GetMetadataRequest{Name: child.Name})
			if err != nil {
				return nil, err
			}
			device := &gen.Device{
				Name:     child.Name,
				Metadata: md,
			}
			if deviceMatchesQuery(request.Query, device) {
				devices = append(devices, device)
				if len(devices) == pageSize {
					// check this child isn't the last child of the page
					if i == len(children)-1 {
						pageToken.ParentPageToken = childrenPage.NextPageToken
						pageToken.PageIndex = 0
						break full
					}

					pageToken.PageIndex += i + 1
					break full
				}
			}
		}

		// we've processed all children in the page without finding enough devices to fill a device page
		pageToken.ParentPageToken = childrenPage.NextPageToken
		pageToken.PageIndex = 0
		if pageToken.ParentPageToken == "" {
			// the response device list isn't full but we've run out of children pages
			break
		}
	}

	res := &gen.ListDevicesResponse{
		Devices: devices,
	}
	if pageToken.ParentPageToken != "" || pageToken.PageIndex > 0 {
		ptData, err := pageToken.encode()
		if err != nil {
			return nil, err
		}
		res.NextPageToken = ptData
	}
	return res, nil
}

func (s *Server) PullDevices(request *gen.PullDevicesRequest, server gen.DevicesApi_PullDevicesServer) error {
	// todo: implement PullDevices
	// I can't currently think of a good way to do this.
	// The simple (in code) way of listing all devices then calling PullMetadata for each one would cause an explosion
	// of go routines. While go can probably handle it there's just no need to do that.
	// Other solutions involve more tight integration with node.Node and the mechanism it uses to update metadata
	return s.UnimplementedDevicesApiServer.PullDevices(request, server)
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
