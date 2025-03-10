package modbusip

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/goburrow/modbus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/driver/modbusip/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
)

type Client struct {
	cli modbus.Client

	handler Handle

	traits.UnimplementedEnergyStorageApiServer
	traits.UnimplementedEmergencyApiServer
	gen.UnimplementedStatusApiServer

	logger *zap.Logger

	fuel      *resource.Value
	faults    *resource.Value
	emergency *resource.Value

	group  *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
}

type Handle interface {
	io.Closer
	Connect() error
	modbus.ClientHandler
}

func NewClient(ctx context.Context, handler Handle) *Client {
	cli := modbus.NewClient(handler)

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	group, ctx := errgroup.WithContext(ctx)

	return &Client{
		cli:     cli,
		handler: handler,

		group:  group,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Client) Connect(pdu config.PDUAddress, resourceType string, address, quantity uint16) error {
	err := c.handler.Connect()

	if err != nil {
		return err
	}

	switch resourceType {
	case fuel:
		c.fuel = resource.NewValue(resource.WithNoDuplicates())
	case faults:
		c.faults = resource.NewValue(resource.WithNoDuplicates())
	case monitor:
		c.emergency = resource.NewValue(resource.WithNoDuplicates())
	default:
		return fmt.Errorf("unknown resource type %s", resourceType)
	}

	// TODO: the below calculations and bit ops should really be configurable based on the device
	// TODO: but we don't have much information about how the modbus registers correspond to trait values
	// TODO: so this is approximated (or hardcoded) for now
	c.group.Go(func() error {
		throttle := time.NewTicker(time.Second)

		for {
			select {
			case <-throttle.C:
				var value uint16
				switch pdu {
				case config.HoldingRegister:
					res, err := c.cli.ReadHoldingRegisters(address, quantity)

					if err != nil {
						c.logger.Error("reading holding registers", zap.Error(err))
						continue
					}

					if res == nil || len(res) == 0 {
						return fmt.Errorf("invalid holding registers response")
					}

					value = uint16(res[0])<<8 | uint16(res[1]) // big endian
				case config.Coil:
					res, err := c.cli.ReadCoils(address, quantity)

					if err != nil {
						c.logger.Error("reading coils", zap.Error(err))
						continue
					}

					if res == nil || len(res) == 0 {
						return fmt.Errorf("invalid coils response")
					}

					value = uint16(res[0])
				case config.InputRegister:
					res, err := c.cli.ReadInputRegisters(address, quantity)

					if err != nil {
						c.logger.Error("reading input registers", zap.Error(err))
						continue
					}

					if res == nil || len(res) == 0 {
						return fmt.Errorf("invalid input registers response")
					}

					value = uint16(res[0])<<8 | uint16(res[1]) // big endian
				case config.DiscreteInput:
					res, err := c.cli.ReadDiscreteInputs(address, quantity)

					if err != nil {
						c.logger.Error("reading discrete inputs", zap.Error(err))
						continue
					}

					if res == nil || len(res) == 0 {
						return fmt.Errorf("invalid discrete inputs response")
					}

					value = uint16(res[0])
				}

				switch resourceType {
				case fuel:
					if _, err := c.fuel.Set(&traits.EnergyLevel{
						Quantity: &traits.EnergyLevel_Quantity{
							Percentage: float32(value),
						},
					}); err != nil {
						c.logger.Error("setting fuel", zap.Error(err))
						continue
					}
				case faults:
					if _, err := c.faults.Set(&gen.StatusLog{
						Level: gen.StatusLog_NOMINAL,
					}); err != nil {
						c.logger.Error("setting faults", zap.Error(err))
						continue
					}
				case monitor:
					if _, err := c.emergency.Set(&traits.Emergency{
						Level: traits.Emergency_OK,
					}); err != nil {
						c.logger.Error("setting monitor", zap.Error(err))
						continue
					}
				}
			case <-c.ctx.Done():
				return c.ctx.Err()
			}
		}
	})

	return nil
}

func (c *Client) Close() error {
	defer c.cancel()
	return c.handler.Close()
}

const (
	fuel    = "fuel"
	faults  = "faults"
	monitor = "monitor"
)

func (c *Client) GetEmergency(_ context.Context, request *traits.GetEmergencyRequest) (*traits.Emergency, error) {
	if c.emergency == nil {
		return nil, status.Errorf(codes.Unimplemented, "%s not implemented", trait.Emergency)
	}
	return c.emergency.Get(resource.WithReadMask(request.GetReadMask())).(*traits.Emergency), nil
}

func (c *Client) UpdateEmergency(_ context.Context, _ *traits.UpdateEmergencyRequest) (*traits.Emergency, error) {
	return nil, status.Error(codes.Unimplemented, "update emergency not implemented")
}

func (c *Client) PullEmergency(request *traits.PullEmergencyRequest, server traits.EmergencyApi_PullEmergencyServer) error {
	if c.emergency == nil {
		return status.Errorf(codes.Unimplemented, "%s not implemented", trait.Emergency)
	}

	changes := c.faults.Pull(server.Context(), resource.WithReadMask(request.GetReadMask()))

	for {
		select {
		case change := <-changes:
			if err := server.Send(&traits.PullEmergencyResponse{
				Changes: []*traits.PullEmergencyResponse_Change{
					{
						Name:       request.GetName(),
						ChangeTime: timestamppb.New(change.ChangeTime),
						Emergency:  change.Value.(*traits.Emergency),
					},
				},
			}); err != nil {
				return err
			}
		case <-server.Context().Done():
			return server.Context().Err()
		}
	}
}

func (c *Client) GetCurrentStatus(_ context.Context, request *gen.GetCurrentStatusRequest) (*gen.StatusLog, error) {
	if c.faults == nil {
		return nil, status.Errorf(codes.Unimplemented, "%s not implemented", statuspb.TraitName)
	}

	return c.faults.Get(resource.WithReadMask(request.GetReadMask())).(*gen.StatusLog), nil
}

func (c *Client) PullCurrentStatus(request *gen.PullCurrentStatusRequest, server gen.StatusApi_PullCurrentStatusServer) error {
	if c.faults == nil {
		return status.Errorf(codes.Unimplemented, "%s not implemented", statuspb.TraitName)
	}

	changes := c.faults.Pull(server.Context(), resource.WithReadMask(request.GetReadMask()))

	for {
		select {
		case <-server.Context().Done():
			return server.Context().Err()
		case change := <-changes:
			if err := server.Send(&gen.PullCurrentStatusResponse{
				Changes: []*gen.PullCurrentStatusResponse_Change{
					{
						Name:          request.GetName(),
						ChangeTime:    timestamppb.New(change.ChangeTime),
						CurrentStatus: change.Value.(*gen.StatusLog),
					},
				}}); err != nil {
				return err
			}
		}
	}
}

func (c *Client) GetEnergyLevel(_ context.Context, request *traits.GetEnergyLevelRequest) (*traits.EnergyLevel, error) {
	if c.fuel == nil {
		return nil, status.Errorf(codes.Unimplemented, "%s not implemented", trait.EnergyStorage)
	}

	return c.fuel.Get(resource.WithReadMask(request.GetReadMask())).(*traits.EnergyLevel), nil
}

func (c *Client) PullEnergyLevel(request *traits.PullEnergyLevelRequest, server traits.EnergyStorageApi_PullEnergyLevelServer) error {
	if c.fuel == nil {
		return status.Errorf(codes.Unimplemented, "%s not implemented", trait.EnergyStorage)
	}

	changes := c.fuel.Pull(server.Context(), resource.WithReadMask(request.GetReadMask()))

	for {
		select {
		case <-server.Context().Done():
			return server.Context().Err()
		case change := <-changes:
			if err := server.Send(&traits.PullEnergyLevelResponse{
				Changes: []*traits.PullEnergyLevelResponse_Change{
					{
						Name:        request.GetName(),
						ChangeTime:  timestamppb.New(change.ChangeTime),
						EnergyLevel: change.Value.(*traits.EnergyLevel),
					},
				},
			}); err != nil {
				return err
			}
		}
	}
}

func (c *Client) Charge(_ context.Context, _ *traits.ChargeRequest) (*traits.ChargeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "charge not implemented")
}
