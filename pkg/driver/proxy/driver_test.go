package proxy

import (
	"testing"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

func Test_proxy_announceChange(t *testing.T) {
	announcer := &testAnnouncer{}
	proxy := &proxy{
		announcer: announcer,
		logger:    zap.NewNop(),
	}

	known := announcedTraits{}
	// new child
	proxy.announceChange(known, &traits.PullChildrenResponse_Change{NewValue: &traits.Child{
		Name: "child01",
		Traits: []*traits.Trait{
			{Name: trait.OnOff.String()},
			{Name: trait.Hail.String()},
		},
	}})

	if _, ok := known[childTrait{name: "child01", trait: trait.OnOff}]; !ok {
		t.Errorf("Expecting child01:OnOff to be remembered, got %v", known)
	}
	if _, ok := known[childTrait{name: "child01", trait: trait.Hail}]; !ok {
		t.Errorf("Expecting child01:Hail to be remembered, got %v", known)
	}

	// child has a new trait
	proxy.announceChange(known, &traits.PullChildrenResponse_Change{NewValue: &traits.Child{
		Name: "child01",
		Traits: []*traits.Trait{
			{Name: trait.OnOff.String()},
			{Name: trait.Hail.String()},
			{Name: trait.Light.String()},
		},
	}})

	if _, ok := known[childTrait{name: "child01", trait: trait.OnOff}]; !ok {
		t.Errorf("Expecting child01:OnOff to be remembered, got %v", known)
	}
	if _, ok := known[childTrait{name: "child01", trait: trait.Hail}]; !ok {
		t.Errorf("Expecting child01:Hail to be remembered, got %v", known)
	}
	if _, ok := known[childTrait{name: "child01", trait: trait.Light}]; !ok {
		t.Errorf("Expecting child01:Light to be remembered, got %v", known)
	}

	// child has a trait removed
	proxy.announceChange(known, &traits.PullChildrenResponse_Change{
		OldValue: &traits.Child{
			Name: "child01",
			Traits: []*traits.Trait{
				{Name: trait.OnOff.String()},
				{Name: trait.Hail.String()},
				{Name: trait.Light.String()},
			},
		},
		NewValue: &traits.Child{
			Name: "child01",
			Traits: []*traits.Trait{
				{Name: trait.OnOff.String()},
				{Name: trait.Light.String()},
			},
		},
	})

	if _, ok := known[childTrait{name: "child01", trait: trait.OnOff}]; !ok {
		t.Errorf("Expecting child01:OnOff to be remembered, got %v", known)
	}
	if _, ok := known[childTrait{name: "child01", trait: trait.Light}]; !ok {
		t.Errorf("Expecting child01:Light to be remembered, got %v", known)
	}
	if _, ok := known[childTrait{name: "child01", trait: trait.Hail}]; ok {
		t.Errorf("Expecting child01:Hail to be forgotten, got %v", known)
	}
}

type testAnnouncer []*announcement

func (t *testAnnouncer) Announce(name string, features ...node.Feature) node.Undo {
	an := &announcement{name: name, features: features}
	*t = append(*t, an)
	return func() {
		an.undone++
	}
}

type announcement struct {
	name     string
	features []node.Feature
	undone   int
}
