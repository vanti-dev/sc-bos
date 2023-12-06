package mqttutil

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func SendJSON(ctx context.Context, client mqtt.Client, topic string, payload any) error {
	payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	token := client.Publish(topic, 1, true, payload)
	select {
	case <-token.Done():
		return token.Error()
	case <-ctx.Done():
		return ctx.Err()
	}
}
