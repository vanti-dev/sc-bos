package node

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

func UnroutedService(desc grpc.ServiceDesc) Function {
	return functionFunc(func(node *Node) {
		err := ensureServiceSupported(node.router, service{
			desc:        desc,
			nameRouting: false,
		})
		if err != nil {
			node.Logger.Error("cannot support unrouted service", zap.Error(err), zap.String("service", desc.ServiceName))
		}
	})
}

func RoutedService(desc grpc.ServiceDesc) Function {
	return functionFunc(func(node *Node) {
		err := ensureServiceSupported(node.router, service{
			desc:        desc,
			nameRouting: true,
		})
		if err != nil {
			node.Logger.Error("cannot support routed service", zap.Error(err), zap.String("service", desc.ServiceName))
		}
	})
}
