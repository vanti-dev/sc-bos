package statusalerts

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/statusalerts/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const AutoName = "statusalerts"

var Factory auto.Factory = factory{}

type factory struct{}

func (f factory) New(services auto.Services) service.Lifecycle {
	a := &autoImpl{Services: services}
	a.Service = service.New(service.MonoApply(a.applyConfig), service.WithParser(config.ReadBytes))
	a.Logger = a.Logger.Named(AutoName)
	return a
}

type autoImpl struct {
	*service.Service[config.Root]
	auto.Services
}

func (a *autoImpl) applyConfig(ctx context.Context, cfg config.Root) error {
	logger := a.Logger
	if cfg.Destination == "" {
		cfg.Destination = a.Node.Name()
	}
	destName := cfg.Destination
	logger = logger.With(zap.String("destination", destName))

	alertAdminClient := gen.NewAlertAdminApiClient(a.Node.ClientConn())
	statusClient := gen.NewStatusApiClient(a.Node.ClientConn())

	if cfg.DelayStart != nil {
		time.Sleep(cfg.DelayStart.Duration)
	}

	var tasks namedTasks
	pullFrom := func(source config.Source) {
		logger := logger.With(zap.String("name", source.Name))
		err := tasks.Run(ctx, source.Name, tasksForSource(source, destName, statusClient, alertAdminClient, logger),
			task.WithRetry(task.RetryUnlimited), task.WithBackoff(time.Millisecond*100, time.Second*10))
		if errors.Is(err, ErrAlreadyRunning) {
			// cool, I guess someone else beat us to it
			return
		}
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("shut down unexpectedly", zap.Error(err))
			return
		}
		logger.Debug("subscription stopped")
	}

	// setup manually configured sources
	for _, name := range cfg.Sources {
		go pullFrom(name)
	}

	// setup discovered sources
	if cfg.DiscoverSources {
		go func() {
			// we use a replacer here because it's an easy and memory efficient way to do prefix matching.
			// strings.Replacer uses a Trie internally.
			ignore := func() *strings.Replacer {
				replacements := make([]string, len(cfg.IgnorePrefixes)*2)
				for _, prefix := range cfg.IgnorePrefixes {
					replacements = append(replacements, prefix, "!")
				}
				return strings.NewReplacer(replacements...)
			}()

			for { // loop in case we get errors
				select {
				case <-ctx.Done():
					return
				default:
				}
				for change := range a.Node.PullAllMetadata(ctx, resource.WithReadPaths(&traits.Metadata{}, "traits", "location", "membership")) {
					if s := ignore.Replace(change.Name); len(s) == 0 || s[0] == '!' {
						continue // ignore
					}
					hadTrait, hasTrait := hasStatusTrait(change.OldValue), hasStatusTrait(change.NewValue)
					switch {
					case hadTrait && !hasTrait: // remove
						err := tasks.Stop(change.Name)
						if err != nil && !errors.Is(err, ErrNotRunning) {
							logger.Debug("error during stop", zap.String("name", change.Name), zap.Error(err))
						}
					case !hadTrait && hasTrait: // add
						source := config.Source{
							Name:      change.Name,
							Floor:     change.NewValue.GetLocation().GetFloor(),
							Zone:      change.NewValue.GetLocation().GetZone(),
							Subsystem: change.NewValue.GetMembership().GetSubsystem(),
						}
						go pullFrom(source)
					}
				}
			}
		}()
	}

	return nil
}

func hasStatusTrait(md *traits.Metadata) bool {
	for _, t := range md.GetTraits() {
		if t.Name == statuspb.TraitName.String() {
			return true
		}
	}
	return false
}
