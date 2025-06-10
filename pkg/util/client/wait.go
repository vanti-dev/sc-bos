package client

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var ErrTransientFailure = errors.New("transient failure")

// WaitForReady waits for the grpc.ClientConn to be in a ready (connected) state or for the connection to fail.
// Transient failures will return ErrTransientFailure, if cc is shutdown it will return grpc.ErrServerStopped,
// otherwise it will block until the connection is ready or the context is done.
func WaitForReady(ctx context.Context, cc *grpc.ClientConn) error {
	// Similar to using grpc.WithBlock during dial
	for {
		s := cc.GetState()
		switch s {
		case connectivity.Idle:
			cc.Connect()
		case connectivity.Ready:
			return nil
		case connectivity.Shutdown:
			return grpc.ErrServerStopped
		case connectivity.Connecting:
		case connectivity.TransientFailure:
			return ErrTransientFailure
		}

		if !cc.WaitForStateChange(ctx, s) {
			return ctx.Err()
		}
	}
}
