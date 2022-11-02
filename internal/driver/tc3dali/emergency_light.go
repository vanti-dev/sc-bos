package tc3dali

import (
	"context"
	"fmt"

	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
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

	rawBatteryLevel, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryBatteryCharge,
		AddressType: dali.Short,
		Address:     s.shortAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("QueryBatteryCharge: %w", err)
	}

	emergencyStatus := decodeEmergencyStatus(uint8(rawEmergencyStatus), uint8(rawEmergencyMode), uint8(rawBatteryLevel))
	return emergencyStatus, nil
}

func decodeEmergencyStatus(rawStatus, rawMode, rawBattery uint8) *rpc.EmergencyStatus {
	const (
		bitInhibit uint8 = 1 << iota
		bitFunctionTestDone
		bitDurationTestDone
		bitBatteryFull
		bitFunctionTestPending
		bitDurationTestPending
		bitIdentificationActive
		bitPhysicallySelected
	)

	dest := &rpc.EmergencyStatus{
		InhibitActive:        rawStatus&bitInhibit != 0,
		IdentificationActive: rawStatus&bitIdentificationActive != 0,
		ActiveModes:          decodeEmergencyMode(rawMode),
	}

	if rawStatus&bitFunctionTestPending != 0 {
		dest.PendingTests = append(dest.PendingTests, rpc.Test_FUNCTION_TEST)
	}
	if rawStatus&bitDurationTestPending != 0 {
		dest.PendingTests = append(dest.PendingTests, rpc.Test_DURATION_TEST)
	}

	if rawStatus&bitFunctionTestDone != 0 {
		dest.ResultsAvailable = append(dest.ResultsAvailable, rpc.Test_FUNCTION_TEST)
	}
	if rawStatus&bitDurationTestDone != 0 {
		dest.ResultsAvailable = append(dest.ResultsAvailable, rpc.Test_DURATION_TEST)
	}

	if level, ok := decodeBatteryLevel(rawBattery); ok {
		dest.BatteryLevelPercent = level
	} else if rawStatus&bitBatteryFull != 0 {
		dest.BatteryLevelPercent = 100
	}

	return dest
}

func decodeEmergencyMode(rawMode uint8) (modes []rpc.EmergencyStatus_Mode) {
	const (
		bitRestActive uint8 = 1 << iota
		bitNormalModeActive
		bitEmergencyModeActive
		bitExtendedEmergencyModeActive
		bitFunctionTestInProgress
		bitDurationTestInProgress
		bitHardwiredInhibit
		bitHardwiredSwitch
	)

	if rawMode&bitRestActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_REST)
	}
	if rawMode&bitNormalModeActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_NORMAL)
	}
	if rawMode&bitEmergencyModeActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_EMERGENCY)
	}
	if rawMode&bitExtendedEmergencyModeActive != 0 {
		modes = append(modes, rpc.EmergencyStatus_EXTENDED_EMERGENCY)
	}
	if rawMode&bitFunctionTestInProgress != 0 {
		modes = append(modes, rpc.EmergencyStatus_FUNCTION_TEST_ACTIVE)
	}
	if rawMode&bitDurationTestInProgress != 0 {
		modes = append(modes, rpc.EmergencyStatus_DURATION_TEST_ACTIVE)
	}
	if rawMode&bitHardwiredInhibit != 0 {
		modes = append(modes, rpc.EmergencyStatus_HARDWIRED_INHIBIT)
	}
	if rawMode&bitHardwiredSwitch != 0 {
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
