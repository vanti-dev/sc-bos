package main

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/smart-core-os/sc-api/go/traits"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/db"
	"github.com/vanti-dev/sc-bos/pkg/util/rpcutil"
)

type PublicationServer struct {
	traits.UnimplementedPublicationApiServer

	logger *zap.Logger
	conn   *pgxpool.Pool
}

func (p *PublicationServer) CreatePublication(
	ctx context.Context, request *traits.CreatePublicationRequest,
) (*traits.Publication, error) {
	logger := rpcutil.ServerLogger(ctx, p.logger)
	input := request.GetPublication()

	pubID := input.GetId()
	if pubID == "" {
		return nil, status.Error(codes.InvalidArgument, "publication ID must be present")
	}
	audience := input.GetAudience().GetName()
	body := input.GetBody()
	mediaType := input.GetMediaType()

	var output *traits.Publication
	err := p.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		// Register the publication
		err := db.CreatePublication(ctx, tx, pubID, audience)
		if err != nil {
			return err
		}

		pubTime := time.Now().UTC()
		// Create the initial version for the publication
		version, err := db.CreatePublicationVersion(ctx, tx, db.PublicationVersion{
			PublicationID: pubID,
			PublishTime:   pubTime,
			Body:          body,
			MediaType:     mediaType,
		})
		if err != nil {
			return err
		}

		output = &traits.Publication{
			Id:          pubID,
			Version:     version,
			Body:        body,
			PublishTime: timestamppb.New(pubTime),
			MediaType:   mediaType,
		}
		if audience != "" {
			output.Audience = &traits.Publication_Audience{
				Name: audience,
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("failed to create publication", zap.Error(err), zap.String("id", input.GetId()))
		return nil, status.Error(codes.Internal, "database error")
	}

	return output, nil
}

func (p *PublicationServer) GetPublication(
	ctx context.Context, request *traits.GetPublicationRequest,
) (*traits.Publication, error) {
	id := request.GetId()
	version := request.GetVersion()

	var output *traits.Publication
	err := p.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		output, err = db.GetPublication(ctx, tx, id, version)
		return err
	})
	return output, err
}

func (p *PublicationServer) UpdatePublication(
	ctx context.Context, request *traits.UpdatePublicationRequest,
) (*traits.Publication, error) {
	if request.GetUpdateMask() != nil {
		return nil, status.Error(codes.Unimplemented, "field mask support not implemented")
	}

	data := db.PublicationVersion{
		PublicationID: request.GetPublication().GetId(),
		PublishTime:   time.Now(),
		Body:          request.GetPublication().GetBody(),
		MediaType:     request.GetPublication().GetMediaType(),
	}
	var updated *traits.Publication

	err := p.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		if request.GetVersion() != "" {
			err = checkLatestVersion(ctx, tx, request.GetPublication().GetId(), request.GetVersion())
			if err != nil {
				return err
			}
		}

		versionID, err := db.CreatePublicationVersion(ctx, tx, data)
		if err != nil {
			return err
		}

		updated, err = db.GetPublication(ctx, tx, data.PublicationID, versionID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (p *PublicationServer) DeletePublication(
	ctx context.Context, request *traits.DeletePublicationRequest,
) (*traits.Publication, error) {
	var pub *traits.Publication

	err := p.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		// if a version is specified, check that it's the latest version
		// this helps clients to avoid racing each other
		if request.GetVersion() != "" {
			err = checkLatestVersion(ctx, tx, request.GetId(), request.GetVersion())
			if err != nil {
				return err
			}
		}

		// grab a copy of the data we're about to delete
		pub, err = db.GetPublication(ctx, tx, request.GetId(), request.GetVersion())
		if errors.Is(err, pgx.ErrNoRows) && request.GetAllowMissing() {
			// the publication doesn't exist, and the client indicates that this is acceptable
			// return an empty publication to indicate this
			pub = &traits.Publication{}
			return nil
		}
		if err != nil {
			return err
		}

		return db.DeletePublication(ctx, tx, request.GetId(), request.GetAllowMissing())
	})

	if err != nil {
		return nil, err
	}
	return pub, nil
}

func (p *PublicationServer) PullPublication(
	_ *traits.PullPublicationRequest, _ traits.PublicationApi_PullPublicationServer,
) error {
	return status.Error(codes.Unimplemented, "PullPublication not implemented")
}

func (p *PublicationServer) ListPublications(
	ctx context.Context, request *traits.ListPublicationsRequest,
) (*traits.ListPublicationsResponse, error) {
	limit := 50
	if request.GetPageSize() > 0 {
		limit = int(request.GetPageSize())
	}

	var (
		publications []*traits.Publication
		nextToken    string
	)
	err := p.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		publications, nextToken, err = db.GetPublicationsPaginated(ctx, tx, request.GetPageToken(), limit)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &traits.ListPublicationsResponse{
		Publications:  publications,
		NextPageToken: nextToken,
	}, nil
}

func (p *PublicationServer) PullPublications(
	_ *traits.PullPublicationsRequest, _ traits.PublicationApi_PullPublicationsServer,
) error {
	return status.Error(codes.Unimplemented, "PullPublications not implemented")
}

func (p *PublicationServer) AcknowledgePublication(
	ctx context.Context, request *traits.AcknowledgePublicationRequest,
) (*traits.Publication, error) {
	var accepted bool
	switch request.GetReceipt() {
	case traits.Publication_Audience_REJECTED:
		accepted = false
	case traits.Publication_Audience_ACCEPTED:
		accepted = true
		if request.GetReceiptRejectedReason() != "" {
			return nil, status.Error(codes.InvalidArgument, "cannot specify rejected_reason when ACCEPTED")
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "REJECTED or ACCEPTED must be specified")
	}

	var updated *traits.Publication
	err := p.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := db.CreatePublicationAcknowledgement(ctx, tx, request.GetId(), request.GetVersion(), time.Now(), accepted,
			request.GetReceiptRejectedReason(), request.GetAllowAcknowledged())
		if err != nil {
			return err
		}

		// get the updated publication
		updated, err = db.GetPublication(ctx, tx, request.GetId(), request.GetVersion())
		return err
	})

	return updated, err
}

func checkLatestVersion(ctx context.Context, tx pgx.Tx, pubID string, expectedVersionID string) error {
	latestVersion, err := db.GetLatestVersionID(ctx, tx, pubID)
	if errors.Is(err, pgx.ErrNoRows) {
		// publication has no versions
		return ErrExpectedVersion
	}
	if err != nil {
		// some other database error
		return err
	}
	if latestVersion != expectedVersionID {
		// publication latest version is not what we expected - concurrent updates detected
		return ErrExpectedVersion
	}

	return nil
}

var ErrExpectedVersion = status.Error(codes.FailedPrecondition, "expected latest version does not match")
