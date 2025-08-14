package node

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/internal/router"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/devicespb"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
	"github.com/vanti-dev/sc-bos/pkg/node/internal/metadatadevices"
	"github.com/vanti-dev/sc-bos/pkg/node/internal/parentdevices"
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
	router *router.Router

	// mu protects writes to devices to avoid concurrent modification errors.
	// The model that backs this field is consistent when accessed concurrently,
	// but returns an error if its modified concurrently instead of waiting.
	// We want it to wait.
	mu      sync.Mutex
	devices *devicespb.Collection

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
		devices: devicespb.NewCollection(resource.WithIDInterceptor(mapID)),
		Logger:  zap.NewNop(),
	}

	// nodes implement the MetadataApi without using the router,
	// the ParentApi is implemented for this nodes name only.
	traits.RegisterMetadataApiServer(node.router, metadatadevices.NewServer(node.devices))
	node.announceLocked(name,
		HasServer(traits.RegisterParentApiServer, traits.ParentApiServer(parentdevices.NewServer(name, node.devices))),
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

	mds := a.metadata
	if !a.noAutoMetadata && len(ts) > 0 {
		md := &traits.Metadata{}
		for _, t := range ts {
			md.Traits = append(md.Traits, &traits.TraitMetadata{Name: string(t.name)})
		}
		mds = append(mds, md)
	}

	md := n.mergeAllMetadata(name, mds...)
	if md != nil {
		_, err := n.devices.Merge(&gen.Device{
			Name:     name,
			Metadata: md,
		}, resource.WithCreateIfAbsent())
		if err != nil {
			log.Errorf("merge metadata %q: %v", name, err)
		} else {
			undo = append(undo, func() {
				n.mu.Lock()
				defer n.mu.Unlock()
				_, _ = n.devices.Delete(name, resource.WithAllowMissing(true))
			})
		}
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
		traitsString = append(traitsString, t.name.Local())
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
