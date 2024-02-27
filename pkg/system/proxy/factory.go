package proxy

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/gentrait/lighttest"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/servicepb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/proxy/config"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"github.com/vanti-dev/sc-bos/pkg/util/slices"
)

const Name = "proxy"

func Factory(holder *lighttest.Holder) system.Factory {
	return &factory{
		server: holder,
	}
}

type factory struct {
	server *lighttest.Holder
}

func (f *factory) New(services system.Services) service.Lifecycle {
	s := &System{
		holder:    f.server,
		self:      services.Node,
		hub:       services.CohortManager,
		ignore:    []string{services.GRPCEndpoint}, // avoid infinite recursion
		tlsConfig: services.ClientTLSConfig,
		announcer: services.Node,
		logger:    services.Logger.Named("proxy"),
	}
	return service.New(service.MonoApply(s.applyConfig))
}

type System struct {
	self      *node.Node
	hub       node.Remote
	ignore    []string
	tlsConfig *tls.Config
	announcer node.Announcer
	logger    *zap.Logger
	holder    *lighttest.Holder
}

// applyConfig runs this system based on the given config.
// This will query the hub for nodes,
// for each node (not ignored) it will query for all children,
// announcing each trait for each child.
func (s *System) applyConfig(ctx context.Context, cfg config.Root) error {
	s.logger.Debug("applying config", zap.Any("config", cfg))
	ignore := append([]string{}, s.ignore...)
	ignore = append(ignore, cfg.Ignore...)

	if cfg.HubMode == "" {
		cfg.HubMode = config.HubModeRemote
	}
	switch cfg.HubMode {
	case config.HubModeRemote:
		hubConn, err := s.hub.Connect(ctx)
		if err != nil {
			return err
		}

		go s.retry(ctx, "announceHub", func(ctx context.Context) (task.Next, error) {
			return s.announceHub(ctx, hubConn)
		})

		go s.retry(ctx, "announceNodes", func(ctx context.Context) (task.Next, error) {
			return s.announceNodes(ctx, hubConn, ignore...)
		})

		s.holder.Fill(gen.NewLightingTestApiClient(hubConn))
	case config.HubModeLocal:
		go s.retry(ctx, "announceNodes", func(ctx context.Context) (task.Next, error) {
			return s.announceLocalNodes(ctx, ignore...)
		})

		var lightTest gen.LightingTestApiClient
		if err := s.self.Client(&lightTest); err != nil {
			s.logger.Warn("no LightingTestApiClient available", zap.Error(err))
		} else {
			s.holder.Fill(lightTest)
		}
	}

	return nil

}

// announceHub adds any routed apis to this node that should be forwarded on to the hub.
// After this you should be able to ask this node to, for example, list alerts on the hub.
func (s *System) announceHub(ctx context.Context, hubConn *grpc.ClientConn) (task.Next, error) {
	// ctx is cancelled when this function returns - i.e. on error
	// this makes sure we're forgetting any announcements in that case.
	// The function will be retried if possible.
	announcer := node.AnnounceContext(ctx, s.announcer)

	// announce any children the hub has
	go s.retry(ctx, "proxyHubChildren", func(ctx context.Context) (task.Next, error) {
		return s.announceNodeChildren(ctx, hubConn)
	})

	// ask the hub what it's name is, and use that for any announcements
	mdClient := traits.NewMetadataApiClient(hubConn)
	stream, err := mdClient.PullMetadata(ctx, &traits.PullMetadataRequest{})
	if err != nil {
		return task.Normal, err
	}

	undo := node.NilUndo // called if the hubName changes to un-announce previous apis
	success := false
	for {
		msg, err := stream.Recv()
		if err != nil {
			if success {
				// at least one request worked so try again immediately
				return task.ResetBackoff, err
			}
			return task.Normal, err
		}
		success = true
		if len(msg.Changes) == 0 {
			continue
		}

		undo()
		lastChange := msg.Changes[len(msg.Changes)-1]
		hubName := lastChange.Metadata.Name

		var undos []node.Undo

		// non-trait routed apis
		undos = append(undos, announcer.Announce(hubName, node.HasMetadata(lastChange.Metadata), node.HasClient(
			gen.NewAlertApiClient(hubConn),
			gen.NewAlertAdminApiClient(hubConn), // Don't do this, we don't want external control of this // SC-469
		)))

		// this is the same logic that you find in announceNodeApis
		undos = append(undos, s.announceServiceApi(announcer, hubConn, hubName))

		// hub traits
		for _, tm := range lastChange.Metadata.Traits {
			traitName := trait.Name(tm.Name)
			if traitName == trait.Metadata {
				continue
			}

			undos = append(undos, s.announceTrait(announcer, hubConn, hubName, traitName))
		}

		undo = node.UndoAll(undos...)
	}
}

// announceNodes fetches all the hubs enrolled nodes and sets up routed apis on this node that proxy those node apis.
func (s *System) announceNodes(ctx context.Context, hubConn *grpc.ClientConn, ignore ...string) (task.Next, error) {
	hubClient := gen.NewHubApiClient(hubConn)
	return s.announceHubNodes(ctx, hubClient, ignore...)
}

func (s *System) announceLocalNodes(ctx context.Context, ignore ...string) (task.Next, error) {
	var hubClient gen.HubApiClient
	if err := s.self.Client(&hubClient); err != nil {
		s.logger.Error("no HubClient available", zap.Error(err))
		return task.Normal, err
	}
	return s.announceHubNodes(ctx, hubClient, ignore...)
}

func (s *System) announceHubNodes(ctx context.Context, hubClient gen.HubApiClient, ignore ...string) (task.Next, error) {
	stream, err := hubClient.PullHubNodes(ctx, &gen.PullHubNodesRequest{})
	if err != nil {
		return task.Normal, err
	}
	knownNodes := make(map[string]context.CancelFunc)
	success := false
	for {
		nodeChanges, err := stream.Recv()
		if err != nil {
			if success {
				return task.ResetBackoff, err
			}
			return task.Normal, err
		}
		success = true

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
				if hubNode.Name == s.self.Name() {
					continue // don't do anything for our own node
				}
				if slices.Contains(hubNode.Address, ignore) {
					continue
				}

				nodeCtx, stopNode := context.WithCancel(ctx)
				knownNodes[hubNode.Address] = stopNode
				nodeConn, err := grpc.DialContext(nodeCtx, hubNode.Address, grpc.WithTransportCredentials(credentials.NewTLS(s.tlsConfig)))
				if err != nil {
					return task.Normal, err
				}

				go s.retry(nodeCtx, "announceNode", func(ctx context.Context) (task.Next, error) {
					return s.announceNode(ctx, hubNode, nodeConn)
				})
			}
		}
	}
}

func (s *System) announceNode(ctx context.Context, hubNode *gen.HubNode, nodeConn *grpc.ClientConn) (task.Next, error) {
	isProxyNode, err := s.isProxy(ctx, nodeConn)
	if err != nil {
		return task.Normal, err
	}

	s.logger.Debug("Proxying node", zap.String("name", hubNode.Name), zap.String("node", hubNode.Address), zap.Bool("isProxy", isProxyNode))
	switch {
	case isProxyNode:
		s.announceProxyNode(ctx, hubNode, nodeConn)
	default:
		s.announceControllerNode(ctx, hubNode, nodeConn)
	}

	<-ctx.Done()
	return task.ResetBackoff, ctx.Err()
}

func (s *System) announceControllerNode(ctx context.Context, hubNode *gen.HubNode, nodeConn *grpc.ClientConn) {
	go s.retry(ctx, "proxyNodeParent", func(ctx context.Context) (task.Next, error) {
		return s.announceNodeParent(ctx, nodeConn, hubNode.Name)
	})
	// proxy any advertised children and child traits
	go s.retry(ctx, "proxyNodeChildren", func(ctx context.Context) (task.Next, error) {
		return s.announceNodeChildren(ctx, nodeConn)
	})

	// proxy any non-trait apis that also use routing
	go s.retry(ctx, "proxyNodeApis", func(ctx context.Context) (task.Next, error) {
		return s.announceNodeApis(ctx, hubNode, nodeConn)
	})
}

func (s *System) announceProxyNode(ctx context.Context, hubNode *gen.HubNode, nodeConn *grpc.ClientConn) {
	go s.retry(ctx, "proxyNodeParent", func(ctx context.Context) (task.Next, error) {
		return s.announceNodeParent(ctx, nodeConn, hubNode.Name)
	})
	// explicitly don't fetch proxy children as they will have the same children as us anyway

	// proxy any non-trait apis that also use routing
	go s.retry(ctx, "proxyNodeApis", func(ctx context.Context) (task.Next, error) {
		return s.announceNodeApis(ctx, hubNode, nodeConn)
	})
}

func (s *System) isProxy(ctx context.Context, nodeConn *grpc.ClientConn) (bool, error) {
	client := gen.NewServicesApiClient(nodeConn)
	req := &gen.ListServicesRequest{Name: "systems"}
	for {
		systems, err := client.ListServices(ctx, req)
		if err != nil {
			return false, err
		}
		for _, sys := range systems.Services {
			if sys.Type == Name {
				return true, nil
			}
		}

		req.PageToken = systems.NextPageToken
		if req.PageToken == "" {
			return false, nil
		}
	}
}

// announceNodeParent discovers all the hub parents trait apis and announces them on this node.
func (s *System) announceNodeParent(ctx context.Context, nodeConn *grpc.ClientConn, name string) (task.Next, error) {
	announcer := node.AnnounceContext(ctx, s.announcer)

	mdClient := traits.NewMetadataApiClient(nodeConn)
	mdStream, err := mdClient.PullMetadata(ctx, &traits.PullMetadataRequest{Name: name})
	if err != nil {
		return task.Normal, err
	}

	undo := func() {}
	success := false
	for {
		mdUpdate, err := mdStream.Recv()
		if err != nil {
			if success {
				return task.ResetBackoff, err // it's an error, but we did succeed with at least one request
			}
			return task.Normal, err
		}
		success = true
		if len(mdUpdate.Changes) == 0 {
			continue
		}

		undo()
		change := mdUpdate.Changes[len(mdUpdate.Changes)-1]

		var undos []node.Undo
		undos = append(undos, announcer.Announce(name, node.HasMetadata(change.Metadata)))
		for _, tm := range change.Metadata.Traits {
			traitName := trait.Name(tm.Name)
			if traitName == trait.Metadata {
				continue
			}
			undos = append(undos, s.announceTrait(announcer, nodeConn, name, traitName))
		}
		undo = node.UndoAll(undos...)
	}
}

// announceNodeChildren discovers and routes all named traits surfaced via nodeConn.
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
	success := false
	for {
		childUpdate, err := childStream.Recv()
		if err != nil {
			if success {
				return task.ResetBackoff, err // it's an error, but we did succeed with at least one request
			}
			return task.Normal, err
		}
		success = true
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
						return s.announceMetadata(mdCtx, nodeConn, child.Name)
					})
					continue
				}
				undos = append(undos, s.announceTrait(announcer, nodeConn, child.Name, traitName))
			}

			if len(undos) == 0 {
				// force the child to exist, even if they don't have any traits
				undos = append(undos, announcer.Announce(child.Name))
			}
			announcedChildren[child.Name] = node.UndoAll(undos...)
		}
	}
}

// announceMetadata pulls the metadata from conn (named name) and updates our nodes local cache of this metadata.
// This makes sure our devices api works locally without having to query all hub nodes.
func (s *System) announceMetadata(ctx context.Context, conn *grpc.ClientConn, name string) (task.Next, error) {
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
	success := false
	for {
		msg, err := stream.Recv()
		if err != nil {
			if success {
				return task.ResetBackoff, err
			}
			return task.Normal, err
		}
		success = true
		for _, change := range msg.Changes {
			lastAnnounce()
			md := change.Metadata
			lastAnnounce = s.announcer.Announce(name, node.HasMetadata(md))
		}
	}
}

// announceNodeApis announces any non-trait based apis that should be routed to a specific node.
// If naming conflicts would occur, name conversion will be performed. For example the services api typically routes
// `drivers`, `automations`, etc which will conflict with each other if we expose them all as is, so we rename to
// `AC-01/drivers`, `AC-01/automations` instead.
func (s *System) announceNodeApis(ctx context.Context, hubNode *gen.HubNode, nodeConn *grpc.ClientConn) (task.Next, error) {
	// ctx is cancelled when this function returns - i.e. on error
	// this makes sure we're forgetting any announcements in that case.
	// The function will be retried if possible.
	announcer := node.AnnounceContext(ctx, s.announcer)

	// service api does name conversion when proxying
	// devices on node AC-01 becomes AC-01/devices on this node
	s.announceServiceApi(announcer, nodeConn, hubNode.Name)

	<-ctx.Done()
	return task.ResetBackoff, ctx.Err()
}

func (s *System) announceServiceApi(announcer node.Announcer, conn *grpc.ClientConn, name string) node.Undo {
	servicesApi := servicepb.RenameApi(gen.NewServicesApiClient(conn), func(n string) string {
		if strings.HasPrefix(n, name+"/") {
			return n[len(name+"/"):]
		}
		return n
	})
	servicesClient := gen.WrapServicesApi(servicesApi)
	var undos []node.Undo
	for _, bucket := range []string{"automations", "drivers", "systems", "zones"} {
		undos = append(undos, announcer.Announce(name+"/"+bucket, node.HasClient(servicesClient)))
	}
	return node.UndoAll(undos...)
}

func (s *System) announceTrait(announcer node.Announcer, nodeConn *grpc.ClientConn, name string, traitName trait.Name) node.Undo {
	var clients []any
	if c := alltraits.APIClient(nodeConn, traitName); c != nil {
		clients = append(clients, c)
	}
	if c := alltraits.HistoryClient(nodeConn, traitName); c != nil {
		clients = append(clients, c)
	}
	if c := alltraits.InfoClient(nodeConn, traitName); c != nil {
		clients = append(clients, c)
	}
	if len(clients) == 0 {
		s.logger.Warn("unable to proxy unknown trait on child",
			zap.String("target", nodeConn.Target()), zap.String("name", name), zap.Stringer("trait", traitName))
	}

	return announcer.Announce(name, node.HasTrait(traitName, node.WithClients(clients...)))
}

func (s *System) retry(ctx context.Context, name string, t task.Task) error {
	attempt := 0
	return task.Run(ctx, func(ctx context.Context) (task.Next, error) {
		attempt++
		next, err := t(ctx)
		if next == task.ResetBackoff {
			// assume some success happened, reset err and attempts
			attempt = 1
		}

		if err == nil {
			return next, err
		}

		switch {
		case attempt == 1:
			s.logger.Warn("failed to run task, will retry", zap.String("task", name), zap.Error(err), zap.Int("attempt", attempt))
		case attempt == 5:
			s.logger.Warn("failed to run task, reducing logging", zap.String("task", name), zap.Error(err), zap.Int("attempt", attempt))
		case attempt%10 == 0:
			s.logger.Debug("failed to run task", zap.String("task", name), zap.Error(err), zap.Int("attempt", attempt))
		}

		return next, err
	}, task.WithRetry(task.RetryUnlimited), task.WithBackoff(10*time.Millisecond, 30*time.Second))
}
