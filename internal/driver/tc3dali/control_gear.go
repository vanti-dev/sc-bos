package tc3dali

import (
	"context"
	"math"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controlGearServer struct {
	traits.UnimplementedLightApiServer
	rpc.UnimplementedDaliApiServer

	bus      dali.Dali
	addr     uint8
	addrType dali.AddressType
	logger   *zap.Logger
}

func (s *controlGearServer) GetBrightness(ctx context.Context, req *traits.GetBrightnessRequest) (*traits.Brightness, error) {
	if s.addrType != dali.Short {
		return nil, status.Error(codes.Unimplemented, "GetBrightness only supported for control gear")
	}

	data, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryActualLevel,
		AddressType: s.addrType,
		Address:     s.addr,
	})
	if err != nil {
		return nil, err
	}

	if data > math.MaxUint8 {
		s.logger.Error("bridge returned out-of-range value for QueryActualLevel", zap.Uint32("data", data))
		return nil, status.Error(codes.Internal, "bridge returned an invalid value")
	}
	percent, ok := daliLevelToPercent(uint8(data))
	if !ok {
		s.logger.Warn("bridge returned MASK value for QueryActualLevel - control gear fault?")
		return nil, status.Error(codes.Unavailable, "light level state unavailable - try again later")
	}

	return &traits.Brightness{
		LevelPercent: percent,
	}, nil
}

func (s *controlGearServer) UpdateBrightness(ctx context.Context, req *traits.UpdateBrightnessRequest) (*traits.Brightness, error) {
	brightness := req.Brightness
	if brightness == nil {
		return nil, status.Error(codes.InvalidArgument, "required field 'brightness' missing")
	}
	if brightness.Preset != nil {
		return nil, status.Error(codes.InvalidArgument, "presets not supported with this driver")
	}
	level, ok := percentToDaliLevel(brightness.LevelPercent)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid 'brightness.level_percent' value")
	}

	_, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.DirectArcPowerControl,
		AddressType: s.addrType,
		Address:     s.addr,
		Data:        level,
	})
	if err != nil {
		return nil, err
	}
	return &traits.Brightness{LevelPercent: brightness.LevelPercent}, nil
}

func (s *controlGearServer) AddToGroup(ctx context.Context, request *rpc.AddToGroupRequest) (*rpc.AddToGroupResponse, error) {
	if request.Group < 0 || request.Group > 15 {
		return nil, status.Errorf(codes.InvalidArgument, "group %d out of range 0-15", request.Group)
	}

	_, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.AddToGroup,
		AddressType: s.addrType,
		Address:     s.addr,
		Data:        byte(request.Group),
	})
	if err != nil {
		return nil, err
	}
	return &rpc.AddToGroupResponse{}, nil
}

func (s *controlGearServer) RemoveFromGroup(ctx context.Context, request *rpc.RemoveFromGroupRequest) (*rpc.RemoveFromGroupResponse, error) {
	if request.Group < 0 || request.Group > 15 {
		return nil, status.Error(codes.InvalidArgument, "group out of range")
	}

	_, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.RemoveFromGroup,
		AddressType: s.addrType,
		Address:     s.addr,
		Data:        byte(request.Group),
	})
	if err != nil {
		return nil, err
	}
	return &rpc.RemoveFromGroupResponse{}, nil
}

func (s *controlGearServer) GetGroupMembership(ctx context.Context, request *rpc.GetGroupMembershipRequest) (*rpc.GetGroupMembershipResponse, error) {
	if s.addrType != dali.Short {
		return nil, status.Errorf(codes.Unimplemented, "query operations only supported for individual control gear")
	}

	data, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:     dali.QueryGroups,
		AddressType: s.addrType,
		Address:     s.addr,
	})
	if err != nil {
		return nil, err
	}

	var groups []int32
	for i := int32(0); i < 16; i++ {
		mask := uint32(1) << i
		if data&mask != 0 {
			groups = append(groups, i)
		}
	}

	return &rpc.GetGroupMembershipResponse{Groups: groups}, nil
}

func daliLevelToPercent(level uint8) (percent float32, ok bool) {
	if level == 255 {
		return 0, false
	}

	return float32(level) * 100.0 / 254.0, true
}

func percentToDaliLevel(percent float32) (level uint8, ok bool) {
	ok = percent >= 0 && percent <= 100
	if !ok {
		return
	}
	level = uint8(percent * 254.0 / 100.0)
	return
}
