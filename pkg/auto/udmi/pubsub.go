package udmi

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PubSub struct {
	Publisher
	Subscriber
}

type Publisher interface {
	Publish(ctx context.Context, topic string, payload any) error
}

type PublisherFunc func(ctx context.Context, topic string, payload any) error

func (p PublisherFunc) Publish(ctx context.Context, topic string, payload any) error {
	return p(ctx, topic, payload)
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, cb mqtt.MessageHandler) error
}

type SubscriberFunc func(ctx context.Context, topic string, cb mqtt.MessageHandler) error

func (p SubscriberFunc) Subscribe(ctx context.Context, topic string, cb mqtt.MessageHandler) error {
	return p(ctx, topic, cb)
}
