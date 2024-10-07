package node

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

func TestNode_Announce_metadata(t *testing.T) {
	n := New("test")
	expectTraits := []*traits.TraitMetadata{{Name: string(trait.Metadata)}}

	n.Announce("d1", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D1"}}))
	got, err := n.allMetadata.GetMetadata("d1")
	if err != nil {
		t.Fatal(err)
	}
	want := &traits.Metadata{Name: "d1", Appearance: &traits.Metadata_Appearance{Title: "D1"}, Traits: expectTraits}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D1 Metadata (-want,+got)\n%s", diff)
	}

	n.Announce("d2", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D2"}, Traits: expectTraits}))
	got, err = n.allMetadata.GetMetadata("d2")
	if err != nil {
		t.Fatal(err)
	}
	want = &traits.Metadata{Name: "d2", Appearance: &traits.Metadata_Appearance{Title: "D2"}, Traits: expectTraits}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D2 Metadata (-want,+got)\n%s", diff)
	}

	// announce again with more metadata
	n.Announce("d2", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Description: "Device 2"}, Traits: expectTraits}))
	got, err = n.allMetadata.GetMetadata("d2")
	if err != nil {
		t.Fatal(err)
	}
	want = &traits.Metadata{Name: "d2", Appearance: &traits.Metadata_Appearance{Title: "D2", Description: "Device 2"}, Traits: expectTraits}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D2 Metadata (-want,+got)\n%s", diff)
	}

	// announce with multiple HasMetadata calls
	n.Announce("d3",
		HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Description: "Device 3"}}),
		HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D3"}}),
	)
	got, err = n.allMetadata.GetMetadata("d3")
	if err != nil {
		t.Fatal(err)
	}
	want = &traits.Metadata{Name: "d3", Appearance: &traits.Metadata_Appearance{Title: "D3", Description: "Device 3"}, Traits: expectTraits}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D3 Metadata (-want,+got)\n%s", diff)
	}
}

func TestNode_ListAllMetadata(t *testing.T) {
	log, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Logger error %v", err)
	}
	newNode := func() *Node {
		n := New("-test") // - means it shows up before other metadata records
		n.Logger = log
		return n
	}

	// the metadata representing the node itself
	nodeMd := &traits.Metadata{Name: "-test", Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}, {Name: string(trait.Parent)}}}

	t.Run("no announce", func(t *testing.T) {
		n := newNode()
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})
	t.Run("default metadata", func(t *testing.T) {
		n := newNode()
		n.Announce("d1")
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
			{Name: "d1", Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}}},
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})
	t.Run("HasMetadata", func(t *testing.T) {
		n := newNode()
		n.Announce("d1", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "Foo"}}))
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
			{Name: "d1", Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}}, Appearance: &traits.Metadata_Appearance{Title: "Foo"}},
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})
	t.Run("HasMetadata unwanted mutation", func(t *testing.T) {
		n := newNode()
		passedMd := &traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "Foo"}}
		n.Announce("d1", HasMetadata(passedMd))
		want := &traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "Foo"}}
		if diff := cmp.Diff(want, passedMd, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("Announce mutated Metadata (-want,+got)\n%s", diff)
		}
	})
	t.Run("HasMetadata merges", func(t *testing.T) {
		n := newNode()
		n.Announce("d1", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "Foo", Description: "Desc"}}))
		n.Announce("d1", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "Bar"}, Membership: &traits.Metadata_Membership{Subsystem: "Tests"}}))
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
			{Name: "d1", Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}}, Appearance: &traits.Metadata_Appearance{Title: "Bar", Description: "Desc"}, Membership: &traits.Metadata_Membership{Subsystem: "Tests"}},
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})

	t.Run("HasMetadata merges traits", func(t *testing.T) {
		n := newNode()
		n.Announce("d1", HasMetadata(&traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "Foo"}}}))
		n.Announce("d1", HasMetadata(&traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "Bar"}}}))
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
			{Name: "d1", Traits: []*traits.TraitMetadata{{Name: "Bar"}, {Name: "Foo"}, {Name: string(trait.Metadata)}}},
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})

	t.Run("HasTrait", func(t *testing.T) {
		n := newNode()
		n.Announce("d1", HasTrait(trait.Light, NoAddChildTrait()))
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
			{Name: "d1", Traits: []*traits.TraitMetadata{{Name: string(trait.Light)}, {Name: string(trait.Metadata)}}},
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})

	t.Run("HasTrait merges", func(t *testing.T) {
		n := newNode()
		n.Announce("d1", HasTrait(trait.Light, NoAddChildTrait()))
		n.Announce("d1", HasTrait(trait.Booking, NoAddChildTrait()))
		got := n.ListAllMetadata()
		want := []*traits.Metadata{
			nodeMd,
			{Name: "d1", Traits: []*traits.TraitMetadata{{Name: string(trait.Booking)}, {Name: string(trait.Light)}, {Name: string(trait.Metadata)}}},
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), protocmp.Transform()); diff != "" {
			t.Fatalf("ListAllMetadata (-want,+got)\n%s", diff)
		}
	})
}
