package opcua

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gopcua/opcua/ua"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/driver/opcua/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/opcua/conv"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Transport struct {
	config.Trait
	gen.UnimplementedTransportApiServer
	gen.UnimplementedTransportInfoServer

	logger    *zap.Logger
	transport *resource.Value // *gen.Transport
	cfg       config.TransportConfig
	scName    string
}

func readTransportConfig(raw []byte) (cfg config.TransportConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

func newTransport(n string, c config.RawTrait, l *zap.Logger) (*Transport, error) {
	cfg, err := readTransportConfig(c.Raw)
	if err != nil {
		return nil, err
	}
	t := &Transport{
		logger:    l,
		transport: resource.NewValue(resource.WithInitialValue(&gen.Transport{}), resource.WithNoDuplicates()),
		cfg:       cfg,
		scName:    n,
	}
	// initialise the doors as we know these from the config
	tp := &gen.Transport{}
	for _, door := range cfg.Doors {
		tp.Doors = append(tp.Doors, &gen.Transport_Door{Title: door.Title})
	}
	_, _ = t.transport.Set(tp)
	return t, nil
}

func (t *Transport) GetTransport(_ context.Context, _ *gen.GetTransportRequest) (*gen.Transport, error) {
	return t.transport.Get().(*gen.Transport), nil
}

func (t *Transport) PullTransport(_ *gen.PullTransportRequest, server gen.TransportApi_PullTransportServer) error {
	for value := range t.transport.Pull(server.Context()) {
		transport := value.Value.(*gen.Transport)
		err := server.Send(&gen.PullTransportResponse{Changes: []*gen.PullTransportResponse_Change{
			{
				Name:       t.scName,
				ChangeTime: timestamppb.New(value.ChangeTime),
				Transport:  transport,
			},
		}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Transport) handleTransportEvent(node *ua.NodeID, value any) {
	old := t.transport.Get().(*gen.Transport)
	if t.cfg.ActualPosition != nil && NodeIdsAreEqual(t.cfg.ActualPosition.NodeId, node) {
		floor, err := conv.ToString(value)
		if err != nil {
			t.logger.Error("failed to convert ActualPosition event", zap.Error(err))
			return
		}
		old.ActualPosition = &gen.Transport_Location{
			Floor: floor,
		}
	}
	if t.cfg.Load != nil && NodeIdsAreEqual(t.cfg.Load.NodeId, node) {
		load, err := conv.Float32Value(value)
		if err != nil {
			t.logger.Error("failed to convert Load event", zap.Error(err))
			return
		}
		old.Load = &load
	}
	if t.cfg.MovingDirection != nil && NodeIdsAreEqual(t.cfg.MovingDirection.NodeId, node) {
		direction, err := conv.ToTraitEnum[gen.Transport_Direction](value, t.cfg.MovingDirection.Enum, gen.Transport_Direction_value)
		if err != nil {
			t.logger.Error("failed to convert MovingDirection to trait enum", zap.Error(err))
			return
		}
		old.MovingDirection = direction
	}
	if t.cfg.NextDestinations != nil {
		for i, dest := range t.cfg.NextDestinations {
			if dest.Type == config.SingleFloor && NodeIdsAreEqual(dest.Source.NodeId, node) {
				floor, err := conv.IntValue(value)
				if err != nil {
					t.logger.Error("failed to convert NextDestinations event", zap.Error(err))
					return
				}
				if i >= len(old.NextDestinations) {
					old.NextDestinations = append(old.NextDestinations, &gen.Transport_Location{
						Floor: strconv.Itoa(floor),
					})
				} else {
					old.NextDestinations[i] = &gen.Transport_Location{
						Floor: strconv.Itoa(floor),
					}
				}
			}
		}
	}
	if t.cfg.Doors != nil {
		for i, door := range t.cfg.Doors {
			if door.Status != nil && NodeIdsAreEqual(door.Status.NodeId, node) {
				status, err := conv.ToTraitEnum[gen.Transport_Door_DoorStatus](value, door.Status.Enum, gen.Transport_Door_DoorStatus_value)
				if err != nil {
					t.logger.Error("failed to convert Door Status to trait enum", zap.Error(err))
					return
				}
				d := &gen.Transport_Door{
					Title: door.Title,
				}
				d.Status = status
				old.Doors[i] = d
			}
		}
	}
	if t.cfg.Speed != nil && NodeIdsAreEqual(t.cfg.Speed.NodeId, node) {
		speed, err := conv.Float32Value(value)
		if err != nil {
			t.logger.Error("failed to convert Speed event", zap.Error(err))
			return
		}
		old.Speed = &speed
	}
	_, _ = t.transport.Set(old)
}

func (t *Transport) DescribeTransport(context.Context, *gen.DescribeTransportRequest) (*gen.TransportSupport, error) {
	return &gen.TransportSupport{
		LoadUnit:  t.cfg.LoadUnit,
		MaxLoad:   t.cfg.MaxLoad,
		SpeedUnit: t.cfg.SpeedUnit,
	}, nil
}
