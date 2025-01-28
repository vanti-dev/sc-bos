package router

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
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

	method, err := l.mr.ResolveMethod(fullMethodName)
	if err != nil {
		return err
	}

	conn, err := method.Resolver.ResolveConn(copyRecver{from: argsProto})
	if err != nil {
		return err
	}

	return conn.Invoke(ctx, fullMethodName, args, reply, opts...)
}

func (l *Loopback) NewStream(ctx context.Context, desc *grpc.StreamDesc, fullMethodName string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	method, err := l.mr.ResolveMethod(fullMethodName)
	if err != nil {
		return nil, err
	}

	return &deferredClientStream{
		ctx:            ctx,
		desc:           desc,
		fullMethodName: fullMethodName,
		opts:           opts,
		resolvedMethod: method,
		ready:          make(chan struct{}),
	}, nil
}

var (
	ErrNonProtoMessage = status.Error(codes.Internal, "non-protobuf messages not supported")
)

// deferredClientStream is a grpc.ClientStream that uses a Method to resolve the real connection when the first message
// is sent using SendMsg, at which point the real stream is opened. After this, calls pass through to the real stream.
// An error resolving the connection from the Method is permanent.
type deferredClientStream struct {
	// NewStream arguments - used to create the real stream
	ctx            context.Context
	desc           *grpc.StreamDesc
	fullMethodName string
	opts           []grpc.CallOption
	// source of the deferred client connection used to create the real stream
	resolvedMethod Method

	// client stream - filled after SendMsg is called
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
	client, err := d.resolvedMethod.Resolver.ResolveConn(copyRecver{from: mProto})
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

// copyRecver is a MsgRecver that copies a known proto.Message into the message provided by the caller of RecvMsg.
type copyRecver struct {
	from proto.Message
}

func (c copyRecver) RecvMsg(m any) error {
	mProto, ok := m.(proto.Message)
	if !ok {
		return ErrNonProtoMessage
	}
	safeProtoMerge(mProto, c.from)
	return nil
}

// safeProtoMerge is like proto.Merge but doesn't panic when using dynamicpb messages.
//
// proto.Merge will panic if the descriptors of src and dest aren't equal according to ==.
// This has problems for us because we're mixing dynamic and explicit messages, which have subtly different descriptors in ways that are implementation dependent.
// To fix this we fall back to marshalling and unmarshalling the messages, which is what proto.Merge is intended to imitate.
func safeProtoMerge(dst, src proto.Message) {
	_, srcDynamic := src.(*dynamicpb.Message)
	_, dstDynamic := dst.(*dynamicpb.Message)
	if !srcDynamic && !dstDynamic {
		proto.Merge(dst, src)
		return
	}

	// fall back to marshalling and unmarshalling,
	// we assume that the descriptors are the same if their names are the same
	if want, got := dst.ProtoReflect().Descriptor().FullName(), src.ProtoReflect().Descriptor().FullName(); want != got {
		panic(fmt.Sprintf("descriptor mismatch: %v != %v", got, want))
	}
	data, err := proto.Marshal(src)
	if err != nil {
		panic(fmt.Sprintf("src marshal: %v", err))
	}
	if err := proto.Unmarshal(data, dst); err != nil {
		panic(fmt.Sprintf("dst unmarshal: %v", err))
	}
}
