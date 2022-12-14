package node

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
