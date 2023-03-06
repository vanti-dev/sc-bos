package prioritypb

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/light"
	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/util/chans"
	"golang.org/x/sync/errgroup"
)

func NewLightPriority(client traits.LightApiClient, realName string, opts ...Option) node.SelfAnnouncer {
	config := readOpts(opts...)
	impl := &Light{
		router:      client,
		realName:    realName,
		defaultName: config.defaultName(realName),

		model: newModel[*traits.Brightness](config.fqns(realName, config.slotNames...)...),
		now:   time.Now,
	}

	return node.SelfAnnouncerFunc(func(a node.Announcer) node.Undo {
		var undos []node.Undo
		// announce the {name}/priority device
		undos = append(undos, a.Announce(config.suffix(realName),
			node.HasMetadata(config.metadata(realName)),
			node.HasClient(gen.WrapPriorityApi(newModelServer(impl.model))),
			node.HasTrait(trait.Light, node.WithClients(light.WrapApi(impl)))))
		for _, name := range impl.names {
			client := light.WrapApi(light.NewMemoryDevice())
			// announce the {name}/priority/{slot} devices
			undos = append(undos, a.Announce(name, node.HasTrait(trait.Light, node.WithClients(client))))
		}

		ctx, stop := context.WithCancel(context.Background())
		undos = append(undos, func() {
			stop()
		})
		go impl.collect(ctx)

		return node.UndoAll(undos...)
	})
}

type Light struct {
	traits.UnimplementedLightApiServer

	router      traits.LightApiClient
	realName    string
	defaultName string

	*model[*traits.Brightness]

	now func() time.Time
}

func (l *Light) collect(ctx context.Context) error {
	// when invoked, these pullers fetch changes from the priority slots
	pullers := make([]pull.Fetcher[*traits.Brightness], len(l.names))
	for i, name := range l.names {
		pullers[i] = &lightBrightnessFetcher{
			now:    l.now,
			client: l.router,
			name:   name,
		}
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return collect(ctx, l.list, pullers...)
	})
	group.Go(func() error {
		for val := range l.model.list.Listen(ctx) {
			if !val.Set {
				continue
			}
			if _, err := l.router.UpdateBrightness(ctx, &traits.UpdateBrightnessRequest{Name: l.realName, Brightness: val.Value}); err != nil {
				return err
			}
		}
		return nil
	})

	return group.Wait()
}

func (l *Light) UpdateBrightness(ctx context.Context, request *traits.UpdateBrightnessRequest) (*traits.Brightness, error) {
	request.Name = l.defaultName
	return l.router.UpdateBrightness(ctx, request)
}

func (l *Light) GetBrightness(ctx context.Context, request *traits.GetBrightnessRequest) (*traits.Brightness, error) {
	request.Name = l.realName
	return l.router.GetBrightness(ctx, request)
}

func (l *Light) PullBrightness(request *traits.PullBrightnessRequest, server traits.LightApi_PullBrightnessServer) error {
	request.Name = l.realName
	stream, err := l.router.PullBrightness(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

type lightBrightnessFetcher struct {
	client traits.LightApiClient
	name   string
	now    func() time.Time
}

func (l *lightBrightnessFetcher) Pull(ctx context.Context, changes chan<- *traits.Brightness) error {
	stream, err := l.client.PullBrightness(ctx, &traits.PullBrightnessRequest{Name: l.name, UpdatesOnly: true})
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range msg.Changes {
			if err := chans.SendContext(ctx, changes, change.Brightness); err != nil {
				return err
			}
		}
	}
}

func (l *lightBrightnessFetcher) Poll(ctx context.Context, changes chan<- *traits.Brightness) error {
	msg, err := l.client.GetBrightness(ctx, &traits.GetBrightnessRequest{Name: l.name})
	if err != nil {
		return err
	}
	return chans.SendContext(ctx, changes, msg)
}
