package statuspb

import (
	"sync"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

// Map tracks the status of multiple named status Model.
// Any names used to update/delete problems will be announced as Status traits using the given announcer.
type Map struct {
	mu        sync.Mutex
	known     map[string]model // keyed by sc name used to announce
	announcer node.Announcer
}

type model struct {
	*Model
	unannounce node.Undo
}

func NewMap(announcer node.Announcer) *Map {
	return &Map{
		known:     make(map[string]model),
		announcer: announcer,
	}
}

func (m *Map) UpdateProblem(name string, problem *gen.StatusLog_Problem) {
	m.getOrCreateModel(name).UpdateProblem(problem)
}

func (m *Map) DeleteProblem(name, problem string) {
	mod, ok := m.getModel(name)
	if !ok {
		return // nothing to do anyway
	}
	mod.DeleteProblem(problem)
}

func (m *Map) Forget(name string) {
	m.mu.Lock()
	mod, ok := m.known[name]
	if !ok {
		m.mu.Unlock()
		return
	}
	delete(m.known, name)
	m.mu.Unlock()
	mod.unannounce()
}

func (m *Map) getOrCreateModel(name string) model {
	m.mu.Lock()
	mod, ok := m.known[name]
	if !ok {
		nm := NewModel()
		client := gen.WrapStatusApi(NewModelServer(nm))
		mod = model{
			Model:      nm,
			unannounce: m.announcer.Announce(name, node.HasTrait(TraitName, node.WithClients(client))),
		}
		m.known[name] = mod
	}
	m.mu.Unlock()
	return mod
}

func (m *Map) getModel(name string) (model, bool) {
	m.mu.Lock()
	mod, ok := m.known[name]
	m.mu.Unlock()
	return mod, ok
}
