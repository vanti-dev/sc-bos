package node

import (
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	// children keeps track of all the names that have been announced to this node.
	// Lazy, initialised when addChildTrait via Announce(HasTrait) or Register are called.
	children *parent.Model
	// routers holds all the APIs this node supports.
	// Populated via Support(Routing).
	routers []router.Router
	// clients holds instances of service clients returned by Client.
	// Typically they are wrappers around each router instance.
	// Populated via Support(Clients).
	clients []any

	Logger *zap.Logger
}

// New creates a new Node node with the given name.
func New(name string) *Node {
	return &Node{
		name:   name,
		Logger: zap.NewNop(),
	}
}

// Name returns the device name for this node, how this node refers to itself.
func (n *Node) Name() string {
	return n.name
}

// Register implements server.GrpcApi and registers all supported routers with s.
func (n *Node) Register(s *grpc.Server) {
	n.parent() // force the parent api to be initialised
	for _, r := range n.routers {
		if api, ok := r.(server.GrpcApi); ok {
			api.Register(s)
		}
	}
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
		log.Debugf("%v now implements %v\n", name, t.name)

		if !t.noAddChildTrait && name != n.name {
			undo = append(undo, n.addChildTrait(a.name, t.name))
		}
		for _, client := range t.clients {
			undo = append(undo, n.addRoute(a.name, client))
		}
		if !t.noAddMetadata {
			md := t.metadata
			if md == nil {
				md = AutoTraitMetadata
			}
			undoMd, err := n.addTraitMetadata(name, t.name, md)
			if err != nil {
				if err != MetadataTraitNotSupported {
					log.Warnf("%v %v: %v", name, t.name, err)
				}
			}
			undo = append(undo, undoMd)
		}
	}
	return UndoAll(undo...)
}

// Support adds new supported functions to this node.
func (n *Node) Support(functions ...Function) {
	for _, function := range functions {
		function.apply(n)
	}
}

func (n *Node) addRouter(r ...router.Router) {
	n.routers = append(n.routers, r...)
}

// addRoute adds name->impl as a route to all routers that support the type impl.
func (n *Node) addRoute(name string, impl interface{}) Undo {
	var undo []Undo
	for _, r := range n.routers {
		if r.HoldsType(impl) {
			r.Add(name, impl)
			undo = append(undo, func() {
				r.Remove(name)
			})
		}
	}
	return UndoAll(undo...)
}

func (n *Node) addChildTrait(name string, traitName ...trait.Name) Undo {
	n.parent().AddChildTrait(name, traitName...)
	return func() {
		// todo: remove child traits from n.parent()
	}
}

func (n *Node) addClient(c ...any) {
	n.clients = append(n.clients, c...)
}

func (n *Node) parent() *parent.Model {
	if n.children == nil {
		// add this model as a device
		n.children = parent.NewModel()
		client := parent.WrapApi(parent.NewModelServer(n.children))
		n.Announce(n.name, HasTrait(trait.Parent, WithClients(client)))
	}
	return n.children
}
