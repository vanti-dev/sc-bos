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
	puller := &udmiMessagePuller{
		client: client,
		name:   name,
	}
	changes := make(chan *gen.PullExportMessagesResponse)

	var tasks []task.Task

	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: subscribe")
		err := subscribe(ctx, name, logger, client, pubsub.Subscriber)
		logger.Debug("task end: subscribe", zap.String("err", errStr(err)))
		return task.Normal, err
	})
	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: pullMessages")
		err := pullMessages(ctx, logger, changes, puller)
		logger.Debug("task end: pullMessages", zap.String("err", errStr(err)))
		return task.Normal, err
	})
	tasks = append(tasks, func(ctx context.Context) (task.Next, error) {
		logger.Debug("task run: publishMessages")
		err := publishMessages(ctx, changes, pubsub.Publisher)
		logger.Debug("task end: publishMessages", zap.String("err", errStr(err)))
		return task.Normal, err
	})

	return tasks
}

// subscribe will fetch the topics for the given name, and for each topic an MQTT subscription is created (via
// Subscriber). Messages received for each of those subscriptions is then passed onto the UdmiService using OnMessage.
func subscribe(ctx context.Context, name string, logger *zap.Logger, client gen.UdmiServiceClient, subscriber Subscriber) error {
	res, err := client.DescribeTopics(ctx, &gen.DescribeTopicsRequest{Name: name})
	if err != nil {
		return err
	}
	for _, topic := range res.Topics {
		err := subscriber.Subscribe(ctx, topic, func(_ mqtt.Client, message mqtt.Message) {
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
		if err != nil {
			return err
		}
	}
	<-ctx.Done()
	return ctx.Err()
}

// pullMessages calls pull (with default backoff/delay) and sends each message on the given channel
func pullMessages(ctx context.Context, logger *zap.Logger, changes chan *gen.PullExportMessagesResponse, puller pull.Fetcher[*gen.PullExportMessagesResponse]) error {
	defer close(changes)
	err := pull.Changes[*gen.PullExportMessagesResponse](ctx, puller, changes, pull.WithLogger(logger))
	if status.Code(err) == codes.Unimplemented {
		return nil
	}
	return err
}

// publishMessages waits for messages on the given channel and sends them to the publisher
// ultimately these end up getting sent as MQTT messages
func publishMessages(ctx context.Context, changes chan *gen.PullExportMessagesResponse, publisher Publisher) error {
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

func errStr(err error) string {
	str := "nil"
	if err != nil {
		str = err.Error()
	}
	return str
}
