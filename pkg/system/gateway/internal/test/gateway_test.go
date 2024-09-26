package test

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/internal/util/grpc/reflectionapi"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/system/gateway/internal/test/shared"
	// make sure that test caching updates based on changes to these files too
	_ "github.com/vanti-dev/sc-bos/pkg/system/gateway/internal/test/ac"
	_ "github.com/vanti-dev/sc-bos/pkg/system/gateway/internal/test/gw"
	_ "github.com/vanti-dev/sc-bos/pkg/system/gateway/internal/test/hub"
)

var skipBuild = flag.Bool("skip-build", false, "skip building and running binaries")
var ignoreEnrolErr = flag.Bool("ignore-enrol-err", false, "ignore enrolment errors")

// TestGateway_e2e tests the gateway by running a cohort of nodes, each a different sc-bos process.
// The test only runs if the -short flag is not set.
func TestGateway_e2e(t *testing.T) {
	// WARNING: This test doesn't play perfectly with go tests caching for a number of reasons:
	// 1. it builds go binaries which are the target of the tests which don't get checked as part of cache invalidation
	// 2. those binaries read files that are also not part of the cache invalidation

	if testing.Short() {
		t.Skip("long test")
	}

	ctx, stop := newCtx(t)
	defer stop()

	if !*skipBuild {
		dir := t.TempDir()

		// First, we need to build each of the different binaries that make up the nodes in the cohort.
		// After this completes we'll have a gw, ac, and hub binary in the tests temp directory.
		t.Logf("Building binaries in %s", dir)
		buildAll(t, dir)

		// Next we start each of the nodes we need for the test
		startCtx, cancelStart := context.WithTimeout(ctx, 30*time.Second)
		defer cancelStart()
		t.Logf("Starting all nodes")
		go runAllNodes(t, startCtx, dir)

		// Wait for the nodes to start up, shouldn't take more than a few seconds on _decent_ hardware.
		t.Logf("Waiting for nodes to start")
		waitForNodes(t, startCtx)
	}

	// Next up we need to configure the cohort
	t.Logf("Configuring cohort")
	configureCohort(t, ctx)

	// Finally we're ready to start checking the setup
	for i, addr := range shared.GWGRPCAddrs {
		addr := addr
		t.Run(fmt.Sprintf("gw%d %s", i+1, addr), func(t *testing.T) {
			// this timeout is long because the GW is using an exponential backoff for retries,
			// capped at 30s, but all attempts before the cohort is configured increase the delay.
			testCtx, stopTests := context.WithTimeout(ctx, 60*time.Second)
			defer stopTests()
			// these func log themselves
			testGW(t, testCtx, addr)
		})
	}
}

func buildAll(t *testing.T, dir string) {
	t.Helper()

	ctx, stop := newCtx(t)
	defer stop()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return build(t, ctx, "gw", dir) })
	g.Go(func() error { return build(t, ctx, "ac", dir) })
	g.Go(func() error { return build(t, ctx, "hub", dir) })
	if err := g.Wait(); err != nil {
		t.Fatal("build failed", err)
	}
}

func build(t *testing.T, ctx context.Context, name, dir string) error {
	t.Helper()

	build := exec.CommandContext(ctx, "go", "build", "-o", filepath.Join(dir, name), "./"+name+"/cmd")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	return build.Run()
}

func runAllNodes(t *testing.T, ctx context.Context, dir string) {
	t.Helper()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return runNode(t, ctx, dir, "hub", shared.HubGRPCAddr, shared.HubHTTPSAddr) })
	for i, addrs := range zip(shared.ACGRPCAddrs, shared.ACHTTPSAddrs) {
		i := i
		addrs := addrs
		g.Go(func() error { return runNode(t, ctx, dir, fmt.Sprintf("ac%d", i+1), addrs[0], addrs[1]) })
	}
	for i, addrs := range zip(shared.GWGRPCAddrs, shared.GWHTTPSAddrs) {
		i := i
		addrs := addrs
		g.Go(func() error { return runNode(t, ctx, dir, fmt.Sprintf("gw%d", i+1), addrs[0], addrs[1]) })
	}
	if err := g.Wait(); err != nil {
		select {
		case <-ctx.Done():
			return
		default:
		}
		t.Fatal("run failed", err)
	}
}

func runNode(t *testing.T, ctx context.Context, dir, name, grpcAddr, httpsAddr string) error {
	t.Helper()

	execName := strings.TrimRight(name, "1234567890")

	node := exec.CommandContext(ctx, filepath.Join(dir, execName),
		"--listen-grpc", grpcAddr,
		"--listen-https", httpsAddr,
		"--policy-mode", "off", // disable policy checking for now
		"--sysconf", filepath.Join("testdata", name, "system.conf.json"),
		"--appconf", filepath.Join("testdata", name, "app.conf.json"),
		"--data", filepath.Join(t.TempDir(), name+"-data"),
	)
	node.Stdout = os.Stdout
	node.Stderr = os.Stderr
	return node.Run()
}

func waitForNodes(t *testing.T, ctx context.Context) {
	t.Helper()

	waitForNode(t, ctx, shared.HubGRPCAddr)
	for _, addr := range shared.ACGRPCAddrs {
		waitForNode(t, ctx, addr)
	}
	for _, addr := range shared.GWGRPCAddrs {
		waitForNode(t, ctx, addr)
	}
}

func waitForNode(t *testing.T, ctx context.Context, addr string) {
	t.Helper()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer conn.Close()

	client := traits.NewMetadataApiClient(conn)
	err = backoff.Retry(func() error {
		_, err := client.GetMetadata(ctx, &traits.GetMetadataRequest{})
		if err != nil {
			t.Logf("failed to poll node %q for liveness: %v", addr, err)
		}
		return err
	}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		t.Fatalf("wait for node %s: %v", addr, err)
	}
}

func configureCohort(t *testing.T, ctx context.Context) {
	t.Helper()

	// todo: use the hubs ca (should be in dir, after the first request) for our client cert checks

	hubConn, err := grpc.DialContext(ctx, shared.HubGRPCAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer hubConn.Close()

	checkErr := func(addr string, err error) {
		t.Helper()
		if *ignoreEnrolErr {
			return
		}
		if err != nil {
			t.Fatalf("enroll %s: %v", addr, err)
		}
	}

	client := gen.NewHubApiClient(hubConn)
	for _, addr := range shared.ACGRPCAddrs {
		_, err := client.EnrollHubNode(ctx, &gen.EnrollHubNodeRequest{Node: &gen.HubNode{Address: addr}})
		checkErr(addr, err)
	}
	for _, addr := range shared.GWGRPCAddrs {
		_, err := client.EnrollHubNode(ctx, &gen.EnrollHubNodeRequest{Node: &gen.HubNode{Address: addr}})
		checkErr(addr, err)
	}
}

func testGW(t *testing.T, ctx context.Context, addr string) {
	t.Helper()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		t.Fatalf("dial %s: %v", addr, err)
	}
	defer conn.Close()

	// Named devices are correctly routed
	deviceNames := []string{
		"ac1/dev1",
		"ac2/dev1",
		"hub/dev1",
	}
	t.Logf("[%s] Waiting for gw to configure gateway system", addr)
	for _, name := range deviceNames {
		waitForDevice(t, ctx, conn, name)
	}

	t.Run("named devices", func(t *testing.T) {
		client := traits.NewOnOffApiClient(conn)
		for _, name := range deviceNames {
			testOnOffApi(t, ctx, addr, name, client)
		}
	})

	t.Run("reflection", func(t *testing.T) {
		testReflection(t, ctx, conn)
	})

	testHubApis(t, ctx, conn)
}

func waitForDevice(t *testing.T, ctx context.Context, conn *grpc.ClientConn, name string) {
	t.Helper()

	client := traits.NewMetadataApiClient(conn)
	err := backoff.Retry(func() error {
		_, err := client.GetMetadata(ctx, &traits.GetMetadataRequest{Name: name})
		return err
	}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		t.Fatalf("wait for device %s: %v", name, err)
	}
}

func testOnOffApi(t *testing.T, ctx context.Context, addr, name string, client traits.OnOffApiClient) {
	t.Helper()

	// useful for cancelling the stream
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// set initial known state: ON
	res, err := client.UpdateOnOff(ctx, &traits.UpdateOnOffRequest{Name: name, OnOff: &traits.OnOff{State: traits.OnOff_ON}})
	if err != nil {
		t.Fatalf("[%s] update onoff %s: %v", addr, name, err)
	}
	if diff := cmp.Diff(&traits.OnOff{State: traits.OnOff_ON}, res, protocmp.Transform()); diff != "" {
		t.Fatalf("[%s] update onoff %s: unexpected response (-want +got):\n%s", addr, name, diff)
	}

	// subscribe
	changes := make(chan *traits.PullOnOffResponse, 1) // we're only expecting 1
	stream, err := client.PullOnOff(ctx, &traits.PullOnOffRequest{Name: name, UpdatesOnly: true})
	if err != nil {
		t.Fatalf("[%s] pull onoff %s: %v", addr, name, err)
	}
	go func() {
		for {
			res, err := stream.Recv()
			if errors.Is(err, io.EOF) || status.Code(err) == codes.Canceled {
				close(changes)
				return
			}
			if err != nil {
				t.Errorf("[%s] pull onoff %s: %v", addr, name, err)
				return
			}
			changes <- res
		}
	}()

	// check initial state
	res, err = client.GetOnOff(ctx, &traits.GetOnOffRequest{Name: name})
	if err != nil {
		t.Fatalf("[%s] get onoff %s: %v", addr, name, err)
	}
	if diff := cmp.Diff(&traits.OnOff{State: traits.OnOff_ON}, res, protocmp.Transform()); diff != "" {
		t.Fatalf("[%s] get onoff %s: unexpected response (-want +got):\n%s", addr, name, diff)
	}

	// Perform update and check for change
	res, err = client.UpdateOnOff(ctx, &traits.UpdateOnOffRequest{Name: name, OnOff: &traits.OnOff{State: traits.OnOff_OFF}})
	if err != nil {
		t.Fatalf("[%s] update onoff %s: %v", addr, name, err)
	}
	if diff := cmp.Diff(&traits.OnOff{State: traits.OnOff_OFF}, res, protocmp.Transform()); diff != "" {
		t.Fatalf("[%s] update onoff %s: unexpected response (-want +got):\n%s", addr, name, diff)
	}
	select {
	case res := <-changes:
		want := &traits.PullOnOffResponse{Changes: []*traits.PullOnOffResponse_Change{
			{
				Name:  name,
				OnOff: &traits.OnOff{State: traits.OnOff_OFF},
			},
		}}
		// clear timestamps to make comparing easier
		for i := range res.Changes {
			res.Changes[i].ChangeTime = nil
		}
		if diff := cmp.Diff(want, res, protocmp.Transform()); diff != "" {
			t.Fatalf("[%s] pull onoff %s: unexpected response (-want +got):\n%s", addr, name, diff)
		}
	}
}

func testReflection(t *testing.T, ctx context.Context, conn *grpc.ClientConn) {
	ctx, stop := context.WithCancel(ctx)
	defer stop()

	client := reflectionpb.NewServerReflectionClient(conn)
	stream, err := client.ServerReflectionInfo(ctx)
	if err != nil {
		t.Fatal("server reflection info:", err)
	}

	services, err := reflectionapi.ListServices(stream)
	if err != nil {
		t.Fatal("list services:", err)
	}
	wantServices := []*reflectionpb.ServiceResponse{
		{Name: "grpc.reflection.v1.ServerReflection"},
		{Name: "grpc.reflection.v1alpha.ServerReflection"},
		{Name: "smartcore.bos.DevicesApi"},
		{Name: "smartcore.bos.EnrollmentApi"},
		{Name: "smartcore.bos.HubApi"},
		{Name: "smartcore.bos.ServicesApi"},
		{Name: "smartcore.traits.MetadataApi"},
		{Name: "smartcore.traits.MetadataInfo"},
		{Name: "smartcore.traits.OnOffApi"},
		{Name: "smartcore.traits.OnOffInfo"},
		{Name: "smartcore.traits.ParentApi"},
		{Name: "smartcore.traits.ParentInfo"},
	}
	if diff := cmp.Diff(wantServices, services, protocmp.Transform()); diff != "" {
		t.Fatalf("services: (-want +got):\n%s", diff)
	}

	types := []string{
		"smartcore.traits.OnOffApi",
		"smartcore.bos.DevicesApi",
	}
	for _, typ := range types {
		_, err = reflectionapi.FileContainingSymbol(stream, typ)
		if err != nil {
			t.Fatalf("file containing symbol %s: %v", typ, err)
		}
	}

	unknownTypes := []string{
		"smartcore.traits.UnknownApiForTesting", // doesn't exist
		// note, all apis that are in the traits or gen packages get loaded at the same time (during package init)
	}
	for _, typ := range unknownTypes {
		_, err = reflectionapi.FileContainingSymbol(stream, typ)
		if status.Code(err) != codes.NotFound {
			t.Fatalf("file containing symbol %s: expected error, got %v", typ, err)
		}
	}
}

func testHubApis(t *testing.T, ctx context.Context, conn *grpc.ClientConn) {
	t.Helper()

	t.Run("HubApi", func(t *testing.T) {
		client := gen.NewHubApiClient(conn)
		res, err := client.ListHubNodes(ctx, &gen.ListHubNodesRequest{})
		if err != nil {
			t.Fatalf("list hub nodes: %v", err)
		}
		wantNames := []string{"ac1", "ac2", "gw1", "gw2"}
		gotNames := make([]string, len(res.Nodes))
		for i, node := range res.Nodes {
			gotNames[i] = node.Name
		}
		sortStrings := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(wantNames, gotNames, sortStrings); diff != "" {
			t.Fatalf("list hub nodes: unexpected response (-want +got):\n%s", diff)
		}
	})

}

func newCtx(t *testing.T) (context.Context, context.CancelFunc) {
	deadline, ok := t.Deadline()
	if !ok {
		return context.WithCancel(context.Background())
	}
	return context.WithDeadline(context.Background(), deadline)
}

func zip[T any](a, b []T) [][2]T {
	if len(a) != len(b) {
		panic("zip: slices have different lengths")
	}
	z := make([][2]T, len(a))
	for i := range a {
		z[i] = [2]T{a[i], b[i]}
	}
	return z
}
