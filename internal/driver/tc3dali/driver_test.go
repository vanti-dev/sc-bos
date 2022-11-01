package tc3dali

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestNetID_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`"1.2.3.4.5.6"`)
	var netID NetID
	err := json.Unmarshal(jsonData, &netID)
	if err != nil {
		t.Error(err)
	}
	expect := NetID{1, 2, 3, 4, 5, 6}
	if diff := cmp.Diff(expect, netID); diff != "" {
		t.Errorf("mismatched NetID (-want +got):\n%s", diff)
	}
}

func TestNetID_MarshalJSON(t *testing.T) {
	netID := NetID{6, 5, 4, 3, 2, 1}
	got, err := json.Marshal(netID)
	if err != nil {
		t.Error(err)
	}
	expect := []byte(`"6.5.4.3.2.1"`)
	if diff := cmp.Diff(expect, got); diff != "" {
		t.Errorf("mismatched JSON (-want +got):\n%s", diff)
	}
}

func TestInitBus(t *testing.T) {
	services := driver.Services{
		Logger: zap.NewNop(),
		Node:   node.New("TestFactory"),
		Tasks:  &task.Group{},
	}
	services.Node.Support(node.Routing(
		parent.NewApiRouter(),
		light.NewApiRouter(),
		occupancysensor.NewApiRouter(),
	))

	config := BusConfig{
		Name: "TestFactory/dali/bus",
		ControlGear: []ControlGearConfig{
			{ShortAddress: 1, Groups: []uint8{0, 1, 2}},
			{ShortAddress: 2, Groups: []uint8{4, 3, 2}},
		},
		ControlDevices: []ControlDeviceConfig{
			{ShortAddress: 3, InstanceTypes: []InstanceType{InstanceTypeOccupancySensor}},
		},
	}

	mockDali := dali.NewMock(zap.NewNop())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := InitBus(ctx, config, mockDali, services)
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Errorf("unexpected InitBus error: %s", err.Error())
		}
	}()

	server := grpc.NewServer()
	defer server.Stop()
	services.Node.Register(server)
	listener := bufconn.Listen(1024 * 1024)
	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("mock server stopped with error: %v", err)
		}
		_ = listener.Close()
	}()

	conn, err := grpc.Dial("",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}

	parentClient := traits.NewParentApiClient(conn)
	res, err := parentClient.ListChildren(context.Background(), &traits.ListChildrenRequest{
		Name:     services.Node.Name(),
		PageSize: 1000,
	})
	if err != nil {
		t.Fatalf("failed to list children: %v", err)
	}

	expected := []*traits.Child{
		{Name: "TestFactory/dali/bus", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/control-device/3", Traits: []*traits.Trait{{Name: string(trait.OccupancySensor)}}},
		{Name: "TestFactory/dali/bus/control-gear/1", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/control-gear/2", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/groups/0", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/groups/1", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/groups/2", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/groups/3", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
		{Name: "TestFactory/dali/bus/groups/4", Traits: []*traits.Trait{{Name: string(trait.Light)}}},
	}
	diff := cmp.Diff(expected, res.Children,
		protocmp.Transform(),
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(x, y *traits.Child) bool {
			return x.GetName() < y.GetName()
		}),
	)
	if diff != "" {
		t.Errorf("ListChildren result unexpected (-want +got):\n%s", diff)
	}
}
