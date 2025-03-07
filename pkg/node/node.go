package node

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadatapb"
	"github.com/smart-core-os/sc-golang/pkg/trait/parentpb"
	"github.com/vanti-dev/sc-bos/internal/router"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
)

// Node represents a smart core node.
// The node has collection of supported apis, represented by router.Router instances and configured via Support.
// When new names are created they should call Announce and with the features relevant to the name.
//
// All supported APIs in a Node are routed based on a name.
// Call Support to add new features to the Node.
// Calling Support after Register will not have any effect on the served apis.
type Node struct {
	name   string
	mu     sync.Mutex // protects all fields below, typically Announce, Support, and methods that rely on that data
	router *router.Router

	// children keeps track of all the names that have been announced to this node.
	// Lazy, initialised when addChildTrait via Announce(HasTrait) or Register are called.
	children *parentpb.Model

	// allMetadata allows users of the node to be notified of any metadata changes via Announce or when
	// that announcement is undone.
	allMetadata *metadatapb.Collection

	Logger *zap.Logger
}

// New creates a new Node with the given name.
func New(name string) *Node {
	mapID := func(requestName string) string {
		if requestName == "" {
			return name
		} else {
			return requestName
		}
	}

	node := &Node{
		name: name,
		router: router.New(router.WithKeyInterceptor(func(key string) (mappedKey string, err error) {
			return mapID(key), nil
		})),
		children:    parentpb.NewModel(),
		Logger:      zap.NewNop(),
		allMetadata: metadatapb.NewCollection(resource.WithIDInterceptor(mapID)),
	}

	// metadata should be supported by default
	traits.RegisterMetadataApiServer(node.router, metadatapb.NewCollectionServer(node.allMetadata))
	_ = node.Announce(name, HasTrait(trait.Metadata))
	node.announceLocked(name,
		HasServer(traits.RegisterParentApiServer, traits.ParentApiServer(parentpb.NewModelServer(node.children))),
		HasTrait(trait.Parent),
	)
	return node
}

// Name returns the device name for this node, how this node refers to itself.
func (n *Node) Name() string {
	return n.name
}

// Announce adds a new name with the given features to this node.
// You may call Announce multiple times with the same name to add additional features, for example new traits.
// You must not Announce the same features on the same name multiple times, until the original announcement of
// those features has been undone.
// Executing the returned Undo will undo any direct changes made, but will not remove support for any services
// from the router.
//
// # A note on undoing
//
// The undo process is not perfect but best effort.
// Hooks and callbacks may have been executed that have side effects that are not undone.
func (n *Node) Announce(name string, features ...Feature) Undo {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.announceLocked(name, features...)
}

func (n *Node) announceLocked(name string, features ...Feature) Undo {
	log := n.Logger.Sugar()
	a := &announcement{name: name}
	for _, feature := range features {
		feature.apply(a)
	}

	var undo []Undo
	undo = append(undo, a.undo...)

	// register all relevant routes with the router
	services := allServices(a, n.Logger)
	for _, s := range services {
		serviceName := s.desc.FullName()
		undoRoute, err := registerDeviceRoute(n.router, name, s)
		if err != nil {
			log.Errorf("cannot register service %s for %q: %v", serviceName, name, err)
		} else {
			undo = append(undo, undoRoute)
		}
	}
	if a.proxyTo != nil {
		undoProxy, err := registerProxyRoute(n.router, name, a.proxyTo)
		if err != nil {
			log.Errorf("cannot register proxy for %q: %v", name, err)
		} else {
			undo = append(undo, undoProxy)
		}
	}

	// unless specifically disabled, all devices support the Metadata trait
	ts := a.traits
	if !a.noAutoMetadata {
		ts = append(ts, traitFeature{name: trait.Metadata})
	}
	for _, t := range ts {
		if !t.noAddChildTrait && name != n.name {
			undo = append(undo, n.addChildTrait(a.name, t.name))
		}
	}

	mds := a.metadata
	if !a.noAutoMetadata && len(ts) > 0 {
		md := &traits.Metadata{}
		for _, t := range ts {
			md.Traits = append(md.Traits, &traits.TraitMetadata{Name: string(t.name)})
		}
		mds = append(mds, md)
	}

	for _, md := range mds {
		// always need to set the name of the device in its metadata
		md.Name = name
		undoMd, err := n.mergeMetadata(name, md)
		if err != nil {
			if errors.Is(err, MetadataTraitNotSupported) {
				log.Warnf("%v metadata: %v", name, err)
			}
			continue
		}
		undo = append(undo, undoMd)
	}

	undo = append(undo, n.logAnnouncement(a, services))

	return UndoAll(undo...)
}

func (n *Node) logAnnouncement(a *announcement, services []service) Undo {
	serviceString := make([]string, 0, len(services))
	for _, s := range services {
		serviceString = append(serviceString, string(s.desc.Name()))
	}
	traitsString := make([]string, 0, len(a.traits))
	for _, t := range a.traits {
		n := t.name.Local()
		if t.noAddChildTrait {
			n = "-" + n
		}
		traitsString = append(traitsString, n)
	}
	var flags []string
	if len(a.metadata) > 0 {
		flags = append(flags, "md")
	}
	if a.noAutoMetadata {
		flags = append(flags, "noAutoMetadata")
	}
	if a.proxyTo != nil {
		flags = append(flags, "proxy")
	}

	log := func(msg string) {
		n.Logger.Debug(msg,
			zap.String("name", a.name),
			zap.Strings("services", serviceString),
			zap.Strings("traits", traitsString),
			zap.Strings("flags", flags),
		)
	}

	log("name announced")
	return func() {
		log("name unannounced")
	}
}

func (n *Node) addChildTrait(name string, traitName ...trait.Name) Undo {
	retryConcurrentOp(func() {
		n.children.AddChildTrait(name, traitName...)
	})
	return func() {
		var child *traits.Child
		parentModel := n.children
		retryConcurrentOp(func() {
			child = parentModel.RemoveChildTrait(name, traitName...)
		})
		// There's a huge assumption here that child was added via AddChildTrait,
		// this should be true but isn't guaranteed
		if child != nil && len(child.Traits) == 0 {
			retryConcurrentOp(func() {
				_, _ = parentModel.RemoveChildByName(child.Name)
			})
		}
	}
}

// retryConcurrentOp runs fn retrying up to 5 times when any panics that isConcurrentUpdateDetectedPanic returns true for.
func retryConcurrentOp(fn func()) (retried bool) {
	var err any
	for i := 0; i < 5; i++ {
		err = catchPanic(fn)
		if isConcurrentUpdateDetectedPanic(err) {
			retried = true
			continue
		}
		if err != nil {
			panic(err) // report other errors
		}
		break // no err
	}
	if err != nil {
		panic(err) // we tried
	}
	return
}

func catchPanic(f func()) (res any) {
	defer func() {
		res = recover()
	}()
	f()
	return
}

func isConcurrentUpdateDetectedPanic(err any) bool {
	e, ok := err.(error)
	return ok && isConcurrentUpdateDetectedError(e)
}

func isConcurrentUpdateDetectedError(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return s.Code() == codes.Aborted && strings.Contains(s.Message(), "concurrent update detected")
}

// Supports s on the router.
// If s has a conn, adds a route for it.
func registerDeviceRoute(r *router.Router, name string, s service) (Undo, error) {
	err := ensureServiceSupported(r, s)
	if err != nil {
		return NilUndo, err
	}
	if s.conn == nil {
		// service just needs to be supported by the router, but don't need to add a route
		return NilUndo, nil
	}

	serviceName := string(s.desc.FullName())
	err = r.AddRoute(serviceName, name, s.conn)
	if err != nil {
		return NilUndo, err
	}

	return func() {
		_ = r.DeleteRoute(serviceName, name)
	}, nil
}

func registerProxyRoute(r *router.Router, name string, conn grpc.ClientConnInterface) (Undo, error) {
	err := r.AddRoute("", name, conn)
	if err != nil {
		return NilUndo, err
	}

	return func() {
		_ = r.DeleteRoute("", name)
	}, nil
}

func ensureServiceSupported(r *router.Router, s service) error {
	serviceName := string(s.desc.FullName())
	if existing := r.GetService(serviceName); existing != nil {
		switch {
		case serviceName == traits.MetadataApi_ServiceDesc.ServiceName:
		case serviceName == traits.MetadataInfo_ServiceDesc.ServiceName:
			// skip, we support metadata specially
		case s.nameRouting && !existing.KeyRoutable():
			// existing service does not support name routing!
			return fmt.Errorf("service %q already exists but does not support name routing", serviceName)
		}
		// already supported, nothing to do
		return nil
	}

	var routerService *router.Service
	if s.nameRouting {
		// smart core traits use the name field to route requests to the right device
		var err error
		routerService, err = router.NewRoutedService(s.desc, "name")
		if err != nil {
			return fmt.Errorf("service %q is not routable by name: %w", serviceName, err)
		}
	} else {
		routerService = router.NewUnroutedService(s.desc)
	}

	// AddService might return ErrServiceExists if another goroutine added support after the GetService check above
	// this is a bit of wasted work but is safe because the service added will be the same
	err := r.AddService(routerService)
	if err != nil && !errors.Is(err, router.ErrServiceExists) {
		return err
	}
	return nil
}

// allServices returns all unique services from a and the traits registered with a.
func allServices(a *announcement, logger *zap.Logger) []service {
	seen := make(map[protoreflect.FullName]struct{})
	var services []service
	for _, s := range a.services {
		if _, ok := seen[s.desc.FullName()]; ok {
			continue
		}
		seen[s.desc.FullName()] = struct{}{}
		services = append(services, s)
	}
	for _, t := range a.traits {
		for _, s := range t.services {
			if _, ok := seen[s.desc.FullName()]; ok {
				continue
			}
			seen[s.desc.FullName()] = struct{}{}
			services = append(services, s)
		}
		tss, err := traitServices(t.name)
		if err != nil {
			logger.Warn("cannot determine services to support for trait", zap.String("trait", string(t.name)), zap.Error(err))
			continue
		}
		for _, s := range tss {
			if _, ok := seen[s.desc.FullName()]; ok {
				continue
			}
			seen[s.desc.FullName()] = struct{}{}
			services = append(services, s)
		}
	}
	return services
}

// returns services that should be supported by the node for the given trait
// (returned services do not contain connections, they are just descriptors)
func traitServices(name trait.Name) ([]service, error) {
	serviceDescs := alltraits.ServiceDesc(name)
	if len(serviceDescs) == 0 {
		return nil, fmt.Errorf("trait %s not recognised", name)
	}

	var services []service
	for _, serviceDesc := range serviceDescs {
		if len(serviceDesc.Methods) == 0 {
			continue // avoid ERROR logs for services without methods, which would act as non-routable
		}
		desc, err := registryDescriptor(serviceDesc.ServiceName)
		if err != nil {
			return nil, err
		}

		services = append(services, service{desc: desc, nameRouting: true})
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("trait %s apis have no rpc methods", name)
	}

	return services, nil
}
