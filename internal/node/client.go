package node

import (
	"errors"
)

// Client returns a new api client of type C which can be used to interact with the named devices in n.
func Client[C any](n *Node) (C, error) {
	for _, client := range n.clients {
		if t, ok := client.(C); ok {
			return t, nil
		}
	}
	var c C
	return c, errors.New("unknown client type")
}

// FindClient places into c a client backed by the named devices in n.
func FindClient[C any](n *Node, c *C) {
	for _, client := range n.clients {
		if t, ok := client.(C); ok {
			*c = t
			return
		}
	}
}
