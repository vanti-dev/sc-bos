package main

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/vanti-dev/bsp-ew/internal/db"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	gen.UnimplementedNodeApiServer

	logger      *zap.Logger
	db          *pgx.Conn
	ca          *enrollment.CA
	managerName string
	managerAddr string
	rootsPEM    []byte
}

func (n *NodeServer) CreateNodeRegistration(ctx context.Context, request *gen.CreateNodeRegistrationRequest) (*gen.NodeRegistration, error) {
	nodeReg := request.GetNodeRegistration()
	if nodeReg == nil {
		return nil, status.Error(codes.InvalidArgument, "node_registration must be supplied")
	}

	en, err := enrollment.EnrollAreaController(ctx, &gen.Enrollment{
		TargetName:     nodeReg.Name,
		TargetAddress:  nodeReg.Address,
		ManagerName:    n.managerAddr,
		ManagerAddress: n.managerName,
		RootCas:        n.rootsPEM,
	}, n.ca)
	if err != nil {
		n.logger.Error("failed to enroll area controller", zap.Error(err),
			zap.String("target_address", nodeReg.Address))
		return nil, status.Error(codes.Unknown, "target refused registration")
	}

	err = n.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return db.AddEnrollment(ctx, tx, db.Enrollment{
			Name:        en.TargetName,
			Description: nodeReg.Description,
			Address:     en.TargetAddress,
			Cert:        en.Certificate,
		})
	})
	if err != nil {
		n.logger.Error("db.AddEnrollment failed", zap.Error(err))
		return nil, status.Error(codes.DataLoss, "failed to save the enrollment - manual intervention required")
	}
	return nodeReg, nil
}

func (n *NodeServer) ListNodeRegistrations(ctx context.Context, request *gen.ListNodeRegistrationsRequest) (*gen.ListNodeRegistrationsResponse, error) {
	var dbEnrollments []db.Enrollment
	err := n.db.BeginFunc(ctx, func(tx pgx.Tx) (err error) {
		dbEnrollments, err = db.ListEnrollments(ctx, tx)
		return
	})
	if err != nil {
		n.logger.Error("db.ListEnrollments failed", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "unable to retrieve enrollments")
	}

	var registrations []*gen.NodeRegistration
	for _, en := range dbEnrollments {
		registrations = append(registrations, &gen.NodeRegistration{
			Name:        en.Name,
			Address:     en.Address,
			Description: en.Description,
		})
	}

	return &gen.ListNodeRegistrationsResponse{NodeRegistrations: registrations}, nil
}
