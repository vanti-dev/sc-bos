// Command proxy-test is an experiment using grpc-proxy to see how applicable it would be for gateways.
package main

import (
	"context"
	"log"
	"net"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func main() {
	realLis, err := net.Listen("tcp", ":")
	if err != nil {
		log.Fatalln(err)
	}
	realServer := grpc.NewServer()
	gen.RegisterTestApiServer(realServer, &testApi{t: &gen.Test{Data: "initial"}})
	go realServer.Serve(realLis)
	defer realServer.Stop()

	log.Printf("Real server hosted on %v", realLis.Addr())
	realConn, err := grpc.Dial(realLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	handler := proxy.TransparentHandler(func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		log.Printf("proxy(%v, %v)", ctx, fullMethodName)
		return ctx, realConn, nil
	})

	proxyLis, err := net.Listen("tcp", ":")
	if err != nil {
		log.Fatalln(err)
	}
	proxyServer := grpc.NewServer(grpc.UnknownServiceHandler(handler),
		grpc.ChainUnaryInterceptor(logUnaryServerCalls),
		grpc.ChainStreamInterceptor(logStreamServerCalls))
	go proxyServer.Serve(proxyLis)
	defer proxyServer.Stop()

	log.Printf("Proxy server hosted on %v", proxyLis.Addr())
	proxyConn, err := grpc.Dial(proxyLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := gen.NewTestApiClient(proxyConn)
	test, err := client.GetTest(context.Background(), &gen.GetTestRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Got response %v", test)
}

func logUnaryServerCalls(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("intercept(%v, %v, %v)", ctx, req, info)
	return handler(ctx, req)
}
func logStreamServerCalls(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("intercept( %v, %v, %v)", srv, ss, info)
	return handler(srv, &loggingServerStream{ss})
}

type loggingServerStream struct {
	grpc.ServerStream
}

func (ss *loggingServerStream) SendMsg(m interface{}) error {
	log.Printf("intercept.SendMsg((%T) {%v})", m, m)
	return ss.ServerStream.SendMsg(m)
}

func (ss *loggingServerStream) RecvMsg(m interface{}) error {
	err := ss.ServerStream.RecvMsg(m)
	log.Printf("intercept.RecvMsg((%T) {%v}) %v", m, m, err)
	return err
}

type testApi struct {
	gen.UnimplementedTestApiServer
	t *gen.Test
}

func (t *testApi) GetTest(ctx context.Context, request *gen.GetTestRequest) (*gen.Test, error) {
	log.Printf("GetTest(%v, %v)", ctx, request)
	return t.t, nil
}

func (t *testApi) UpdateTest(ctx context.Context, request *gen.UpdateTestRequest) (*gen.Test, error) {
	log.Printf("UpdateTest(%v, %v)", ctx, request)
	t.t = request.Test
	return t.t, nil
}

func (t *testApi) mustEmbedUnimplementedTestApiServer() {
	// TODO implement me
	panic("implement me")
}
