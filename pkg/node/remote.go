package node

import (
	"context"
	"io"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

// Remote represents a remote smart core node
type Remote interface {
	io.Closer
	Target() string
	Connect(ctx context.Context) (*grpc.ClientConn, error)
}

// Dial calls grpc.NewClient and returns it as a Remote.
func Dial(ctx context.Context, target string, opts ...grpc.DialOption) Remote {
	conn, err := dial(target, opts...)
	return &eagerRemote{
		target:  target,
		conn:    conn,
		dialErr: err,
	}
}

// eagerRemote implements Remote by requiring an existing grpc.ClientConn.
type eagerRemote struct {
	target string

	conn    *grpc.ClientConn
	dialErr error
}

func (r *eagerRemote) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

func (r *eagerRemote) Target() string {
	return r.target
}

func (r *eagerRemote) Connect(_ context.Context) (*grpc.ClientConn, error) {
	return r.conn, r.dialErr
}

// DialChan returns a Remote that connects to the last received target from targets.
// Remote.Connect will return the same grpc.ClientConn, the underlying connection will be replaced each time targets
// emits a value.
//
// Do not use grpc.WithBlock option with DialChan.
func DialChan(ctx context.Context, targets <-chan string, opts ...grpc.DialOption) Remote {
	// A note about scheme names:
	// The manual resolver will return the scheme name exactly as we give it, but when looking up resolvers in dial
	// grpc explicitly lower cases the parsed target scheme so if we don't use lower case here then it won't find it!
	res := manual.NewBuilderWithScheme("dialchan")
	// This is required, even though it's empty, as any operation on the conn will block until at least one call to
	// res.CC.UpdateState. This does this for us when the resolver builders Build method is called during dial.
	res.InitialState(resolver.State{})
	opts = append(opts, grpc.WithResolvers(res))
	conn, err := dial("dialchan:ignored", opts...)
	remote := &chanRemote{
		targets:  targets,
		resolver: res,
		conn:     conn,
		dialErr:  err,

		closed: make(chan struct{}),
	}
	if err == nil { // only start if the conn was a success
		conn.Connect() // connect causes the resolver to be built
		remote.start()
	}
	return remote
}

type chanRemote struct {
	targets <-chan string

	resolver *manual.Resolver
	conn     *grpc.ClientConn
	dialErr  error

	targetMu sync.Mutex
	target   string

	closed    chan struct{}
	closeErr  error
	closeOnce sync.Once
}

func (c *chanRemote) start() {
	go func() {
		for {
			select {
			case target, ok := <-c.targets:
				if !ok {
					// we don't close here, instead we're saying if the target chan closes we make no more updates
					// to the address.
					return
				}

				c.targetMu.Lock()
				c.target = target
				c.targetMu.Unlock()

				// let the conn know that the address has changed
				c.resolver.UpdateState(resolver.State{Addresses: []resolver.Address{
					{Addr: target},
				}})
			case <-c.closed:
				return
			}
		}
	}()
}

func (c *chanRemote) Close() error {
	c.closeOnce.Do(func() {
		if c.conn != nil {
			c.closeErr = c.conn.Close()
		}
		close(c.closed)
	})
	return c.closeErr
}

func (c *chanRemote) Target() string {
	c.targetMu.Lock()
	defer c.targetMu.Unlock()
	return c.target
}

func (c *chanRemote) Connect(_ context.Context) (*grpc.ClientConn, error) {
	return c.conn, c.dialErr
}

func dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOpts := processOpts(opts)

	conn, err := grpc.NewClient(target, opts...)
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
	return conn, err
}

func processOpts(opts []grpc.DialOption) *dialOptions {
	dialOpts := &dialOptions{}
	for _, opt := range opts {
		if nOpt, ok := opt.(DialOption); ok {
			nOpt.applyOpts(dialOpts)
		}
	}
	return dialOpts
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
