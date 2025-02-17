// Package wordpress provides a simple and flexible way to automate
// the process of posting recordings to a WordPress REST API endpoint.
// It supports scheduled or regular posting of recordings
//
// This package is designed to be extensible, enabling future
// support for other CMS platforms or REST APIs. It abstracts
// the complexities of API interactions
package wordpress

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress/config"
	wordpressHttp "github.com/vanti-dev/sc-bos/pkg/auto/wordpress/http"
	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress/job"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const AutoName = "wordpress"

var Factory auto.Factory = factory{}

type factory struct{}

type autoImpl struct {
	*service.Service[config.Root]
	auto.Services
}

func (f factory) New(services auto.Services) service.Lifecycle {
	a := &autoImpl{Services: services}
	a.Service = service.New(service.MonoApply(a.applyConfig), service.WithParser(config.ReadBytes))
	a.Logger = a.Logger.Named(AutoName)
	return a
}

func (a *autoImpl) applyConfig(ctx context.Context, cfg config.Root) error {
	logger := a.Logger.Named(cfg.Name).With(zap.String("baseUrl", cfg.BaseUrl))

	jobs := job.FromConfig(cfg, logger)

	if len(jobs) < 1 {
		return nil
	}

	// set up
	for _, jb := range jobs {
		for _, cli := range jb.GetClients() {
			if err := a.Node.Client(cli); err != nil {
				logger.Warn(fmt.Sprintf("failed to create %s client", jb.GetName()), zap.Error(err))

				return err
			}
		}
	}

	var opts []wordpressHttp.Option

	for _, auth := range cfg.Auths {
		switch auth.Type {
		case config.AuthenticationBearer:
			opts = append(opts, wordpressHttp.WithAuthorizationBearer(auth.Token))
		case config.AuthenticationBasic:
			opts = append(opts, wordpressHttp.WithAuthorizationBasic(auth.Token))
		default:
			return fmt.Errorf("authentication type %s not supported", auth.Type)
		}
	}

	opts = append(opts, wordpressHttp.WithLogger(cfg.Logs, logger))
	client := wordpressHttp.New(opts...)

	go func() {
		mulpx := job.Multiplex(jobs...)

		// tear down
		defer func() {
			for _, jb := range jobs {
				jb.Stop()
			}
			mulpx.WaitForDone()
		}()

		// run
		for {
			select {
			case jb := <-mulpx.C:
				if err := jb.Do(ctx, client.Post); err != nil {
					logger.Warn(fmt.Sprintf("failed to run %s", jb.GetName()), zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
