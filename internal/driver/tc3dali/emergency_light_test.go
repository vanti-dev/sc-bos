package tc3dali

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestEmergencyLightServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	mock := dali.NewMock(logger)
	server := &emergencyLightServer{
		bus:       mock,
		shortAddr: 0,
	}
	client := rpc.WrapDaliApi(server)

	// Check initial state
	actualStatus, err := client.GetEmergencyStatus(context.Background(), &rpc.GetEmergencyStatusRequest{})
	if err != nil {
		t.Fatal(err)
	}
	expectedStatus := &rpc.EmergencyStatus{
		ActiveModes: []rpc.EmergencyStatus_Mode{
			rpc.EmergencyStatus_NORMAL,
		},
		BatteryLevelPercent: 100.0,
	}

	diff := cmp.Diff(expectedStatus, actualStatus, protocmp.Transform())
	if diff != "" {
		t.Errorf("unexpected emergency status (-want +got):\n%s\n", diff)
	}

	// Expect no function duration test result available
	actualResult, err := client.GetTestResult(context.Background(), &rpc.GetTestResultRequest{Test: rpc.Test_DURATION_TEST})
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("expected NotFound error, but got %v with result\n%s", err, protojson.Format(actualResult))
	}

	// Start duration test & check that it's running
	_, err = client.StartTest(context.Background(), &rpc.StartTestRequest{Test: rpc.Test_DURATION_TEST})
	if err != nil {
		t.Errorf("StartTest: %v", err)
	}
	actualStatus, err = client.GetEmergencyStatus(context.Background(), &rpc.GetEmergencyStatusRequest{})
	if err != nil {
		t.Fatalf("GetEmergencyStatus: %v", err)
	}
	if !contains(rpc.EmergencyStatus_DURATION_TEST_ACTIVE, actualStatus.ActiveModes) {
		t.Fatalf("DURATION_TEST_ACTIVE not in ActiveModes; found %v", actualStatus.ActiveModes)
	}

	// mark test as completed successfully, check this is reflected
	const testResult = 123
	if ok := mock.ControlGear[0].CompleteDurationTest(true, testResult); !ok {
		t.Fatal("could not mark duration test as complete")
	}
	actualResult, err = client.GetTestResult(context.Background(), &rpc.GetTestResultRequest{Test: rpc.Test_DURATION_TEST})
	if err != nil {
		t.Fatalf("GetTestResult: %v", err)
	}
	expectedResult := &rpc.TestResult{
		Test:     rpc.Test_DURATION_TEST,
		Pass:     true,
		Duration: durationpb.New(testResult * 2 * time.Minute),
	}
	if diff := cmp.Diff(expectedResult, actualResult, protocmp.Transform()); diff != "" {
		t.Errorf("unexpected TestResult (-want +got):\n%s", diff)
	}

	// clear the test result, check this is reflected
	_, err = client.DeleteTestResult(context.Background(), &rpc.DeleteTestResultRequest{Test: rpc.Test_DURATION_TEST})
	if err != nil {
		t.Fatalf("DeleteTestResult: %v", err)
	}
	if mock.ControlGear[0].TestInProgress() {
		t.Error("Test is still running")
	}
	actualResult, err = client.GetTestResult(context.Background(), &rpc.GetTestResultRequest{Test: rpc.Test_DURATION_TEST})
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("expected NotFound error, but got %v with result\n%s", err, protojson.Format(actualResult))
	}
}

func contains[T comparable](target T, in []T) bool {
	for _, item := range in {
		if item == target {
			return true
		}
	}
	return false
}
