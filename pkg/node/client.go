package node

import (
	"fmt"

	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
)

// Client implements Clienter backed by the node's router.
func (n *Node) Client(p any) error {
	if !alltraits.NewClient(p, n.router) {
		return fmt.Errorf("no client of type %T", p)
	}
	return nil
}

// Clienter represents a type that can respond with an API client.
type Clienter interface {
	// Client sets into the pointer p a client, if one is available, or returns an error.
	// Argument p should be a pointer to a variable of the required client type.
	//
	// Example
	//
	//	var client traits.OnOffApiClient
	//	err := n.Client(&client)
	Client(p any) error
}

// ClientFunc adapts a func of the correct signature to implement Clienter.
type ClientFunc func(p any) error

func (c ClientFunc) Client(p any) error {
	return c(p)
}
