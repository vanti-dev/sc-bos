package tc3dali

import (
	"context"
	"fmt"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali/dali202"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type emergencyLightServer struct {
	rpc.UnimplementedDaliApiServer

	bus       dali.Dali
	shortAddr uint8
}

func (s *emergencyLightServer) Identify(ctx context.Context, request *rpc.IdentifyRequest) (*rpc.IdentifyResponse, error) {
	_, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.StartIdentification202,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, err
	}
	return &rpc.IdentifyResponse{}, nil
}

func (s *emergencyLightServer) GetEmergencyStatus(ctx context.Context, request *rpc.GetEmergencyStatusRequest) (*rpc.EmergencyStatus, error) {
	rawEmergencyStatus, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryEmergencyStatus,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryEmergencyStatus: %w", err)
	}

	rawEmergencyMode, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryEmergencyMode,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryEmergencyMode: %w", err)
	}

	rawFailure, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryFailureStatus,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryFailureStatus: %w", err)
	}

	rawBatteryLevel, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryBatteryCharge,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryBatteryCharge: %w", err)
	}

	emergencyStatus := decodeEmergencyStatus(uint8(rawEmergencyStatus), uint8(rawEmergencyMode), uint8(rawBatteryLevel),
		uint8(rawFailure))
	return emergencyStatus, nil
}

func (s *emergencyLightServer) StartTest(ctx context.Context, request *rpc.StartTestRequest) (*rpc.StartTestResponse, error) {
	var command dali.Command
	switch request.GetTest() {
	case rpc.Test_FUNCTION_TEST:
		command = dali.StartFunctionTest
	case rpc.Test_DURATION_TEST:
		command = dali.StartDurationTest
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid test specified")
	}

	_, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     command,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, err
	}
	return &rpc.StartTestResponse{}, nil
}

func (s *emergencyLightServer) StopTest(ctx context.Context, request *rpc.StopTestRequest) (*rpc.StopTestResponse, error) {
	_, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.StopTest,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, err
	}
	return &rpc.StopTestResponse{}, nil
}

var errInvalidTest = status.Error(codes.InvalidArgument, "invalid test specified")

func (s *emergencyLightServer) GetTestResult(ctx context.Context, request *rpc.GetTestResultRequest) (*rpc.TestResult, error) {
	// set up some data for the following requests, based on what kind of test result is needed
	var (
		doneMask        uint8
		failureMask     uint8
		requestDuration bool
	)
	switch request.GetTest() {
	case rpc.Test_FUNCTION_TEST:
		doneMask = dali202.StatusBitFunctionTestDone
		failureMask = dali202.FailureBitFunctionTest
	case rpc.Test_DURATION_TEST:
		doneMask = dali202.StatusBitDurationTestDone
		failureMask = dali202.FailureBitDurationTest
		requestDuration = true
	default:
		return nil, errInvalidTest
	}

	// work out if the test requested has been completed and has data
	rawStatus, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryEmergencyStatus,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryEmergencyStatus: %w", err)
	}
	// get the results of the test i.e. did it succeed or fail, and how long did the battery last for duration tests
	rawFailure, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryFailureStatus,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryFailureStatus: %w", err)
	}

	var result *rpc.TestResult
	if uint8(rawFailure)&failureMask != 0 {
		// the test failed, report failure
		result = &rpc.TestResult{
			Test: request.GetTest(),
			Pass: false,
		}
	} else if uint8(rawStatus)&doneMask != 0 {
		result = &rpc.TestResult{
			Test: request.GetTest(),
			Pass: true,
		}
	} else {
		// test not complete, no data to return
		return nil, status.Error(codes.NotFound, "test results not present - a test must be run first")
	}

	if requestDuration {
		// fetch how long the duration test lasted, in units of two minutes
		rawDuration, err := s.bus.ExecuteCommand(ctx, dali.Request{
			Command:     dali.QueryDurationTestResult,
			AddressType: dali.Short,
			Address:     s.shortAddr,
		})
		if err != nil {
			return nil, fmt.Errorf("QueryDurationTestResult: %w", err)
		}

		result.Duration = durationpb.New(time.Duration(rawDuration) * 2 * time.Minute)
	}

	return result, nil
}

func (s *emergencyLightServer) DeleteTestResult(ctx context.Context, request *rpc.DeleteTestResultRequest) (*rpc.TestResult, error) {
	result, err := s.GetTestResult(ctx, &rpc.GetTestResultRequest{
		Name: request.GetName(),
		Test: request.GetTest(),
	})
	if err != nil {
		return nil, err
	}
	// we can't make the emergency light forget a test failure, so don't try
	// the only way to clear the failed test result is to run another test which passes
	if !result.Pass {
		return nil, status.Error(codes.FailedPrecondition, "test failure cannot be cleared until a test passes")
	}

	var command dali.Command
	switch request.GetTest() {
	case rpc.Test_FUNCTION_TEST:
		command = dali.ResetFunctionTestDoneFlag
	case rpc.Test_DURATION_TEST:
		command = dali.ResetDurationTestDoneFlag
	default:
		return nil, errInvalidTest
	}

	_, err = s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     command,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("Reset...TestDoneFlag: %w", err)
	}

	return result, nil
}

func decodeEmergencyStatus(rawStatus, rawMode, rawBattery, rawFailure uint8) *rpc.EmergencyStatus {
	dest := &rpc.EmergencyStatus{
		InhibitActive:        rawStatus&dali202.StatusBitInhibit != 0,
		IdentificationActive: rawStatus&dali202.StatusBitIdentificationActive != 0,
		ActiveModes:          decodeEmergencyMode(rawMode),
		Failures:             decodeFailures(rawFailure),
	}

	if rawStatus&dali202.StatusBitFunctionTestPending != 0 {
		dest.PendingTests = append(dest.PendingTests, rpc.Test_FUNCTION_TEST)
	}
	if rawStatus&dali202.StatusBitDurationTestPending != 0 {
		dest.PendingTests = append(dest.PendingTests, rpc.Test_DURATION_TEST)
	}

	if rawStatus&dali202.StatusBitFunctionTestDone != 0 {
		dest.ResultsAvailable = append(dest.ResultsAvailable, rpc.Test_FUNCTION_TEST)
	}
	if rawStatus&dali202.StatusBitDurationTestDone != 0 {
		dest.ResultsAvailable = append(dest.ResultsAvailable, rpc.Test_DURATION_TEST)
	}

	if rawFailure&dali202.FailureBitFunctionMaxDelayExceeded != 0 {
		dest.OverdueTests = append(dest.OverdueTests, rpc.Test_FUNCTION_TEST)
	}
	if rawFailure&dali202.FailureBitDurationMaxDelayExceeded != 0 {
		dest.OverdueTests = append(dest.OverdueTests, rpc.Test_DURATION_TEST)
	}

	if level, ok := decodeBatteryLevel(rawBattery); ok {
		dest.BatteryLevelPercent = level
	} else if rawStatus&dali202.StatusBitBatteryFull != 0 {
		dest.BatteryLevelPercent = 100
	}

	return dest
}

func decodeEmergencyMode(rawMode uint8) (modes []rpc.EmergencyStatus_Mode) {
	if rawMode&dali202.ModeBitRestActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_REST)
	}
	if rawMode&dali202.ModeBitNormalModeActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_NORMAL)
	}
	if rawMode&dali202.ModeBitEmergencyModeActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_EMERGENCY)
	}
	if rawMode&dali202.ModeBitExtendedEmergencyModeActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_EXTENDED_EMERGENCY)
	}
	if rawMode&dali202.ModeBitFunctionTestInProgress != 0 {
		modes = append(modes, rpc.EmergencyStatus_FUNCTION_TEST_ACTIVE)
	}
	if rawMode&dali202.ModeBitDurationTestInProgress != 0 {
		modes = append(modes, rpc.EmergencyStatus_DURATION_TEST_ACTIVE)
	}
	if rawMode&dali202.ModeBitHardwiredInhibit != 0 {
		modes = append(modes, rpc.EmergencyStatus_HARDWIRED_INHIBIT)
	}
	if rawMode&dali202.ModeBitHardwiredSwitch != 0 {
		modes = append(modes, rpc.EmergencyStatus_HARDWIRED_SWITCH)
	}
	return
}

func decodeBatteryLevel(rawLevel uint8) (percent float32, ok bool) {
	if rawLevel == 255 {
		return 0, false
	}
	return float32(rawLevel) * 100.0 / 254.0, true
}

func decodeFailures(rawFailure uint8) (failures []rpc.EmergencyStatus_Failure) {
	if rawFailure&dali202.FailureBitCircuit != 0 {
		failures = append(failures, rpc.EmergencyStatus_CIRCUIT_FAILURE)
	}
	if rawFailure&dali202.FailureBitBatteryDuration != 0 {
		failures = append(failures, rpc.EmergencyStatus_BATTERY_DURATION_FAILURE)
	}
	if rawFailure&dali202.FailureBitBattery != 0 {
		failures = append(failures, rpc.EmergencyStatus_BATTERY_FAILURE)
	}
	if rawFailure&dali202.FailureBitEmergencyLamp != 0 {
		failures = append(failures, rpc.EmergencyStatus_LAMP_FAILURE)
	}
	if rawFailure&dali202.FailureBitFunctionTest != 0 {
		failures = append(failures, rpc.EmergencyStatus_FUNCTION_TEST_FAILED)
	}
	if rawFailure&dali202.FailureBitDurationTest != 0 {
		failures = append(failures, rpc.EmergencyStatus_DURATION_TEST_FAILED)
	}
	return
}
