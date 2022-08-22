package app

import (
	"context"
	"crypto"
	"crypto/tls"

	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"github.com/vanti-dev/bsp-ew/internal/util/pki"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func ServeEnrollment(ctx context.Context, logger *zap.Logger, server *enrollment.Server, key crypto.PrivateKey, listenGRPC string) error {
	certSource, err := pki.NewSelfSignedCertSource(key, logger)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		GetCertificate: certSource.TLSConfigGetCertificate,
	}

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
