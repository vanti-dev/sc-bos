package status

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/zone"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/status/config"
)

var Feature = zone.FactoryFunc(func(services zone.Services) service.Lifecycle {
	f := &feature{
		announce: services.Node,
		devices:  services.Devices,
		clients:  services.Node,
		logger:   services.Logger,
	}
	f.Service = service.New(service.MonoApply(f.applyConfig))
	return f
})

type feature struct {
	*service.Service[config.Root]
	announce node.Announcer
	devices  *zone.Devices
	clients  node.Clienter
	logger   *zap.Logger
}

func (f *feature) applyConfig(ctx context.Context, cfg config.Root) error {
	announce := node.AnnounceContext(ctx, f.announce)
	logger := f.logger.With(zap.String("zone", cfg.Name))

	if len(cfg.StatusLogs) > 0 || cfg.StatusLogAll {
		var client gen.StatusApiClient
		if err := f.clients.Client(&client); err != nil {
			return err
		}

		f.devices.Add(cfg.StatusLogs...)
		if cfg.StatusLogAll {
			go func() {
				select {
				case <-ctx.Done():
					return
				case <-f.devices.Frozen():
					names := f.namesThatImplementStatus(ctx, f.devices.Names()...)
					group := &Group{
						client: client,
						names:  names,
						logger: logger,
					}
					announce.Announce(cfg.Name, node.HasTrait(statuspb.TraitName, node.WithClients(gen.WrapStatusApi(group))))
				}
			}()
		} else {
			group := &Group{
				client: client,
				names:  cfg.StatusLogs,
				logger: logger,
			}
			announce.Announce(cfg.Name, node.HasTrait(statuspb.TraitName, node.WithClients(gen.WrapStatusApi(group))))
		}
	}

	return nil
}

func (f *feature) namesThatImplementStatus(ctx context.Context, names ...string) []string {
	var mdClient traits.MetadataApiClient
	if err := f.clients.Client(&mdClient); err != nil {
		f.logger.Warn("cannot discover status devices, metadata api client not supported")
		return nil
	}

	res := make([]string, len(names))
	var wg sync.WaitGroup
	wg.Add(len(names))
	for i, name := range names {
		i, name := i, name
		ctx, stop := context.WithTimeout(ctx, 10*time.Second)
		go func() {
			defer stop()
			defer wg.Done()
			if f.retryNameImplementsStatus(ctx, mdClient, name) {
				res[i] = name
			}
		}()
	}
	wg.Wait()

	var noEmpty []string
	for _, r := range res {
		if r != "" {
			noEmpty = append(noEmpty, r)
		}
	}

	return noEmpty
}

func (f *feature) retryNameImplementsStatus(ctx context.Context, client traits.MetadataApiClient, name string) bool {
	delay := 10 * time.Millisecond
	const inc = 2
	for {
		md, err := client.GetMetadata(ctx, &traits.GetMetadataRequest{Name: name, ReadMask: &fieldmaskpb.FieldMask{Paths: []string{"traits"}}})
		if err != nil {
			select {
			case <-ctx.Done():
				return false
			default:
				select {
				case <-time.After(delay):
					delay *= inc
					continue
				case <-ctx.Done():
					return false
				}
			}

		}
		for _, tmd := range md.Traits {
			if tmd.Name == string(statuspb.TraitName) {
				return true
			}
		}

		return false
	}
}
