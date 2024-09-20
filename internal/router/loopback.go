package router

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Loopback implements grpc.ClientConnInterface using a MethodResolver to find the real connection for each request.
//
// Resolution of the real connection is deferred until the first message is sent.
type Loopback struct {
	mr MethodResolver
}

var _ grpc.ClientConnInterface = (*Loopback)(nil)

func NewLoopback(mr MethodResolver) *Loopback {
	return &Loopback{mr: mr}
}

func (l *Loopback) Invoke(ctx context.Context, fullMethodName string, args any, reply any, opts ...grpc.CallOption) error {
	argsProto, ok := args.(proto.Message)
	if !ok {
		return ErrNonProtoMessage
	}

	method, ok := l.mr.ResolveMethod(fullMethodName)
	if !ok {
		return ErrUnknownMethod
	}

	conn, err := method.Resolver.Resolve(messageCopyRecver{msg: argsProto})
	if err != nil {
		return err
	}

	return conn.Invoke(ctx, fullMethodName, args, reply, opts...)
}

func (l *Loopback) NewStream(ctx context.Context, desc *grpc.StreamDesc, fullMethodName string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	method, ok := l.mr.ResolveMethod(fullMethodName)
	if !ok {
		return nil, ErrUnknownMethod
	}

	return &deferredClientStream{
		ctx:            ctx,
		desc:           desc,
		fullMethodName: fullMethodName,
		opts:           opts,
		resolvedMethod: method,
	}, nil
}

// messageCopyRecver is a MsgRecver that copies a known proto.Message into the message provided by the caller of RecvMsg.
type messageCopyRecver struct {
	msg proto.Message
}

func (m messageCopyRecver) RecvMsg(msg any) error {
	msgProto, ok := msg.(proto.Message)
	if !ok {
		return ErrNonProtoMessage
	}
	proto.Merge(msgProto, m.msg)
	return nil
}

var (
	ErrNonProtoMessage = status.Error(codes.Internal, "non-protobuf messages not supported")
	ErrUnknownMethod   = status.Error(codes.Unimplemented, "unknown service method")
)

type deferredClientStream struct {
	// NewStream arguments
	ctx            context.Context
	desc           *grpc.StreamDesc
	fullMethodName string
	opts           []grpc.CallOption

	// source of the deferred client stream
	resolvedMethod Method

	// client stream
	ready  chan struct{}
	stream grpc.ClientStream
	err    error
}

func (d *deferredClientStream) Header() (metadata.MD, error) {
	select {
	case <-d.ready:
		if d.err != nil {
			return nil, d.err
		} else {
			return d.stream.Header()
		}
	case <-d.ctx.Done():
		return nil, nil
	}
}

func (d *deferredClientStream) Trailer() metadata.MD {
	select {
	case <-d.ready:
		if d.err != nil {
			return nil
		} else {
			return d.stream.Trailer()
		}
	default:
		// The interface prohibits calling Trailer() before RecvMsg() (or CloseAndRecv) has returned.
		// By the time the first RecvMsg has returned, the stream must be ready.
		panic("Trailer() called before RecvMsg() has returned")
	}
}

func (d *deferredClientStream) CloseSend() error {
	select {
	case <-d.ready:
		if d.err != nil {
			return d.err
		} else {
			return d.stream.CloseSend()
		}
	case <-d.ctx.Done():
		return d.ctx.Err()
	}
}

func (d *deferredClientStream) Context() context.Context {
	// the stream context is not the same as d.ctx!
	select {
	case <-d.ready:
		return d.stream.Context()
	default:
		// The interface prohibits using the stream in this way
		panic("Context() called before Header() / RecvMsg() has returned")
	}
}

func (d *deferredClientStream) SendMsg(m any) error {
	// concurrent calls to SendMsg are not allowed, so it doesn't matter that there's no lock held between
	// checking d.ready and closing d.ready
	select {
	case <-d.ready:
		if d.err != nil {
			return d.err
		} else {
			return d.stream.SendMsg(m)
		}
	default:
	}

	// resolve the stream - we only get one attempt at this so any error will be permanent
	defer close(d.ready)
	mProto, ok := m.(proto.Message)
	if !ok {
		err := ErrNonProtoMessage
		d.err = err
		return err
	}
	client, err := d.resolvedMethod.Resolver.Resolve(copyRecver{from: mProto})
	if err != nil {
		d.err = err
		return err
	}
	// create the real stream and send the captured message
	stream, err := client.NewStream(d.ctx, d.desc, d.fullMethodName, d.opts...)
	if err != nil {
		d.err = err
		return err
	}
	err = stream.SendMsg(m)
	if err != nil {
		d.err = err
		return err
	}
	d.stream = stream
	return nil
}

func (d *deferredClientStream) RecvMsg(m any) error {
	select {
	case <-d.ready:
		if d.err != nil {
			return d.err
		} else {
			return d.stream.RecvMsg(m)
		}
	case <-d.ctx.Done():
		return d.ctx.Err()
	}
}

type copyRecver struct {
	from proto.Message
}

func (c copyRecver) RecvMsg(m any) error {
	mProto, ok := m.(proto.Message)
	if !ok {
		return ErrNonProtoMessage
	}
	proto.Merge(mProto, c.from)
	return nil
}
