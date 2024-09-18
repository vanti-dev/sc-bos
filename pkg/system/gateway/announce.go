package gateway

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/internal/router"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/servicepb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
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
		conns:       newSyncMap[*grpc.ClientConn](),
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
	var serviceChanges <-chan rx.Change[remoteService]
	var stopServiceSub context.CancelFunc
	var undoServices tasks
	setupServiceSub := func() {
		var serviceCtx context.Context
		serviceCtx, stopServiceSub = context.WithCancel(ctx)
		var services *slices.Sorted[remoteService]
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
				undoServices.remove(c.Old.name)
			}
			if c.Type != rx.Remove {
				undoServices[c.New.name] = a.announceRemoteService(c.New)
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
	servicesClient := gen.WrapServicesApi(servicesApi)
	var undos []node.Undo
	for _, bucket := range []string{"automations", "drivers", "systems", "zones"} {
		undos = append(undos, a.Announce(name+"/"+bucket, node.HasClient(servicesClient)))
	}
	return node.UndoAll(undos...)
}

// announceName allows this node to answer requests aimed at the given remoteDesc.
// This includes named RPCs (like trait requests) and DeviceApi / ParentApi requests that would include this name in their responses.
func (a *announcer) announceName(d remoteDesc) node.Undo {
	if d.name == "" {
		return node.NilUndo
	}

	var services []grpc.ServiceDesc
	for _, t := range d.md.GetTraits() {
		traitName := trait.Name(t.Name)

		if traitName == trait.Metadata {
			// skip, the node will handle the implementation of this for us
			continue
		}

		services = append(services, alltraits.ServiceDesc(traitName)...)
	}

	if old, replaced := a.table.conns.Set(d.name, a.node.conn); replaced {
		if old != a.node.conn {
			a.logger.Warn("name already registered by another party, it has been replaced",
				zap.String("name", d.name), zap.String("old", old.Target()), zap.String("new", a.node.addr))
		}
	}

	return node.UndoAll(
		a.Announce(d.name,
			node.HasMetadata(d.md),
			node.HasClientConn(a.node.conn, services...),
		),
		func() { a.table.conns.Del(d.name) },
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
func (a *announcer) announceRemoteService(rs remoteService) node.Undo {
	if a.table.services.Inc(rs.name) == 1 {
		undo := a.announceRemoteServiceApis(rs)
		if undo != nil { // can be nil if there's nothing to undo
			a.table.serviceUndo.Set(rs.name, undo)
		}
	}

	return func() {
		if a.table.services.Dec(rs.name) == 0 {
			// we were the last to remove the service, clean everything up
			undo, ok := a.table.serviceUndo.Del(rs.name)
			if ok {
				undo()
			}
		}
	}
}

// announceRemoteServices updates this node to respond to requests for each remoteService in seq.
func (a *announcer) announceRemoteServices(seq seq2[int, remoteService]) tasks {
	dst := tasks{}
	seq(func(_ int, rs remoteService) bool {
		dst[rs.name] = a.announceRemoteService(rs)
		return true
	})
	return dst
}

func (a *announcer) announceRemoteServiceApis(rs remoteService) node.Undo {
	var keyFuncs []router.KeyFunc
	for _, method := range rs.methods {
		keyFunc, err := router.NameKey(method.Input())
		if err != nil {
			continue // we don't care about err, just that this method doesn't have a name key
		}
		keyFuncs = append(keyFuncs, keyFunc)
	}

	// which type of proxying should each method use?
	var newMethod func(int, protoreflect.MethodDescriptor) router.Method
	switch {
	case len(keyFuncs) == len(rs.methods):
		// all methods have a key func, we can enable routing for this service as a whole
		newMethod = func(i int, method protoreflect.MethodDescriptor) router.Method {
			return routedMethod(method, keyFuncs[i], a.table.conns)
		}
	case a.node.isHub:
		// route everything else to the hub
		newMethod = func(_ int, method protoreflect.MethodDescriptor) router.Method {
			return fixedMethod(method, a.node.conn)
		}
	default:
		// Found a non-routable service on a non-hub node.
		// We didn't think we had any of these so log it to remind us.
		a.logger.Warn("non-routable service on non-hub node", zap.String("service", rs.name))
		return nil
	}

	var undos []node.Undo
	for i, method := range rs.methods {
		name := fullNameToRpcPath(method.FullName())
		if a.System.methods.Add(name, newMethod(i, method)) {
			undos = append(undos, func() {
				a.System.methods.Delete(name)
			})
		}
	}

	// check if someone else has registered the service methods before us, and let them
	if len(undos) < len(rs.methods) {
		for _, undo := range undos {
			undo()
		}
		undos = nil
		a.logger.Warn("service already registered by another party", zap.String("service", rs.name))
		// note: we don't register with a.table.serviceUndo as we've already undone everything
		return nil
	}

	return node.UndoAll(undos...)
}

// table tracks state across all remote nodes.
type table struct {
	conns       *syncMap[*grpc.ClientConn] // keyed by rpc method path
	services    *counts                    // keyed by service full name, counts how many times we've seen services across all remote nodes
	serviceUndo *syncMap[node.Undo]        // keyed by service full name
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

// fixedMethod returns an unknown.Method that resolves to the given conn.
func fixedMethod(method protoreflect.MethodDescriptor, conn *grpc.ClientConn) router.Method {
	return router.Method{
		StreamDesc: grpc.StreamDesc{
			ServerStreams: method.IsStreamingServer(),
			ClientStreams: method.IsStreamingClient(),
		},
		Resolver: router.NewFixedResolver(conn),
	}
}

// routedMethod returns an unknown.Method that resolves to a conn based on the keyFunc and given table.
func routedMethod(method protoreflect.MethodDescriptor, keyFunc router.KeyFunc, table *syncMap[*grpc.ClientConn]) router.Method {
	return router.Method{
		StreamDesc: grpc.StreamDesc{
			ServerStreams: method.IsStreamingServer(),
			ClientStreams: method.IsStreamingClient(),
		},
		Resolver: router.ResolverFunc(func(mr router.MsgRecver) (grpc.ClientConnInterface, error) {
			key, err := keyFunc(mr)
			if err != nil {
				return nil, err
			}
			conn, ok := table.Get(key)
			if !ok {
				return nil, status.Errorf(codes.NotFound, "name not known: %q", key)
			}
			return conn, nil
		}),
	}
}

// fullNameToRpcPath maps proto full names (package.Service.Method) to rpc paths (/package.Service/Method).
func fullNameToRpcPath(fullName protoreflect.FullName) string {
	return fmt.Sprintf("/%s/%s", fullName.Parent(), fullName.Name())
}

// seq2 is like iter.Seq2 but before we've updated to go1.22.
type seq2[T1, T2 any] func(yield func(T1, T2) bool)
