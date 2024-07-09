package node

import (
	"context"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
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

// SelfAnnouncerFunc allows adapting a func of the correct signature to implement SelfAnnouncer
type SelfAnnouncerFunc func(a Announcer) Undo

func (sa SelfAnnouncerFunc) AnnounceSelf(a Announcer) Undo {
	return sa(a)
}

// AnnounceWithNamePrefix returns an Announcer whose Announce method acts like `Announce(prefix+name, features...)`
func AnnounceWithNamePrefix(prefix string, a Announcer) Announcer {
	return AnnouncerFunc(func(name string, features ...Feature) Undo {
		return a.Announce(prefix+name, features...)
	})
}

// AnnounceContext returns a new Announcer that undoes any announcements when ctx is Done.
// This leaks a go routine if ctx is never done.
func AnnounceContext(ctx context.Context, a Announcer) Announcer {
	mu := sync.Mutex{}
	var undos []Undo
	undone := make(chan struct{})
	go func() {
		defer close(undone)
		<-ctx.Done()
		mu.Lock()
		defer mu.Unlock()
		for _, undo := range undos {
			undo()
		}
		undos = nil
	}()
	return AnnouncerFunc(func(name string, features ...Feature) Undo {
		select {
		case <-ctx.Done():
			return NilUndo
		default:
		}

		undo := a.Announce(name, features...)
		mu.Lock()
		defer mu.Unlock()

		select {
		case <-undone:
			// undos have been called already, we missed the boat
			undo()
			return NilUndo
		default:
		}

		undo = UndoOnce(undo)
		undos = append(undos, undo)
		return undo
	})
}

// AnnounceFeatures returns an Announcer that acts like `Announce(name, [moreFeatures..., features...])`
func AnnounceFeatures(a Announcer, moreFeatures ...Feature) Announcer {
	return AnnouncerFunc(func(name string, features ...Feature) Undo {
		// do it this way to avoid modifying modeFeatures if it happens to have capacity
		allFeatures := append([]Feature{}, moreFeatures...)
		allFeatures = append(allFeatures, features...)
		return a.Announce(name, allFeatures...)
	})
}

type announcement struct {
	name           string
	traits         []traitFeature
	clients        []any
	metadata       []*traits.Metadata
	noAutoMetadata bool
	undo           []Undo
}

type traitFeature struct {
	name     trait.Name
	clients  []any
	metadata map[string]string

	noAddChildTrait bool
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
func HasClient(clients ...any) Feature {
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

// HasMetadata merges the given metadata into any existing metadata held against the device name.
// Merging metadata is not commutative, two calls to Announce causes the second call to overwrite the first for all common fields.
// Announce accepts multiple HasMetadata features acting as if Announce were called in sequence with each HasMetadata feature.
//
// See metadata.Model.MergeMetadata for more details.
// If md is nil, does not adjust the announcement.
func HasMetadata(md *traits.Metadata) Feature {
	return featureFunc(func(a *announcement) {
		if md == nil {
			return // do nothing, helps callers to avoid nil checks
		}
		// We clone because if this is the first time the name has been associated with metadata,
		// then the passed md is used as is instead of cloning which can cause unexpected mutation
		// from the pov of the caller.
		a.metadata = append(a.metadata, proto.Clone(md).(*traits.Metadata))
	})
}

// HasNoAutoMetadata indicates that announcing the device should not announce automatic metadata for the device.
// Announcing will normally inspect all traits and announce basic metadata for them, this feature turns that off.
// HasMetadata will still be applied.
func HasNoAutoMetadata() Feature {
	return featureFunc(func(a *announcement) {
		a.noAutoMetadata = true
	})
}

// TraitOption controls how a Node behaves when presented with a new device trait.
type TraitOption func(t *traitFeature)

// WithClients indicates that the trait is implemented by these client instances.
// The clients will be added to the relevant routers when the trait is announced.
func WithClients(client ...any) TraitOption {
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
