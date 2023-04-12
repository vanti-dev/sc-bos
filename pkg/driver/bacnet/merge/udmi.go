package merge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

const UdmiMergeName = "udmi"

type UdmiMergeConfig struct {
	config.Trait
	TopicPrefix string                         `json:"topicPrefix,omitempty"`
	Points      map[string]*config.ValueSource `json:"points"`
}

func readUdmiMergeConfig(raw []byte) (cfg UdmiMergeConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

// udmiMerge implements the UdmiService and will merge multiple BACnet objects into one UDMI payload
// BACnet objects are polled for changes, and any changes sent as UDMI events
// control is implemented via OnMessage, only points present in the config are controllable.
type udmiMerge struct {
	gen.UnimplementedUdmiServiceServer
	client *gobacnet.Client
	known  known.Context
	logger *zap.Logger

	config UdmiMergeConfig
	bus    minibus.Bus[*gen.PullExportMessagesResponse]

	pollTask *task.Intermittent
	// protect the points value
	pointsLock sync.Mutex
	points     udmi.PointsEvent
}

func newUdmiMerge(client *gobacnet.Client, ctx known.Context, config config.RawTrait, logger *zap.Logger) (*udmiMerge, error) {
	cfg, err := readUdmiMergeConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	f := &udmiMerge{
		client: client,
		known:  ctx,
		config: cfg,
		logger: logger,
	}
	f.pollTask = task.NewIntermittent(f.startPoll)
	return f, nil
}

func (f *udmiMerge) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(f.config.Name, node.HasTrait(udmipb.TraitName, node.WithClients(gen.WrapUdmiService(f))))
}

func (f *udmiMerge) PullControlTopics(request *gen.PullControlTopicsRequest, server gen.UdmiService_PullControlTopicsServer) error {
	err := server.Send(&gen.PullControlTopicsResponse{
		Name:   f.config.Name,
		Topics: []string{f.config.TopicPrefix + "/config"},
	})
	if err != nil {
		return err
	}
	ctx := server.Context()
	<-ctx.Done()
	return ctx.Err()
}

func (f *udmiMerge) OnMessage(ctx context.Context, request *gen.OnMessageRequest) (*gen.OnMessageResponse, error) {
	if request.Message == nil {
		return nil, status.Error(codes.InvalidArgument, "no message")
	}
	var msg udmi.ConfigMessage
	err := json.Unmarshal([]byte(request.Message.Payload), &msg)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid config message: %s", err)
	}
	for point, value := range msg.PointSet.Points {
		if cfg, ok := f.config.Points[point]; ok {
			set := value.SetValue
			if floatValue, isFloat := value.SetValue.(float64); isFloat {
				// let's assume for now values shouldn't be float64, true for the YABE room simulator at least
				set = float32(floatValue)
			}
			err = writeProperty(ctx, f.client, f.known, *cfg, set, 0)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to write point %s: %s", point, err)
			}
		}
	}
	return &gen.OnMessageResponse{Name: f.config.Name}, nil
}

func (f *udmiMerge) PullExportMessages(request *gen.PullExportMessagesRequest, server gen.UdmiService_PullExportMessagesServer) error {
	_ = f.pollTask.Attach(server.Context())
	for msg := range f.bus.Listen(server.Context()) {
		err := server.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *udmiMerge) startPoll(init context.Context) (stop task.StopFn, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(f.config.PollPeriodDuration())
	go func() {
		for {
			select {
			case <-ticker.C:
				err := f.pollPeer(ctx)
				if err != nil { // todo: should this return?
					f.logger.Warn("pollPeer error", zap.String("err", err.Error()))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel, nil
}

// pollPeer fetches data from the peer device, save locally, and fire a change if there is one
func (f *udmiMerge) pollPeer(ctx context.Context) error {
	events := make(udmi.PointsEvent)
	var errs []error
	for key, cfg := range f.config.Points {
		value, err := readProperty(ctx, f.client, f.known, *cfg)
		if err != nil {
			errs = append(errs, fmt.Errorf("read property %q: %w", key, err))
			continue
		}
		events[key] = udmi.PointValue{PresentValue: value}
	}
	if len(errs) == len(f.config.Points) {
		return multierr.Combine(errs...)
	}
	if len(errs) > 0 {
		f.logger.Debug("ignoring some errors", zap.Errors("errs", errs))
	}

	f.pointsLock.Lock()
	isEqual := f.points.Equal(events)
	hasUpdate := !isEqual
	if hasUpdate {
		f.points = events
	}
	f.pointsLock.Unlock()
	if hasUpdate {
		// send the update
		b, err := json.Marshal(f.points)
		if err != nil {
			return err
		}
		f.bus.Send(ctx, &gen.PullExportMessagesResponse{
			Name: f.config.Name,
			Message: &gen.MqttMessage{
				Topic:   f.config.TopicPrefix + "/event/pointset/points",
				Payload: string(b),
			},
		})
	}
	return nil
}
