package node

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/vanti-dev/sc-bos/internal/router"
	"github.com/vanti-dev/sc-bos/pkg/node/alltraits"
)

// Client implements Clienter backed by the node's router.
//
// Deprecated: Use ClientConn() to acquire a connection and construct clients directly.
func (n *Node) Client(p any) error {
	if !alltraits.NewClient(p, n.ClientConn()) {
		return fmt.Errorf("no client of type %T", p)
	}
	return nil
}

// ClientConn returns a connection to the Node's router.
func (n *Node) ClientConn() grpc.ClientConnInterface {
	return router.NewLoopback(n.router)
}

func (n *Node) ServerHandler() grpc.StreamHandler {
	return router.StreamHandler(n.router)
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
