package securityevent

import (
	"context"
	"strconv"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var subSystems = [2]string{"access control", "cctv"}

type Model struct {
	mu                sync.Mutex // guards allSecurityEvents and genId
	allSecurityEvents []*gen.SecurityEvent
	genId             int

	lastSecurityEvent *resource.Value // of *gen.SecurityEvent
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.SecurityEvent{})}
	opts = append(defaultOpts, opts...)

	m := &Model{
		lastSecurityEvent: resource.NewValue(opts...),
	}

	// let's add some events to start with so we can test the list method without waiting
	startTime := time.Now().Add(-100 * time.Minute)
	for i := 0; i < 100; i++ {
		_, _ = m.GenerateSecurityEvent(timestamppb.New(startTime))
		startTime = startTime.Add(time.Minute)
	}

	return m
}

// AddSecurityEvent manually add a security event to the model
func (m *Model) AddSecurityEvent(se *gen.SecurityEvent, opts ...resource.WriteOption) (*gen.SecurityEvent, error) {
	v, err := m.lastSecurityEvent.Set(se, opts...)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allSecurityEvents = append(m.allSecurityEvents, se)
	return v.(*gen.SecurityEvent), nil
}

// GenerateSecurityEvent generates a security event and adds it to the model
func (m *Model) GenerateSecurityEvent(time *timestamppb.Timestamp) (*gen.SecurityEvent, error) {
	m.mu.Lock()
	se := &gen.SecurityEvent{
		SecurityEventTime: time,
		Description:       "Generated security event",
		Id:                strconv.Itoa(m.genId),
		Source: &gen.SecurityEvent_Source{
			Subsystem: subSystems[m.genId%2],
		},
		EventType: gen.SecurityEvent_EventType(m.genId % 22),
	}
	m.genId++
	m.mu.Unlock()
	return m.AddSecurityEvent(se)
}

func (m *Model) ListSecurityEvents(req *gen.ListSecurityEventsRequest) (*gen.ListSecurityEventsResponse, error) {
	// page token is just the index of where we left off (if any)
	// this works with the current basic implementation because we only support a list of all events without filtering/sorting
	// and the events are stored in ascending chronological order. If this either of these things change, this will need to be rethought
	pageToken := req.GetPageToken()
	m.mu.Lock()
	defer m.mu.Unlock()
	startIndex := len(m.allSecurityEvents)
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
		resp.SecurityEvents = append(resp.SecurityEvents, m.allSecurityEvents[i])
		if len(resp.SecurityEvents) >= int(count) {
			resp.NextPageToken = strconv.Itoa(i - 1)
			break
		}
	}
	resp.TotalSize = int32(len(m.allSecurityEvents))
	return resp, nil
}

func (m *Model) PullSecurityEventsWrapper(request *gen.PullSecurityEventsRequest, server gen.SecurityEventApi_PullSecurityEventsServer) error {
	if !request.UpdatesOnly {
		m.mu.Lock()
		i := len(m.allSecurityEvents) - 50
		if i < 0 {
			i = 0
		}
		for ; i < len(m.allSecurityEvents)-1; i++ {
			change := &gen.PullSecurityEventsResponse_Change{
				Name:       request.Name,
				NewValue:   m.allSecurityEvents[i],
				ChangeTime: m.allSecurityEvents[i].SecurityEventTime,
				Type:       types.ChangeType_ADD,
			}
			if err := server.Send(&gen.PullSecurityEventsResponse{Changes: []*gen.PullSecurityEventsResponse_Change{change}}); err != nil {
				m.mu.Unlock()
				return err
			}
		}
		m.mu.Unlock()
	}
	for change := range m.PullSecurityEvents(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullSecurityEventsResponse{}
		msg.Changes = append(msg.Changes, change)
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) PullSecurityEvents(ctx context.Context, opts ...resource.ReadOption) <-chan *gen.PullSecurityEventsResponse_Change {
	send := make(chan *gen.PullSecurityEventsResponse_Change)
	recv := m.lastSecurityEvent.Pull(ctx, opts...)
	go func() {
		defer close(send)
		for change := range recv {
			value := change.Value.(*gen.SecurityEvent)
			send <- &gen.PullSecurityEventsResponse_Change{
				OldValue:   nil,
				NewValue:   value, // the mock driver only generates new security events and does not delete them
				ChangeTime: value.SecurityEventTime,
				Type:       types.ChangeType_ADD,
			}
		}
	}()

	return send
}

func (m *Model) unlock() {
	m.mu.Unlock()
}

func (m *Model) lock() {
	m.mu.Lock()
}
