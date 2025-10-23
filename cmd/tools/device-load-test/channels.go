package main

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type channelCounter struct {
	openUnary    atomic.Int64
	openStream   atomic.Int64
	openChannels atomic.Int64

	maxUnary    atomic.Int64
	maxStream   atomic.Int64
	maxChannels atomic.Int64
}

func (c *channelCounter) startChannel() {
	current := c.openChannels.Add(1)
	for {
		oldMax := c.maxChannels.Load()
		if current <= oldMax {
			break
		}
		if c.maxChannels.CompareAndSwap(oldMax, current) {
			break
		}
	}
}

func (c *channelCounter) stopChannel() {
	c.openChannels.Add(-1)
}

func (c *channelCounter) startUnary() {
	c.startChannel()
	current := c.openUnary.Add(1)
	for {
		oldMax := c.maxUnary.Load()
		if current <= oldMax {
			break
		}
		if c.maxUnary.CompareAndSwap(oldMax, current) {
			break
		}
	}
}

func (c *channelCounter) stopUnary() {
	c.openUnary.Add(-1)
	c.stopChannel()
}

func (c *channelCounter) startStream() {
	c.startChannel()
	current := c.openStream.Add(1)
	for {
		oldMax := c.maxStream.Load()
		if current <= oldMax {
			break
		}
		if c.maxStream.CompareAndSwap(oldMax, current) {
			break
		}
	}
}

func (c *channelCounter) stopStream() {
	c.openStream.Add(-1)
	c.stopChannel()
}

func (c *channelCounter) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		c.startUnary()
		defer c.stopUnary()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *channelCounter) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return stream, err
		}
		return &streamWrapper{
			ClientStream: stream,
			c:            c,
		}, nil
	}
}

func (c *channelCounter) CurrentCounts() channelCounts {
	return channelCounts{
		Unary:   c.openUnary.Load(),
		Stream:  c.openStream.Load(),
		Channel: c.openChannels.Load(),
	}
}

func (c *channelCounter) MaxCounts() channelCounts {
	return channelCounts{
		Unary:   c.maxUnary.Load(),
		Stream:  c.maxStream.Load(),
		Channel: c.maxChannels.Load(),
	}
}

type channelCounts struct {
	Unary   int64
	Stream  int64
	Channel int64
}

type streamWrapper struct {
	grpc.ClientStream
	c           *channelCounter
	initialised bool
}

func (s *streamWrapper) init() {
	if s.initialised {
		return
	}
	s.initialised = true
	s.c.startStream()
	_ = context.AfterFunc(s.Context(), s.c.stopStream)
}

func (s *streamWrapper) Header() (metadata.MD, error) {
	md, err := s.ClientStream.Header()
	s.init()
	return md, err
}

func (s *streamWrapper) RecvMsg(m any) error {
	err := s.ClientStream.RecvMsg(m)
	s.init()
	return err
}
