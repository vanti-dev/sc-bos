package pgxalerts

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// selectAlertSQL selects fields in the order expected by scanAlert.
const selectAlertSQL = `SELECT id, description, severity, create_time, resolve_time, floor, zone, source, federation, ack_time, ack_author_id, ack_author_name, ack_author_email FROM alerts`

type QueryRower interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func insertAlert(ctx context.Context, q QueryRower, alert *gen.Alert) (*gen.Alert, error) {
	var createTime time.Time
	err := q.QueryRow(ctx,
		`INSERT INTO alerts (description, severity, floor, zone, source, federation) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, create_time`,
		alert.Description, alert.Severity, alert.Floor, alert.Zone, alert.Source, alert.Federation,
	).Scan(&alert.Id, &createTime)
	if err != nil {
		return nil, err
	}

	alert.CreateTime = timestamppb.New(createTime)
	return alert, nil
}

func readAlertById(ctx context.Context, tx QueryRower, id string, dst *gen.Alert) error {
	row := tx.QueryRow(ctx,
		selectAlertSQL+` WHERE id=$1`,
		id,
	)
	return scanAlert(row, dst)
}

func scanAlert(scanner pgx.Row, dst *gen.Alert) error {
	var createTime, resolveTime, ackTime *time.Time
	var floor, zone, source, federation *string
	var ackAuthorId, ackAuthorName, ackAuthorEmail *string
	err := scanner.Scan(&dst.Id, &dst.Description, &dst.Severity, &createTime, &resolveTime, &floor, &zone, &source, &federation, &ackTime, &ackAuthorId, &ackAuthorName, &ackAuthorEmail)
	if err != nil {
		return err
	}
	if floor != nil {
		dst.Floor = *floor
	}
	if zone != nil {
		dst.Zone = *zone
	}
	if source != nil {
		dst.Source = *source
	}
	if federation != nil {
		dst.Federation = *federation
	}
	if createTime != nil {
		dst.CreateTime = timestamppb.New(*createTime)
	}
	if resolveTime != nil {
		dst.ResolveTime = timestamppb.New(*resolveTime)
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
