package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ServeHTTPS(ctx context.Context, server *http.Server, timeout time.Duration, logger *zap.Logger) error {
	serveErr := make(chan error, 1)
	go func() {
		// assume the cert and key are provided in server.TLSConfig
		serveErr <- server.ListenAndServeTLS("", "")
	}()
	logger.Info("now serving HTTPS", zap.String("addr", server.Addr))

	select {
	case err := <-serveErr:
		// we received an error from ListenAndServe earlier than expected - there must have been some problem starting
		// the server
		logger.Error("server stopped unexpectedly", zap.Error(err))
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	logger.Info("attempting to stop HTTPS server", zap.Duration("timeout", timeout))
	err := server.Shutdown(shutdownCtx)
	if errors.Is(err, context.DeadlineExceeded) {
		// the server didn't shut down in time, so force it to stop serving immediately
		logger.Warn("timeout expired - forcing HTTPS server to stop now")
		_ = server.Close()
	}
	return ctx.Err()
}

func ServeGRPC(ctx context.Context, server *grpc.Server, listen string, timeout time.Duration, logger *zap.Logger) error {
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()
	logger.Info("now serving gRPC", zap.String("addr", listen))

	// The server may return an error early if there was a problem starting up.
	select {
	case err := <-serveErr:
		logger.Error("server stopped unexpectedly", zap.Error(err))
		return err
	case <-ctx.Done():
	}

	// Try to stop the server gracefully. If it takes too long, force it to stop.
	logger.Info("attempting to stop gRPC server", zap.Duration("timeout", timeout))
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		server.GracefulStop()
	}()

	select {
	case <-stopped:
	case <-time.After(timeout):
		logger.Warn("timeout expired - forcing gRPC server to stop now")
		server.Stop()
	}

	return multierr.Combine(<-serveErr, ctx.Err())
}
