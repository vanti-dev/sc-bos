package pgxalerts

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

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

// initAlertMetadata seeds s.md and keeps it up to date with changes to the underlying store
// if s.md hasn't already been seeded.
// initAlertMetadata blocks until either s.md is seeded or an error occurs or ctx expires.
// If an error occurs seeding s.md then it is returned
func (s *Server) initAlertMetadata(ctx context.Context) error {
	s.mdMu.Lock()
	defer s.mdMu.Unlock()
	wait := func() error {
		select {
		case <-s.mdC:
			return s.mdErr
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if s.mdC != nil {
		// init is either in progress or has completed
		return wait()
	}

	// need to init the alert md
	s.mdC = make(chan struct{})
	ctx, s.mdStop = context.WithCancel(context.Background())
	go func() {
		// no need to lock here as we already know we're the only ones writing to mdC or mdErr
		defer close(s.mdC)

		// These are populated in the db transaction to make sure we don't miss any events.
		// There is a race condition here, though very unlikely (I think).
		// There is space between our transaction and other RPC transactions (like AckAlert)
		// for events to be missed or duplicated.
		// I apologise to any future maintainer who find this comment and has to fix it, likely under pressure.
		var events <-chan *gen.PullAlertsResponse_Change
		var eventsCtx context.Context
		var eventsCancel context.CancelFunc

		// Collect initial stats from the DB
		md := &gen.AlertMetadata{
			AcknowledgedCounts:   make(map[bool]uint32),
			FloorCounts:          make(map[string]uint32),
			ZoneCounts:           make(map[string]uint32),
			SeverityCounts:       make(map[int32]uint32),
			ResolvedCounts:       make(map[bool]uint32),
			NeedsAttentionCounts: make(map[string]uint32),
		}
		err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
			// totals
			err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts`).Scan(&md.TotalCount)
			if err != nil {
				return fmt.Errorf("count all %w", err)
			}

			// ack and nak counts
			var ackCount uint32
			err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE ack_time IS NOT NULL`).Scan(&ackCount)
			if err != nil {
				return fmt.Errorf("nak count %w", err)
			}
			md.AcknowledgedCounts = map[bool]uint32{
				true:  ackCount,
				false: md.TotalCount - ackCount,
			}
			// resolved and unresolved counts
			var resolveCounts uint32
			err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE resolve_time IS NOT NULL`).Scan(&resolveCounts)
			if err != nil {
				return fmt.Errorf("resolve count %w", err)
			}
			md.ResolvedCounts = map[bool]uint32{
				true:  resolveCounts,
				false: md.TotalCount - resolveCounts,
			}

			// needs attention counts
			var needsAttentionCount uint32
			err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE ack_time IS NOT NULL AND resolve_time IS NOT NULL`).Scan(&needsAttentionCount)
			if err != nil {
				return fmt.Errorf("ack and resolved count %w", err)
			}
			md.NeedsAttentionCounts["ack_resolved"] = needsAttentionCount
			err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE ack_time IS NOT NULL AND resolve_time IS NULL`).Scan(&needsAttentionCount)
			if err != nil {
				return fmt.Errorf("ack and unresolved count %w", err)
			}
			md.NeedsAttentionCounts["ack_unresolved"] = needsAttentionCount
			err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE ack_time IS NULL AND resolve_time IS NOT NULL`).Scan(&needsAttentionCount)
			if err != nil {
				return fmt.Errorf("ack and resolved count %w", err)
			}
			md.NeedsAttentionCounts["nack_resolved"] = needsAttentionCount
			err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE ack_time IS NULL AND resolve_time IS NULL`).Scan(&needsAttentionCount)
			if err != nil {
				return fmt.Errorf("ack and unresolved count %w", err)
			}
			md.NeedsAttentionCounts["nack_unresolved"] = needsAttentionCount

			// floors
			floors, err := queryGroupCounts(ctx, tx, "floor", s.Floors)
			if err != nil {
				return fmt.Errorf("floors %w", err)
			}
			md.FloorCounts = floors
			// zones
			zones, err := queryGroupCounts(ctx, tx, "zone", s.Floors)
			if err != nil {
				return fmt.Errorf("zones %w", err)
			}
			md.ZoneCounts = zones
			// severity
			severity, err := queryGroupCounts(ctx, tx, "severity", s.Severity)
			if err != nil {
				return fmt.Errorf("severity %w", err)
			}
			md.SeverityCounts = severityCountMapToProto(severity)

			if eventsCancel != nil {
				eventsCancel()
			}
			eventsCtx, eventsCancel = context.WithCancel(ctx)
			events = s.bus.Listen(eventsCtx)

			return nil
		})
		if err != nil {
			s.mdErr = err
			s.mdC = nil // allow another call to attempt again
			return
		}
		s.md = resource.NewValue(resource.WithInitialValue(md))

		// setup listeners for changes to the DB so we can track those changes in metadata
		go func() {
			defer eventsCancel()
			for {
				select {
				case <-ctx.Done():
					return // the server is closing
				case e := <-events:
					applyMdDelta(s.md, e)
				}
			}
		}()
	}()
	return wait()
}

func applyMdDelta(md *resource.Value, e *gen.PullAlertsResponse_Change) error {
	// note: AlertMetadata uses uint32 so we can't use it to store our deltas
	_, err := md.Set(&gen.AlertMetadata{}, resource.InterceptBefore(func(old, new proto.Message) {
		oldMd, newMd := old.(*gen.AlertMetadata), new.(*gen.AlertMetadata)
		proto.Merge(newMd, oldMd)
		// proto.Merge doesn't set maps that are empty but non-nil!
		// Without this explicit copy we'd end up with nil maps that we then assign to, which panics
		if newMd.FloorCounts == nil {
			newMd.FloorCounts = make(map[string]uint32)
		}
		if newMd.ZoneCounts == nil {
			newMd.ZoneCounts = make(map[string]uint32)
		}
		if newMd.SeverityCounts == nil {
			newMd.SeverityCounts = make(map[int32]uint32)
		}
		if newMd.AcknowledgedCounts == nil {
			newMd.AcknowledgedCounts = make(map[bool]uint32)
		}
		if newMd.ResolvedCounts == nil {
			newMd.ResolvedCounts = make(map[bool]uint32)
		}
		if newMd.NeedsAttentionCounts == nil {
			newMd.NeedsAttentionCounts = make(map[string]uint32)
		}

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
		boolMapDelta(e.OldValue, e.NewValue, newMd.AcknowledgedCounts, func(v *gen.Alert) bool {
			return v.GetAcknowledgement().GetAcknowledgeTime() != nil
		})
		// resolved/unresolved
		boolMapDelta(e.OldValue, e.NewValue, newMd.ResolvedCounts, func(v *gen.Alert) bool {
			return v.GetResolveTime() != nil
		})

		// floors, zones, and severity
		mapDelta(e.GetOldValue().GetFloor(), e.GetNewValue().GetFloor(), newMd.FloorCounts)
		mapDelta(e.GetOldValue().GetZone(), e.GetNewValue().GetZone(), newMd.ZoneCounts)
		mapDelta(int32(e.GetOldValue().GetSeverity()), int32(e.GetNewValue().GetSeverity()), newMd.SeverityCounts)

		// needs attention
		needsAttentionMap(e.GetOldValue(), e.GetNewValue(), newMd.NeedsAttentionCounts)
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

func boolMapDelta(o, n *gen.Alert, m map[bool]uint32, f func(*gen.Alert) bool) {
	if o == nil && n != nil {
		if f(n) {
			m[true]++
		} else {
			m[false]++
		}
	} else if o != nil && n == nil {
		if f(o) {
			mapSub(true, m)
		} else {
			mapSub(false, m)
		}
	} else {
		oOK, nOK := f(o), f(n)
		if oOK && !nOK {
			m[true]++
			mapSub(false, m)
		} else if !oOK && nOK {
			mapSub(true, m)
			m[false]++
		}
	}
}

func mapSub[K comparable](k K, m map[K]uint32) {
	if m[k] > 0 {
		m[k]--
	}
}

func mapCmpBool[K comparable](k K, o, n bool, m map[K]uint32) {
	switch {
	case !o && n:
		m[k]++
	case o && !n:
		mapSub(k, m)
	}
}

func needsAttentionMap(o, n *gen.Alert, m map[string]uint32) {
	switch {
	case o == nil && n != nil:
		ackResolved, ackUnresolved, nackResolved, nackUnresolved := needsAttentionFlags(n)
		if ackResolved {
			m["ack_resolved"]++
		}
		if ackUnresolved {
			m["ack_unresolved"]++
		}
		if nackResolved {
			m["nack_resolved"]++
		}
		if nackUnresolved {
			m["nack_unresolved"]++
		}
	case o != nil && n == nil:
		ackResolved, ackUnresolved, nackResolved, nackUnresolved := needsAttentionFlags(o)
		if ackResolved {
			mapSub("ack_resolved", m)
		}
		if ackUnresolved {
			mapSub("ack_unresolved", m)
		}
		if nackResolved {
			mapSub("nack_resolved", m)
		}
		if nackUnresolved {
			mapSub("nack_unresolved", m)
		}
	case o != nil && n != nil:
		oAckResolved, oAckUnresolved, oNackResolved, oNackUnresolved := needsAttentionFlags(o)
		nAckResolved, nAckUnresolved, nNackResolved, nNackUnresolved := needsAttentionFlags(n)
		mapCmpBool("ack_resolved", oAckResolved, nAckResolved, m)
		mapCmpBool("ack_unresolved", oAckUnresolved, nAckUnresolved, m)
		mapCmpBool("nack_resolved", oNackResolved, nNackResolved, m)
		mapCmpBool("nack_unresolved", oNackUnresolved, nNackUnresolved, m)
	}
}

func needsAttentionFlags(a *gen.Alert) (ackResolved, ackUnresolved, nackResolved, nackUnresolved bool) {
	if a == nil {
		return
	}
	ack := a.GetAcknowledgement().GetAcknowledgeTime() != nil
	resolved := a.GetResolveTime() != nil
	ackResolved = ack && resolved
	ackUnresolved = ack && !resolved
	nackResolved = !ack && resolved
	nackUnresolved = !ack && !resolved
	return
}

func queryGroupCounts[K comparable](ctx context.Context, tx pgx.Tx, col string, seed []K) (map[K]uint32, error) {
	all := make(map[K]uint32)
	for _, floor := range seed {
		all[floor] = 0
	}
	query, err := tx.Query(ctx, fmt.Sprintf(`SELECT %s, COUNT(id) FROM alerts GROUP BY %s`, col, col))
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}
	for query.Next() {
		var id *K
		var count uint32
		err = query.Scan(&id, &count)
		if err != nil {
			return nil, fmt.Errorf("row scan %w", err)
		}
		if id == nil {
			var zero K
			id = &zero
		}
		all[*id] = count
	}
	return all, nil
}

func severityCountMapToProto(m map[gen.Alert_Severity]uint32) map[int32]uint32 {
	dst := make(map[int32]uint32, len(m))
	for k, v := range m {
		dst[int32(k)] = v
	}
	return dst
}
