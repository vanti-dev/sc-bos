package node

import (
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Local represents a local smart core node.
// The node has collection of supported apis, represented by router.Router instances and configured via Support.
// When new names are created they should call Announce and with the features relevant to the name.
//
// All supported APIs in a Local are routed based on a name.
// Call Support to add new features to the Local.
// Calling Support after Register will not have any effect on the served apis.
type Local struct {
	name string

	// children keeps track of all the names that have been announced to this node.
	// Lazy, initialised when addChildTrait via Announce(HasTrait) or Register are called.
	children *parent.Model
	// routers holds all the APIs this node supports.
	// Populated via Support(WithRouter).
	routers []router.Router

	Logger *zap.Logger
}

// New creates a new Local node with the given name.
func New(name string) *Local {
	return &Local{
		name:   name,
		Logger: zap.NewNop(),
	}
}

// Name returns the device name for this node, how this node refers to itself.
func (n *Local) Name() string {
	return n.name
}

// Register implements server.GrpcApi and registers all supported routers with s.
func (n *Local) Register(s *grpc.Server) {
	n.parent() // force the parent api to be initialised
	for _, r := range n.routers {
		if api, ok := r.(server.GrpcApi); ok {
			api.Register(s)
		}
	}
}

// Announce adds a new name with the given features to this node.
// You may call Announce multiple times with the same name to add additional features, for example new traits.
func (n *Local) Announce(name string, features ...Feature) {
	a := &announcement{name: name}
	for _, feature := range features {
		feature.apply(a)
	}
	for _, client := range a.clients {
		n.addRoute(name, client)
	}
	log := n.Logger.Sugar()
	for _, t := range a.traits {
		log.Debugf("%v now implements %v\n", name, t.name)

		if !t.noAddChildTrait && name != n.name {
			n.addChildTrait(a.name, t.name)
		}
		for _, client := range t.clients {
			n.addRoute(a.name, client)
		}
		if !t.noAddMetadata {
			md := t.metadata
			if md == nil {
				md = AutoTraitMetadata
			}
			if err := n.addTraitMetadata(name, t.name, md); err != nil {
				if err != MetadataTraitNotSupported {
					log.Warnf("%v %v: %v", name, t.name, err)
				}
			}
		}
	}
}

// Support adds new supported functions to this node.
func (n *Local) Support(functions ...Function) {
	for _, function := range functions {
		function.apply(n)
	}
}

func (n *Local) addRouter(r ...router.Router) {
	n.routers = append(n.routers, r...)
}

// addRoute adds name->impl as a route to all routers that support the type impl.
func (n *Local) addRoute(name string, impl interface{}) (added bool) {
	for _, r := range n.routers {
		if r.HoldsType(impl) {
			r.Add(name, impl)
			added = true
		}
	}
	return
}

func (n *Local) addChildTrait(name string, traitName ...trait.Name) {
	n.parent().AddChildTrait(name, traitName...)
}

func (n *Local) parent() *parent.Model {
	if n.children == nil {
		// add this model as a device
		n.children = parent.NewModel()
		client := parent.WrapApi(parent.NewModelServer(n.children))
		n.Announce(n.name, HasTrait(trait.Parent, WithClients(client)))
	}
	return n.children
}
