package dali

import (
	"context"
	"errors"
	"fmt"

	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali/dali202"
	"go.uber.org/zap"
)

var (
	MockErrInvalidAddressType = errors.New("invalid address type")
	MockErrNoSuchControlGear  = errors.New("no control gear with the given ID")
)

type Mock struct {
	ControlGear [64]*MockControlGear
	logger      *zap.Logger
}

func NewMock(logger *zap.Logger) *Mock {
	mock := &Mock{logger: logger}
	for i := range mock.ControlGear {
		mock.ControlGear[i] = newMockControlGear(uint8(i))
	}
	return mock
}

func (m *Mock) ExecuteCommand(ctx context.Context, request Request) (data uint32, err error) {
	switch request.Command {
	case QueryActualLevel:
		level, err := m.queryActualLevel(request)
		return uint32(level), err
	case DirectArcPowerControl:
		return 0, m.directArcPowerControl(request)
	case GoToScene:
		return 0, m.goToScene(request)
	case SetFadeTime:
		// we don't implement this
		return 0, nil
	case QueryGroups:
		groups, err := m.queryGroups(request)
		return uint32(groups), err
	case AddToGroup:
		return 0, m.addToGroup(request)
	case RemoveFromGroup:
		return 0, m.removeFromGroup(request)
	case IdentifyDevice102, IdentifyDevice103, StartIdentification202:
		m.identify(request)
		return 0, nil
	case QueryEmergencyStatus:
		cg, err := m.getSingleControlGear(request)
		if err != nil {
			return 0, err
		}
		return uint32(cg.EmergencyStatus), nil
	case QueryEmergencyMode:
		cg, err := m.getSingleControlGear(request)
		if err != nil {
			return 0, err
		}
		return uint32(cg.EmergencyMode), nil
	case QueryFailureStatus:
		cg, err := m.getSingleControlGear(request)
		if err != nil {
			return 0, err
		}
		return uint32(cg.EmergencyFailure), nil
	case QueryDurationTestResult:
		cg, err := m.getSingleControlGear(request)
		if err != nil {
			return 0, err
		}
		return uint32(cg.DurationTestResult), nil
	case QueryBatteryCharge:
		cg, err := m.getSingleControlGear(request)
		if err != nil {
			return 0, err
		}
		return uint32(cg.BatteryLevel), nil
	case StartFunctionTest:
		m.startTest(request, dali202.ModeBitFunctionTestInProgress)
		return 0, nil
	case StartDurationTest:
		m.startTest(request, dali202.ModeBitDurationTestInProgress)
		return 0, nil
	case StopTest:
		m.stopTest(request)
		return 0, nil
	case ResetFunctionTestDoneFlag:
		m.clearFunctionTest(request)
		return 0, nil
	case ResetDurationTestDoneFlag:
		m.clearDurationTest(request)
		return 0, nil
	}

	return 0, ErrCommandUnimplemented
}

func (m *Mock) queryActualLevel(request Request) (actualLevel byte, err error) {
	cg, err := m.getSingleControlGear(request)
	if err != nil {
		return 0, err
	}
	return cg.ActualLevel, nil
}

func (m *Mock) directArcPowerControl(request Request) error {
	targets := m.getTargetControlGear(request)

	for _, target := range targets {
		target.ActualLevel = request.Data
		m.logger.Debug("Direct Arc Power Control",
			zap.Uint8("shortAddress", target.shortAddr),
			zap.Uint8("level", target.ActualLevel))
	}
	return nil
}

func (m *Mock) goToScene(request Request) error {
	targets := m.getTargetControlGear(request)
	scene := request.Data

	for _, target := range targets {
		if target.SceneLevels[scene] != 255 {
			target.ActualLevel = target.SceneLevels[scene]
			m.logger.Debug("Go To Scene",
				zap.Uint8("scene", scene),
				zap.Uint8("level", target.SceneLevels[scene]))
		} else {
			m.logger.Debug("Go To Scene - MASK",
				zap.Uint8("scene", scene))
		}
	}

	return nil
}

func (m *Mock) queryGroups(request Request) (groups uint16, err error) {
	gear, err := m.getSingleControlGear(request)
	if err != nil {
		return 0, err
	}
	return gear.GroupMembership, nil
}

func (m *Mock) addToGroup(request Request) error {
	gear := m.getTargetControlGear(request)
	for _, g := range gear {
		g.GroupMembership |= 1 << request.Data
	}
	return nil
}

func (m *Mock) removeFromGroup(request Request) error {
	gear := m.getTargetControlGear(request)
	for _, g := range gear {
		g.GroupMembership &= ^(1 << request.Data)
	}
	return nil
}

func (m *Mock) identify(request Request) {
	gear := m.getTargetControlGear(request)
	addrs := make([]uint8, 0, len(gear))
	for _, g := range gear {
		addrs = append(addrs, g.shortAddr)
	}
	m.logger.Info("identify control gear", zap.Uint8s("shortAddrs", addrs))
}

func (m *Mock) clearFunctionTest(request Request) {
	gear := m.getTargetControlGear(request)
	for _, g := range gear {
		g.EmergencyStatus &= ^dali202.StatusBitFunctionTestDone
		g.EmergencyFailure &= ^dali202.FailureBitFunctionTest
	}
}

func (m *Mock) clearDurationTest(request Request) {
	gear := m.getTargetControlGear(request)
	for _, g := range gear {
		g.EmergencyStatus &= ^dali202.StatusBitDurationTestDone
		g.EmergencyFailure &= ^dali202.FailureBitDurationTest
		g.DurationTestResult = 0
	}
}

func (m *Mock) stopTest(request Request) {
	gear := m.getTargetControlGear(request)
	for _, g := range gear {
		if g.TestInProgress() {
			g.EmergencyMode = dali202.ModeBitNormalModeActive
		} else {
			m.logger.Warn("can't stop test: no test running",
				zap.Uint8("shortAddr", g.shortAddr))
		}
	}
}

// pass either dali202.ModeBitFunctionTestInProgress or dali202.ModeBitDurationTestInProgress as testBit
func (m *Mock) startTest(request Request, testBit uint8) {
	for _, g := range m.getTargetControlGear(request) {
		if g.EmergencyMode&dali202.ModeBitNormalModeActive == 0 {
			m.logger.Warn("can't start a test because control gear not in normal mode",
				zap.Uint8("shortAddr", g.shortAddr),
				zap.String("mode", fmt.Sprintf("%08b", g.EmergencyMode)),
			)
			continue
		}

		g.EmergencyMode = testBit
	}
}

func (m *Mock) getTargetControlGear(request Request) []*MockControlGear {
	var targets []*MockControlGear
	switch request.AddressType {
	case Short:
		cg := m.ControlGear[request.Address]
		if cg != nil {
			targets = []*MockControlGear{cg}
		}
	case Group:
		for _, cg := range m.ControlGear {
			if cg != nil && cg.IsGroupMember(request.Address) {
				targets = append(targets, cg)
			}
		}
	case Broadcast:
		targets = m.ControlGear[:]
	default:
		panic("unsupported address type")
	}
	if len(targets) == 0 {
		m.logger.Warn("no control gear selected by request")
	}
	return targets
}

func (m *Mock) getSingleControlGear(request Request) (*MockControlGear, error) {
	if request.AddressType != Short {
		return nil, MockErrInvalidAddressType
	}
	if m.ControlGear[request.Address] == nil {
		return nil, MockErrNoSuchControlGear
	}
	return m.ControlGear[request.Address], nil
}

func (m *Mock) EnableInputEventListener(params InputEventParameters) error {
	panic("input devices not implemented")
}

func (m *Mock) OnInputEvent(handler InputEventHandler) error {
	panic("input devices not implemented")
}

func (m *Mock) Close() error {
	return nil
}

var _ Dali = (*Mock)(nil)

func newMockControlGear(shortAddr uint8) *MockControlGear {
	mcg := &MockControlGear{shortAddr: shortAddr}
	for i := range mcg.SceneLevels {
		mcg.SceneLevels[i] = 255 // MASK means 'no level programmed'
	}
	// by default make every control gear a member of group 0
	mcg.GroupMembership |= 1 << 0
	mcg.EmergencyMode = dali202.ModeBitNormalModeActive
	mcg.BatteryLevel = 254
	return mcg
}

type MockControlGear struct {
	shortAddr          uint8
	ActualLevel        uint8
	GroupMembership    uint16
	SceneLevels        [16]uint8
	EmergencyStatus    uint8
	EmergencyMode      uint8
	EmergencyFailure   uint8
	DurationTestResult uint8
	BatteryLevel       uint8
}

func (m *MockControlGear) SetBatteryLevel(level uint8) {
	m.BatteryLevel = level
	if level == 255 {
		m.EmergencyStatus |= dali202.StatusBitBatteryFull
	} else {
		m.EmergencyStatus &= ^dali202.StatusBitBatteryFull
	}
}

func (m *MockControlGear) CompleteDurationTest(success bool, result uint8) (ok bool) {
	if m.EmergencyMode&dali202.ModeBitDurationTestInProgress == 0 {
		// duration test must be in progress to complete test
		return false
	}
	m.EmergencyMode = dali202.ModeBitNormalModeActive
	m.EmergencyStatus |= dali202.StatusBitDurationTestDone
	if success {
		m.EmergencyFailure &= ^dali202.FailureBitDurationTest
	} else {
		m.EmergencyFailure |= dali202.FailureBitDurationTest
	}
	m.DurationTestResult = result

	return true
}

func (m *MockControlGear) CompleteFunctionTest(success bool) (ok bool) {
	if m.EmergencyMode&dali202.ModeBitFunctionTestInProgress == 0 {
		// duration test must be in progress to complete test
		return false
	}
	m.EmergencyMode = dali202.ModeBitNormalModeActive
	m.EmergencyStatus |= dali202.StatusBitFunctionTestDone
	if success {
		m.EmergencyFailure &= ^dali202.FailureBitFunctionTest
	} else {
		m.EmergencyFailure |= dali202.FailureBitFunctionTest
	}

	return true
}

func (m *MockControlGear) IsGroupMember(groupNum uint8) bool {
	if !isValidGroup(groupNum) {
		panic("invalid group")
	}

	return (m.GroupMembership & (1 << groupNum)) != 0
}

func (m *MockControlGear) TestInProgress() bool {
	testMask := dali202.ModeBitFunctionTestInProgress | dali202.ModeBitDurationTestInProgress
	return m.EmergencyMode&testMask != 0
}

func isValidGroup(groupNum uint8) bool {
	return groupNum < 64
}
