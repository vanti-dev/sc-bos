package gateway

import (
	"cmp"
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/internal/util/grpc/reflectionapi"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

// scanRemoteHub calls scanRemoteNode for the hub and all enrolled nodes placing the results in a cohort.
// This blocks until ctx is done.
func (s *System) scanRemoteHub(ctx context.Context, c *cohort, hubConn *grpc.ClientConn) {
	hubClient := gen.NewHubApiClient(hubConn)
	hubNode := newRemoteNode(hubConn.Target(), hubConn)
	hubNode.isHub = true
	c.Nodes.Set(hubNode)

	s.scanRemoteNode(ctx, "hub", hubNode)

	s.retry(ctx, "pull enrolled nodes", func(ctx context.Context) (task.Next, error) {
		return s.pullEnrolledNodes(ctx, hubClient, c)
	}, zap.String("node", hubNode.addr))
}

// scanLocalHub calls scanRemoteNode for all enrolled nodes placing the results in a cohort.
// This blocks until ctx is done.
func (s *System) scanLocalHub(ctx context.Context, c *cohort, hubClient gen.HubApiClient) {
	s.retry(ctx, "pull enrolled nodes", func(ctx context.Context) (task.Next, error) {
		return s.pullEnrolledNodes(ctx, hubClient, c)
	}, zap.String("node", "local"))
}

// scanRemoteNode starts background tasks to collect information about a remote node.
// Errors encountered during tasks are logged and the tasks are retried until ctx is done.
// nodeName is a human-readable name for the node, used in logs to distinguish between remote nodes.
func (s *System) scanRemoteNode(ctx context.Context, nodeName string, n *remoteNode) {
	go s.poll(ctx, func(ctx context.Context) error {
		return s.reflectNode(ctx, n)
	}, zap.String("task", "reflect"), zap.String("nodeName", nodeName), zap.String("node", n.addr))

	go s.retry(ctx, "pull metadata", func(ctx context.Context) (task.Next, error) {
		return s.pullMetadata(ctx, "", n)
	}, zap.String("node", n.addr), zap.String("nodeName", nodeName))

	go s.retry(ctx, "pull systems", func(ctx context.Context) (task.Next, error) {
		return s.pullSystems(ctx, n)
	}, zap.String("node", n.addr), zap.String("nodeName", nodeName))

	go s.retry(ctx, "pull children", func(ctx context.Context) (task.Next, error) {
		return s.pullChildren(ctx, "", n)
	}, zap.String("node", n.addr), zap.String("nodeName", nodeName))
}

// pullMetadata updates node with metadata about name.
// If name is empty, metadata about node itself is fetched and updated,
// otherwise metadata about a child is fetched and updated.
func (s *System) pullMetadata(ctx context.Context, name string, node *remoteNode) (task.Next, error) {
	mdClient := traits.NewMetadataApiClient(node.conn)
	mdStream, err := mdClient.PullMetadata(ctx, &traits.PullMetadataRequest{Name: name})
	if err != nil {
		return neverSucceeded, err
	}

	for msgReceived := false; ; msgReceived = true {
		cs, err := mdStream.Recv()
		if err != nil {
			if msgReceived {
				return someSuccess, err
			}
			return neverSucceeded, err
		}

		if len(cs.Changes) == 0 {
			continue
		}
		c := cs.Changes[len(cs.Changes)-1]
		md := c.Metadata

		// are we fetching metadata for the remoteNode itself or a child?
		self := name == ""
		if self {
			node.Self.Set(remoteDesc{name: md.Name, md: md})
		} else {
			node.Children.Set(remoteDesc{name: name, md: md})
		}
	}
}

// pullSystems populates node.Systems using the ServicesApi of node.
func (s *System) pullSystems(ctx context.Context, node *remoteNode) (task.Next, error) {
	client := gen.NewServicesApiClient(node.conn)
	stream, err := client.PullServices(ctx, &gen.PullServicesRequest{Name: "systems"})
	if err != nil {
		return neverSucceeded, err
	}

	var systems remoteSystems

	for msgReceived := false; ; msgReceived = true {
		msg, err := stream.Recv()
		if err != nil {
			if msgReceived {
				return someSuccess, err
			}
			return neverSucceeded, err
		}

		oldSystems := systems
		for _, c := range msg.Changes {
			id := cmp.Or(c.GetNewValue().GetId(), c.GetOldValue().GetId())
			if id == "" {
				continue // not sure what happened, all services should have an id
			}
			switch id {
			case Name:
				systems.proxy = c.GetNewValue()
			}
		}

		// while we aren't using proto.Equal here, that's fine as the pointers would have changed if we got an update
		if oldSystems != systems {
			node.Systems.Set(systems)
		}
	}
}

// pullChildren uses node's ParentApi to collect the list of children and metadata about them.
// pullChildren blocks while the ParentApi stream is active.
func (s *System) pullChildren(ctx context.Context, name string, node *remoteNode) (task.Next, error) {
	parentClient := traits.NewParentApiClient(node.conn)
	childStream, err := parentClient.PullChildren(ctx, &traits.PullChildrenRequest{Name: name})
	if err != nil {
		return neverSucceeded, err
	}

	// tasks tracks the information we query for each child.
	tasks := tasks{}
	defer tasks.callAll()

	for msgReceived := false; ; msgReceived = true {
		cs, err := childStream.Recv()
		if err != nil {
			if msgReceived {
				return someSuccess, err
			}
			return neverSucceeded, err
		}

		for _, c := range cs.Changes {
			// for anything that isn't an add stop the existing task for the child
			if c.OldValue != nil {
				tasks.remove(c.OldValue.Name)
			}
			if c.NewValue == nil {
				continue // was a deletion
			}

			child := c.NewValue
			childCtx, stop := context.WithCancel(ctx)
			tasks[c.NewValue.Name] = stop
			go s.retry(childCtx, "pull child metadata", func(ctx context.Context) (task.Next, error) {
				return s.pullMetadata(ctx, child.Name, node)
			}, zap.String("name", child.Name), zap.String("node", node.addr))
		}
	}
}

// pullEnrolledNodes calls scanRemoteNode for all enrolled nodes in the hub.
func (s *System) pullEnrolledNodes(ctx context.Context, hubClient gen.HubApiClient, cohort *cohort) (task.Next, error) {
	stream, err := hubClient.PullHubNodes(ctx, &gen.PullHubNodesRequest{})
	if err != nil {
		return neverSucceeded, err
	}

	// keyed by enrolled node address, just like in the cohort
	tasks := tasks{}
	defer tasks.callAll()

	for msgReceived := false; ; msgReceived = true {
		msg, err := stream.Recv()
		if err != nil {
			if msgReceived {
				return someSuccess, err
			}
			return neverSucceeded, err
		}

		for _, c := range msg.Changes {
			// todo: handle the update case instead of treating it as delete+add
			if c.OldValue != nil {
				tasks.remove(c.OldValue.Address)
			}
			if c.NewValue == nil {
				continue
			}

			var stops []func()
			stopAll := func() {
				for _, f := range stops {
					f()
				}
			}

			node := c.NewValue

			// check whether we should process this node at all
			if node.Name == s.self.Name() {
				continue // this is us, don't re-scan ourselves
			}
			if cohort.ShouldIgnore(node.Address) {
				continue
			}

			nodeCtx, stop := context.WithCancel(ctx)
			stops = append(stops, stop)

			conn, err := grpc.DialContext(nodeCtx, node.Address, grpc.WithTransportCredentials(credentials.NewTLS(s.tlsConfig)))
			if err != nil {
				// An error here means there's something wrong with the params we passed to grpc.DialContext.
				// It's unlikely that retrying will help, so skip this node and log an error.
				s.logger.Error("failed to dial remote node", zap.String("address", node.Address), zap.Error(err))
				stopAll()
				continue
			}
			stops = append(stops, func() { conn.Close() })

			remoteNode := newRemoteNode(node.Address, conn)
			// Notify others about the remote node before we scan it.
			// We do it this way to give others a chance to subscribe to the node before we start updating it
			cohort.Nodes.Set(remoteNode)
			s.scanRemoteNode(nodeCtx, cmp.Or(node.Name, node.Address), remoteNode)

			stops = append(stops, func() { cohort.Nodes.Remove(remoteNode) })
			tasks[node.Address] = stopAll
		}
	}
}

// reflectNode uses the reflection service to collect information about the services the node implements.
func (s *System) reflectNode(ctx context.Context, node *remoteNode) error {
	client := reflectionpb.NewServerReflectionClient(node.conn)

	streamCtx, stopStream := context.WithCancel(ctx)
	defer stopStream()
	stream, err := client.ServerReflectionInfo(streamCtx)
	if err != nil {
		return err
	}

	services, err := reflectionapi.ListServices(stream)
	if err != nil {
		return err
	}

	// collect all the type information for the services the node implements
	seenFiles := make(map[string]struct{})
	fileSet := &descriptorpb.FileDescriptorSet{}
	for _, svc := range services {
		serviceFiles, err := reflectionapi.FileContainingSymbol(stream, svc.Name)
		if err != nil {
			return err
		}
		for _, fileBytes := range serviceFiles {
			fileDescPB := &descriptorpb.FileDescriptorProto{}
			err := proto.Unmarshal(fileBytes, fileDescPB)
			if err != nil {
				return fmt.Errorf("unmarshal file descriptor: %w", err)
			}
			if _, ok := seenFiles[fileDescPB.GetName()]; ok {
				continue // skip duplicates
			}
			seenFiles[fileDescPB.GetName()] = struct{}{}
			fileSet.File = append(fileSet.File, fileDescPB)
		}
	}

	stopStream() // don't need the stream anymore, we have all the info we need

	files, err := protodesc.NewFiles(fileSet)
	if err != nil {
		return fmt.Errorf("create file set: %w", err)
	}

	remoteServices := make([]protoreflect.ServiceDescriptor, 0, len(services))
	for _, svc := range services {
		desc, err := files.FindDescriptorByName(protoreflect.FullName(svc.Name))
		if err != nil {
			// shouldn't have an error here as files is computed from the same set of services
			return fmt.Errorf("service %q not found in files: %w", svc.Name, err)
		}

		serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
		if !ok {
			return fmt.Errorf("descriptor %q is not a service, got %T", desc.FullName(), desc)
		}

		remoteServices = append(remoteServices, serviceDesc)
	}

	// we now can fully describe the remote service
	node.Services.Replace(remoteServices)

	return nil
}

// retry semantics for long-running tasks that can succeed and fail at the same time
const (
	neverSucceeded = task.Normal
	someSuccess    = task.ResetBackoff
)
