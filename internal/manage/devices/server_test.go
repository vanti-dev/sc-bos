package devices

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func Test_deviceMatchesQuery(t *testing.T) {
	// simple validation test, exhaustive tests are at a lower level
	t.Run("matches", func(t *testing.T) {
		query := &gen.Device_Query{
			Conditions: []*gen.Device_Query_Condition{
				{Field: "metadata.membership.subsystem", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "Lighting"}},
				{Field: "metadata.location.more.floor", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "4"}},
			},
		}
		device := &gen.Device{
			Name: "Light on floor 4",
			Metadata: &traits.Metadata{
				Membership: &traits.Metadata_Membership{
					Subsystem: "Lighting",
				},
				Location: &traits.Metadata_Location{
					More: map[string]string{
						"floor": "4",
					},
				},
			},
		}
		got := deviceMatchesQuery(query, device)
		if !got {
			t.Fatalf("deviceMatchesQuery want true, got false")
		}
	})

	t.Run("not matches", func(t *testing.T) {
		query := &gen.Device_Query{
			Conditions: []*gen.Device_Query_Condition{
				{Field: "metadata.membership.subsystem", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "Lighting"}},
				{Field: "metadata.location.more.floor", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "4"}},
			},
		}
		device := &gen.Device{
			Name: "Light on floor 4",
			Metadata: &traits.Metadata{
				Membership: &traits.Metadata_Membership{
					Subsystem: "Lighting",
				},
				Location: &traits.Metadata_Location{
					More: map[string]string{
						"floor": "5",
					},
				},
			},
		}
		got := deviceMatchesQuery(query, device)
		if got {
			t.Fatalf("deviceMatchesQuery want false, got true")
		}
	})
}

func TestServer_ListDevices(t *testing.T) {
	n := node.New("test")
	{
		r := parent.NewApiRouter()
		n.Support(node.Routing(r), node.Clients(parent.WrapApi(r)))
	}
	{
		r := metadata.NewApiRouter(metadata.WithMetadataApiClientFactory(func(name string) (traits.MetadataApiClient, error) {
			return metadata.WrapApi(metadata.NewModelServer(metadata.NewModel(resource.WithInitialValue(&traits.Metadata{
				Name: name,
			})))), nil
		}))
		n.Support(node.Routing(r), node.Clients(metadata.WrapApi(r)))
	}
	// create 40 devices, half are lights, half are hvac,
	// they're created interleaved to try and avoid page 1 being all lights and page 2 being all hvac
	for i := 0; i < 20; i++ {
		n.Announce(fmt.Sprintf("device/%02d/light", i),
			node.HasMetadata(&traits.Metadata{
				Membership: &traits.Metadata_Membership{
					Subsystem: "Lighting",
				},
			}),
			node.HasTrait(trait.Light))
		n.Announce(fmt.Sprintf("device/%02d/hvac", i),
			node.HasMetadata(&traits.Metadata{
				Membership: &traits.Metadata_Membership{
					Subsystem: "HVAC",
				},
			}),
			node.HasTrait(trait.AirTemperature))
	}

	server := &Server{parentName: "test", node: n}
	server.ChildPageSize = 5 // force multiple pages to be read from the parent

	mdNamed := func(name string) *traits.Metadata {
		return &traits.Metadata{
			Name:       name,
			Membership: &traits.Metadata_Membership{Subsystem: "Lighting"},
			Traits: []*traits.TraitMetadata{
				{Name: trait.Light.String()},
			},
		}
	}

	wantPage1 := []*gen.Device{
		{Name: "device/00/light", Metadata: mdNamed("device/00/light")},
		{Name: "device/01/light", Metadata: mdNamed("device/01/light")},
		{Name: "device/02/light", Metadata: mdNamed("device/02/light")},
		{Name: "device/03/light", Metadata: mdNamed("device/03/light")},
		{Name: "device/04/light", Metadata: mdNamed("device/04/light")},
		{Name: "device/05/light", Metadata: mdNamed("device/05/light")},
		{Name: "device/06/light", Metadata: mdNamed("device/06/light")},
	}
	wantPage2 := []*gen.Device{
		{Name: "device/07/light", Metadata: mdNamed("device/07/light")},
		{Name: "device/08/light", Metadata: mdNamed("device/08/light")},
		{Name: "device/09/light", Metadata: mdNamed("device/09/light")},
		{Name: "device/10/light", Metadata: mdNamed("device/10/light")},
		{Name: "device/11/light", Metadata: mdNamed("device/11/light")},
		{Name: "device/12/light", Metadata: mdNamed("device/12/light")},
		{Name: "device/13/light", Metadata: mdNamed("device/13/light")},
	}
	wantPage3 := []*gen.Device{
		{Name: "device/14/light", Metadata: mdNamed("device/14/light")},
		{Name: "device/15/light", Metadata: mdNamed("device/15/light")},
		{Name: "device/16/light", Metadata: mdNamed("device/16/light")},
		{Name: "device/17/light", Metadata: mdNamed("device/17/light")},
		{Name: "device/18/light", Metadata: mdNamed("device/18/light")},
		{Name: "device/19/light", Metadata: mdNamed("device/19/light")},
	}

	// PAGE 1 - should return a full page
	devices, err := server.ListDevices(context.Background(), &gen.ListDevicesRequest{
		PageSize: 7,
		Query: &gen.Device_Query{Conditions: []*gen.Device_Query_Condition{
			{Field: "metadata.membership.subsystem", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "Lighting"}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantPage1, devices.Devices, protocmp.Transform()); diff != "" {
		t.Fatalf("Page 1 (-want,+got)\n%s", diff)
	}
	if devices.NextPageToken == "" {
		t.Fatalf("Expecting a NextPageToken, but got none")
	}

	// PAGE 2 - should also return a full page
	devices, err = server.ListDevices(context.Background(), &gen.ListDevicesRequest{
		PageSize: 7,
		Query: &gen.Device_Query{Conditions: []*gen.Device_Query_Condition{
			{Field: "metadata.membership.subsystem", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "Lighting"}},
		}},
		PageToken: devices.NextPageToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantPage2, devices.Devices, protocmp.Transform()); diff != "" {
		t.Fatalf("Page 2 (-want,+got)\n%s", diff)
	}
	if devices.NextPageToken == "" {
		t.Fatalf("Expecting a NextPageToken, but got none")
	}

	// PAGE 3 - is not a full page, and should have no page token
	devices, err = server.ListDevices(context.Background(), &gen.ListDevicesRequest{
		PageSize: 7,
		Query: &gen.Device_Query{Conditions: []*gen.Device_Query_Condition{
			{Field: "metadata.membership.subsystem", Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "Lighting"}},
		}},
		PageToken: devices.NextPageToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantPage3, devices.Devices, protocmp.Transform()); diff != "" {
		t.Fatalf("Page 3 (-want,+got)\n%s", diff)
	}
	if devices.NextPageToken != "" {
		t.Fatalf("Expecting no NextPageToken, but got %s", devices.NextPageToken)
	}
}

func assertAllLights(t *testing.T, listResponse *gen.ListDevicesResponse) {
	for i, device := range listResponse.Devices {
		if device.GetMetadata().GetMembership().GetSubsystem() != "Lighting" {
			t.Fatalf("device %d is not a light: %+v", i, device)
		}
	}
}
