package node

type Scoped struct {
	announcer Announcer
	undo      Undo
}

// Reset clears old scopes and creates a new scope with the announced name.
// The caller should never call Reset concurrently.
func (s *Scoped) Reset(name string) {
	// undo the previous scope if it exists
	if s.undo != nil {
		s.undo()
	}
	// set up the new scope
	var announcer Announcer
	announcer, s.undo = AnnounceScope(s.announcer)
	// use the new scope to announce our name
	announcer.Announce(name)
}

func ExampleAnnounceScope_scoped() {
	scoped := &Scoped{announcer: New("test")}

	// set up the first scope
	scoped.Reset("a")

	// some time later, set up a new scope
	scoped.Reset("b")
}
