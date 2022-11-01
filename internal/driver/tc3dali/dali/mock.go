package dali

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	MockErrInvalidAddressType = errors.New("invalid address type")
	MockErrNoSuchControlGear  = errors.New("no control gear with the given ID")
	MockErrUnimplemented      = errors.New("command unimplemented")
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
	}

	return 0, MockErrUnimplemented
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
	return mcg
}

type MockControlGear struct {
	shortAddr       uint8
	ActualLevel     uint8
	GroupMembership uint16
	SceneLevels     [16]uint8
}

func (m *MockControlGear) IsGroupMember(groupNum uint8) bool {
	if !isValidGroup(groupNum) {
		panic("invalid group")
	}

	return (m.GroupMembership & (1 << groupNum)) != 0
}

func isValidGroup(groupNum uint8) bool {
	return groupNum < 64
}
