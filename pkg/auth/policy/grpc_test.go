package policy

import (
	"context"
	"net"
	"testing"

	"github.com/open-policy-agent/opa/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
)

func TestInterceptor(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)

	compiler, err := ast.CompileModules(regoFiles)
	if err != nil {
		t.Fatal(err)
	}
	interceptor := NewInterceptor(&static{compiler: compiler})
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.GRPCUnaryInterceptor()),
		grpc.ChainStreamInterceptor(interceptor.GRPCStreamingInterceptor()),
	)
	traits.RegisterOnOffApiServer(server, onoff.NewModelServer(onoff.NewModel()))
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Errorf("server stopped with error: %v", err)
		}
	}()

	t.Cleanup(func() {
		if err := lis.Close(); err != nil {
			t.Logf("failed to close listener: %v", err)
		}
		server.Stop()
	})

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}

	client := traits.NewOnOffApiClient(conn)

	// check simple name based auth, global for all smartcore.* apis
	_, err = client.GetOnOff(ctx, &traits.GetOnOffRequest{Name: "allow"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, err = client.GetOnOff(ctx, &traits.GetOnOffRequest{Name: "deny"})
	if err == nil {
		t.Error("expected error")
	}
	if c := status.Code(err); c != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied, got %v", err)
	}

	// check action based auth, specific to this trait
	_, err = client.UpdateOnOff(ctx, &traits.UpdateOnOffRequest{Name: "any", OnOff: &traits.OnOff{State: traits.OnOff_ON}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, err = client.UpdateOnOff(ctx, &traits.UpdateOnOffRequest{Name: "any", OnOff: &traits.OnOff{State: traits.OnOff_OFF}})
	if err == nil {
		t.Error("expected error")
	}
	if c := status.Code(err); c != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied, got %v", err)
	}
}

var regoFiles = map[string]string{
	"smartcore.rego": `package smartcore

# This simple rule allows any request whose name is "allow", all other requests are denied
allow {
	input.request.name == "allow"
}
`,
	"smartcore.traits.OnOffApi.rego": `package smartcore.traits.OnOffApi

# This rule allows people to turn any device on (but not off)
allow {
	input.method == "UpdateOnOff"
	input.request.onOff.state == "ON"
}
`,
}
