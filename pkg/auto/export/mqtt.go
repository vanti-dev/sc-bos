package export

import (
	"context"
	"fmt"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/export/config"
	"github.com/vanti-dev/sc-bos/pkg/auto/export/source"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

var MQTTFactory auto.Factory = factory{}

type factory struct{}

func (_ factory) New(services auto.Services) service.Lifecycle {
	return NewMQTTExport(services)
}

func (_ factory) AddSupport(supporter node.Supporter) {
	r := gen.NewMqttServiceRouter()
	supporter.Support(node.Routing(r), node.Clients(gen.WrapMqttService(r)))
}

func NewMQTTExport(services auto.Services) service.Lifecycle {
	e := &mqttExport{services: services}
	e.Service = service.New(service.MonoApply(e.applyConfig))
	e.services.Logger = services.Logger.Named("export.mqtt")
	return e
}

type mqttExport struct {
	*service.Service[config.Root]
	services auto.Services
}

func (e *mqttExport) applyConfig(ctx context.Context, cfg config.Root) error {
	client, err := newMqttClient(cfg)
	if err != nil {
		return err
	}
	services := source.Services{
		Services:  e.services,
		Publisher: mqttPublisher(client, 0, false),
	}

	connected := client.Connect()
	connected.Wait()
	if connected.Error() != nil {
		return connected.Error()
	}

	go func() {
		<-ctx.Done()
		client.Disconnect(5000)
	}()

	return configureSources(ctx, services, cfg.Sources)
}

func newMqttClient(cfg config.Root) (mqtt.Client, error) {
	options, err := cfg.Broker.ClientOptions()
	if err != nil {
		return nil, err
	}
	return mqtt.NewClient(options), nil
}

var supportedSources = map[string]func(source.Services) task.Starter{
	"bacnet":     source.NewBacnet,
	"mqtt":       source.NewMqtt,
	"smart-core": source.NewSmartCore,
}

func configureSources(ctx context.Context, services source.Services, cfgs []config.RawSource) error {
	var started []task.Starter
	go func() {
		<-ctx.Done()
		var err error
		for _, impl := range started {
			if task.Stoppable(impl) {
				err = multierr.Append(err, task.Stop(impl))
			}
		}
		if err != nil {
			services.Logger.Warn("Failed to cleanly stop after ctx done", zap.Error(err))
		}
	}()

	var allErrs error
	for _, cfg := range cfgs {
		f, ok := supportedSources[cfg.Type]
		if !ok {
			allErrs = multierr.Append(allErrs, fmt.Errorf("unsupported type %v", cfg.Type))
			continue
		}
		impl := f(services)
		if err := impl.Start(ctx); err != nil {
			allErrs = multierr.Append(allErrs, fmt.Errorf("start %s %w", cfg.Name, err))
			continue
		}
		// keep track so we can stop them if ctx ends
		started = append(started, impl)

		if task.Configurable(impl) {
			if err := task.Configure(impl, cfg.Raw); err != nil {
				allErrs = multierr.Append(allErrs, fmt.Errorf("configure %s %w", cfg.Name, err))
			}
		}
	}
	return allErrs
}

func mqttPublisher(client mqtt.Client, qos byte, retained bool) source.Publisher {
	return source.PublisherFunc(func(ctx context.Context, topic string, payload any) error {
		token := client.Publish(topic, qos, retained, payload)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-token.Done():
			return token.Error()
		}
	})
}
