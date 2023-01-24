package devices

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadata"
	"github.com/smart-core-os/sc-golang/pkg/trait/parent"
	"google.golang.org/protobuf/testing/protocmp"

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

	wantPage1 := []*gen.Device{
		{Name: "device/00/light", Metadata: &traits.Metadata{Name: "device/00/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/01/light", Metadata: &traits.Metadata{Name: "device/01/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/02/light", Metadata: &traits.Metadata{Name: "device/02/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/03/light", Metadata: &traits.Metadata{Name: "device/03/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/04/light", Metadata: &traits.Metadata{Name: "device/04/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/05/light", Metadata: &traits.Metadata{Name: "device/05/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/06/light", Metadata: &traits.Metadata{Name: "device/06/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
	}
	wantPage2 := []*gen.Device{
		{Name: "device/07/light", Metadata: &traits.Metadata{Name: "device/07/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/08/light", Metadata: &traits.Metadata{Name: "device/08/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/09/light", Metadata: &traits.Metadata{Name: "device/09/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/10/light", Metadata: &traits.Metadata{Name: "device/10/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/11/light", Metadata: &traits.Metadata{Name: "device/11/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/12/light", Metadata: &traits.Metadata{Name: "device/12/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/13/light", Metadata: &traits.Metadata{Name: "device/13/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
	}
	wantPage3 := []*gen.Device{
		{Name: "device/14/light", Metadata: &traits.Metadata{Name: "device/14/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/15/light", Metadata: &traits.Metadata{Name: "device/15/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/16/light", Metadata: &traits.Metadata{Name: "device/16/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/17/light", Metadata: &traits.Metadata{Name: "device/17/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/18/light", Metadata: &traits.Metadata{Name: "device/18/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
		{Name: "device/19/light", Metadata: &traits.Metadata{Name: "device/19/light", Membership: &traits.Metadata_Membership{Subsystem: "Lighting"}}},
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
