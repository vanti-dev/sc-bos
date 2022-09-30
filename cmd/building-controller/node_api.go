package main

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vanti-dev/bsp-ew/internal/db"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/util/rpcutil"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	gen.UnimplementedNodeApiServer

	logger        *zap.Logger
	db            *pgxpool.Pool
	managerName   string
	managerAddr   string
	authority     pki.Source // trust authority for the cohort of smart core nodes
	testTLSConfig *tls.Config
}

func (n *NodeServer) GetNodeRegistration(ctx context.Context, request *gen.GetNodeRegistrationRequest) (*gen.NodeRegistration, error) {
	logger := rpcutil.ServerLogger(ctx, n.logger)
	var dbEnrollment db.Enrollment
	err := n.db.BeginFunc(ctx, func(tx pgx.Tx) (err error) {
		dbEnrollment, err = db.GetEnrollment(ctx, tx, request.GetNodeName())
		return
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "no node registration with specified node_name")
	} else if err != nil {
		logger.Error("GetEnrollment failed", zap.Error(err), zap.String("name", request.GetNodeName()))
		return nil, status.Error(codes.Internal, "failed to retrieve enrollment")
	}

	return &gen.NodeRegistration{
		Name:        dbEnrollment.Name,
		Address:     dbEnrollment.Address,
		Description: dbEnrollment.Description,
	}, nil
}

func (n *NodeServer) CreateNodeRegistration(ctx context.Context, request *gen.CreateNodeRegistrationRequest) (*gen.NodeRegistration, error) {
	logger := rpcutil.ServerLogger(ctx, n.logger)
	nodeReg := request.GetNodeRegistration()
	if nodeReg == nil {
		return nil, status.Error(codes.InvalidArgument, "node_registration must be supplied")
	}

	en, err := enrollment.EnrollAreaController(ctx, &gen.Enrollment{
		TargetName:     nodeReg.Name,
		TargetAddress:  nodeReg.Address,
		ManagerName:    n.managerName,
		ManagerAddress: n.managerAddr,
	}, n.authority)
	if err != nil {
		logger.Error("failed to enroll area controller", zap.Error(err),
			zap.String("target_address", nodeReg.Address))
		return nil, status.Error(codes.Unknown, "target refused registration")
	}

	err = n.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return db.CreateEnrollment(ctx, tx, db.Enrollment{
			Name:        en.TargetName,
			Description: nodeReg.Description,
			Address:     en.TargetAddress,
			Cert:        en.Certificate,
		})
	})
	if err != nil {
		logger.Error("db.CreateEnrollment failed", zap.Error(err))
		return nil, status.Error(codes.DataLoss, "failed to save the enrollment - manual intervention required")
	}
	return nodeReg, nil
}

func (n *NodeServer) ListNodeRegistrations(ctx context.Context, request *gen.ListNodeRegistrationsRequest) (*gen.ListNodeRegistrationsResponse, error) {
	logger := rpcutil.ServerLogger(ctx, n.logger)
	var dbEnrollments []db.Enrollment
	err := n.db.BeginFunc(ctx, func(tx pgx.Tx) (err error) {
		dbEnrollments, err = db.ListEnrollments(ctx, tx)
		return
	})
	if err != nil {
		logger.Error("db.ListEnrollments failed", zap.Error(err))
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

func (n *NodeServer) TestNodeCommunication(ctx context.Context, request *gen.TestNodeCommunicationRequest) (*gen.TestNodeCommunicationResponse, error) {
	reg, err := n.GetNodeRegistration(ctx, &gen.GetNodeRegistrationRequest{NodeName: request.GetNodeName()})
	if err != nil {
		return nil, err
	}

	logger := n.logger.With(zap.String("node_address", reg.Address))

	conn, err := grpc.DialContext(ctx, reg.Address, grpc.WithTransportCredentials(credentials.NewTLS(n.testTLSConfig)))
	if err != nil {
		logger.Error("failed to connect to node", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "unable to connect to target node")
	}

	client := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := client.ServerReflectionInfo(ctx)
	if err != nil {
		logger.Error("failed to reflect target node", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "target node reflection failed")
	}

	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host:           reg.Address,
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{ListServices: ""},
	})
	if err != nil {
		logger.Error("node reflection: send request", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "target node reflection failed")
	}

	res, err := stream.Recv()
	if err != nil {
		logger.Error("node reflection: receive response", zap.Error(err))
		return nil, status.Error(codes.Unavailable, "target node reflection failed")
	}

	var serviceNames []string
	for _, service := range res.GetListServicesResponse().GetService() {
		serviceNames = append(serviceNames, service.Name)
	}

	return &gen.TestNodeCommunicationResponse{Services: serviceNames}, nil
}
