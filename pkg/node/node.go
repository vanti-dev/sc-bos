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

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/vanti-dev/sc-bos/internal/router"
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
	mu     sync.Mutex // protects all fields below, typically Announce, Support, and methods that rely on that data

	// children keeps track of all the names that have been announced to this node.
	// Lazy, initialised when addChildTrait via Announce(HasTrait) or Register are called.
	children *parent.Model

	// allMetadata allows users of the node to be notified of any metadata changes via Announce or when
	// that announcement is undone.
	allMetadata *metadata.Collection

	Logger *zap.Logger
}

// New creates a new Node with the given name.
func New(name string) *Node {
	node := &Node{
		name:        name,
		router:      router.New(),
		Logger:      zap.NewNop(),
		allMetadata: metadata.NewCollection(),
	}

	// metadata should be supported by default
	metadataService, err := LocalService(traits.MetadataApi_ServiceDesc, metadata.NewCollectionServer(node.allMetadata))
	if err != nil {
		// should be impossible because the metadata ServiceDesc is valid
		panic("metadata service not valid")
	}
	// can't error because there's no services/routes to conflict with
	_, _ = node.AnnounceService(metadataService)

	node.parentLocked()
	return node
}

// Name returns the device name for this node, how this node refers to itself.
func (n *Node) Name() string {
	return n.name
}

// Announce adds a new name with the given features to this node.
// You may call Announce multiple times with the same name to add additional features, for example new traits.
// Executing the returned Undo will undo any direct changes made.
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
	var services []service
	services = append(services, a.services...)
	for _, t := range a.traits {
		services = append(services, t.services...)
	}
	for _, s := range services {
		undoRoute, err := registerDeviceRoute(n.router, name, s)
		if err != nil {
			log.Errorf("cannot register %s route for %q: %v", s.desc.FullName(), name, err)
		} else {
			undo = append(undo, undoRoute)
		}
	}
	if a.proxyTo != nil {
		undoProxy, err := registerProxyRoute(n.router, name, a.proxyTo)
		if err != nil {
			log.Errorf("cannot register proxy route for %q: %v", name, err)
		} else {
			undo = append(undo, undoProxy)
		}
	}

	// unless specifically disabled, all devices support the Metadata trait
	if !a.noAutoMetadata {
		a.traits = append(a.traits, traitFeature{name: trait.Metadata})
	}
	for _, t := range a.traits {
		log.Debugf("%v now implements %v", name, t.name)
		undo = append(undo, func() {
			log.Debugf("%v no longer implements %v", name, t.name)
		})

		if !t.noAddChildTrait && name != n.name {
			undo = append(undo, n.addChildTrait(a.name, t.name))
		}
	}

	mds := a.metadata
	if !a.noAutoMetadata && len(a.traits) > 0 {
		md := &traits.Metadata{}
		for _, t := range a.traits {
			md.Traits = append(md.Traits, &traits.TraitMetadata{Name: string(t.name)})
		}
		mds = append(mds, md)
	}
	// always need to set the name of the device in its metadata
	mds = append(mds, &traits.Metadata{Name: name})

	for _, md := range mds {
		undoMd, err := n.mergeMetadata(name, md)
		if err != nil {
			if errors.Is(err, MetadataTraitNotSupported) {
				log.Warnf("%v metadata: %v", name, err)
			}
			continue
		}
		undo = append(undo, undoMd)
	}

	return UndoAll(undo...)
}

// Support adds new supported functions to this node.
func (n *Node) Support(functions ...Function) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, function := range functions {
		function.apply(n)
	}
}

func (n *Node) addChildTrait(name string, traitName ...trait.Name) Undo {
	retryConcurrentOp(func() {
		n.parentLocked().AddChildTrait(name, traitName...)
	})
	return func() {
		var child *traits.Child
		parentModel := n.parent()
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

func (n *Node) parent() *parent.Model {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.parentLocked()
}

func (n *Node) parentLocked() *parent.Model {
	if n.children == nil {
		// add this model as a device
		n.children = parent.NewModel()
		n.announceLocked(n.name,
			HasServer(traits.RegisterParentApiServer, traits.ParentApiServer(parent.NewModelServer(n.children))),
			HasTrait(trait.Parent),
		)
	}
	return n.children
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

func registerDeviceRoute(r *router.Router, name string, s service) (Undo, error) {
	err := ensureServiceSupported(r, s)
	if err != nil {
		return NilUndo, err
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
		if s.nameRouting && !existing.KeyRoutable() {
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

	// AddService might return an error if another goroutine added support after the GetService check above
	// this is a bit of wasted work but is safe because the service added will be the same
	_ = r.AddService(routerService)
	return nil
}
