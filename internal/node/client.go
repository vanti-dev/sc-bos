package node

import "fmt"

// Client returns a new api client of type C which can be used to interact with the named devices in n.
func Client[C any](n *Node) (C, error) {
	for _, client := range n.clients {
		if t, ok := client.(C); ok {
			return t, nil
		}
	}
	var c C
	return c, fmt.Errorf("no such client type %T", c)
}
