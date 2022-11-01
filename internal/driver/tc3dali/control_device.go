package tc3dali

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/bridge"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// controlDeviceServer implements the OccupancySensor trait using DALI part 303 occupancy sensors.
// See https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc3_dali/6777329803.html&id=5128453449526025647
// for details on the operation of such sensors.
type controlDeviceServer struct {
	traits.UnimplementedOccupancySensorApiServer
	bus       dali.Dali
	shortAddr uint8

	occupancy *resource.Value
	logger    *zap.Logger

	enableEventsOnce *once
}

func (s *controlDeviceServer) GetOccupancy(ctx context.Context, _ *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	data, err := s.bus.ExecuteCommand(ctx, dali.Request{
		Command:             dali.QueryInputValue,
		AddressType:         dali.Short,
		Address:             s.shortAddr,
		InstanceAddressType: dali.IATInstanceType,
		InstanceAddress:     dali.InstanceTypeOccupancy,
	})
	if err != nil {
		return nil, err
	}

	if data > math.MaxUint8 {
		return nil, status.Error(codes.Internal, "bridge returned an out-of-range value")
	}
	state := occupancyFromInputValue(uint8(data))
	occupancy := &traits.Occupancy{
		State: state,
	}
	// save the value in the local cache
	// filter it though the resource.Value to respect any transformations it performs
	protoOccupancy, err := s.occupancy.Set(occupancy)
	if err != nil {
		s.logger.Error("can't save occupancy value", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to cache occupancy value")
	}

	return protoOccupancy.(*traits.Occupancy), nil
}

func (s *controlDeviceServer) PullOccupancy(req *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	err := s.ensureEventsEnabled(ctx)
	if err != nil {
		s.logger.Error("failed to enable occupancy events for sensor", zap.Error(err))
		return status.Error(codes.Unavailable, "cannot communicate with occupancy sensor")
	}

	changes := s.occupancy.Pull(ctx, resource.WithBackpressure(false))
	for change := range changes {
		err := server.Send(&traits.PullOccupancyResponse{
			Changes: []*traits.PullOccupancyResponse_Change{
				{
					Name:       req.Name,
					ChangeTime: timestamppb.New(change.ChangeTime),
					Occupancy:  change.Value.(*traits.Occupancy),
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *controlDeviceServer) handleInputEvent(event dali.InputEvent, err error) {
	if err != nil || event.Err != nil {
		// this event doesn't contain any useful data
	}

	// only process events if they are for this control device
	if event.Scheme != dali.EventSchemeDevice || event.DeviceShortAddress() != s.shortAddr ||
		event.InstanceType() != bridge.InstanceType(dali.InstanceTypeOccupancy) {

		return
	}

	// only process occupancy events
	occupied := extractEventOccupancy(event)
	var occupancy *traits.Occupancy
	if occupied {
		occupancy = &traits.Occupancy{State: traits.Occupancy_OCCUPIED}
	} else {
		occupancy = &traits.Occupancy{State: traits.Occupancy_UNOCCUPIED}
	}

	_, err = s.occupancy.Set(occupancy)
	if err != nil {
		s.logger.Warn("failed to update occupancy resource after event received",
			zap.Error(err), zap.String("state", occupancy.State.String()))
	}
}

func (s *controlDeviceServer) ensureEventsEnabled(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err, done := s.enableEventsOnce.Do(ctx, func() error {
		// Enable the Occupancy instance on the device.
		// The device won't send any Occupancy-related events until the Occupancy instance has been enabled.
		// This is idempotent.
		// This one needs to be done first, the other commands can happen in any order.
		_, err := s.bus.ExecuteCommand(ctx, dali.Request{
			Command:             dali.EnableInstance,
			AddressType:         dali.Short,
			Address:             s.shortAddr,
			InstanceAddressType: dali.IATInstanceType,
			InstanceAddress:     dali.InstanceTypeOccupancy,
		})
		if err != nil {
			return fmt.Errorf("EnableInstance: %w", err)
		}
		// Set the event filter. The event filter determines which kinds of events will be sent to the DALI bus.
		// We enable events when the area becomes occupied and unoccupied, plus repeat events which are fired occasionally
		// when there have been no changes for a while - this helps in case we miss the original event.
		// Other events are not enabled because they cause bus congestion for no benefit.
		// This is idempotent.
		_, err = s.bus.ExecuteCommand(ctx, dali.Request{
			Command:             dali.SetEventFilter,
			AddressType:         dali.Short,
			Address:             s.shortAddr,
			InstanceAddressType: dali.IATInstanceType,
			InstanceAddress:     dali.InstanceTypeOccupancy,
			Data:                eventFilterOccupied | eventFilterVacant | eventFilterRepeat,
		})
		if err != nil {
			return fmt.Errorf("SetEventFilter: %w", err)
		}
		// Tell the PLC program to listen for events from this device. When the PLC receives a matching event,
		// it will place the event in the queue so we can detect it from Go. This is idempotent.
		err = s.bus.EnableInputEventListener(dali.InputEventParameters{
			Scheme:       dali.EventSchemeDevice,
			AddressInfo1: s.shortAddr,
			AddressInfo2: dali.InstanceTypeOccupancy,
		})
		if err != nil {
			return fmt.Errorf("EnableInputEventListener: %w", err)
		}
		// The callback will be called when the bridge.Dali instance receives any new event from the Go program.
		// This is not idempotent - if we register twice we'll receive each event twice.
		err = s.bus.OnInputEvent(s.handleInputEvent)
		if err != nil {
			return fmt.Errorf("register event handler: %w", err)
		}
		return nil
	})
	if !done {
		return ctx.Err()
	}
	return err
}

// Determines if this event, received from an occupancy sensor, indicates the sensor detects an occupied space
// or an unoccupied space.
// See https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc3_dali/6777329803.html&id=5128453449526025647
// for documentation of the bit fields in this event.
func extractEventOccupancy(event dali.InputEvent) (occupied bool) {
	return (event.Data & (1 << 1)) != 0
}

func occupancyFromInputValue(inputValue uint8) traits.Occupancy_State {
	switch inputValue {
	case 0x00, 0x55:
		return traits.Occupancy_UNOCCUPIED
	case 0xAA, 0xFF:
		return traits.Occupancy_OCCUPIED
	}
	return traits.Occupancy_STATE_UNSPECIFIED
}

// The values of the event filter bit fields that an occupancy sensor control device will send to the DALI bus.
// Documented in the "Event filter" table at
// https://infosys.beckhoff.com/content/1033/tcplclib_tc3_dali/6777329803.html?id=5128453449526025647
const (
	eventFilterOccupied   uint8 = 1 << 0
	eventFilterVacant     uint8 = 1 << 1
	eventFilterRepeat     uint8 = 1 << 2
	eventFilterMovement   uint8 = 1 << 3
	eventFilterNoMovement uint8 = 1 << 4
)
