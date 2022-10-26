package node

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

// Announcer defines the Announce method.
// Calling Announce signals that a given name exists and has the collection of features provided.
// Typically announcing the same name more than once will combine the features of the existing name with the new features.
type Announcer interface {
	// Announce signals that the name exists and has the given features.
	Announce(name string, features ...Feature) Undo
}

// SelfAnnouncer is a complement to node.Announcer allowing a type to announce itself.
type SelfAnnouncer interface {
	AnnounceSelf(a Announcer) Undo
}

// AnnouncerFunc allows adapting a func of the correct signature to implement Announcer
type AnnouncerFunc func(name string, features ...Feature) Undo

func (a AnnouncerFunc) Announce(name string, features ...Feature) Undo {
	return a(name, features...)
}

// AnnounceWithNamePrefix returns an Announcer whose Announce method acts like `Announce(prefix+name, features...)`
func AnnounceWithNamePrefix(prefix string, a Announcer) Announcer {
	return AnnouncerFunc(func(name string, features ...Feature) Undo {
		return a.Announce(prefix+name, features...)
	})
}

type announcement struct {
	name    string
	traits  []traitFeature
	clients []interface{}
	undo    []Undo
}

type traitFeature struct {
	name     trait.Name
	clients  []interface{}
	metadata map[string]string

	noAddChildTrait bool
	noAddMetadata   bool
}

// Feature describes some aspect of a named device.
type Feature interface{ apply(a *announcement) }

// EmptyFeature is a Feature that means nothing.
// It can be embedded in custom Feature types to allow them to extend the capabilities of the Feature system.
type EmptyFeature struct{}

func (e EmptyFeature) apply(a *announcement) {
	// do nothing
}

type featureFunc func(a *announcement)

func (f featureFunc) apply(a *announcement) {
	f(a)
}

// HasClient indicates that the name implements non-trait apis as defined by these clients.
// The clients are still added to routers and all requests on the clients should accept a Name.
func HasClient(clients ...interface{}) Feature {
	return featureFunc(func(a *announcement) {
		a.clients = append(a.clients, clients...)
	})
}

// HasTrait indicates that the device implements the named trait.
func HasTrait(name trait.Name, opt ...TraitOption) Feature {
	return featureFunc(func(a *announcement) {
		feature := traitFeature{name: name}
		for _, option := range opt {
			option(&feature)
		}
		a.traits = append(a.traits, feature)
	})
}

// TraitOption controls how a Node behaves when presented with a new device trait.
type TraitOption func(t *traitFeature)

// WithClients indicates that the trait is implemented by these client instances.
// The clients will be added to the relevant routers when the trait is announced.
func WithClients(client ...interface{}) TraitOption {
	return func(t *traitFeature) {
		t.clients = append(t.clients, client...)
	}
}

// NoAddChildTrait instructs the Node not to add the trait to the nodes parent.Model.
func NoAddChildTrait() TraitOption {
	return func(t *traitFeature) {
		t.noAddChildTrait = true
	}
}

// NoAddMetadata instructs the Node not to add the trait to the nodes traits.Metadata.
func NoAddMetadata() TraitOption {
	return func(t *traitFeature) {
		t.noAddMetadata = true
	}
}

// WithTraitMetadata instructs the Node to use the given metadata when adding the trait to the nodes traits.Metadata.
// Metadata maps will be merged together, with conflicting keys in later calls overriding existing keys.
func WithTraitMetadata(md map[string]string) TraitOption {
	return func(t *traitFeature) {
		if t.metadata == nil {
			t.metadata = make(map[string]string)
		}
		for k, v := range md {
			t.metadata[k] = v
		}
	}
}
