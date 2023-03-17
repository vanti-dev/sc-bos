package hubalerts

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/util/once"
)

// Server backs off to a remove server, modifying requests to filter by federation parameter.
type Server struct {
	gen.UnimplementedAlertApiServer
	gen.UnimplementedAlertAdminApiServer

	remoteName string // can be empty
	federation string
	remoteNode node.Remote

	alertsOnce once.RetryError
	alerts     gen.AlertApiClient
	alertAdmin gen.AlertAdminApiClient

	// to support alert metadata, which we have to track ourselves
	mdOnce once.RetryError
	md     *resource.Value // of *gen.AlertMetadata, used to track changes
	mdStop func()          // closes any go routines that are listening for changes
}

func NewServer(remoteName, localName string, remoteNode node.Remote) *Server {
	return &Server{
		remoteName: remoteName,
		federation: localName,
		remoteNode: remoteNode,
	}
}

func (s *Server) CreateAlert(ctx context.Context, request *gen.CreateAlertRequest) (*gen.Alert, error) {
	if err := s.initConn(ctx); err != nil {
		return nil, err
	}
	request.Name = s.remoteName
	request.Alert.Federation = s.federation
	return s.alertAdmin.CreateAlert(ctx, request)
}

func (s *Server) UpdateAlert(ctx context.Context, request *gen.UpdateAlertRequest) (*gen.Alert, error) {
	if err := s.initConn(ctx); err != nil {
		return nil, err
	}
	request.Name = s.remoteName
	request.Alert.Federation = s.federation
	return s.alertAdmin.UpdateAlert(ctx, request)
}

func (s *Server) DeleteAlert(ctx context.Context, request *gen.DeleteAlertRequest) (*gen.DeleteAlertResponse, error) {
	if err := s.initConn(ctx); err != nil {
		return nil, err
	}
	request.Name = s.remoteName
	return s.alertAdmin.DeleteAlert(ctx, request)
}

func (s *Server) ListAlerts(ctx context.Context, request *gen.ListAlertsRequest) (*gen.ListAlertsResponse, error) {
	if err := s.initConn(ctx); err != nil {
		return nil, err
	}
	request.Name = s.remoteName
	if request.Query == nil {
		request.Query = &gen.Alert_Query{}
	}
	request.Query.Federation = s.federation
	return s.alerts.ListAlerts(ctx, request)
}

func (s *Server) PullAlerts(request *gen.PullAlertsRequest, server gen.AlertApi_PullAlertsServer) error {
	if err := s.initConn(server.Context()); err != nil {
		return err
	}
	request.Name = s.remoteName
	if request.Query == nil {
		request.Query = &gen.Alert_Query{}
	}
	request.Query.Federation = s.federation
	stream, err := s.alerts.PullAlerts(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (s *Server) AcknowledgeAlert(ctx context.Context, request *gen.AcknowledgeAlertRequest) (*gen.Alert, error) {
	if err := s.initConn(ctx); err != nil {
		return nil, err
	}
	request.Name = s.remoteName
	return s.alerts.AcknowledgeAlert(ctx, request)
}

func (s *Server) UnacknowledgeAlert(ctx context.Context, request *gen.AcknowledgeAlertRequest) (*gen.Alert, error) {
	if err := s.initConn(ctx); err != nil {
		return nil, err
	}
	request.Name = s.remoteName
	return s.alerts.UnacknowledgeAlert(ctx, request)
}

func (s *Server) GetAlertMetadata(ctx context.Context, request *gen.GetAlertMetadataRequest) (*gen.AlertMetadata, error) {
	if err := s.initAlertMetadata(ctx); err != nil {
		return nil, err
	}
	return s.md.Get(resource.WithReadMask(request.ReadMask)).(*gen.AlertMetadata), nil
}

func (s *Server) PullAlertMetadata(request *gen.PullAlertMetadataRequest, server gen.AlertApi_PullAlertMetadataServer) error {
	if err := s.initAlertMetadata(server.Context()); err != nil {
		return err
	}
	for change := range s.md.Pull(server.Context(), resource.WithReadMask(request.GetReadMask()), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		err := server.Send(&gen.PullAlertMetadataResponse{Changes: []*gen.PullAlertMetadataResponse_Change{
			{
				Name:       request.Name,
				ChangeTime: timestamppb.New(change.ChangeTime),
				Metadata:   change.Value.(*gen.AlertMetadata),
			},
		}})
		if err != nil {
			return err
		}
	}

	return nil
}

// initConn sets up the remote clients alerts and alertAdmin.
func (s *Server) initConn(ctx context.Context) error {
	return s.alertsOnce.Do(ctx, func() error {
		conn, err := s.remoteNode.Connect(ctx)
		if err != nil {
			return err
		}
		s.alerts = gen.NewAlertApiClient(conn)
		s.alertAdmin = gen.NewAlertAdminApiClient(conn)
		return nil
	})
}

// initAlertMetadata seeds s.md and keeps it up to date with changes to the underlying store
// if s.md hasn't already been seeded.
// initAlertMetadata blocks until either s.md is seeded or an error occurs or ctx expires.
// If an error occurs seeding s.md then it is returned
func (s *Server) initAlertMetadata(ctx context.Context) error {
	return s.mdOnce.Do(ctx, func() error {
		var ctx context.Context
		ctx, s.mdStop = context.WithCancel(context.Background())

		// Collect initial stats from the DB
		md := &gen.AlertMetadata{
			AcknowledgedCounts: make(map[bool]uint32),
			FloorCounts:        make(map[string]uint32),
			ZoneCounts:         make(map[string]uint32),
			SeverityCounts:     make(map[int32]uint32),
		}
		val := resource.NewValue(resource.WithInitialValue(md))

		// We need to pull and list the alerts from the server.
		// Normally we'd be able to do this via a single pull call but the pgx server doesn't support that :(
		// For now we'll have to live with this race condition
		query := &gen.Alert_Query{Federation: s.federation}
		stream, err := s.alerts.PullAlerts(ctx, &gen.PullAlertsRequest{Name: s.remoteName, Query: query})
		if err != nil {
			return err
		}
		listReq := &gen.ListAlertsRequest{Name: s.remoteName, Query: query}
		for {
			list, err := s.alerts.ListAlerts(ctx, listReq)
			if err != nil {
				return err
			}
			listReq.PageToken = list.NextPageToken

			for _, alert := range list.Alerts {
				err := applyMdDelta(val, &gen.PullAlertsResponse_Change{
					Type:     types.ChangeType_ADD,
					NewValue: alert,
				})
				if err != nil {
					return err
				}
			}

			if listReq.PageToken == "" {
				break
			}
		}

		s.md = val

		// setup listeners for changes so we can track those changes in metadata
		go func() {
			for {
				msg, err := stream.Recv()
				if err != nil {
					stream, err = s.pullAlertsAgain(ctx, &gen.PullAlertsRequest{Name: s.remoteName, Query: query})
					if err != nil {
						return // ctx done, aka server stopped
					}
				}

				for _, change := range msg.Changes {
					applyMdDelta(s.md, change)
				}
			}
		}()
		return nil
	})
}

func applyMdDelta(md *resource.Value, e *gen.PullAlertsResponse_Change) error {
	// note: AlertMetadata uses uint32 so we can't use it to store our deltas
	_, err := md.Set(&gen.AlertMetadata{}, resource.InterceptBefore(func(old, new proto.Message) {
		oldMd, newMd := old.(*gen.AlertMetadata), new.(*gen.AlertMetadata)
		proto.Merge(newMd, oldMd)

		// total
		if e.OldValue == nil && e.NewValue != nil {
			newMd.TotalCount++
		} else if e.OldValue != nil && e.NewValue == nil {
			newMd.TotalCount--
			if newMd.TotalCount < 0 { // shouldn't be needed, but just in case
				newMd.TotalCount = 0
			}
		}

		// ack/nak
		if e.OldValue == nil && e.NewValue != nil {
			acked := e.NewValue.GetAcknowledgement().GetAcknowledgeTime() != nil
			if acked {
				newMd.AcknowledgedCounts[true]++
			} else {
				newMd.AcknowledgedCounts[false]++
			}
		} else if e.OldValue != nil && e.NewValue == nil {
			acked := e.OldValue.GetAcknowledgement().GetAcknowledgeTime() != nil
			if acked {
				newMd.AcknowledgedCounts[true]--
			} else {
				newMd.AcknowledgedCounts[false]--
			}
		} else {
			oldAck, newAck := e.GetOldValue().GetAcknowledgement().GetAcknowledgeTime(),
				e.GetNewValue().GetAcknowledgement().GetAcknowledgeTime()
			if oldAck == nil && newAck != nil {
				newMd.AcknowledgedCounts[true]++
				if newMd.AcknowledgedCounts[false] > 0 { // just in case
					newMd.AcknowledgedCounts[false]--
				}
			} else if oldAck != nil && newAck == nil {
				if newMd.AcknowledgedCounts[true] > 0 { // just in case
					newMd.AcknowledgedCounts[true]--
				}
				newMd.AcknowledgedCounts[false]++
			}
		}

		// floors, zones, and severity
		mapDelta(e.GetOldValue().GetFloor(), e.GetNewValue().GetFloor(), newMd.FloorCounts)
		mapDelta(e.GetOldValue().GetZone(), e.GetNewValue().GetZone(), newMd.ZoneCounts)
		mapDelta(int32(e.GetOldValue().GetSeverity()), int32(e.GetNewValue().GetSeverity()), newMd.SeverityCounts)
	}))

	return err
}

func mapDelta[T comparable](o, n T, m map[T]uint32) {
	if o == n {
		return
	}
	var zero T
	if o != zero {
		if c, ok := m[o]; ok && c > 0 {
			m[o] = c - 1
		}
	}
	if n != zero {
		m[n]++
	}
}

// pullAlertsAgain calls s.alerts.PullAlerts, retrying on failure until ctx is done.
func (s *Server) pullAlertsAgain(ctx context.Context, req *gen.PullAlertsRequest) (gen.AlertApi_PullAlertsClient, error) {
	scale := 1.1
	max := 20 * time.Second
	delay := 10 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			// this is the only non-successful return we have
			return nil, ctx.Err()
		default:
		}

		stream, err := s.alerts.PullAlerts(ctx, req)
		if err == nil {
			return stream, nil
		}

		time.Sleep(delay)
		delay = time.Duration(float64(delay) * scale)
		if delay > max {
			delay = max
		}
	}
}
