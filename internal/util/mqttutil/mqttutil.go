package mqttutil

import (
	"bytes"
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WaitToken(ctx context.Context, token mqtt.Token) (completed bool, err error) {
	select {
	case <-token.Done():
		return true, token.Error()
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func Connect(ctx context.Context, options *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(options)
	_, err := WaitToken(ctx, client.Connect())
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SendJSON converts payload to JSON and sends on topic to client.
// payload will be converted to JSON using the following rules:
//
//  1. if payload is a []byte, string, or bytes.Buffer, it will be sent as-is
//  2. if payload is a proto.Message, it will be converted using protojson.Marshal
//  3. otherwise, it will be converted using json.Marshal
func SendJSON(ctx context.Context, client mqtt.Client, topic string, payload any) error {
	switch v := payload.(type) {
	case []byte, string, bytes.Buffer: // payload is already as expected
	case proto.Message:
		var err error
		payload, err = protojson.Marshal(v)
		if err != nil {
			return err
		}
	default:
		var err error
		payload, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	token := client.Publish(topic, 1, true, payload)
	select {
	case <-token.Done():
		return token.Error()
	case <-ctx.Done():
		return ctx.Err()
	}
}
