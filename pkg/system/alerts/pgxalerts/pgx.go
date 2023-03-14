package pgxalerts

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

//go:embed schema.sql
var schemaSql string

func SetupDB(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, schemaSql)
		return err
	})
}

func NewServer(ctx context.Context, connStr string) (*Server, error) {
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("connect %w", err)
	}

	return NewServerFromPool(ctx, pool)
}

func NewServerFromPool(ctx context.Context, pool *pgxpool.Pool) (*Server, error) {
	err := SetupDB(ctx, pool)
	if err != nil {
		return nil, fmt.Errorf("setup %w", err)
	}

	return &Server{
		pool: pool,
		bus:  &minibus.Bus[*gen.PullAlertsResponse_Change]{},
		Severity: []gen.Alert_Severity{
			gen.Alert_INFO,
			gen.Alert_WARNING,
			gen.Alert_SEVERE,
			gen.Alert_LIFE_SAFETY,
		},
	}, nil
}

type Server struct {
	gen.UnimplementedAlertApiServer
	gen.UnimplementedAlertAdminApiServer

	// Floors, if set, is used to pre-populate AlertMetadata with zero values for cases when no alerts appear on a floor.
	Floors []string
	// Zones, if set, is used to pre-populate AlertMetadata with zero values for cases when no alerts appear in a zone.
	Zones []string
	// Severity, if set, is used to pre-populate AlertMetadata with zero values for cases when no alerts have a severity.
	Severity []gen.Alert_Severity

	pool *pgxpool.Pool
	bus  *minibus.Bus[*gen.PullAlertsResponse_Change]

	// to support alert metadata
	mdMu   sync.Mutex      // guards the following md fields
	md     *resource.Value // of *gen.AlertMetadata, used to track changes
	mdC    chan struct{}   // nil if needs init, blocked if init-ing, closed if done
	mdErr  error           // non-nil if mdC is done and completed with error
	mdStop func()          // closes any go routines that are listening for changes
}

func (s *Server) Close() error {
	s.mdMu.Lock()
	if s.mdStop != nil {
		s.mdStop()
	}
	s.mdMu.Unlock()

	return nil
}

func (s *Server) CreateAlert(ctx context.Context, request *gen.CreateAlertRequest) (*gen.Alert, error) {
	alert := request.Alert
	if alert.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description empty")
	}
	if alert.Severity == 0 {
		alert.Severity = gen.Alert_WARNING
	}

	var createTime time.Time
	err := s.pool.QueryRow(ctx,
		`INSERT INTO alerts (description, severity, floor, zone, source) VALUES ($1, $2, $3, $4, $5) RETURNING id, create_time`,
		alert.Description, alert.Severity, alert.Floor, alert.Zone, alert.Source,
	).Scan(&alert.Id, &createTime)
	if err != nil {
		return nil, dbErrToStatus(err)
	}

	alert.CreateTime = timestamppb.New(createTime)

	// notify
	go s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
		Name:       request.Name,
		Type:       types.ChangeType_ADD,
		ChangeTime: alert.CreateTime,
		NewValue:   alert,
	})

	return alert, nil
}

func (s *Server) UpdateAlert(ctx context.Context, request *gen.UpdateAlertRequest) (*gen.Alert, error) {
	alert := request.Alert
	if alert.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}

	var fields []string
	var values []any
	if shouldUpdateField(request.UpdateMask, "description", alert.Description) {
		fields = append(fields, "description")
		values = append(values, alert.Description)
	}
	if shouldUpdateField(request.UpdateMask, "severity", alert.Severity) {
		fields = append(fields, "severity")
		values = append(values, alert.Severity)
	}
	if shouldUpdateField(request.UpdateMask, "floor", alert.Floor) {
		fields = append(fields, "floor")
		values = append(values, alert.Floor)
	}
	if shouldUpdateField(request.UpdateMask, "zone", alert.Zone) {
		fields = append(fields, "zone")
		values = append(values, alert.Zone)
	}
	if shouldUpdateField(request.UpdateMask, "source", alert.Source) {
		fields = append(fields, "source")
		values = append(values, alert.Source)
	}

	if len(fields) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	original := &gen.Alert{}
	updated := &gen.Alert{}

	setStr := ""
	for i, field := range fields {
		if setStr != "" {
			setStr += ","
		}
		setStr += fmt.Sprintf("%s=$%d", field, i+2) // +2 to make it 1-based leaving room for $1 to be the id
	}
	args := append([]any{alert.Id}, values...)
	sql := fmt.Sprintf(`UPDATE alerts SET %s WHERE id=$1`, setStr)
	err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// record the original value which will be used for event notifications
		if err := readAlertById(ctx, tx, alert.Id, original); err != nil {
			return err
		}
		res, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		rows := res.RowsAffected()
		if rows == 0 {
			return status.Error(codes.NotFound, "id not found")
		}

		// get alert so we can return it
		return readAlertById(ctx, tx, alert.Id, updated)
	})

	if err != nil {
		return nil, dbErrToStatus(err)
	}

	// notify
	go s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
		Name:       request.Name,
		Type:       types.ChangeType_UPDATE,
		ChangeTime: timestamppb.Now(),
		OldValue:   original,
		NewValue:   updated,
	})

	return updated, nil
}

func (s *Server) DeleteAlert(ctx context.Context, request *gen.DeleteAlertRequest) (*gen.DeleteAlertResponse, error) {
	if request.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}
	existing := &gen.Alert{}
	err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// We do extra work to get the old value so we can include it in bus events.
		// Without this any filtered PullAlerts call wouldn't be able to correctly include the event in responses.
		err := readAlertById(ctx, tx, request.Id, existing)
		if err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `DELETE FROM alerts WHERE id=$1`, request.Id)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			if request.AllowMissing {
				return nil
			}
			return status.Error(codes.NotFound, request.Id)
		}
		return nil
	})

	if err != nil {
		err := dbErrToStatus(err)
		if status.Code(err) == codes.NotFound && request.AllowMissing {
			return &gen.DeleteAlertResponse{}, nil
		}
		return nil, err
	}

	// notify
	go s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
		Name:       request.Name,
		Type:       types.ChangeType_REMOVE,
		ChangeTime: timestamppb.Now(),
		OldValue:   existing,
	})

	return &gen.DeleteAlertResponse{}, nil
}

func (s *Server) ListAlerts(ctx context.Context, request *gen.ListAlertsRequest) (*gen.ListAlertsResponse, error) {
	var args []any
	var argIdx int
	var where []string // combined with AND
	if request.PageToken != "" {
		pt, err := DecodePageToken(request.PageToken)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "bad page token")
		}
		where = append(where, fmt.Sprintf(`create_time>$%d OR (create_time=$%d AND id>$%d)`, argIdx+1, argIdx+1, argIdx+2))
		args = append(args, pt.LastCreateTime, pt.LastID)
		argIdx += 2
	}
	if request.Query != nil {
		q := request.Query

		if q.CreatedNotBefore != nil {
			where = append(where, fmt.Sprintf(`create_time>=$%d`, argIdx+1))
			args = append(args, q.CreatedNotBefore.AsTime())
			argIdx += 1
		}
		if q.CreatedNotAfter != nil {
			where = append(where, fmt.Sprintf(`create_time<=$%d`, argIdx+1))
			args = append(args, q.CreatedNotAfter.AsTime())
			argIdx += 1
		}
		if q.SeverityNotBelow != 0 {
			where = append(where, fmt.Sprintf(`severity>=$%d`, argIdx+1))
			args = append(args, q.SeverityNotBelow)
			argIdx += 1
		}
		if q.SeverityNotAbove != 0 {
			where = append(where, fmt.Sprintf(`severity<=$%d`, argIdx+1))
			args = append(args, q.SeverityNotAbove)
			argIdx += 1
		}
		if q.Floor != "" {
			where = append(where, fmt.Sprintf(`floor=$%d`, argIdx+1))
			args = append(args, q.Floor)
			argIdx += 1
		}
		if q.Zone != "" {
			where = append(where, fmt.Sprintf(`zone=$%d`, argIdx+1))
			args = append(args, q.Zone)
			argIdx += 1
		}
		if q.Source != "" {
			where = append(where, fmt.Sprintf(`source=$%d`, argIdx+1))
			args = append(args, q.Source)
			argIdx += 1
		}
		if q.Acknowledged != nil {
			if *q.Acknowledged {
				where = append(where, `ack_time IS NOT NULL`)
			} else {
				where = append(where, `ack_time IS NULL`)
			}
		}
	}

	sql := selectAlertSQL
	if len(where) > 0 {
		sql += " WHERE (" + strings.Join(where, ") AND (") + ")"
	}
	sql += " ORDER BY create_time DESC, id DESC"
	// +1 so we read past the end of the page to know if theres another page or not
	pageSize := normalizePageSize(request.PageSize)
	sql += fmt.Sprintf(" LIMIT %d", pageSize+1)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, dbErrToStatus(err)
	}
	defer rows.Close()

	var alerts []*gen.Alert
	for rows.Next() {
		alert := &gen.Alert{}
		err := scanAlert(rows, alert)
		if err != nil {
			return nil, dbErrToStatus(err)
		}
		alerts = append(alerts, alert)
	}

	res := &gen.ListAlertsResponse{}
	if len(alerts) <= pageSize {
		// no more pages, remember we selected pageSize+1 rows so we could be sure.
		res.Alerts = alerts
		return res, nil
	}

	// else there are more pages
	// calc the next page token
	res.Alerts = alerts[:pageSize]
	lastRecord := alerts[len(alerts)-1]
	npt := PageToken{
		LastCreateTime: lastRecord.CreateTime.AsTime(),
		LastID:         lastRecord.Id,
	}
	nptStr, err := npt.Encode()
	if err != nil {
		return nil, err
	}
	res.NextPageToken = nptStr

	filter := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))
	for _, alert := range res.Alerts {
		filter.Filter(alert)
	}
	return res, nil
}

func (s *Server) PullAlerts(request *gen.PullAlertsRequest, server gen.AlertApi_PullAlertsServer) error {
	filter := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))
	for change := range s.bus.Listen(server.Context()) {
		change := convertChangeForQuery(request.Query, change)
		if change == nil {
			continue
		}
		if change.OldValue != nil {
			change.OldValue = filter.FilterClone(change.OldValue).(*gen.Alert)
		}
		if change.NewValue != nil {
			change.NewValue = filter.FilterClone(change.NewValue).(*gen.Alert)
		}
		err := server.Send(&gen.PullAlertsResponse{Changes: []*gen.PullAlertsResponse_Change{change}})
		if err != nil {
			return err
		}
	}
	return server.Context().Err()
}

func (s *Server) AcknowledgeAlert(ctx context.Context, request *gen.AcknowledgeAlertRequest) (*gen.Alert, error) {
	if request.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}

	existing := &gen.Alert{}
	updated := &gen.Alert{}
	err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := readAlertById(ctx, tx, request.Id, existing)
		if err != nil {
			return err
		}

		if existing.Acknowledgement != nil {
			if request.AllowAcknowledged {
				updated = existing
				return nil
			}
			return status.Error(codes.Aborted, "already acknowledged")
		}

		var ackAuthorId, ackAuthorName, ackAuthorEmail *string
		if request.Author != nil {
			a := request.Author
			if a.Id != "" {
				ackAuthorId = &a.Id
			}
			if a.DisplayName != "" {
				ackAuthorName = &a.DisplayName
			}
			if a.Email != "" {
				ackAuthorEmail = &a.Email
			}
		}
		_, err = tx.Exec(ctx, `UPDATE alerts SET ack_time=now(), ack_author_id=$2, ack_author_name=$3, ack_author_email=$4 WHERE id=$1`,
			request.Id, ackAuthorId, ackAuthorName, ackAuthorEmail)
		if err != nil {
			return err
		}

		return readAlertById(ctx, tx, request.Id, updated)
	})

	if err != nil {
		err := dbErrToStatus(err)
		if status.Code(err) == codes.NotFound && request.AllowMissing {
			return &gen.Alert{}, nil
		}
		return nil, err
	}

	if !proto.Equal(existing, updated) {
		// notify
		go s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
			Name:       request.Name,
			Type:       types.ChangeType_UPDATE,
			ChangeTime: timestamppb.Now(),
			OldValue:   existing,
			NewValue:   updated,
		})
	}

	return updated, nil
}

func (s *Server) UnacknowledgeAlert(ctx context.Context, request *gen.AcknowledgeAlertRequest) (*gen.Alert, error) {
	if request.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}

	existing := &gen.Alert{}
	updated := &gen.Alert{}
	err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := readAlertById(ctx, tx, request.Id, existing)
		if err != nil {
			return err
		}

		if existing.Acknowledgement == nil {
			if request.AllowAcknowledged {
				updated = existing
				return nil
			}
			return status.Error(codes.Aborted, "already not acknowledged")
		}

		_, err = tx.Exec(ctx, `UPDATE alerts SET ack_time=null, ack_author_id=null, ack_author_name=null, ack_author_email=null WHERE id=$1`, request.Id)
		if err != nil {
			return err
		}

		return readAlertById(ctx, tx, request.Id, updated)
	})

	if err != nil {
		err := dbErrToStatus(err)
		if status.Code(err) == codes.NotFound && request.AllowMissing {
			return &gen.Alert{}, nil
		}
		return nil, err
	}

	if !proto.Equal(existing, updated) {
		// notify
		go s.bus.Send(context.Background(), &gen.PullAlertsResponse_Change{
			Name:       request.Name,
			Type:       types.ChangeType_UPDATE,
			ChangeTime: timestamppb.Now(),
			OldValue:   existing,
			NewValue:   updated,
		})
	}

	return updated, nil
}

func readAlertById(ctx context.Context, tx pgx.Tx, id string, dst *gen.Alert) error {
	row := tx.QueryRow(ctx,
		selectAlertSQL+` WHERE id=$1`,
		id,
	)
	return scanAlert(row, dst)
}

// selectAlertSQL selects fields in the order expected by scanAlert.
const selectAlertSQL = `SELECT id, description, severity, create_time, floor, zone, source, ack_time, ack_author_id, ack_author_name, ack_author_email FROM alerts`

func scanAlert(scanner pgx.Row, dst *gen.Alert) error {
	var createTime, ackTime *time.Time
	var ackAuthorId, ackAuthorName, ackAuthorEmail *string
	err := scanner.Scan(&dst.Id, &dst.Description, &dst.Severity, &createTime, &dst.Floor, &dst.Zone, &dst.Source, &ackTime, &ackAuthorId, &ackAuthorName, &ackAuthorEmail)
	if err != nil {
		return err
	}
	if createTime != nil {
		dst.CreateTime = timestamppb.New(*createTime)
	}
	if ackTime != nil {
		dst.Acknowledgement = &gen.Alert_Acknowledgement{
			AcknowledgeTime: timestamppb.New(*ackTime),
		}

		// ack author details, we assume there is an author only if the author has any information
		var hasAuthor bool
		author := &gen.Alert_Acknowledgement_Author{}
		if ackAuthorId != nil {
			hasAuthor = true
			author.Id = *ackAuthorId
		}
		if ackAuthorName != nil {
			hasAuthor = true
			author.DisplayName = *ackAuthorName
		}
		if ackAuthorEmail != nil {
			hasAuthor = true
			author.Email = *ackAuthorEmail
		}
		if hasAuthor {
			dst.Acknowledgement.Author = author
		}
	}
	return nil
}

func dbErrToStatus(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, "alert")
	}
	if pgerr, ok := err.(*pgconn.PgError); ok {
		if strings.HasPrefix(pgerr.Code, "22") { // data error
			// If we're getting data errors involving uuids then someone is giving us an id string that doesn't match
			// the uuid format, we don't care that it doesn't match, but we should return NotFound in this case.
			if strings.Contains(pgerr.Message, "uuid") {
				return status.Error(codes.NotFound, "id")
			}
			return status.Error(codes.InvalidArgument, pgerr.Error())
		}
	}
	// todo: convert err to a more meaningful status code
	return err
}

// shouldUpdateField returns true if the given field value should be written to the db given the field mask.
// If m is nil then only non-zero values should be written, otherwise fieldMaskIncludesPath is consulted.
func shouldUpdateField[T comparable](m *fieldmaskpb.FieldMask, f string, v T) bool {
	if m == nil {
		var zero T
		return v != zero
	}
	return fieldMaskIncludesPath(m, f)
}

// fieldMaskIncludesPath returns true if the given FieldMask indicates that p should be included.
// This happens if m is nil, m.Paths includes p, or p is a more specific path for one in m.
func fieldMaskIncludesPath(m *fieldmaskpb.FieldMask, p string) bool {
	if m == nil {
		return true
	}
	i := fieldmaskpb.Intersect(m, &fieldmaskpb.FieldMask{Paths: []string{p}})
	return len(i.Paths) > 0
}

func convertChangeForQuery(q *gen.Alert_Query, change *gen.PullAlertsResponse_Change) *gen.PullAlertsResponse_Change {
	if q == nil {
		return change
	}

	res := proto.Clone(change).(*gen.PullAlertsResponse_Change)
	if change.OldValue != nil && !alertMatchesQuery(q, change.OldValue) {
		res.OldValue = nil
	}
	if change.NewValue != nil && !alertMatchesQuery(q, change.NewValue) {
		res.NewValue = nil
	}
	if res.OldValue == nil && res.NewValue == nil {
		return nil // change should be ignored
	}
	if res.NewValue == nil {
		// delete
		res.Type = types.ChangeType_REMOVE
		return res
	}
	if res.OldValue == nil {
		// create
		res.Type = types.ChangeType_ADD
		return res
	}
	// else update
	res.Type = types.ChangeType_UPDATE
	return res
}

func alertMatchesQuery(q *gen.Alert_Query, a *gen.Alert) bool {
	if q == nil {
		return true
	}
	if a == nil {
		return false
	}

	if q.CreatedNotBefore != nil {
		if a.CreateTime == nil {
			return false
		}
		if a.CreateTime.AsTime().Before(q.CreatedNotBefore.AsTime()) {
			return false
		}
	}
	if q.CreatedNotAfter != nil {
		if a.CreateTime == nil {
			return false
		}
		if a.CreateTime.AsTime().After(q.CreatedNotAfter.AsTime()) {
			return false
		}
	}
	if q.SeverityNotBelow != 0 && int32(a.Severity) < q.SeverityNotBelow {
		return false
	}
	if q.SeverityNotAbove != 0 && int32(a.Severity) > q.SeverityNotAbove {
		return false
	}
	if q.Floor != "" && q.Floor != a.Floor {
		return false
	}
	if q.Zone != "" && q.Zone != a.Zone {
		return false
	}
	if q.Source != "" && q.Source != a.Source {
		return false
	}

	if q.Acknowledged != nil {
		wantAck := *q.Acknowledged
		if wantAck && a.Acknowledgement == nil {
			return false
		}
		if !wantAck && a.Acknowledgement != nil {
			return false
		}
	}

	return true
}
