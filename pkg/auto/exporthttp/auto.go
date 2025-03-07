// Package exporthttp provides a simple and flexible way to automate
// the process of POSTing recordings to a REST API endpoint.
// It supports scheduled or regular posting of recordings
//
// This package is designed to be extensible, enabling future
// support for other CMS platforms or REST APIs. It abstracts
// the complexities of API interactions
package exporthttp

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/exporthttp/config"
	exportHttp "github.com/vanti-dev/sc-bos/pkg/auto/exporthttp/http"
	"github.com/vanti-dev/sc-bos/pkg/auto/exporthttp/job"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const AutoName = "exporthttp"

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

	jobs := job.FromConfig(cfg, logger, a.Node)

	if len(jobs) < 1 {
		return nil
	}

	var client *exportHttp.Client

	switch cfg.Auth.Type {
	case config.AuthenticationBearer:
		client = exportHttp.New(exportHttp.WithAuthorizationBearer(cfg.Auth.Token), exportHttp.WithLogger(cfg.Logs, logger))
	default:
		return fmt.Errorf("authentication type %s not supported", cfg.Auth.Type)
	}

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
