package node

import (
	"sync"
)

// Undo allows a statement to be undone.
type Undo func()

// NilUndo does nothing.
func NilUndo() {}

// UndoAll creates an Undo that undoes all the given Undo in order.
func UndoAll(undo ...Undo) Undo {
	return func() {
		for _, u := range undo {
			u()
		}
	}
}

// UndoOnce returns an Undo that only calls undo once.
func UndoOnce(undo Undo) Undo {
	var once sync.Once
	return func() {
		once.Do(undo)
	}
}
