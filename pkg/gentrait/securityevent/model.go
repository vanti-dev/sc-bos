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
	lastSecurityEvent *resource.Value // of *gen.SecurityEvent
	allSecurityEvents []*gen.SecurityEvent
	genId             int
	mu                sync.Mutex
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
	m.mu.Unlock()
	m.allSecurityEvents = append(m.allSecurityEvents, se)
	return v.(*gen.SecurityEvent), nil
}

// GenerateSecurityEvent generates a security event and adds it to the model
func (m *Model) GenerateSecurityEvent(time *timestamppb.Timestamp) (*gen.SecurityEvent, error) {
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
	return m.AddSecurityEvent(se)
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
