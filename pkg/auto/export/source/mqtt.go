package source

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/vanti-dev/sc-bos/pkg/auto/export/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/util/pull"
)

func NewMqtt(services Services) task.Starter {
	r := &mqtt{services: services}
	r.Lifecycle = task.NewLifecycle(r.applyConfig)
	r.Logger = services.Logger.Named("smart-core")
	return r
}

type mqtt struct {
	*task.Lifecycle[config.MqttServiceSource]
	services Services
}

func (m *mqtt) applyConfig(ctx context.Context, cfg config.MqttServiceSource) error {
	clients := m.services.Node

	var client gen.MqttServiceClient
	err := clients.Client(&client)
	if err != nil {
		return err
	}

	sent := allowDuplicates()
	if cfg.Duplicates.TrackDuplicates() {
		sent = trackDuplicates(cfg.Duplicates.Cmp())
	}

	grp, ctx := errgroup.WithContext(ctx)

	for _, name := range cfg.RpcNames {
		name := name // save for go routine usage
		puller := &mqttMessagePuller{
			client: client,
			name:   name,
		}
		changes := make(chan *gen.PullMessagesResponse)
		grp.Go(func() error {
			defer close(changes)
			err := pull.Changes[*gen.PullMessagesResponse](ctx, puller, changes, pull.WithLogger(m.Logger.Named(name)))
			if status.Code(err) == codes.Unimplemented {
				m.Logger.Debug("read not supported")
				return nil
			}
			return err
		})
		grp.Go(func() error {
			for change := range changes {
				if commit, publish := sent.Changed(name, change); publish {
					data, err := protojson.MarshalOptions{
						EmitUnpopulated: true,
					}.Marshal(change)
					if err != nil {
						return err
					}
					err = m.services.Publisher.Publish(ctx, name, string(data))
					if err != nil {
						return err
					}
					commit()
				}
			}
			return nil
		})
	}

	go func() {
		err := grp.Wait()
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			m.Logger.Warn("source shut down", zap.Error(err))
		} else {
			m.Logger.Debug("source shut down")
		}
	}()

	return nil
}

type mqttMessagePuller struct {
	client gen.MqttServiceClient
	name   string
}

func (m *mqttMessagePuller) Pull(ctx context.Context, changes chan<- *gen.PullMessagesResponse) error {
	stream, err := m.client.PullMessages(ctx, &gen.PullMessagesRequest{Name: m.name})
	if err != nil {
		return err
	}

	for {
		change, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case changes <- change:
		}
	}
}

func (m *mqttMessagePuller) Poll(ctx context.Context, changes chan<- *gen.PullMessagesResponse) error {
	return status.Error(codes.Unimplemented, "not supported")
}
