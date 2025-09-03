package healthpb

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/vanti-dev/sc-bos/internal/manage/devices"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/devicespb"
)

// ExampleRegistry_devicesApi shows how to implement the [gen.DevicesApiServer] using a [Registry].
func ExampleRegistry_devicesApi() {
	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})

	// todo: make this more ergonomic
	registry := &Registry{
		onCheckUpdate: func(name string, c *gen.HealthCheck) {
			_, _ = devs.Update(&gen.Device{Name: name}, resource.WithMerger(func(mask *masks.FieldUpdater, dst, src proto.Message) {
				dstDev := dst.(*gen.Device)
				dstDev.HealthChecks = mergeChecks(mask.Merge, dstDev.HealthChecks, c)
			}), resource.WithCreateIfAbsent())
		},
		onCheckDelete: func(name, id string) {
			_, _ = devs.Update(&gen.Device{Name: name}, resource.WithMerger(func(mask *masks.FieldUpdater, dst, src proto.Message) {
				dstDev := dst.(*gen.Device)
				dstDev.HealthChecks = removeCheck(dstDev.HealthChecks, id)
			}), resource.WithAllowMissing(true))
		},
	}
	exampleChecks := registry.ForOwner("example")

	// create the device with some metadata
	_, _ = devs.Update(&gen.Device{Name: "device1", Metadata: &traits.Metadata{
		Appearance: &traits.Metadata_Appearance{Title: "Example Device 1"},
	}}, resource.WithCreateIfAbsent())
	// prepare a health check for the device
	dev1Check, _ := exampleChecks.NewErrorCheck("device1", &gen.HealthCheck{
		DisplayName: "Is it working",
	})
	// report on the health of the device
	dev1Check.UpdateError(errors.New("malfunction"))

	// use the server, typically via gRPC from the client
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))
	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	if err != nil {
		panic(err)
	}
	for _, d := range deviceList.Devices {
		var states []gen.HealthCheck_Check_State
		for _, c := range d.GetHealthChecks() {
			states = append(states, c.GetCheck().GetState())
		}
		fmt.Printf("Device %q is %v", d.GetMetadata().GetAppearance().GetTitle(), states)
	}
	// Output: Device "Example Device 1" is [ABNORMAL]
}

type devicesServerModel struct {
	devices.Collection
}

func (m devicesServerModel) ClientConn() grpc.ClientConnInterface {
	return nil
}
