package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Servers struct {
	Logger          *zap.Logger
	ShutdownTimeout time.Duration // How long to wait for connections to close when the context is cancelled
	GRPC            *grpc.Server
	GRPCAddress     string // Address to pass to net.Listen for the gRPC server
	HTTP            *http.Server
}

func (s *Servers) Serve(ctx context.Context) error {
	// softCtx is cancelled when we want to stop
	// after softCtx ends, we wait for ShutdownTimeout then cancel hardCtx as well
	group, softCtx := errgroup.WithContext(ctx)
	hardCtx, hardStop := context.WithCancel(softCtx)
	defer hardStop()
	go func() {
		<-softCtx.Done()
		s.Logger.Info("waiting for servers to stop", zap.Duration("timeout", s.ShutdownTimeout))
		time.Sleep(s.ShutdownTimeout)
		s.Logger.Warn("forcing servers to stop now")
		hardStop()
	}()

	if s.HTTP != nil {
		group.Go(func() error {
			return s.serveHTTPS(softCtx, hardCtx)
		})
	}

	if s.GRPC != nil {
		group.Go(func() error {
			return s.serveGRPC(softCtx, hardCtx)
		})
	}

	return group.Wait()
}

func (s *Servers) serveHTTPS(softCtx context.Context, hardCtx context.Context) error {
	var group errgroup.Group

	group.Go(func() error {
		// cert and key are provided by s.HTTP.TLSConfig
		return s.HTTP.ListenAndServeTLS("", "")
	})

	group.Go(func() error {
		<-softCtx.Done()
		return s.HTTP.Shutdown(hardCtx)
	})

	group.Go(func() error {
		<-hardCtx.Done()
		return s.HTTP.Close()
	})

	return group.Wait()
}

func (s *Servers) serveGRPC(softCtx context.Context, hardCtx context.Context) error {
	var group errgroup.Group

	grpcListener, err := net.Listen("tcp", s.GRPCAddress)
	if err != nil {
		return err
	}

	group.Go(func() error {
		return s.GRPC.Serve(grpcListener)
	})

	group.Go(func() error {
		<-softCtx.Done()
		s.GRPC.GracefulStop()
		return nil
	})

	group.Go(func() error {
		<-hardCtx.Done()
		s.GRPC.Stop()
		return nil
	})

	return group.Wait()
}
