package helvarnet

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-bos/internal/manage/devices"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/devicespb"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/healthpb"
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

func setup(devs *devicespb.Collection) *healthpb.Registry {

	return healthpb.NewRegistry(
		healthpb.WithOnCheckCreate(func(name string, c *gen.HealthCheck) *gen.HealthCheck {
			_, _ = devs.Update(&gen.Device{Name: name}, resource.WithMerger(func(mask *masks.FieldUpdater, dst, src proto.Message) {
				dstDev := dst.(*gen.Device)
				dstDev.HealthChecks = healthpb.MergeChecks(mask.Merge, dstDev.HealthChecks, c)
			}), resource.WithCreateIfAbsent(), resource.WithExpectAbsent())
			return nil
		}),
		healthpb.WithOnCheckUpdate(func(name string, c *gen.HealthCheck) {
			_, _ = devs.Update(&gen.Device{Name: name}, resource.WithMerger(func(mask *masks.FieldUpdater, dst, src proto.Message) {
				dstDev := dst.(*gen.Device)
				dstDev.HealthChecks = healthpb.MergeChecks(mask.Merge, dstDev.HealthChecks, c)
			}))
		}),
		healthpb.WithOnCheckDelete(func(name, id string) {
			_, _ = devs.Update(&gen.Device{Name: name}, resource.WithMerger(func(mask *masks.FieldUpdater, dst, src proto.Message) {
				dstDev := dst.(*gen.Device)
				dstDev.HealthChecks = healthpb.RemoveCheck(dstDev.HealthChecks, id)
			}), resource.WithAllowMissing(true))
		}),
	)
}

func TestDeviceNoError(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	status := int64(0)
	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))
	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, 0, len(checks[0].GetFaults().CurrentFaults))
	}
}

func TestDeviceDeviceOffline(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	status := int64(DeviceOfflineCode)
	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))
	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, strconv.Itoa(DeviceOfflineCode), checks[0].GetFaults().CurrentFaults[0].Code.Code)
	}
}

func TestDeviceBadResponse(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	status := int64(BadResponseCode)
	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))
	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, strconv.Itoa(BadResponseCode), checks[0].GetFaults().CurrentFaults[0].Code.Code)
	}
}

func TestSingleHelvarnetFault(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	status := int64(0x00000001)
	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))
	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, strconv.Itoa(0x00000001), checks[0].GetFaults().CurrentFaults[0].Code.Code)
	}
}

func TestDoubleHelvarnetFault(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	status := int64(0x00000011)
	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))
	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, 2, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, strconv.Itoa(0x00000001), checks[0].GetFaults().CurrentFaults[0].Code.Code)
		require.Equal(t, strconv.Itoa(0x00000010), checks[0].GetFaults().CurrentFaults[1].Code.Code)
	}
}

func TestAddFaultThenClear(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))

	// Step 1: Add a fault
	status := int64(0x00000001) // Disabled
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, strconv.Itoa(0x00000001), checks[0].GetFaults().CurrentFaults[0].Code.Code)
	}

	// Verify fault is tracked
	require.True(t, raisedFaults[0x00000001])

	// Step 2: Clear the fault
	status = int64(0)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err = client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_NORMAL, checks[0].Normality)
		require.Equal(t, 0, len(checks[0].GetFaults().CurrentFaults))
	}

	// Verify fault tracking is cleared
	require.False(t, raisedFaults[0x00000001])
}

func TestAddMultipleFaultsThenClear(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))

	// Step 1: Add multiple faults
	status := int64(0x00000001 | 0x00000002 | 0x00000004) // Disabled | LampFailure | Missing
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)
		require.Equal(t, 3, len(checks[0].GetFaults().CurrentFaults))
	}

	// Verify all faults are tracked
	require.True(t, raisedFaults[0x00000001])
	require.True(t, raisedFaults[0x00000002])
	require.True(t, raisedFaults[0x00000004])

	// Step 2: Clear all faults
	status = int64(0)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err = client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_NORMAL, checks[0].Normality)
		require.Equal(t, 0, len(checks[0].GetFaults().CurrentFaults))
	}

	// Verify all fault tracking is cleared
	require.False(t, raisedFaults[0x00000001])
	require.False(t, raisedFaults[0x00000002])
	require.False(t, raisedFaults[0x00000004])
}

func TestAddFaultThenPartialClear(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))

	// Step 1: Add multiple faults
	status := int64(0x00000001 | 0x00000002 | 0x00000004) // Disabled | LampFailure | Missing
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)
		require.Equal(t, 3, len(checks[0].GetFaults().CurrentFaults))
	}

	// Step 2: Clear some faults but keep one
	status = int64(0x00000001) // Keep only Disabled
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err = client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, strconv.Itoa(0x00000001), checks[0].GetFaults().CurrentFaults[0].Code.Code)
	}

	// Verify fault tracking is updated correctly
	require.True(t, raisedFaults[0x00000001])
	require.False(t, raisedFaults[0x00000002])
	require.False(t, raisedFaults[0x00000004])

	// Step 3: Now clear the remaining fault
	status = int64(0)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err = client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_NORMAL, checks[0].Normality)
		require.Equal(t, 0, len(checks[0].GetFaults().CurrentFaults))
	}

	// Verify all fault tracking is cleared
	require.False(t, raisedFaults[0x00000001])
}

func TestDeviceOfflineThenClear(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))

	// Step 1: Set device offline
	status := int64(DeviceOfflineCode)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, "Device Offline", checks[0].GetFaults().CurrentFaults[0].SummaryText)
		require.Equal(t, gen.HealthCheck_Reliability_NO_RESPONSE, checks[0].Reliability.State)
	}

	// Step 2: Device comes back online (status 0 = no faults)
	status = int64(0)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err = client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_NORMAL, checks[0].Normality)
		require.Equal(t, 0, len(checks[0].GetFaults().CurrentFaults))
	}
}

func TestBadResponseThenClear(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))

	// Step 1: Set bad response
	status := int64(BadResponseCode)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)
		require.Equal(t, 1, len(checks[0].GetFaults().CurrentFaults))
		require.Equal(t, "Bad Response", checks[0].GetFaults().CurrentFaults[0].SummaryText)
		require.Equal(t, gen.HealthCheck_Reliability_BAD_RESPONSE, checks[0].Reliability.State)
	}

	// Step 2: Device recovers (status 0 = no faults)
	status = int64(0)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err = client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_NORMAL, checks[0].Normality)
		require.Equal(t, 0, len(checks[0].GetFaults().CurrentFaults))
	}
}

func TestUnknownNegativeStatus(t *testing.T) {

	devs := devicespb.NewCollection()
	server := devices.NewServer(devicesServerModel{Collection: devs})
	deviceName := "helvarnet-device-1"
	reg := setup(devs)
	exampleChecks := reg.ForOwner("example")

	_, _ = devs.Update(&gen.Device{Name: deviceName}, resource.WithCreateIfAbsent())

	check := getDeviceHealthCheck()
	fc, err := exampleChecks.NewFaultCheck(deviceName, check)
	require.NoError(t, err)
	defer fc.Dispose()

	raisedFaults := make(map[int64]bool)
	ctx := context.Background()
	client := gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, server))

	// Test with an unknown negative status code (not DeviceOfflineCode or BadResponseCode)
	status := int64(-99)
	updateDeviceFaults(ctx, status, fc, raisedFaults)

	deviceList, err := client.ListDevices(context.TODO(), &gen.ListDevicesRequest{})
	require.NoError(t, err)

	for _, d := range deviceList.Devices {
		checks := d.GetHealthChecks()
		require.Equal(t, 1, len(checks))
		require.Equal(t, gen.HealthCheck_ABNORMAL, checks[0].Normality)

		// Should have created an "Internal Driver Error" fault
		faults := checks[0].GetFaults().CurrentFaults
		require.Equal(t, 1, len(faults))
		require.Equal(t, "Internal Driver Error", faults[0].SummaryText)
		require.Equal(t, "The device has an unrecognised internal status code", faults[0].DetailsText)
		require.Equal(t, strconv.Itoa(UnrecognisedErrorCode), faults[0].Code.Code)
		require.Equal(t, SystemName, faults[0].Code.System)

		// Should have set reliability to UNRELIABLE
		require.Equal(t, gen.HealthCheck_Reliability_UNRELIABLE, checks[0].Reliability.State)
		require.NotNil(t, checks[0].Reliability.UnreliableTime)
	}
}

type devicesServerModel struct {
	devices.Collection
}

func (m devicesServerModel) ClientConn() grpc.ClientConnInterface {
	return nil
}
