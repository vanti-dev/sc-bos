package udmi

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

// tasksForSource returns an array of tasks to run for each UdmiService source/name
// all of these need to be run for the implementation to work
func tasksForSource(name string, logger *zap.Logger, client gen.UdmiServiceClient, pubsub *PubSub) []task.Task {
	messageChanges := make(chan *gen.PullExportMessagesResponse)
	topicChanges := make(chan *gen.PullControlTopicsResponse)

	var tasks []task.Task

	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: pullTopics")
		err := pullTopics(ctx, name, logger, client, topicChanges)
		logger.Debug("task end: pullTopics", zap.Error(err))
		return task.Normal, err
	})
	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: handleTopicChanges")
		err := handleTopicChanges(ctx, name, logger, client, topicChanges, pubsub.Subscriber)
		logger.Debug("task end: handleTopicChanges", zap.Error(err))
		return task.Normal, err
	})
	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: pullMessages")
		err := pullMessages(ctx, name, logger, client, messageChanges)
		logger.Debug("task end: pullMessages", zap.Error(err))
		return task.Normal, err
	})
	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: handleMessages")
		err := handleMessages(ctx, messageChanges, pubsub.Publisher)
		logger.Debug("task end: handleMessages", zap.Error(err))
		return task.Normal, err
	})

	return tasks
}

// pullTopics calls pull for control topics (with default backoff/delay) and sends each message on the given channel
func pullTopics(ctx context.Context, name string, logger *zap.Logger, client gen.UdmiServiceClient, changes chan<- *gen.PullControlTopicsResponse) error {
	puller := &udmiControlTopicsPuller{
		client: client,
		name:   name,
	}
	defer close(changes)
	err := pull.Changes[*gen.PullControlTopicsResponse](ctx, puller, changes, pull.WithLogger(logger))
	if status.Code(err) == codes.Unimplemented {
		return nil
	}
	return err
}

// handleTopicChanges will wait for topic messages on the channel, and for each topic an MQTT subscription is created (via
// Subscriber). Messages received for each of those subscriptions is then passed onto the UdmiService using OnMessage.
func handleTopicChanges(ctx context.Context, name string, logger *zap.Logger, client gen.UdmiServiceClient, changes <-chan *gen.PullControlTopicsResponse, subscriber Subscriber) error {
	subscribeTopic := func(ctx context.Context, topic string) error {
		return subscriber.Subscribe(ctx, topic, func(_ mqtt.Client, message mqtt.Message) {
			payload := string(message.Payload())
			logger.Debug("received MQTT message", zap.String("topic", topic), zap.String("payload", payload))
			_, err := client.OnMessage(ctx, &gen.OnMessageRequest{
				Name: name,
				Message: &gen.MqttMessage{
					Topic:   message.Topic(),
					Payload: payload,
				},
			})
			if err != nil {
				logger.Warn("unable to call OnMessage", zap.Error(err))
			} else {
				logger.Debug("forwarded MQTT message to UDMI service", zap.String("topic", topic), zap.String("payload", payload))
			}
		})
	}

	current := func() {}
	defer func() {
		current()
	}()
	for change := range changes {
		current() // cancel previous subscriptions
		ctx, cancel := context.WithCancel(ctx)
		current = cancel
		// todo: work out topic changes, rather than just restart all
		for _, topic := range change.Topics {
			err := subscribeTopic(ctx, topic)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// pullMessages calls pull for export messages (with default backoff/delay) and sends each message on the given channel
func pullMessages(ctx context.Context, name string, logger *zap.Logger, client gen.UdmiServiceClient, changes chan<- *gen.PullExportMessagesResponse) error {
	puller := &udmiExportMessagePuller{
		client: client,
		name:   name,
	}
	defer close(changes)
	err := pull.Changes[*gen.PullExportMessagesResponse](ctx, puller, changes, pull.WithLogger(logger))
	if status.Code(err) == codes.Unimplemented {
		return nil
	}
	return err
}

// handleMessages waits for messages on the given channel and sends them to the publisher
// ultimately these end up getting sent as MQTT messages
func handleMessages(ctx context.Context, changes <-chan *gen.PullExportMessagesResponse, publisher Publisher) error {
	for change := range changes {
		if change.Message == nil {
			continue
		}
		err := publisher.Publish(ctx, change.Message.Topic, change.Message.Payload)
		if err != nil {
			return err
		}
	}
	return nil
}
