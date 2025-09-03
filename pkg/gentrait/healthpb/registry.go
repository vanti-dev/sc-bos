package healthpb

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var (
	// ErrAlreadyExists is returned when attempting to create a health check with the same name and ID as an existing check.
	ErrAlreadyExists = errors.New("already exists")
	// ErrInvalid is returned when arguments to Registry or Checks methods are invalid.
	ErrInvalid = errors.New("invalid")
)

// Registry manages a collection of health checks.
type Registry struct {
	mu     sync.RWMutex
	byName map[string]*namedChecks

	onCheckUpdate func(name string, c *gen.HealthCheck)
	onCheckDelete func(name, id string)
}

// ForOwner returns a Checks instance that can create checks owned by the given owner.
// The owner should be unique within the node.
func (r *Registry) ForOwner(owner string) *Checks {
	return &Checks{
		r:     r,
		owner: owner,
	}
}

// addCheck adds c, identified by id, as a health check against name.
// ErrAlreadyExists is returned if id already exists for name.
func (r *Registry) addCheck(name, id string, c *checkBase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.byName == nil {
		r.byName = make(map[string]*namedChecks)
	}
	nc, ok := r.byName[name]
	if !ok {
		nc = &namedChecks{
			n:    name,
			byId: make(map[string]*checkBase),
		}
		r.byName[name] = nc
	}
	_, exists := nc.byId[id]
	if exists {
		return fmt.Errorf("%w: %s[%s]", ErrAlreadyExists, name, id)
	}
	c.onCommit = func(c *gen.HealthCheck) {
		if r.onCheckUpdate != nil {
			r.onCheckUpdate(name, c)
		}
	}
	c.onDispose = func(c *gen.HealthCheck) {
		// clean up our own state
		r.mu.Lock()
		delete(nc.byId, id)
		if len(nc.byId) == 0 {
			delete(r.byName, name)
		}
		r.mu.Unlock()
		if r.onCheckDelete != nil {
			r.onCheckDelete(name, id)
		}
	}
	nc.byId[id] = c
	return nil
}

// Checks provides a mechanism to create managed health checks against named devices.
//
// For each factory method:
//   - [gen.HealthCheck.Id] must be unique per name per [Checks] instance, or [ErrAlreadyExists] is returned. An absent id is equivalent to the empty string.
//   - The [gen.HealthCheck] instance must not be modified after the call.
//   - [ErrInvalid] is returned if the passed [gen.HealthCheck] is not valid for the type of check being created.
type Checks struct {
	r     *Registry
	owner string
}

// NewBoundsCheck creates a new health check that updates state by comparing a value against bounds.
func (hc *Checks) NewBoundsCheck(name string, c *gen.HealthCheck) (*BoundsCheck, error) {
	hc.adjustId(c)
	check, err := newBoundsCheck(c)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if err := hc.r.addCheck(name, c.Id, check.checkBase); err != nil {
		return nil, err
	}
	return check, nil
}

// NewErrorCheck creates a new health check that tracks state via an error value.
func (hc *Checks) NewErrorCheck(name string, c *gen.HealthCheck) (*ErrorCheck, error) {
	hc.adjustId(c)
	check, err := newErrorCheck(c)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if err := hc.r.addCheck(name, c.Id, check.checkBase); err != nil {
		return nil, err
	}
	return check, nil
}

func (hc *Checks) adjustId(c *gen.HealthCheck) {
	switch {
	case c.Id != "":
		c.Id = hc.owner + ":" + c.Id
	default:
		c.Id = hc.owner
	}
}

type namedChecks struct {
	n    string
	byId map[string]*checkBase
}
