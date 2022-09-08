package node

import "github.com/smart-core-os/sc-golang/pkg/router"

// Supporter is a type that can have its supported functions changed at runtime.
type Supporter interface {
	Support(functions ...Function)
}

// Function represents something that is supported by a type.
// For example an API might be represented as a Function and added to a server.
type Function interface {
	apply(node *Local)
}

// EmptyFunction does not change the functions a Local supports.
// Can be embedded in custom Function types to allow extending Local support.
type EmptyFunction struct{}

func (e EmptyFunction) apply(_ *Local) {
	// Do nothing
}

type functionFunc func(node *Local)

func (f functionFunc) apply(node *Local) {
	f(node)
}

// WithRouter adds the given routers to the supported API of the node.
func WithRouter(r ...router.Router) Function {
	return functionFunc(func(node *Local) {
		node.addRouter(r...)
	})
}
