package node

import (
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
)

// Supporter is a type that can have its supported functions changed at runtime.
type Supporter interface {
	Support(functions ...Function)
}

// SelfSupporter inverts the receiver for Supporter.
type SelfSupporter interface {
	AddSupport(supporter Supporter)
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

// Api instructs the node to register the given apis with the grpc.Server when Node.Register is called.
func Api(apis ...server.GrpcApi) Function {
	return functionFunc(func(node *Node) {
		node.addApi(apis...)
	})
}

// Routing adds the given routers to the supported API of the node.
func Routing(rs ...router.Router) Function {
	return functionFunc(func(node *Node) {
		node.addRouter(rs...)
		// Special case, if the router implements GrpcApi, act as though they called Support(Router(r), Api(r)).
		// This is mostly for backwards compatibility reasons, we didn't use to have addApi.
		for _, r := range rs {
			if api, ok := r.(server.GrpcApi); ok {
				node.addApi(api)
			}
		}
	})
}

// Clients adds the given clients, which should be proto service clients, to a node.
// Code can access a nodes clients via Client.
// Typically, these are associated with the nodes routers via server->client conversion.
func Clients(c ...any) Function {
	return functionFunc(func(node *Node) {
		node.addClient(c...)
	})
}
