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
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
)

// Node represents a smart core node.
// The node has collection of supported apis, represented by router.Router instances and configured via Support.
// When new names are created they should call Announce and with the features relevant to the name.
//
// All supported APIs in a Node are routed based on a name.
// Call Support to add new features to the Node.
// Calling Support after Register will not have any effect on the served apis.
type Node struct {
	name string
	mu   sync.Mutex // protects all fields below, typically Announce, Support, and methods that rely on that data

	// children keeps track of all the names that have been announced to this node.
	// Lazy, initialised when addChildTrait via Announce(HasTrait) or Register are called.
	children *parent.Model
	// routers holds all the APIs this node supports.
	// Populated via Support(Routing).
	routers []router.Router
	// clients holds instances of service clients returned by Client.
	// Typically, they are wrappers around each router instance.
	// Populated via Support(Clients).
	clients []any
	// apis holds each of the APIs that this node registers with a grpc.Server.
	// Populated via Support(Api) explicitly, or Support(Routing) if the router implements server.GrpcApi.
	apis []server.GrpcApi

	// allMetadata allows users of the node to be notified of any metadata changes via Announce or when
	// that announcement is undone.
	allMetadata *resource.Collection // of *traits.Metadata

	Logger *zap.Logger
}

// New creates a new Node with the given name.
func New(name string) *Node {
	return &Node{
		name:        name,
		Logger:      zap.NewNop(),
		allMetadata: resource.NewCollection(),
	}
}

// Name returns the device name for this node, how this node refers to itself.
func (n *Node) Name() string {
	return n.name
}

// Register implements server.GrpcApi and registers all supported routers with s.
func (n *Node) Register(s *grpc.Server) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, api := range n.apis {
		api.Register(s)
	}
	n.registerRequiredLocked(s)
	n.parentLocked() // force the parent api to be initialised
}

// registerRequiredLocked registers the required APIs with s if they haven't already been registered.
func (n *Node) registerRequiredLocked(s *grpc.Server) {
	// Note: we don't use alltraits here to avoid importing and initialising all traits when we only need a few.
	registeredServices := s.GetServiceInfo()
	if _, ok := registeredServices["smartcore.traits.MetadataApi"]; !ok {
		r := metadata.NewApiRouter()
		registerRouter(n, s, r, metadata.WrapApi(r))
	}
	if _, ok := registeredServices["smartcore.traits.MetadataInfo"]; !ok {
		r := metadata.NewInfoRouter()
		registerRouter(n, s, r, metadata.WrapInfo(r))
	}
	if _, ok := registeredServices["smartcore.traits.ParentApi"]; !ok {
		r := parent.NewApiRouter()
		registerRouter(n, s, r, parent.WrapApi(r))
	}
	if _, ok := registeredServices["smartcore.traits.ParentInfo"]; !ok {
		r := parent.NewInfoRouter()
		registerRouter(n, s, r, parent.WrapInfo(r))
	}
}

func registerRouter[R interface {
	router.Router
	server.GrpcApi
}](n *Node, s *grpc.Server, r R, c any) {
	n.addRouter(r)
	n.addApi(r)
	n.addClient(c)
	r.Register(s)
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
	a := &announcement{name: name}
	for _, feature := range features {
		feature.apply(a)
	}

	var undo []Undo
	undo = append(undo, a.undo...)

	for _, client := range a.clients {
		undo = append(undo, n.addRoute(name, client))
	}
	log := n.Logger.Sugar()
	for _, t := range a.traits {
		log.Debugf("%v now implements %v", name, t.name)
		undo = append(undo, func() {
			log.Debugf("%v no longer implements %v", name, t.name)
		})

		// adding clients must happen before addChildTrait because addChildTrait can call announce and might look
		// for the metadata client in order to add the Parent trait to it.
		for _, client := range t.clients {
			undo = append(undo, n.addRoute(a.name, client))
		}
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

func (n *Node) addRouter(rs ...router.Router) {
	n.routers = append(n.routers, rs...)
}

func (n *Node) addApi(apis ...server.GrpcApi) {
	n.apis = append(n.apis, apis...)
}

// addRoute adds name->impl as a route to all routers that support the type impl.
func (n *Node) addRoute(name string, c client) Undo {
	var undo []Undo
	var addCount int
	for _, r := range n.routers {
		if r.HoldsType(c.impl) {
			addCount++
			r.Add(name, c.impl)
			r := r
			undo = append(undo, func() {
				r.Remove(name)
			})
		}
	}
	if addCount == 0 && !c.allowUnsupported {
		n.Logger.Warn(fmt.Sprintf("no router for %s typed %T", name, c.impl))
	}
	return UndoAll(undo...)
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

func (n *Node) addClient(c ...any) {
	n.clients = append(n.clients, c...)
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
		client := parent.WrapApi(parent.NewModelServer(n.children))
		n.announceLocked(n.name, HasTrait(trait.Parent, WithClients(client)))
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
