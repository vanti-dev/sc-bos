package iothub

import (
	"context"
	"fmt"
	"io"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/vanti-dev/sc-bos/internal/util/mqttutil"
	"github.com/vanti-dev/sc-bos/pkg/auto/azureiot/auth"
)

type Conn interface {
	SendOutputMessage(ctx context.Context, telemetry any) error
	io.Closer
}

type conn struct {
	mqtt           mqtt.Client
	telemetryTopic string
}

func Dial(ctx context.Context, params ConnectionParameters) (Conn, error) {
	mqttOpts, err := MQTTClientOptions(params.HostName, params.DeviceID, "", &auth.LocalSigner{Secret: params.SharedAccessKey})
	if err != nil {
		return nil, err
	}
	mqttClient, err := mqttutil.Connect(ctx, mqttOpts)
	if err != nil {
		return nil, err
	}

	return &conn{
		mqtt:           mqttClient,
		telemetryTopic: fmt.Sprintf("devices/%s/messages/events/", params.DeviceID),
	}, nil
}

func (c *conn) SendOutputMessage(ctx context.Context, telemetry any) error {
	return mqttutil.SendJSON(ctx, c.mqtt, c.telemetryTopic, telemetry)
}

func (c *conn) Close() error {
	c.mqtt.Disconnect(0)
	return nil
}
