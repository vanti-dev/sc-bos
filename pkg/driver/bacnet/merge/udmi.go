package merge

import (
	"context"
	"encoding/json"
	"math"
	"sync"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	status2 "github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
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
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	config UdmiMergeConfig
	bus    minibus.Bus[*gen.PullExportMessagesResponse]

	pollTask *task.Intermittent
	// protect the points value
	pointsLock sync.Mutex
	points     udmi.PointsEvent
}

func newUdmiMerge(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*udmiMerge, error) {
	cfg, err := readUdmiMergeConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	f := &udmiMerge{
		client:   client,
		known:    devices,
		statuses: statuses,
		config:   cfg,
		logger:   logger,
	}
	f.pollTask = task.NewIntermittent(f.startPoll)
	initTraitStatus(statuses, cfg.Name, "UDMI")
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
			err = comm.WriteProperty(ctx, f.client, f.known, *cfg, set, 0)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to write point %s: %s", point, err)
			}
		}
	}
	return &gen.OnMessageResponse{Name: f.config.Name}, nil
}

func (f *udmiMerge) GetExportMessage(ctx context.Context, request *gen.GetExportMessageRequest) (*gen.MqttMessage, error) {
	pollCtx, cleanup := context.WithTimeout(ctx, f.config.PollTimeoutDuration()/4)
	defer cleanup()
	events := f.bus.Listen(pollCtx)
	_ = f.pollTask.Attach(pollCtx)
	select {
	case <-pollCtx.Done():
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		f.pointsLock.Lock()
		points := f.points
		f.pointsLock.Unlock()
		if len(points) == 0 {
			return nil, status.Error(codes.Unavailable, "no recent events")
		}
		return f.pointsToPointSet(f.config.TopicPrefix, points)
	case msg := <-events:
		return msg.Message, nil
	}
}

func (f *udmiMerge) PullExportMessages(request *gen.PullExportMessagesRequest, server gen.UdmiService_PullExportMessagesServer) error {
	events := f.bus.Listen(server.Context())
	_ = f.pollTask.Attach(server.Context())

	// initial value
	if request.IncludeLast {
		f.pointsLock.Lock()
		points := f.points
		f.pointsLock.Unlock()
		if len(points) > 0 {
			msg, err := f.pointsToPointSet(f.config.TopicPrefix, points)
			if err != nil {
				return err
			}
			err = server.Send(&gen.PullExportMessagesResponse{
				Name:    request.Name,
				Message: msg,
			})
			if err != nil {
				return err
			}
		}
	}
	for msg := range events {
		err := server.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *udmiMerge) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "udmi", f.config.PollPeriodDuration(), f.config.PollTimeoutDuration(), f.logger, f.pollPeer)
}

// pollPeer fetches data from the peer device, save locally, and fire a change if there is one
func (f *udmiMerge) pollPeer(ctx context.Context) error {
	events := make(udmi.PointsEvent)
	var errs []error
	requestValues := make([]config.ValueSource, 0, len(f.config.Points))
	keys := make([]string, 0, len(f.config.Points))
	for key, cfg := range f.config.Points {
		requestValues = append(requestValues, *cfg)
		keys = append(keys, key)
	}
	for i, result := range comm.ReadPropertiesChunked(ctx, f.client, f.known, f.config.ChunkSize, requestValues...) {
		switch e := result.(type) {
		case error:
			errs = append(errs, comm.ErrReadProperty{Prop: keys[i], Cause: e})
		default:
			events[keys[i]] = udmi.PointValue{PresentValue: e}
		}
	}

	status2.UpdatePollErrorStatus(f.statuses, f.config.Name, "UDMI", keys, errs)
	if len(errs) == len(f.config.Points) {
		err := multierr.Combine(errs...)
		return err
	}
	if len(errs) > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
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
		msg, err := f.pointsToPointSet(f.config.TopicPrefix, events)
		if err != nil {
			return err
		}
		f.bus.Send(ctx, &gen.PullExportMessagesResponse{
			Name:    f.config.Name,
			Message: msg,
		})
	}
	return nil
}

func sanitise(points udmi.PointsEvent) {
	for k, v := range points {
		if pv, ok := v.PresentValue.(float64); ok {
			if math.IsNaN(pv) || math.IsInf(pv, 0) {
				points[k] = udmi.PointValue{PresentValue: nil}
			}
		}
	}
}

func (f *udmiMerge) pointsToPointSet(topicPrefix string, points udmi.PointsEvent) (*gen.MqttMessage, error) {

	sanitise(points)

	b, err := json.Marshal(points)
	if err != nil {
		return nil, err
	}
	return &gen.MqttMessage{
		Topic:   topicPrefix + "/event/pointset/points",
		Payload: string(b),
	}, nil
}
