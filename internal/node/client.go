package node

import (
	"errors"
	"fmt"
	"reflect"
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

// Client implements Clienter backed by clients configured using Support(Clients).
func (n *Node) Client(p any) error {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("%T is not a pointer", p)
	}
	elem := v.Elem()
	et := elem.Type()
	if !elem.CanSet() {
		return fmt.Errorf("%T can not be set", p)
	}

	for _, client := range n.clients {
		if reflect.TypeOf(client).AssignableTo(et) {
			elem.Set(reflect.ValueOf(client))
			return nil
		}
	}

	return fmt.Errorf("no client of type %v", elem)
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
