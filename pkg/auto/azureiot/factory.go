package azureiot

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/internal/iothub"
	"github.com/vanti-dev/sc-bos/internal/iothub/auth"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverType = "azureiot"

const minPollInterval = 5 * time.Second

var Factory auto.Factory = factory{}

type factory struct{}

func (f factory) New(services auto.Services) service.Lifecycle {
	services.Logger = services.Logger.Named(DriverType)
	a := &Auto{services: services}
	a.Service = service.New(service.MonoApply(a.applyConfig), service.WithParser(ParseConfig))
	return a
}

type Auto struct {
	*service.Service[Config]
	services auto.Services
}

func (a *Auto) applyConfig(ctx context.Context, cfg Config) error {
	if len(cfg.Devices) == 0 {
		a.services.Logger.Warn("no devices configured; no polling will happen")
		return nil
	}

	// load the group key from string or disk, but only if a device will need it
	var needsGroupKey bool // group key is required if a device lacks a connection string
	for _, deviceCfg := range cfg.Devices {
		if deviceCfg.ConnectionString == "" {
			needsGroupKey = true
		}
	}
	var groupKey auth.SASKey
	if needsGroupKey {
		var err error
		groupKey, err = loadGroupKey(cfg)
		if err != nil {
			return fmt.Errorf("failed to load group key: %w", err)
		}
	}

	// the group of tasks that pull from device traits and publish to IoT Hub
	grp, ctx := errgroup.WithContext(ctx)
	for _, device := range cfg.Devices {
		devDialler, err := diallerFromConfig(device, cfg.IDScope, groupKey)
		if err != nil {
			a.services.Logger.Error("device not initialised due to invalid configuration", zap.String("device", device.Name), zap.Error(err))
			continue
		}

		stream := make(chan proto.Message)
		grp.Go(func() error {
			defer close(stream)
			return a.pullDevice(ctx, stream, device)
		})
		grp.Go(func() error {
			a.sendOutputMessages(ctx, devDialler, stream)
			return nil
		})
	}

	go func() {
		err := grp.Wait()
		switch {
		case err == nil:
			// if everything completed without error, then nothing did anything anyway
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			// context errors are expected when the service is shutting down
		default:
			a.services.Logger.Error("one or more background tasks stopped", zap.Error(err))
		}
	}()

	return nil
}

func (a *Auto) sendOutputMessages(ctx context.Context, dialler dialler, stream <-chan proto.Message) {
	var conn iothub.Conn
	buf := list.New()
	notSent := func(m any) {
		for buf.Len() > DefaultBacklogSize {
			drop := buf.Front()
			buf.Remove(drop)
			a.services.Logger.Warn("dropping message due to full buffer", zap.Any("message", drop.Value))
		}
		a.services.Logger.Debug("buffering message", zap.Any("message", m))
		buf.PushBack(m)
	}
	sendBacklog := func(m any) error {
		for m := buf.Front(); m != nil; m = m.Next() {
			err := conn.SendOutputMessage(ctx, m.Value)
			if err != nil {
				return err
			}
			buf.Remove(m)
		}
		return nil
	}
	for msg := range stream {
		if conn == nil {
			var err error
			conn, err = dialler.Dial(ctx)
			if err != nil {
				a.services.Logger.Error("failed to dial IoT Hub", zap.Error(err))
				notSent(msg)
				continue
			}
		}
		err := sendBacklog(msg)
		if err != nil {
			a.services.Logger.Warn("failed to send backlog message to IoT Hub", zap.Error(err))
			notSent(msg)
			continue
		}
		err = conn.SendOutputMessage(ctx, msg)
		if err != nil {
			a.services.Logger.Warn("failed to send message to IoT Hub", zap.Error(err))
			notSent(msg)
			continue
		}
	}
}

func loadGroupKey(cfg Config) (auth.SASKey, error) {
	if cfg.GroupKey != "" {
		return auth.ParseSASKey(cfg.GroupKey)
	}

	raw, err := os.ReadFile(cfg.GroupKeyFile)
	if err != nil {
		return nil, err
	}
	return auth.ParseSASKey(string(raw))
}
