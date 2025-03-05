package temperaturepb

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestGetTemperature_priority(t *testing.T) {
	// getting without specifying a priority should read from the peer
	t.Run("get without priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}
		ts.getTemperature.Returns(want)

		ps := &PriorityServer{TemperatureApiServer: ts}

		got, err := ps.GetTemperature(context.Background(), &gen.GetTemperatureRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
			t.Errorf("GetTemperature unexpected response (-want +got):\n%s", diff)
		}
	})

	// getting a priority level when there aren't any priorities stored should return nil
	t.Run("get no priorities", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)

		ps := &PriorityServer{TemperatureApiServer: ts}

		_, ctx := newTestServerTransportStream(context.Background())
		got, err := ps.GetTemperature(ctx, &gen.GetTemperatureRequest{Priority: 10})
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff((*gen.Temperature)(nil), got, protocmp.Transform()); diff != "" {
			t.Errorf("GetTemperature unexpected response (-want +got):\n%s", diff)
		}
	})

	// getting a priority level that hasn't been set should return nil
	t.Run("get priority unset", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)

		ps := &PriorityServer{TemperatureApiServer: ts}
		ps.priorityArray[20] = &gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}}
		assertPriorityCall(t, nil, 20, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.GetTemperature(ctx, &gen.GetTemperatureRequest{Priority: 10})
		})
	})

	// gets using a priority should return the value stored at that priority
	t.Run("get priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.getTemperature.Returns(&gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}})

		ps := &PriorityServer{TemperatureApiServer: ts}
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}
		ps.priorityArray[10] = &gen.UpdateTemperatureRequest{Temperature: want}
		assertPriorityCall(t, want, 10, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.GetTemperature(ctx, &gen.GetTemperatureRequest{Priority: 10})
		})
	})

	// writes should use the default priority if none is given
	t.Run("update no priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.updateTemperature.ReturnsFunc(func(ctx context.Context, req *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
			return req.Temperature, nil
		})

		ps := &PriorityServer{TemperatureApiServer: ts}
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}
		assertPriorityCall(t, want, gen.Priority_DEFAULT, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: want,
			})
		})
		ts.updateTemperature.AssertCalledWith(&gen.UpdateTemperatureRequest{Temperature: want})
		ts.updateTemperature.AssertNoMoreCalls()
	})

	// writing a lower priority to those stored already should not write to the peer
	t.Run("update low priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.updateTemperature.ReturnsFunc(func(ctx context.Context, req *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
			return req.Temperature, nil
		})

		ps := &PriorityServer{TemperatureApiServer: ts}
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}
		ps.priorityArray[20] = &gen.UpdateTemperatureRequest{Temperature: want}
		assertPriorityCall(t, want, 20, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 25}},
				Priority:    30,
			})
		})
		ts.updateTemperature.AssertNoMoreCalls()
	})

	// writing at a higher priority to those already stored should write to the peer
	t.Run("update high priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.updateTemperature.ReturnsFunc(func(ctx context.Context, req *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
			return req.Temperature, nil
		})

		ps := &PriorityServer{TemperatureApiServer: ts}
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 25}}
		ps.priorityArray[20] = &gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}}
		assertPriorityCall(t, want, 10, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: want,
				Priority:    10,
			})
		})
		ts.updateTemperature.AssertCalledWith(&gen.UpdateTemperatureRequest{Temperature: want, Priority: 10})
		ts.updateTemperature.AssertNoMoreCalls()
	})

	// writing a delta update should result in a concrete write to the peer
	t.Run("update delta", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.getTemperature.Returns(&gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}})
		ts.updateTemperature.ReturnsFunc(func(ctx context.Context, req *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
			return req.Temperature, nil
		})

		ps := &PriorityServer{TemperatureApiServer: ts}
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 23}}
		assertPriorityCall(t, want, 10, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 2}},
				Delta:       true,
				Priority:    10,
			})
		})
		ts.updateTemperature.AssertCalledWith(&gen.UpdateTemperatureRequest{Temperature: want, Priority: 10})
		ts.updateTemperature.AssertNoMoreCalls()
	})

	// writing a delta update with a lower priority should store a concrete write in the priority array
	t.Run("update delta low priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.getTemperature.Returns(&gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}})
		ts.updateTemperature.ReturnsFunc(func(ctx context.Context, req *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
			return req.Temperature, nil
		})

		ps := &PriorityServer{TemperatureApiServer: ts}
		ps.priorityArray[20] = &gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}}
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 21}}
		assertPriorityCall(t, want, 20, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 2}},
				Delta:       true,
				Priority:    30,
			})
		})
		ts.updateTemperature.AssertNoMoreCalls()
		got := ps.priorityArray[30]
		want = &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 23}}
		if diff := cmp.Diff(&gen.UpdateTemperatureRequest{Temperature: want, Priority: 30}, got, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected priorityArray[30] (-want +got):\n%s", diff)
		}
	})

	// deleting a priority should write the next priority level available
	t.Run("delete priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		ts.updateTemperature.ReturnsFunc(func(ctx context.Context, req *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
			return req.Temperature, nil
		})

		ps := &PriorityServer{TemperatureApiServer: ts}
		ps.priorityArray[20] = &gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 20}}}
		ps.priorityArray[30] = &gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 30}}}

		assertPriorityCall(t, &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 30}}, 30, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: nil,
				Priority:    20,
			})
		})
		ts.updateTemperature.AssertCalledWith(&gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 30}}})
		ts.updateTemperature.AssertNoMoreCalls()
	})

	// deleting a priority when there are no more priorities should write nothing but return the peer value
	t.Run("delete only priority", func(t *testing.T) {
		ts := newTestTemperatureApiServer(t)
		want := &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 20}}
		ts.getTemperature.Returns(want)

		ps := &PriorityServer{TemperatureApiServer: ts}
		ps.priorityArray[20] = &gen.UpdateTemperatureRequest{Temperature: &gen.Temperature{SetPoint: &types.Temperature{ValueCelsius: 22}}}

		assertPriorityCall(t, want, 0, func(ctx context.Context) (*gen.Temperature, error) {
			return ps.UpdateTemperature(ctx, &gen.UpdateTemperatureRequest{
				Temperature: nil,
				Priority:    20,
			})
		})
		ts.updateTemperature.AssertNoMoreCalls()
	})
}

// assertPriorityCall calls fn then asserts that the response matches wantRes and a header is set on the response with the priority wantPriority.
func assertPriorityCall[Res proto.Message](t *testing.T, wantRes Res, wantPriority gen.Priority_Level, fn func(ctx context.Context) (Res, error)) {
	t.Helper()
	stream, ctx := newTestServerTransportStream(context.Background())
	gotRes, err := fn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRes, gotRes, protocmp.Transform()); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
	if wantPriority == 0 {
		if vals, ok := stream.HasHeader("PriorityLevel", ""); ok {
			t.Errorf("expected no PriorityLevel header, got %v", vals)
		}
	} else {
		if vals, ok := stream.HasHeader("PriorityLevel", strconv.Itoa(int(wantPriority))); !ok {
			t.Errorf("expected PriorityLevel header == %d, got %v", wantPriority, vals)
		}
	}
}

type testServer struct {
	gen.UnimplementedTemperatureApiServer
	getTemperature    *testRPC[*gen.GetTemperatureRequest, *gen.Temperature]
	updateTemperature *testRPC[*gen.UpdateTemperatureRequest, *gen.Temperature]
}

func newTestTemperatureApiServer(t *testing.T) *testServer {
	return &testServer{
		getTemperature:    newTestRPC[*gen.GetTemperatureRequest, *gen.Temperature](t, "GetTemperature"),
		updateTemperature: newTestRPC[*gen.UpdateTemperatureRequest, *gen.Temperature](t, "UpdateTemperature"),
	}
}

func (t *testServer) GetTemperature(ctx context.Context, request *gen.GetTemperatureRequest) (*gen.Temperature, error) {
	return t.getTemperature.Do(ctx, request)
}

func (t *testServer) UpdateTemperature(ctx context.Context, request *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
	return t.updateTemperature.Do(ctx, request)
}

type testRPC[Req, Res proto.Message] struct {
	t         *testing.T
	name      string
	requests  []rpcArgs[Req]
	responses rpcRes[Req, Res]
}

func newTestRPC[Req, Res proto.Message](t *testing.T, name string) *testRPC[Req, Res] {
	return &testRPC[Req, Res]{t: t, name: name}
}

func (rpc *testRPC[Req, Res]) Do(ctx context.Context, req Req) (Res, error) {
	if rpc.responses == nil {
		var zero Res
		return zero, fmt.Errorf("unexpected call to %s(%v, %v)", rpc.name, ctx, req)
	}
	ctx = grpc.NewContextWithServerTransportStream(ctx, &testServerTransportStream{method: rpc.name})
	rpc.requests = append(rpc.requests, rpcArgs[Req]{rpc.t, ctx, req})
	res, err, next := rpc.responses(ctx, req)
	rpc.responses = next
	return res, err
}

func (rpc *testRPC[Req, Res]) Errors(err error) {
	rpc.responses = rpc.responses.Append(func(ctx context.Context, req Req) (Res, error, rpcRes[Req, Res]) {
		var zero Res
		return zero, err, nil
	})
}

func (rpc *testRPC[Req, Res]) Returns(res Res) {
	rpc.responses = rpc.responses.Append(func(ctx context.Context, req Req) (Res, error, rpcRes[Req, Res]) {
		return res, nil, nil
	})
}

func (rpc *testRPC[Req, Res]) ReturnsFunc(f func(ctx context.Context, req Req) (Res, error)) {
	rpc.responses = rpc.responses.Append(func(ctx context.Context, req Req) (Res, error, rpcRes[Req, Res]) {
		r, err := f(ctx, req)
		return r, err, nil
	})
}

func (rpc *testRPC[Req, Res]) AssertCalled(fn func(ctx context.Context, req Req)) {
	rpc.t.Helper()
	if len(rpc.requests) == 0 {
		rpc.t.Errorf("expected %s to be called, but it was not", rpc.name)
		return
	}
	req := rpc.requests[0]
	rpc.requests = rpc.requests[1:]
	fn(req.ctx, req.req)
}

func (rpc *testRPC[Req, Res]) AssertCalledWith(req Req) {
	rpc.t.Helper()
	rpc.AssertCalled(func(ctx context.Context, got Req) {
		rpc.t.Helper()
		if diff := cmp.Diff(req, got, protocmp.Transform()); diff != "" {
			rpc.t.Errorf("unexpected %s request (-want +got):\n%s", rpc.name, diff)
		}
	})
}

func (rpc *testRPC[Req, Res]) AssertNoMoreCalls() {
	rpc.t.Helper()
	if len(rpc.requests) != 0 {
		rpc.t.Errorf("expected %s to not be called, but it was", rpc.name)
	}
}

func (rpc *testRPC[Req, Res]) nextResponse(ctx context.Context, req Req) (Res, error) {
	rpc.t.Helper()
	if rpc.responses == nil {
		rpc.t.Errorf("unexpected request %s(%v, %v)", rpc.name, ctx, req)
		var zero Res
		return zero, fmt.Errorf("unexpected call to %s", rpc.name)
	}

	res, err, next := rpc.responses(ctx, req)
	rpc.responses = next
	return res, err

}

type rpcArgs[T proto.Message] struct {
	t   *testing.T
	ctx context.Context
	req T
}

type rpcRes[Req proto.Message, Res proto.Message] func(ctx context.Context, req Req) (res Res, err error, next rpcRes[Req, Res])

func (rpc rpcRes[Req, Res]) Append(next rpcRes[Req, Res]) rpcRes[Req, Res] {
	cur := rpc
	var f rpcRes[Req, Res]
	f = func(ctx context.Context, req Req) (Res, error, rpcRes[Req, Res]) {
		if cur == nil {
			return next(ctx, req)
		}
		var (
			res Res
			err error
		)
		res, err, cur = cur(ctx, req)
		return res, err, f
	}
	return f
}

type testServerTransportStream struct {
	method  string
	header  metadata.MD
	trailer metadata.MD
	sent    bool
}

func newTestServerTransportStream(ctx context.Context) (*testServerTransportStream, context.Context) {
	stream := &testServerTransportStream{}
	ctx = grpc.NewContextWithServerTransportStream(ctx, stream)
	return stream, ctx
}

func (t *testServerTransportStream) Method() string {
	return t.method
}

func (t *testServerTransportStream) SetHeader(md metadata.MD) error {
	t.header = metadata.Join(t.header, md)
	return nil
}

func (t *testServerTransportStream) SendHeader(md metadata.MD) error {
	t.header = metadata.Join(t.header, md)
	t.sent = true
	return nil
}

func (t *testServerTransportStream) SetTrailer(md metadata.MD) error {
	t.trailer = metadata.Join(t.trailer, md)
	return nil
}

func (t testServerTransportStream) HasHeader(key, val string) ([]string, bool) {
	vals := t.header.Get(key)
	for _, v := range vals {
		if v == val {
			return nil, true
		}
	}
	return vals, false
}
