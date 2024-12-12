package gateway

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alphapb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/servicepb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system/gateway/internal/rx"
	"github.com/vanti-dev/sc-bos/pkg/util/slices"
)

// announceCohort announces information about the cohort as if it were present on this system.
// This includes announcing remote names and apis via the local DevicesApi and other apis.
// This blocks until ctx is done.
func (s *System) announceCohort(ctx context.Context, c *cohort) {
	tasks := tasks{}
	defer tasks.callAll()

	table := &table{
		services:    &counts{m: make(map[string]int)},
		serviceUndo: newSyncMap[node.Undo](),
	}

	runAnnouncer := func(n *remoteNode) {
		nodeCtx, stop := context.WithCancel(ctx)
		tasks[n.addr] = stop
		a := &announcer{
			System: s,
			logger: s.logger.With(
				zap.String("remoteAddr", n.addr),
				zap.Bool("isHub", n.isHub),
			),
			node:      n,
			Announcer: node.AnnounceContext(nodeCtx, s.announcer),
			table:     table,
		}
		go a.announceRemoteNode(nodeCtx)
	}

	nodes, nodeChanges := c.Nodes.Sub(ctx)
	nodes.All(func(_ int, n *remoteNode) bool {
		runAnnouncer(n)
		return true
	})

	for nodeChange := range nodeChanges {
		if nodeChange.Old != nil {
			tasks.remove(nodeChange.Old.addr)
		}
		if nodeChange.New == nil {
			continue // was a deletion
		}
		runAnnouncer(nodeChange.New)
	}
}

// announcer provides the functionality for announcing parts of a cohort.
// Typically, the entrypoint is announceRemoteNode.
type announcer struct {
	*System
	logger *zap.Logger
	table  *table
	node   *remoteNode
	node.Announcer
}

// announceRemoteNode announces the remote node as if it were present on this system.
// This blocks until ctx is done.
// All side effects are undone when the method returns.
func (a *announcer) announceRemoteNode(ctx context.Context) {
	self, selfChanges := a.node.Self.Sub(ctx)
	undoSelf := a.announceName(self)

	// Remote nodes that are a proxy will likely also be announcing all the same children and apis we will.
	// To avoid circular routing we don't announce advertised children for remote nodes that are a proxy.
	// To avoid extra dynamic proxying we also don't announce reflected services for remote nodes that are a proxy.
	// As being a proxy is something that can change,
	// we track and update our advertised children and services when the proxy status changes.
	systems, systemChanges := a.node.Systems.Sub(ctx)
	isProxy := func() bool {
		// we are intentionally ignoring the loading state of the proxy system,
		// under the assumption that whether the remote proxy system has loaded or not it's
		// still intending to be a proxy eventually.
		return systems.proxy.GetActive()
	}
	wasProxy := isProxy()

	// The services we proxy are dependent on whether the remote node is a proxy or not.
	// These track and update our advertised services when the proxy status changes.
	var serviceChanges <-chan rx.Change[protoreflect.ServiceDescriptor]
	var stopServiceSub context.CancelFunc
	var undoServices tasks
	setupServiceSub := func() {
		var serviceCtx context.Context
		serviceCtx, stopServiceSub = context.WithCancel(ctx)
		var services *slices.Sorted[protoreflect.ServiceDescriptor]
		services, serviceChanges = a.node.Services.Sub(serviceCtx)
		undoServices = a.announceRemoteServices(services.All)
	}
	closeServiceSub := func() {
		if stopServiceSub != nil {
			stopServiceSub()
		}
		undoServices.callAll()
		serviceChanges = nil
		undoServices = nil
	}
	defer closeServiceSub()

	// The children we proxy are dependent on whether the remote node is a proxy or not.
	// These track and update our advertised children when the proxy status changes.
	var childChanges <-chan rx.Change[remoteDesc]
	var stopChildSub context.CancelFunc
	var undoChildren tasks
	setupChildSub := func() {
		var childCtx context.Context
		childCtx, stopChildSub = context.WithCancel(ctx)
		var children *slices.Sorted[remoteDesc]
		children, childChanges = a.node.Children.Sub(childCtx)
		undoChildren = a.announceMetadataTraitsSet(children.All)
	}
	closeChildSub := func() {
		if stopChildSub != nil {
			stopChildSub()
		}
		undoChildren.callAll()
		childChanges = nil
		undoChildren = nil
	}
	defer closeChildSub()

	// The types and services we return via reflection depends on whether the remote node is a proxy or not.
	// These update the reflection server when the proxy status changes.
	setupReflection := func() {
		a.reflection.Add(a.node.conn)
	}
	closeReflection := func() {
		a.reflection.Remove(a.node.conn)
	}
	defer closeReflection()

	// helper for switching between proxy and non-proxy mode
	switchProxyMode := func(isProxy bool) {
		if isProxy {
			closeChildSub()
			closeServiceSub()
			closeReflection()
		} else {
			setupChildSub()
			setupServiceSub()
			setupReflection()
		}
	}

	switchProxyMode(isProxy())

	for {
		select {
		case <-ctx.Done():
			// a done ctx should clean up all the subscriptions and announced names
			return
		case self = <-selfChanges:
			undoSelf()
			undoSelf = node.UndoAll(
				a.announceName(self),
				a.announceServiceApi(self.name),
			)
		case c, ok := <-serviceChanges:
			if !ok {
				continue // we stopped watching
			}
			if c.Type != rx.Add {
				undoServices.remove(string(c.Old.FullName()))
			}
			if c.Type != rx.Remove {
				undoServices[string(c.New.FullName())] = a.announceRemoteService(c.New)
			}
		case systems = <-systemChanges:
			isProxy := isProxy()
			if isProxy != wasProxy {
				switchProxyMode(isProxy)
				wasProxy = isProxy
			}
		case c, ok := <-childChanges:
			if !ok {
				continue // we stopped watching
			}
			if c.Type != rx.Add {
				undoChildren.remove(c.Old.name)
			}
			if c.Type != rx.Remove {
				undoChildren[c.New.name] = a.announceName(c.New)
			}
		}
	}
}

// announceServiceApi adds proxying for the ServicesApi to a.node.
// As services were historically named `drivers`, `automations`, `systems`, and `zones`,
// we rename them to `{name}/drivers`, `{name}/automations`, etc. to avoid conflicts with
// the same names from other remote nodes.
func (a *announcer) announceServiceApi(name string) node.Undo {
	servicesApi := servicepb.RenameApi(gen.NewServicesApiClient(a.node.conn), func(n string) string {
		if strings.HasPrefix(n, name+"/") {
			return n[len(name+"/"):]
		}
		return n
	})
	var undos []node.Undo
	for _, bucket := range []string{"automations", "drivers", "systems", "zones"} {
		undos = append(undos, a.Announce(name+"/"+bucket,
			node.HasServer(gen.RegisterServicesApiServer, servicesApi),
			node.HasNoAutoMetadata(),
		))
	}
	return node.UndoAll(undos...)
}

// announceName allows this node to answer requests aimed at the given remoteDesc.
// This includes named RPCs (like trait requests) and DeviceApi / ParentApi requests that would include this name in their responses.
func (a *announcer) announceName(d remoteDesc) node.Undo {
	if d.name == "" {
		return node.NilUndo
	}

	return a.Announce(d.name,
		node.HasMetadata(d.md),
		node.HasProxy(a.node.conn),
	)
}

// announceMetadataTraitsSet announces the metadata traits for each remoteDesc in seq.
func (a *announcer) announceMetadataTraitsSet(seq seq2[int, remoteDesc]) tasks {
	dst := tasks{}
	seq(func(_ int, d remoteDesc) bool {
		dst[d.name] = a.announceName(d)
		return true
	})
	return dst
}

// announceRemoteService updates this node to respond to requests for the given remoteService.
func (a *announcer) announceRemoteService(rs protoreflect.ServiceDescriptor) node.Undo {
	if a.ignoreRemoteService(rs) {
		return node.NilUndo
	}
	name := string(rs.FullName())
	if a.table.services.Inc(name) == 1 {
		undo := a.announceRemoteServiceApis(rs)
		if undo != nil { // can be nil if there's nothing to undo
			a.table.serviceUndo.Set(name, undo)
		}
	}

	return func() {
		if a.table.services.Dec(name) == 0 {
			// we were the last to remove the service, clean everything up
			undo, ok := a.table.serviceUndo.Del(name)
			if ok {
				undo()
			}
		}
	}
}

// announceRemoteServices updates this node to respond to requests for each remoteService in seq.
func (a *announcer) announceRemoteServices(seq seq2[int, protoreflect.ServiceDescriptor]) tasks {
	dst := tasks{}
	seq(func(_ int, rs protoreflect.ServiceDescriptor) bool {
		name := string(rs.FullName())
		dst[name] = a.announceRemoteService(rs)
		return true
	})
	return dst
}

func (a *announcer) announceRemoteServiceApis(rs protoreflect.ServiceDescriptor) node.Undo {
	srv := node.ReflectedConnService(rs, a.node.conn)
	var undos []node.Undo

	// which type of proxying should each method use?
	switch {
	case srv.NameRoutable():
		// routes will be added by device announcements as this service is routable by name
		err := a.self.SupportService(srv)
		if err != nil {
			a.logger.Warn("failed to announce routable service",
				zap.String("service", string(rs.FullName())),
				zap.Error(err),
			)
			break
		}
		a.logger.Debug("routable service announced",
			zap.String("service", string(rs.FullName())),
		)
	case a.node.isHub:
		// route everything to the hub
		undo, err := a.self.AnnounceService(srv)
		if err != nil {
			a.logger.Warn("failed to announce non-routable hub service",
				zap.String("service", string(rs.FullName())),
				zap.Error(err),
			)
			break
		}
		a.logger.Debug("non-routable hub service announced",
			zap.String("service", string(rs.FullName())),
		)
		undos = append(undos, undo)
	default:
		// Found a non-routable service on a non-hub node.
		// We didn't think we had any of these so log it to remind us.
		a.logger.Warn("unable to announce non-routable service on non-hub node", zap.String("service", string(rs.FullName())))
		return nil
	}

	return node.UndoAll(undos...)
}

// ignoreRemoteService returns true for service descriptors that we shouldn't proxy automatically.
// All services that return true are proxied in some other way.
func (a *announcer) ignoreRemoteService(rs protoreflect.ServiceDescriptor) bool {
	name := string(rs.FullName())
	switch name {
	// services that are handled explicitly via other mechanisms
	case
		gen.DevicesApi_ServiceDesc.ServiceName,                // handled by the node outside the gateway service
		gen.EnrollmentApi_ServiceDesc.ServiceName,             // handled by the app controller during boot
		gen.ServicesApi_ServiceDesc.ServiceName,               // see announceServiceApi
		reflectionpb.ServerReflection_ServiceDesc.ServiceName, // see setup/closeReflection in announceRemoteNode
		reflectionv1alphapb.ServerReflection_ServiceDesc.ServiceName:
		return true
	}
	return false
}

// table tracks state across all remote nodes.
type table struct {
	services    *counts             // keyed by service full name, counts how many times we've seen services across all remote nodes
	serviceUndo *syncMap[node.Undo] // keyed by service full name
}

// syncMap is a simple synchronised map with string keys.
type syncMap[T any] struct {
	mu sync.RWMutex
	m  map[string]T
}

func newSyncMap[T any]() *syncMap[T] {
	return &syncMap[T]{m: make(map[string]T)}
}

func (r *syncMap[T]) Get(k string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.m[k]
	return v, ok
}

// Set adds or replaces a value in the map, returning the old value and true if the key was already present.
func (r *syncMap[T]) Set(k string, v T) (T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	old, ok := r.m[k]
	r.m[k] = v
	return old, ok
}

// Del removes a value from the map, returning the old value and true if the key was present.
func (r *syncMap[T]) Del(k string) (T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.m[k]
	if ok {
		delete(r.m, k)
	}
	return v, ok
}

// counts acts like a synchronous multiset of strings.
type counts struct {
	mu sync.Mutex
	m  map[string]int
}

// Inc adds one to the count of k, returning the new count.
func (r *counts) Inc(k string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[k]++
	return r.m[k]
}

// Dec subtracts one from the count of k, returning the new count.
// Dec will never cause the count to become less than zero.
func (r *counts) Dec(k string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[k]--
	if r.m[k] <= 0 {
		delete(r.m, k)
	}
	return r.m[k]
}

// seq2 is like iter.Seq2 but before we've updated to go1.22.
type seq2[T1, T2 any] func(yield func(T1, T2) bool)
