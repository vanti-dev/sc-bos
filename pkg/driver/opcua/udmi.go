package opcua

import (
	"context"
	"encoding/json"

	"github.com/gopcua/opcua/ua"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver/opcua/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

type Udmi struct {
	config.Trait
	gen.UnimplementedUdmiServiceServer

	logger *zap.Logger
	// the points that have been configured to be monitored. node id -> value source
	monitoredPoints map[string]*config.ValueSource
	pointEvents     map[string]any
	scName          string
	udmiBus         minibus.Bus[*gen.PullExportMessagesResponse]
	udmiConfig      config.UdmiConfig
}

func readUdmiConfig(raw []byte) (cfg config.UdmiConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

func newUdmi(n string, c config.RawTrait, l *zap.Logger) (*Udmi, error) {
	cfg, err := readUdmiConfig(c.Raw)
	if err != nil {
		return nil, err
	}
	udmi := &Udmi{
		logger:          l,
		monitoredPoints: make(map[string]*config.ValueSource),
		pointEvents:     make(map[string]any),
		scName:          n,
		udmiConfig:      cfg,
	}
	// the points in the config are defined as a map of name -> value source to match the pattern we use
	// we get a callback with the node id, so convert to a map of node id -> value source for fast lookup
	for _, p := range cfg.Points {
		udmi.monitoredPoints[p.NodeId] = p
	}
	return udmi, nil
}

func (u *Udmi) sendUdmiMessage(ctx context.Context, node *ua.NodeID, value any) {

	if p, ok := u.monitoredPoints[node.String()]; ok {

		pointName := p.Name
		u.pointEvents[pointName] = p.GetValueFromIntKey(value)

		body, err := json.Marshal(u.pointEvents)
		if err != nil {
			u.logger.Error("failed to marshal points event", zap.Error(err))
			return
		}

		u.udmiBus.Send(ctx, &gen.PullExportMessagesResponse{
			Name: u.scName,
			Message: &gen.MqttMessage{
				Topic:   u.udmiConfig.TopicPrefix + config.PointsEventTopicSuffix,
				Payload: string(body),
			},
		})
	}
}

func (u *Udmi) PullControlTopics(_ *gen.PullControlTopicsRequest, topicsServer gen.UdmiService_PullControlTopicsServer) error {
	// we don't have any control topics, yet
	<-topicsServer.Context().Done()
	return nil
}

func (u *Udmi) OnMessage(context.Context, *gen.OnMessageRequest) (*gen.OnMessageResponse, error) {
	// we don't support doing anything here, yet
	return &gen.OnMessageResponse{}, nil
}

func (u *Udmi) PullExportMessages(_ *gen.PullExportMessagesRequest, server gen.UdmiService_PullExportMessagesServer) error {
	for msg := range u.udmiBus.Listen(server.Context()) {
		err := server.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
