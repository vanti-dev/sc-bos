package udmi

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/task/service"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/udmi/config"
)

const AutoType = "udmi"

var Factory auto.Factory = factory{}

type factory struct{}

func (_ factory) New(services auto.Services) service.Lifecycle {
	return NewUDMI(services)
}

func NewUDMI(services auto.Services) service.Lifecycle {
	e := &udmiAuto{services: services}
	e.Service = service.New(service.MonoApply(e.applyConfig))
	e.services.Logger = services.Logger.Named(AutoType)
	return e
}

type udmiAuto struct {
	*service.Service[config.Root]
	services auto.Services
}

func (e *udmiAuto) applyConfig(ctx context.Context, cfg config.Root) error {
	var udmiClient gen.UdmiServiceClient
	err := e.services.Node.Client(&udmiClient)
	if err != nil {
		return err
	}

	client, err := newMqttClient(cfg)
	if err != nil {
		return err
	}

	pubSub := &PubSub{
		Publisher:  mqttPublisher(client, 0, false),
		Subscriber: mqttSubscriber(client, 0),
	}

	connected := client.Connect()
	connected.Wait()
	if connected.Error() != nil {
		return connected.Error()
	}
	e.services.Logger.Debug("connected")

	go func() {
		<-ctx.Done()
		client.Disconnect(5000)
	}()

	grp, ctx := errgroup.WithContext(ctx)
	var tasks []task.Task

	// todo: check if we've already setup the sources
	// todo: how do we stop one source?
	for _, name := range cfg.Sources {
		tasks = append(tasks, tasksForSource(name, e.services.Logger.Named(name), udmiClient, pubSub)...)
	}

	for _, t := range tasks {
		t := t // save for go routine usage
		grp.Go(func() error {
			return task.Run(ctx, t, task.WithRetry(task.RetryUnlimited), task.WithBackoff(time.Millisecond*100, time.Second*10))
		})
	}

	go func() {
		err := grp.Wait()
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			e.services.Logger.Warn("shut down", zap.Error(err))
		} else {
			e.services.Logger.Debug("shut down")
		}
	}()
	return nil
}
