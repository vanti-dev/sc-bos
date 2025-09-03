// Package devices provides mechanisms for querying devices on a node.
// If you want to find all the lighting devices on floor 2, this is the package you need.
package devices

import (
	"context"
	"net/url"
	"sort"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

type Server struct {
	gen.UnimplementedDevicesApiServer

	// ChildPageSize overrides the default page size used when querying the parent trait for children
	ChildPageSize int32

	node *node.Node
	now  func() time.Time

	downloadUrlBase      url.URL // defaults to /dl/devices
	downloadTokenWriter  DownloadTokenWriter
	downloadTokenReader  DownloadTokenReader
	downloadKey          func() ([]byte, error)
	downloadExpiry       time.Duration // defaults to 1 hour
	downloadExpiryLeeway time.Duration // defaults to 1 minute
	downloadPageTimeout  time.Duration // defaults to 10 seconds, applies to get and history cursor calls
}

func NewServer(n *node.Node, opts ...Option) *Server {
	s := &Server{
		node:                 n,
		now:                  time.Now,
		downloadUrlBase:      url.URL{Path: "/dl/devices"},
		downloadExpiry:       time.Hour,
		downloadExpiryLeeway: time.Minute,
		downloadKey:          newHMACKeyGen(64), // todo: replace with something that works between nodes
		downloadPageTimeout:  10 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Register(server *grpc.Server) {
	gen.RegisterDevicesApiServer(server, s)
}

func (s *Server) ListDevices(_ context.Context, request *gen.ListDevicesRequest) (*gen.ListDevicesResponse, error) {
	if err := validateQuery(request.GetQuery()); err != nil {
		return nil, err
	}

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

	// note: allDevices is already sorted by name
	allDevices := s.node.ListDevices(
		resource.WithReadMask(request.ReadMask),
		resource.WithInclude(func(id string, item proto.Message) bool {
			if item == nil {
				return false
			}
			device := item.(*gen.Device)
			return deviceMatchesQuery(request.Query, device)
		}),
	)
	nextIndex := 0
	if pageToken.LastName != "" {
		nextIndex = sort.Search(len(allDevices), func(i int) bool {
			return allDevices[i].Name >= pageToken.LastName
		})
		if nextIndex < len(allDevices) && allDevices[nextIndex].Name == pageToken.LastName {
			nextIndex++
		}
		pageToken.LastName = ""
	}

	var devices []*gen.Device
	for _, device := range allDevices[nextIndex:] {
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
		TotalSize: int32(len(allDevices)),
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
	if err := validateQuery(request.GetQuery()); err != nil {
		return err
	}

	changes := s.node.PullDevices(server.Context(),
		resource.WithUpdatesOnly(request.UpdatesOnly),
		resource.WithReadMask(request.ReadMask),
		resource.WithInclude(func(id string, item proto.Message) bool {
			if item == nil {
				return false
			}
			device := item.(*gen.Device)
			return deviceMatchesQuery(request.Query, device)
		}),
	)
	for change := range changes {
		resChange := &gen.PullDevicesResponse_Change{
			Name:       change.Id,
			ChangeTime: timestamppb.New(change.ChangeTime),
			Type:       change.ChangeType,
			OldValue:   change.OldValue,
			NewValue:   change.NewValue,
		}
		res := &gen.PullDevicesResponse{Changes: []*gen.PullDevicesResponse_Change{resChange}}
		if err := server.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) GetDevicesMetadata(_ context.Context, request *gen.GetDevicesMetadataRequest) (*gen.DevicesMetadata, error) {
	devices := s.node.ListDevices()
	var res *gen.DevicesMetadata
	col := newMetadataCollector(request.GetIncludes().GetFields()...)
	for _, device := range devices {
		res = col.add(device)
	}
	return res, nil
}

func (s *Server) PullDevicesMetadata(request *gen.PullDevicesMetadataRequest, server gen.DevicesApi_PullDevicesMetadataServer) error {
	var md *gen.DevicesMetadata
	send := func(msg *gen.DevicesMetadata, t time.Time) error {
		if proto.Equal(md, msg) {
			return nil
		}
		md = proto.Clone(msg).(*gen.DevicesMetadata)
		change := &gen.PullDevicesMetadataResponse_Change{
			ChangeTime:      timestamppb.New(t),
			DevicesMetadata: msg,
		}
		return server.Send(&gen.PullDevicesMetadataResponse{Changes: []*gen.PullDevicesMetadataResponse_Change{
			change,
		}})
	}

	// do this before getting initial values
	changes := s.node.PullDevices(server.Context())

	// send initial values.
	// Note we recalculate the metadata for the initial value and for updates separately. We can't guarantee data
	// consistency between the Get and Pull calls, at least this way the data should be accurate.
	if !request.UpdatesOnly {
		md, err := s.GetDevicesMetadata(server.Context(), &gen.GetDevicesMetadataRequest{Includes: request.Includes})
		if err != nil {
			return err
		}
		if err := send(md, s.now()); err != nil {
			return err
		}
	}

	// watch for and send updates to metadata
	col := newMetadataCollector(request.Includes.Fields...)
	seeding := true
	for change := range changes {
		var md *gen.DevicesMetadata
		if change.OldValue != nil {
			md = col.remove(change.OldValue)
		}
		if change.NewValue != nil {
			md = col.add(change.NewValue)
		}

		if change.LastSeedValue {
			seeding = false // technically not true as this is the last seed value, but we'll deal with that below
		}
		if seeding || (change.LastSeedValue && request.UpdatesOnly) {
			// avoid sending an update for each of the initial values that seed the collector
			continue
		}
		if err := send(md, change.ChangeTime); err != nil {
			return err
		}
	}

	return nil
}

func validateQuery(q *gen.Device_Query) error {
	if q == nil {
		return nil
	}
	for _, c := range q.Conditions {
		if in := c.GetStringIn(); in != nil {
			if len(in.Strings) > 100 {
				return status.Errorf(codes.InvalidArgument, "condition string_in len > 100")
			}
		}
		if in := c.GetStringInFold(); in != nil {
			if len(in.Strings) > 100 {
				return status.Errorf(codes.InvalidArgument, "condition string_in_fold len > 100")
			}
		}
	}
	return nil
}
