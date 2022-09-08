package node

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"io"
)

// Remote represents a remote smart core node
type Remote interface {
	io.Closer
	Target() string
	Connect(ctx context.Context) (*grpc.ClientConn, error)
}

// Dial calls grpc.Dial and returns it as a Remote.
// The call to grpc.Dial may happen upon first call to Remote.Connect or may happen as part of this call.
func Dial(ctx context.Context, target string, opts ...grpc.DialOption) (Remote, error) {
	dialOpts := &dialOptions{}
	for _, opt := range opts {
		if nOpt, ok := opt.(DialOption); ok {
			nOpt.applyOpts(dialOpts)
		}
	}

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}

	if dialOpts.stateChangeLogger != nil {
		logger := dialOpts.stateChangeLogger.Sugar()
		go func() {
			// watch state changes, returns when the conn is closed
			var state connectivity.State
			for state != connectivity.Shutdown && conn.WaitForStateChange(context.Background(), state) {
				state = conn.GetState()
				logger.Infof("%v is %v", target, state)
			}
		}()
	}

	return &remoteNode{
		target:  target,
		conn:    conn,
		dialErr: err,
	}, nil
}

type dialOptions struct {
	stateChangeLogger *zap.Logger
}

type DialOption interface {
	grpc.DialOption
	applyOpts(node *dialOptions)
}

// EmptyDialOption does not change how a dial will be performed.
// Useful for embedding into custom DialOption types to extend the dial api.
type EmptyDialOption struct {
	grpc.EmptyDialOption
}

func (e EmptyDialOption) applyOpts(_ *dialOptions) {}

type dialOptionFunc struct {
	grpc.DialOption
	f func(req *dialOptions)
}

func (d dialOptionFunc) applyOpts(node *dialOptions) {
	d.f(node)
}

func newDialOption(f func(n *dialOptions)) grpc.DialOption {
	return dialOptionFunc{grpc.EmptyDialOption{}, f}
}

func WithLogStateChange(logger *zap.Logger) grpc.DialOption {
	return newDialOption(func(n *dialOptions) {
		n.stateChangeLogger = logger
	})
}

type remoteNode struct {
	target string

	conn    *grpc.ClientConn
	dialErr error
}

func (r *remoteNode) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

func (r *remoteNode) Target() string {
	return r.target
}

func (r *remoteNode) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	return r.conn, r.dialErr
}
