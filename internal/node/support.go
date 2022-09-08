package node

import (
	"github.com/smart-core-os/sc-golang/pkg/router"
)

// Supporter is a type that can have its supported functions changed at runtime.
type Supporter interface {
	Support(functions ...Function)
}

// Function represents something that is supported by a type.
// For example an API might be represented as a Function and added to a server.
type Function interface {
	apply(node *Node)
}

// EmptyFunction does not change the functions a Node supports.
// Can be embedded in custom Function types to allow extending Node support.
type EmptyFunction struct{}

func (e EmptyFunction) apply(_ *Node) {
	// Do nothing
}

type functionFunc func(node *Node)

func (f functionFunc) apply(node *Node) {
	f(node)
}

// Routing adds the given routers to the supported API of the node.
func Routing(r ...router.Router) Function {
	return functionFunc(func(node *Node) {
		node.addRouter(r...)
	})
}

// Clients adds the given clients, which should be proto service clients, to a node.
// Code can access a nodes clients via Client.
// Typically these are associated with the nodes routers via server->client conversion.
func Clients(c ...any) Function {
	return functionFunc(func(node *Node) {
		node.addClient(c...)
	})
}
