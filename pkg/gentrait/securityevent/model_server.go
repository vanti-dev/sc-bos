package securityevent

import (
	"context"
	"strconv"

	"google.golang.org/grpc"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedSecurityEventApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterSecurityEventApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) ListSecurityEvents(_ context.Context, req *gen.ListSecurityEventsRequest) (*gen.ListSecurityEventsResponse, error) {

	// page token is just the index of where we left off (if any)
	// this works with the current basic implementation because we only support a list of all events without filtering/sorting
	// and the events are stored in ascending chronological order. If this either of these things change, this will need to be rethought
	pageToken := req.GetPageToken()
	m.model.mu.Lock()
	defer m.model.mu.Unlock()
	startIndex := len(m.model.allSecurityEvents)
	if pageToken != "" {
		_, err := strconv.Atoi(req.GetPageToken())
		if err != nil {
			return nil, err
		}
		startIndex, _ = strconv.Atoi(pageToken)
	}

	count := req.PageSize
	if count == 0 {
		count = 50
	} else if count > 1000 {
		count = 1000
	}

	resp := &gen.ListSecurityEventsResponse{}

	// reverse to retrieve the latest events first
	for i := startIndex - 1; i >= 0; i-- {
		resp.SecurityEvents = append(resp.SecurityEvents, m.model.allSecurityEvents[i])
		if len(resp.SecurityEvents) >= int(count) {
			resp.NextPageToken = strconv.Itoa(i - 1)
			break
		}
	}
	resp.TotalSize = int32(len(m.model.allSecurityEvents))
	return resp, nil
}

// PullSecurityEvents returns a channel of security events
// If updatesOnly is false, only the previous 50 events will be sent before any new events
// For historical events use ListSecurityEvents
func (m *ModelServer) PullSecurityEvents(request *gen.PullSecurityEventsRequest, server gen.SecurityEventApi_PullSecurityEventsServer) error {
	if !request.UpdatesOnly {
		m.model.mu.Lock()
		i := len(m.model.allSecurityEvents) - 50
		if i < 0 {
			i = 0
		}
		for ; i < len(m.model.allSecurityEvents)-1; i++ {
			change := &gen.PullSecurityEventsResponse_Change{
				Name:       request.Name,
				NewValue:   m.model.allSecurityEvents[i],
				ChangeTime: m.model.allSecurityEvents[i].SecurityEventTime,
				Type:       types.ChangeType_ADD,
			}
			if err := server.Send(&gen.PullSecurityEventsResponse{Changes: []*gen.PullSecurityEventsResponse_Change{change}}); err != nil {
				m.model.mu.Unlock()
				return err
			}
		}
		m.model.mu.Unlock()
	}
	for change := range m.model.PullSecurityEvents(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		var msg gen.PullSecurityEventsResponse
		msg.Changes = append(msg.Changes, change)
		if err := server.Send(&msg); err != nil {
			return err
		}
	}
	return nil
}
