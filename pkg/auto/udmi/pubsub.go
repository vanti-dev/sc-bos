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
	// Subscribe starts a new subscription for the given topic. This will return a nil error if the subscription
	// has been successful, and it has been acknowledged.
	Subscribe(ctx context.Context, topic string, cb mqtt.MessageHandler) error
}

type SubscriberFunc func(ctx context.Context, topic string, cb mqtt.MessageHandler) error

func (p SubscriberFunc) Subscribe(ctx context.Context, topic string, cb mqtt.MessageHandler) error {
	return p(ctx, topic, cb)
}
