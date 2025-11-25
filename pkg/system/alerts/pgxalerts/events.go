package pgxalerts

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-bos/pkg/gen"
)

func (s *Server) notifyAdd(name string, alert *gen.Alert) {
	// notify
	s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
		Name:       name,
		Type:       types.ChangeType_ADD,
		ChangeTime: alert.CreateTime,
		NewValue:   alert,
	})
}

func (s *Server) notifyUpdate(name string, original *gen.Alert, updated *gen.Alert) int {
	return s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
		Name:       name,
		Type:       types.ChangeType_UPDATE,
		ChangeTime: timestamppb.Now(),
		OldValue:   original,
		NewValue:   updated,
	})
}

func (s *Server) notifyRemove(name string, existing *gen.Alert) int {
	return s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
		Name:       name,
		Type:       types.ChangeType_REMOVE,
		ChangeTime: timestamppb.Now(),
		OldValue:   existing,
	})
}
