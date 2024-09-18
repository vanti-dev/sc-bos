package router

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
)

func TestStreamHandler(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	// downstream nodes
	n1Client, err := newNode(t, ctx, "n1")
	if err != nil {
		t.Fatalf("newNode(n1) = %v", err)
	}
	n2Client, err := newNode(t, ctx, "n2")
	if err != nil {
		t.Fatalf("newNode(n2) = %v", err)
	}

	reg := New()
	registerService(reg, traits.OnOffApi_ServiceDesc.ServiceName, n1Client, n2Client)

	proxyServer := grpc.NewServer(grpc.UnknownServiceHandler(StreamHandler(reg)))
	proxyLis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { proxyServer.Stop() })
	go func() {
		if err := proxyServer.Serve(proxyLis); err != nil {
			t.Errorf("proxyServer.Serve(proxyLis) = %v", err)
		}
	}()

	proxyConn, err := bufConn(ctx, proxyLis)
	if err != nil {
		t.Fatalf("bufConn(proxyLis) = %v", err)
	}
	t.Cleanup(func() { proxyConn.Close() })

	onOffClient := traits.NewOnOffApiClient(proxyConn)
	t.Run("downstream", func(t *testing.T) {
		testDownstream(t, ctx, onOffClient, "n1")
		testDownstream(t, ctx, onOffClient, "n2")
	})

	t.Run("unknown key", func(t *testing.T) {
		// known api, unknown key
		_, err = onOffClient.GetOnOff(ctx, &traits.GetOnOffRequest{Name: "missing"})
		if code := status.Code(err); code != codes.NotFound {
			t.Fatalf("onOffClient.GetOnOff(missing) want NotFound, got: %v", err)
		}
	})

	t.Run("unknown api", func(t *testing.T) {
		// unknown api
		client2 := traits.NewOnOffInfoClient(proxyConn)
		_, err = client2.DescribeOnOff(ctx, &traits.DescribeOnOffRequest{Name: "n1"})
		if code := status.Code(err); code != codes.Unimplemented {
			t.Fatalf("client2.DescribeOnOff(n1) want Unimplemented, got: %v", err)
		}
	})
}

func testDownstream(t *testing.T, ctx context.Context, client traits.OnOffApiClient, name string) {
	t.Helper()

	ctx, stop := context.WithTimeout(ctx, time.Second)
	defer stop() // also cancels the stream

	stream, err := client.PullOnOff(ctx, &traits.PullOnOffRequest{Name: name, UpdatesOnly: true})
	if err != nil {
		t.Fatalf("client.PullOnOff(%s) = %v", name, err)
	}

	// check initial state
	res, err := client.GetOnOff(ctx, &traits.GetOnOffRequest{Name: name})
	if err != nil {
		t.Fatalf("client.GetOnOff(%s) = %v", name, err)
	}
	if diff := cmp.Diff(res, &traits.OnOff{State: traits.OnOff_OFF}, protocmp.Transform()); diff != "" {
		t.Fatalf("client.GetOnOff(%s) mismatch (-want +got):\n%s", name, diff)
	}

	// update state
	res, err = client.UpdateOnOff(ctx, &traits.UpdateOnOffRequest{Name: name, OnOff: &traits.OnOff{State: traits.OnOff_ON}})
	if err != nil {
		t.Fatalf("client.UpdateOnOff(%s) = %v", name, err)
	}
	if diff := cmp.Diff(res, &traits.OnOff{State: traits.OnOff_ON}, protocmp.Transform()); diff != "" {
		t.Fatalf("client.UpdateOnOff(%s) mismatch (-want +got):\n%s", name, diff)
	}

	// check stream
	event, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv(%s) = %v", name, err)
	}
	// clear timestamps to make comparison easier
	for i := range event.Changes {
		event.Changes[i].ChangeTime = nil
	}
	wantEvent := &traits.PullOnOffResponse{Changes: []*traits.PullOnOffResponse_Change{
		{
			Name:  name,
			OnOff: &traits.OnOff{State: traits.OnOff_ON},
		},
	}}
	if diff := cmp.Diff(event, wantEvent, protocmp.Transform()); diff != "" {
		t.Fatalf("stream.Recv(%s) mismatch (-want +got):\n%s", name, diff)
	}
}

type namedMessage interface {
	proto.Message
	GetName() string
}

func nameKey[T namedMessage](m T) KeyFunc {
	return func(mr MsgRecver) (string, error) {
		req := m.ProtoReflect().New().Interface().(T) // new instance of T
		if err := mr.RecvMsg(req); err != nil {
			return "", err
		}
		return req.GetName(), nil
	}
}

func registerService(reg *Router, service string, clients ...*grpc.ClientConn) {
	reg.SupportService(NewUnroutedService(serviceDescriptor(service)))

	for i, client := range clients {
		name := fmt.Sprintf("n%d", i+1)
		err := reg.AddRoute(service, name, client)
		if err != nil {
			panic(err)
		}
	}
}

func newNode(t *testing.T, ctx context.Context, name string) (*grpc.ClientConn, error) {
	lis := bufconn.Listen(1024 * 1024)
	s := nodeServer(name)
	c, err := bufConn(ctx, lis)

	t.Cleanup(func() { c.Close() })
	t.Cleanup(func() { s.Stop() })

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("%s.Serve() = %v", name, err)
		}
	}()
	return c, err
}

// nodeServer returns a *grpc.Server that implements the OnOffApi service that responds to the given name.
func nodeServer(name string) *grpc.Server {
	s := grpc.NewServer()
	r := onoff.NewApiRouter() // use a router to force NOT_FOUND for unknown names
	m := onoff.NewModel(onoff.WithInitialOnOff(&traits.OnOff{State: traits.OnOff_OFF}))
	r.Add(name, onoff.WrapApi(onoff.NewModelServer(m)))
	traits.RegisterOnOffApiServer(s, r)
	return s
}

func bufConn(ctx context.Context, buf *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return buf.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func serviceDescriptor(name string) protoreflect.ServiceDescriptor {
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil {
		panic(err)
	}
	sd, ok := desc.(protoreflect.ServiceDescriptor)
	if !ok {
		panic("not a service descriptor")
	}
	return sd
}
