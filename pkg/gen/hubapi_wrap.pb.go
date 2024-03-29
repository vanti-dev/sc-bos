// Code generated by protoc-gen-wrapper. DO NOT EDIT.

package gen

import (
	context "context"
	wrap "github.com/smart-core-os/sc-golang/pkg/wrap"
	grpc "google.golang.org/grpc"
)

// WrapHubApi	adapts a HubApiServer	and presents it as a HubApiClient
func WrapHubApi(server HubApiServer) HubApiClient {
	return &hubApiWrapper{server}
}

type hubApiWrapper struct {
	server HubApiServer
}

// compile time check that we implement the interface we need
var _ HubApiClient = (*hubApiWrapper)(nil)

// UnwrapServer returns the underlying server instance.
func (w *hubApiWrapper) UnwrapServer() HubApiServer {
	return w.server
}

// Unwrap implements wrap.Unwrapper and returns the underlying server instance as an unknown type.
func (w *hubApiWrapper) Unwrap() any {
	return w.UnwrapServer()
}

func (w *hubApiWrapper) GetHubNode(ctx context.Context, req *GetHubNodeRequest, _ ...grpc.CallOption) (*HubNode, error) {
	return w.server.GetHubNode(ctx, req)
}

func (w *hubApiWrapper) ListHubNodes(ctx context.Context, req *ListHubNodesRequest, _ ...grpc.CallOption) (*ListHubNodesResponse, error) {
	return w.server.ListHubNodes(ctx, req)
}

func (w *hubApiWrapper) PullHubNodes(ctx context.Context, in *PullHubNodesRequest, opts ...grpc.CallOption) (HubApi_PullHubNodesClient, error) {
	stream := wrap.NewClientServerStream(ctx)
	server := &pullHubNodesHubApiServerWrapper{stream.Server()}
	client := &pullHubNodesHubApiClientWrapper{stream.Client()}
	go func() {
		err := w.server.PullHubNodes(in, server)
		stream.Close(err)
	}()
	return client, nil
}

type pullHubNodesHubApiClientWrapper struct {
	grpc.ClientStream
}

func (w *pullHubNodesHubApiClientWrapper) Recv() (*PullHubNodesResponse, error) {
	m := new(PullHubNodesResponse)
	if err := w.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type pullHubNodesHubApiServerWrapper struct {
	grpc.ServerStream
}

func (s *pullHubNodesHubApiServerWrapper) Send(response *PullHubNodesResponse) error {
	return s.ServerStream.SendMsg(response)
}

func (w *hubApiWrapper) InspectHubNode(ctx context.Context, req *InspectHubNodeRequest, _ ...grpc.CallOption) (*HubNodeInspection, error) {
	return w.server.InspectHubNode(ctx, req)
}

func (w *hubApiWrapper) EnrollHubNode(ctx context.Context, req *EnrollHubNodeRequest, _ ...grpc.CallOption) (*HubNode, error) {
	return w.server.EnrollHubNode(ctx, req)
}

func (w *hubApiWrapper) RenewHubNode(ctx context.Context, req *RenewHubNodeRequest, _ ...grpc.CallOption) (*HubNode, error) {
	return w.server.RenewHubNode(ctx, req)
}

func (w *hubApiWrapper) TestHubNode(ctx context.Context, req *TestHubNodeRequest, _ ...grpc.CallOption) (*TestHubNodeResponse, error) {
	return w.server.TestHubNode(ctx, req)
}

func (w *hubApiWrapper) ForgetHubNode(ctx context.Context, req *ForgetHubNodeRequest, _ ...grpc.CallOption) (*ForgetHubNodeResponse, error) {
	return w.server.ForgetHubNode(ctx, req)
}
