package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/pborman/uuid"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
	"github.com/vanti-dev/sc-bos/pkg/util/maps"
)

var (
	ErrNotFound = errors.New("not found")
)

// Map tracks multiple Record with the ability to create, delete, and listen for changes to the tracked records.
type Map struct {
	mu    sync.Mutex
	known map[string]*Record
	bus   *minibus.Bus[*Change]

	idFunc IdFunc
	create CreateFunc

	now func() time.Time
}

// CreateFunc returns a new Lifecycle for the given kind.
// CreateFunc is called during Map.Create.
type CreateFunc func(kind string) (Lifecycle, error)

// IdFunc generates a new id given the passed parameters.
// The function should return an error if no id can be found for which exists(id) returns false.
// IdFunc is called during Map.Create to mint new IDs where they are not provided as part of the call.
type IdFunc func(kind string, exists func(id string) bool) (string, error)

// NewMap creates and returns a new empty Map using the given create funcs.
func NewMap(createFunc CreateFunc, idFunc IdFunc) *Map {
	return &Map{
		known:  make(map[string]*Record),
		bus:    &minibus.Bus[*Change]{},
		idFunc: idFunc,
		create: createFunc,
		now:    time.Now,
	}
}

var ErrImmutable = errors.New("immutable")

// NewMapOf creates an immutable map containing only the given known services.
// Services will have an ID assigned based on the index in known.
// The returned map will return an error for Create and Delete.
func NewMapOf(known []Lifecycle) *Map {
	knownMap := make(map[string]*Record, len(known))
	for i, lifecycle := range known {
		id := strconv.FormatInt(int64(i), 10)
		kind := kindFromType(lifecycle)
		knownMap[id] = &Record{Id: id, Kind: kind, Service: lifecycle}
	}
	return &Map{
		known: knownMap,
		bus:   &minibus.Bus[*Change]{},
		now:   time.Now,
	}
}

func kindFromType(t any) string {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	return rt.String()
}

// MapSetNow sets the now func on m returning a func that undoes the set.
// Use for testing.
func MapSetNow(m *Map, now func() time.Time) func() {
	old := m.now
	m.now = now
	return func() {
		m.now = old
	}
}

var (
	// IdIsKind is an IdFunc that attempts to use the kind as the id.
	// This implies that only one service of a given type can exist in the map.
	IdIsKind = IdFunc(func(kind string, exists func(id string) bool) (string, error) {
		if exists(kind) {
			return "", errors.New("id exists")
		}
		return kind, nil
	})
	IdIsUUID = IdFunc(func(kind string, exists func(id string) bool) (string, error) {
		max := 10
		for i := 0; i < max; i++ {
			id := uuid.New()
			if exists(id) {
				continue
			}
			return id, nil
		}
		return "", errors.New("exhausted attempts finding unique uuid")
	})
	IdIsRequired = IdFunc(func(_ string, _ func(id string) bool) (string, error) {
		return "", errors.New("id is required")
	})
)

type Record struct {
	Id, Kind string
	Service  Lifecycle
}

// Create creates and adds a new record to m returning the new ID and the records service State.
// The kind argument is required, but id is optional. If absent the IdFunc will be used to mint a new ID.
// Only State.Active and State.Config are optionally used in the passed state to either Start or Configure the created
// Lifecycle. If either are present and the corresponding Lifecycle call returns an error, then creating the new record
// will be aborted and that error will be returned.
func (m *Map) Create(id, kind string, state State) (string, State, error) {
	r, err := m.createRecord(id, kind)
	if err != nil {
		return "", State{}, err
	}

	outState := r.Service.State()
	defer func() {
		if err != nil {
			m.mu.Lock()
			defer m.mu.Unlock()

			// cleanup the known map
			got, ok := m.known[r.Id]
			if ok && got == r {
				delete(m.known, r.Id)
			}

			// cleanup the service
			if outState.Active {
				_, _ = r.Service.Stop()
			}
		}
	}()

	if state.Active {
		outState, err = r.Service.Start()
	}
	if len(state.Config) > 0 {
		outState, err = r.Service.Configure(state.Config)
	}

	if err != nil {
		return "", State{}, err
	}

	change := &Change{
		ChangeTime: m.now(),
		ChangeType: types.ChangeType_ADD,
		NewValue:   r,
	}
	go m.bus.Send(context.Background(), change)
	return r.Id, outState, nil
}

// Delete stops and removes the record with the given ID from m.
// If id is not found, returns ErrNotFound.
// Delete does not error if the record is already stopped.
func (m *Map) Delete(id string) (State, error) {
	r, err := m.deleteRecord(id)
	if err != nil {
		return State{}, err
	}

	// stop before sending the event
	state, err := r.Service.Stop()
	if errors.Is(err, ErrAlreadyStopped) {
		err = nil
	}

	go m.bus.Send(context.Background(), &Change{
		ChangeTime: m.now(),
		ChangeType: types.ChangeType_REMOVE,
		OldValue:   r,
	})
	return state, err
}

// Get returns the record associated with the given id, or nil if not found.
func (m *Map) Get(id string) *Record {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.known[id]
}

// Values returns all known values in an indeterminate order.
func (m *Map) Values() []*Record {
	m.mu.Lock()
	defer m.mu.Unlock()
	return maps.Values(m.known)
}

// Listen emits changes to m on the returns chan until ctx is done.
func (m *Map) Listen(ctx context.Context) <-chan *Change {
	return m.bus.Listen(ctx)
}

func (m *Map) createRecord(id, kind string) (*Record, error) {
	if m.create == nil {
		return nil, ErrImmutable
	}
	s, err := m.create(kind)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	// gen id if needed
	if id == "" {
		var err error
		id, err = m.idFunc(kind, m.idExists)
		if err != nil {
			return nil, fmt.Errorf("bad id %w", err)
		}
	}

	r := &Record{
		Id:      id,
		Kind:    kind,
		Service: s,
	}
	m.known[id] = r
	return r, nil
}

func (m *Map) deleteRecord(id string) (*Record, error) {
	if m.create == nil {
		return nil, ErrImmutable
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.known[id]
	if !ok {
		return nil, ErrNotFound
	}
	delete(m.known, id)
	return r, nil
}

// idExists returns whether the given id exists in m.known.
// m.mu lock must be held when calling this method
func (m *Map) idExists(id string) bool {
	return maps.Exists(m.known, id)
}

// Change represents a change to a Map.
type Change struct {
	ChangeTime time.Time
	ChangeType types.ChangeType
	OldValue   *Record
	NewValue   *Record
}
