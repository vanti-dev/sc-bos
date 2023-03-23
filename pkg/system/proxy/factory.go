package proxy

import (
	"context"
	"crypto/tls"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/proxy/config"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/util/slices"
)

const Name = "proxy"

var Factory = system.FactoryFunc(func(services system.Services) service.Lifecycle {
	s := &System{
		hub:       services.CohortManager,
		ignore:    []string{services.GRPCEndpoint}, // avoid infinite recursion
		tlsConfig: services.ClientTLSConfig,
		announcer: services.Node,
		logger:    services.Logger.Named("proxy"),
	}
	return service.New(service.MonoApply(s.applyConfig))
})

type System struct {
	hub       node.Remote
	ignore    []string
	tlsConfig *tls.Config
	announcer node.Announcer
	logger    *zap.Logger
}

// applyConfig runs this system based on the given config.
// This will query the hub for nodes,
// for each node (not ignored) it will query for all children,
// announcing each trait for each child.
func (s *System) applyConfig(ctx context.Context, cfg config.Root) error {
	hubConn, err := s.hub.Connect(ctx)
	if err != nil {
		return err
	}

	ignore := append([]string{}, s.ignore...)
	ignore = append(ignore, cfg.Ignore...)

	hubClient := gen.NewHubApiClient(hubConn)
	go s.retry(ctx, "pullHubNodes", func(ctx context.Context) (task.Next, error) {
		stream, err := hubClient.PullHubNodes(ctx, &gen.PullHubNodesRequest{})
		if err != nil {
			return task.Normal, err
		}

		knownNodes := make(map[string]context.CancelFunc)
		for {
			nodeChanges, err := stream.Recv()
			if err != nil {
				return task.ResetBackoff, err
			}

			for _, change := range nodeChanges.Changes {
				// todo: this could be more efficient but I'm low on time right now
				if change.OldValue != nil {
					if stop, ok := knownNodes[change.OldValue.Address]; ok {
						stop()
						delete(knownNodes, change.OldValue.Address)
					}
				}
				if change.NewValue != nil {
					hubNode := change.NewValue
					// todo: consider not processing nodes that are also proxies
					if slices.Contains(hubNode.Address, ignore) {
						continue
					}
					s.logger.Debug("Proxying all devices from node", zap.String("node", hubNode.Address))
					nodeCtx, stopNode := context.WithCancel(ctx)
					knownNodes[hubNode.Address] = stopNode
					nodeConn, err := grpc.DialContext(nodeCtx, hubNode.Address, grpc.WithTransportCredentials(credentials.NewTLS(s.tlsConfig)))
					if err != nil {
						return task.Normal, err
					}
					go s.retry(nodeCtx, "collectChildren", func(ctx context.Context) (task.Next, error) {
						return s.announceNodeChildren(ctx, nodeConn)
					})
				}
			}
		}
	})
	return nil
}

func (s *System) announceNodeChildren(ctx context.Context, nodeConn *grpc.ClientConn) (task.Next, error) {
	// ctx is cancelled when this function returns - i.e. on error
	// this makes sure we're forgetting any announcements in that case.
	// The function will be retried if possible.
	announcer := node.AnnounceContext(ctx, s.announcer)

	parentClient := traits.NewParentApiClient(nodeConn)
	childStream, err := parentClient.PullChildren(ctx, &traits.PullChildrenRequest{})
	if err != nil {
		return task.Normal, err
	}

	announcedChildren := make(map[string]node.Undo)
	for {
		childUpdate, err := childStream.Recv()
		if err != nil {
			return task.ResetBackoff, err // it's an error, but we did succeed with at least one request
		}
		for _, change := range childUpdate.Changes {
			if change.OldValue != nil {
				if undo, ok := announcedChildren[change.OldValue.Name]; ok {
					delete(announcedChildren, change.OldValue.Name)
					undo()
				}
			}
			if change.NewValue == nil {
				continue
			}

			child := change.NewValue
			var undos []node.Undo
			for _, childTrait := range child.Traits {
				traitName := trait.Name(childTrait.Name)
				if traitName == trait.Metadata {
					// treat metadata differently to other traits, we want to proactively get it so the devices api works
					mdCtx, stopMdCtx := context.WithCancel(ctx)
					undos = append(undos, func() {
						stopMdCtx()
					})
					go s.retry(mdCtx, "watchMetadata", func(ctx context.Context) (task.Next, error) {
						return s.announceChildMetadata(mdCtx, nodeConn, child.Name)
					})
					continue
				}
				traitClient := alltraits.APIClient(nodeConn, traitName)
				if traitClient == nil {
					s.logger.Warn("unable to proxy unknown trait on child",
						zap.String("target", nodeConn.Target()), zap.String("name", child.Name), zap.String("trait", childTrait.Name))
					continue
				}

				// todo: we need to better support metadata so our DevicesApi works as expected!
				undo := announcer.Announce(child.Name, node.HasTrait(traitName, node.WithClients(traitClient)))
				undos = append(undos, undo)
			}

			if len(undos) == 0 {
				// force the child to exist, even if they don't have any traits
				undos = append(undos, announcer.Announce(child.Name))
			}
			announcedChildren[child.Name] = node.UndoAll(undos...)
		}
	}
}

func (s *System) announceChildMetadata(ctx context.Context, conn *grpc.ClientConn, name string) (task.Next, error) {
	mdClient := traits.NewMetadataApiClient(conn)
	stream, err := mdClient.PullMetadata(ctx, &traits.PullMetadataRequest{Name: name})
	if err != nil {
		return task.Normal, err
	}
	// we aren't using node.AnnounceContext here because we want to undo both when ctx is done and if the md is updated.
	lastAnnounce := node.NilUndo
	go func() {
		<-ctx.Done()
		lastAnnounce()
	}()
	for {
		msg, err := stream.Recv()
		if err != nil {
			return task.ResetBackoff, err
		}
		for _, change := range msg.Changes {
			lastAnnounce()
			md := change.Metadata
			lastAnnounce = s.announcer.Announce(name, node.HasMetadata(md))
		}
	}
}

func (s *System) retry(ctx context.Context, name string, t task.Task) error {
	return task.Run(ctx, t, task.WithRetry(task.RetryUnlimited), task.WithBackoff(10*time.Millisecond, 30*time.Second), task.WithErrorLogger(s.logger.With(zap.String("task", name))))
}
