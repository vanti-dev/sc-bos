package node

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
