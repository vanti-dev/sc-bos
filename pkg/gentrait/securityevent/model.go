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

func (m *Model) GetSecurityEventCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.allSecurityEvents)
}

func (m *Model) ListSecurityEvents(start, count int) []*gen.SecurityEvent {

	var events []*gen.SecurityEvent
	// reverse to retrieve the latest events first
	for i := start - 1; i >= 0; i-- {
		events = append(events, m.allSecurityEvents[i])
		if len(events) >= count {
			break
		}
	}
	return events
}

func (m *Model) pullSecurityEventsWrapper(request *gen.PullSecurityEventsRequest, server gen.SecurityEventApi_PullSecurityEventsServer) error {
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
