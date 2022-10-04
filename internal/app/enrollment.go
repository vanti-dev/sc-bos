package app

import (
	"context"
	"errors"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/internal/util/pki/expire"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func ServeEnrollment(ctx context.Context, logger *zap.Logger, server *enrollment.Server, key pki.PrivateKey, listenGRPC string) error {
	certSource := pki.CacheSource(pki.SelfSignedSource(key, pki.WithIfaces()), expire.AfterProgress(0.5))
	tlsConfig := pki.TLSConfig(certSource)

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	reflection.Register(grpcServer)
	gen.RegisterEnrollmentApiServer(grpcServer, server)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := ServeGRPC(ctx, grpcServer, listenGRPC, 0, logger.Named("server.grpc"))
		if err != nil {
			logger.Warn("server stopped", zap.Error(err))
		}
	}()
	logger.Info("gRPC serving; waiting for enrollment")

	_, ok := server.Wait(ctx)
	if ok {
		logger.Info("controller is now enrolled")
	} else {
		logger.Error("server stopped without an enrollment")
		return enrollment.ErrNotEnrolled
	}
	return nil
}

func RemoteManager(server *enrollment.Server, opts ...grpc.DialOption) (connFunc func() (*grpc.ClientConn, error), closeFunc func() error) {
	var lastManagerAddress string
	var lastConn *grpc.ClientConn
	var lastErr error
	closeConn := func() error {
		if lastConn != nil {
			err := lastConn.Close()
			lastConn = nil
			lastErr = nil
			return err
		}
		return nil
	}
	return func() (*grpc.ClientConn, error) {
		e, ok := server.Enrollment()
		if !ok {
			_ = closeConn()
			return nil, errors.New("not enrolled")
		}

		if e.ManagerAddress != lastManagerAddress {
			_ = closeConn()
			lastManagerAddress = e.ManagerAddress
		}

		if lastConn == nil {
			lastConn, lastErr = grpc.Dial(lastManagerAddress, opts...)
		}

		return lastConn, lastErr
	}, closeConn
}
