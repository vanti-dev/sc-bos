package node

import (
	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
)

func TestNode_Announce_metadata(t *testing.T) {
	n := New("test")
	d1Metadata := metadata.NewModel()
	d2Metadata := metadata.NewModel()

	r := metadata.NewApiRouter()
	r.Add("d1", metadata.WrapApi(metadata.NewModelServer(d1Metadata)))
	r.Add("d2", metadata.WrapApi(metadata.NewModelServer(d2Metadata)))

	n.Support(Routing(r))

	n.Announce("d1", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D1"}}))
	got, err := d1Metadata.GetMetadata()
	if err != nil {
		t.Fatal(err)
	}
	want := &traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D1"}}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D1 Metadata (-want,+got)\n%s", diff)
	}

	n.Announce("d2", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D2"}}))
	got, err = d2Metadata.GetMetadata()
	if err != nil {
		t.Fatal(err)
	}
	want = &traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D2"}}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D2 Metadata (-want,+got)\n%s", diff)
	}

	// announce again with more metadata
	n.Announce("d2", HasMetadata(&traits.Metadata{Appearance: &traits.Metadata_Appearance{Description: "Device 2"}}))
	got, err = d2Metadata.GetMetadata()
	if err != nil {
		t.Fatal(err)
	}
	want = &traits.Metadata{Appearance: &traits.Metadata_Appearance{Title: "D2", Description: "Device 2"}}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("D2 Metadata (-want,+got)\n%s", diff)
	}
}
