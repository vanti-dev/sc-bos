package merge

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/vanti-dev/gobacnet"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

const udmiMergeName = "udmi"

type udmiMergeConfig struct {
	config.Trait
	TopicPrefix string                         `json:"topicPrefix,omitempty"`
	Points      map[string]*config.ValueSource `json:"points"`
}

func readUdmiMergeConfig(raw []byte) (cfg udmiMergeConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

// udmiMerge implements the UdmiService and will merge multiple BACnet objects into one UDMI payload
// BACnet objects are polled for changes, and any changes sent as UDMI events
// control is not yet implemented
type udmiMerge struct {
	gen.UnimplementedUdmiServiceServer
	client *gobacnet.Client
	known  known.Context
	logger *zap.Logger

	config udmiMergeConfig
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
	return a.Announce(f.config.Name, node.HasClient(gen.WrapUdmiService(f)))
}

func (f *udmiMerge) PullControlTopics(request *gen.PullControlTopicsRequest, server gen.UdmiService_PullControlTopicsServer) error {
	f.logger.Debug("PullControlTopics")
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
	f.logger.Debug("OnMessage")
	// todo: implement control
	return &gen.OnMessageResponse{Name: f.config.Name}, nil
}

func (f *udmiMerge) PullExportMessages(request *gen.PullExportMessagesRequest, server gen.UdmiService_PullExportMessagesServer) error {
	f.logger.Debug("PullExportMessages")
	_ = f.pollTask.Attach(server.Context())
	for msg := range f.bus.Listen(server.Context()) {
		f.logger.Debug("sending export message", zap.Any("message", msg))
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
				break
			}
		}
	}()
	return cancel, nil
}

// pollPeer fetches data from the peer device, save locally, and fire a change if there is one
func (f *udmiMerge) pollPeer(ctx context.Context) error {
	events := make(udmi.PointsEvent)
	for key, cfg := range f.config.Points {
		value, err := readProperty(ctx, f.client, f.known, *cfg)
		if err != nil {
			return err
		}
		events[key] = &udmi.PointValue{PresentValue: value}
	}
	f.pointsLock.Lock()
	defer f.pointsLock.Unlock()
	if !f.points.Equal(events) {
		f.points = events
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
