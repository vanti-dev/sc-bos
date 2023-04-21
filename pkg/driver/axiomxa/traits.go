package axiomxa

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olebedev/emitter"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/mps"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

func (d *Driver) announceTraits(ctx context.Context, cfg config.Root, announcer node.Announcer, mpsMessages *emitter.Emitter, names *devices) error {
	if len(cfg.Devices) == 0 {
		return nil // no devices, no traits!
	}

	udmiServer := &udmiServer{
		msgs:         mpsMessages,
		udmiMessages: &minibus.Bus[udmiMessage]{},
		names:        names,
		log:          d.logger.Named("udmi"),
	}
	udmiServer.pipe = task.NewIntermittent(udmiServer.mpsToUDMIMessages)

	udmiClient := gen.WrapUdmiService(udmiServer)
	for _, device := range cfg.Devices {
		announcer.Announce(device.Name, node.HasMetadata(device.Metadata), node.HasTrait(udmipb.TraitName, node.WithClients(udmiClient)))
	}

	return nil
}

// udmiServer implements gen.UdmiServiceServer backed by the given emitter events.
// We only support telemetry, not control. These methods just don't do anything, rather than returning UNIMPLEMENTED.
type udmiServer struct {
	gen.UnimplementedUdmiServiceServer
	msgs         *emitter.Emitter          // original message ports
	pipe         *task.Intermittent        // task that converts msgs to udmiMessages
	udmiMessages *minibus.Bus[udmiMessage] // emits converted and processed message ports

	names *devices
	log   *zap.Logger
}

func (u *udmiServer) PullControlTopics(request *gen.PullControlTopicsRequest, topicsServer gen.UdmiService_PullControlTopicsServer) error {
	// we don't have any control topics
	<-topicsServer.Context().Done()
	return nil
}

func (u *udmiServer) OnMessage(ctx context.Context, request *gen.OnMessageRequest) (*gen.OnMessageResponse, error) {
	// we don't support doing anything here
	return &gen.OnMessageResponse{}, nil
}

func (u *udmiServer) PullExportMessages(request *gen.PullExportMessagesRequest, messagesServer gen.UdmiService_PullExportMessagesServer) error {
	if err := u.pipe.Attach(messagesServer.Context()); err != nil {
		return err
	}

	for msg := range u.udmiMessages.Listen(messagesServer.Context()) {
		if msg.name != request.Name {
			continue
		}

		var points *udmi.PointsEvent
		switch msg.topic {
		case KeyAccessGranted:
			points = u.toCardReaderPoints(msg.data, "access granted")
		case KeyAccessDenied:
			points = u.toCardReaderPoints(msg.data, "access denied")
		case KeySecure:
			points = u.toDoorPoints(msg.data, "not open")
		case KeyDoorHeldOpen, KeyForcedEntry:
			points = u.toDoorPoints(msg.data, "held open")
		default:
			continue
		}

		if err := u.sendPoints(messagesServer, points, msg.udmiTopic, msg.name); err != nil {
			return err
		}
	}

	return nil
}

type udmiMessage struct {
	topic     string
	name      string
	udmiTopic string
	data      mps.Fields
}

func (u *udmiServer) mpsToUDMIMessages(_ context.Context) (task.StopFn, error) {
	mpsMsgs := u.msgs.On("*")
	ctx, stop := context.WithCancel(context.Background())
	go func() {
		defer u.msgs.Off("*", mpsMsgs)

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-mpsMsgs:
				data := msg.Args[0].(mps.Fields)
				name, topic, ok := u.lookupNames(data)
				if !ok {
					continue
				}
				u.udmiMessages.Send(context.Background(), udmiMessage{
					topic:     msg.OriginalTopic,
					name:      name,
					udmiTopic: topic,
					data:      data,
				})
			}
		}
	}()
	return stop, nil
}

func (u *udmiServer) lookupNames(fields mps.Fields) (name, topic string, ok bool) {
	name, ok = u.names.SmartCoreName(fields)
	if !ok {
		u.log.Debug("msg port has no smart core name", zap.Any("msg", fields))
		return "", "", false
	}

	topic, ok = u.names.UDMITopicPrefix(fields)
	if !ok {
		u.log.Debug("msg port has no udmi topic", zap.Any("msg", fields))
		return "", "", false
	}
	return name, topic, true
}

func (u *udmiServer) sendPoints(stream gen.UdmiService_PullExportMessagesServer, points *udmi.PointsEvent, name string, topic string) error {
	body, err := json.Marshal(points)
	if err != nil {
		u.log.Warn("Failed to marshal UDMI payload", zap.Any("points", points), zap.Error(err))
		return nil
	}
	toSend := &gen.PullExportMessagesResponse{Name: name, Message: &gen.MqttMessage{
		Topic:   topic + "/event/pointset/points",
		Payload: string(body),
	}}
	return stream.Send(toSend)
}

func (u *udmiServer) toCardReaderPoints(data mps.Fields, state string) *udmi.PointsEvent {
	return &udmi.PointsEvent{
		"CardRderType":       udmi.PointValue{PresentValue: "CardRder"},
		"CardRderLstRdState": udmi.PointValue{PresentValue: state},
		"CardRderLstUserID":  udmi.PointValue{PresentValue: fmt.Sprintf("%d", data.CardID)},
		"CardRderLstRdTm":    udmi.PointValue{PresentValue: data.Timestamp.Format(time.RFC3339)},
	}
}

func (u *udmiServer) toDoorPoints(data mps.Fields, state string) *udmi.PointsEvent {
	return &udmi.PointsEvent{
		"DrType":    udmi.PointValue{PresentValue: "Dr"},
		"DrState":   udmi.PointValue{PresentValue: state},
		"DrLstOpTm": udmi.PointValue{PresentValue: data.Timestamp.Format(time.RFC3339)},
	}
}
