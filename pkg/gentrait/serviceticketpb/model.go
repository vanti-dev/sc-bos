package serviceticketpb

import (
	"errors"

	"github.com/pborman/uuid"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Model struct {
	tickets map[string]*gen.Ticket // TicketId -> Ticket
	support *gen.TicketSupport
}

func NewModel(opts ...resource.Option) *Model {
	m := &Model{}
	m.tickets = make(map[string]*gen.Ticket)
	return m
}

func (m Model) addTicket(ticket *gen.Ticket) *gen.Ticket {
	id := uuid.NewUUID().String()
	ticket.Id = id
	m.tickets[id] = ticket
	return ticket
}

func (m Model) updateTicket(ticket *gen.Ticket) (*gen.Ticket, error) {
	if _, ok := m.tickets[ticket.Id]; !ok {
		return nil, errors.New("ticket not found")
	}
	m.tickets[ticket.Id] = ticket
	return ticket, nil
}

// SetSupport sets the gen.TicketSupport to use in the ServiceTicketInfoServer.
func (m Model) SetSupport(s *gen.TicketSupport) {
	m.support = s
}
