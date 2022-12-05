package pgxalerts

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

//go:embed schema.sql
var schemaSql string

// selectAlertSQL selects fields in the order expected by scanAlert.
const selectAlertSQL = `SELECT id, description, severity, create_time, ack_time, floor, zone, source FROM alerts`

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

	return &Server{pool: pool}, nil
}

type Server struct {
	gen.UnimplementedAlertApiServer
	gen.UnimplementedAlertAdminApiServer

	pool *pgxpool.Pool
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
		res, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return dbErrToStatus(err)
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

	return updated, nil
}

func (s *Server) DeleteAlert(ctx context.Context, request *gen.DeleteAlertRequest) (*gen.DeleteAlertResponse, error) {
	if request.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM alerts WHERE id=$1`, request.Id)
	if err != nil {
		err := dbErrToStatus(err)
		if status.Code(err) == codes.NotFound && request.AllowMissing {
			return &gen.DeleteAlertResponse{}, nil
		}
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		if request.AllowMissing {
			return &gen.DeleteAlertResponse{}, nil
		}
		return nil, status.Error(codes.NotFound, request.Id)
	}
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
	// todo: PullAlerts
	return s.UnimplementedAlertApiServer.PullAlerts(request, server)
}

func (s *Server) AcknowledgeAlert(ctx context.Context, request *gen.AcknowledgeAlertRequest) (*gen.Alert, error) {
	if request.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}

	updated := &gen.Alert{}
	err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		existing := &gen.Alert{}
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

		_, err = tx.Exec(ctx, `UPDATE alerts SET ack_time=now() WHERE id=$1`, request.Id)
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

	return updated, nil
}

func (s *Server) UnacknowledgeAlert(ctx context.Context, request *gen.AcknowledgeAlertRequest) (*gen.Alert, error) {
	if request.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id empty")
	}

	updated := &gen.Alert{}
	err := s.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		existing := &gen.Alert{}
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

		_, err = tx.Exec(ctx, `UPDATE alerts SET ack_time=null WHERE id=$1`, request.Id)
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

	return updated, nil
}

func readAlertById(ctx context.Context, tx pgx.Tx, id string, dst *gen.Alert) error {
	row := tx.QueryRow(ctx,
		selectAlertSQL+` WHERE id=$1`,
		id,
	)
	return scanAlert(row, dst)
}

func scanAlert(scanner pgx.Row, dst *gen.Alert) error {
	var createTime, ackTime *time.Time
	err := scanner.Scan(&dst.Id, &dst.Description, &dst.Severity, &createTime, &ackTime, &dst.Floor, &dst.Zone, &dst.Source)
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
