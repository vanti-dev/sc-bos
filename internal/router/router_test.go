package router

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

// tests overall type behavior: registering services, adding routes, and routing requests with correct priority.
func TestRouter(t *testing.T) {
	r := New()
	r.SupportService(routedRegistryService(t, traits.OnOffApi_ServiceDesc.ServiceName, "name"))
	r.SupportService(routedRegistryService(t, traits.OccupancySensorApi_ServiceDesc.ServiceName, "name"))
	r.SupportService(routedRegistryService(t, traits.AirQualitySensorApi_ServiceDesc.ServiceName, "name"))

	fooModel := onoff.NewModel(resource.WithInitialValue(&traits.OnOff{State: traits.OnOff_OFF}))
	defaultModel := onoff.NewModel(resource.WithInitialValue(&traits.OnOff{State: traits.OnOff_ON}))
	occupancyModel := occupancysensor.NewModel(resource.WithInitialValue(&traits.Occupancy{State: traits.Occupancy_OCCUPIED}))

	// register a specific route for "foo"
	err := r.AddRoute("", "foo",
		wrap.ServerToClient(traits.OnOffApi_ServiceDesc, onoff.NewModelServer(fooModel)))
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}
	// register a specific route for "foo" for the occupancy service - this should have higher priority
	err = r.AddRoute(traits.OccupancySensorApi_ServiceDesc.ServiceName, "foo",
		wrap.ServerToClient(traits.OccupancySensorApi_ServiceDesc, occupancysensor.NewModelServer(occupancyModel)))
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}
	// add a catch-all for all OnOffApi requests that are not to "foo"
	err = r.AddRoute(traits.OnOffApi_ServiceDesc.ServiceName, "",
		wrap.ServerToClient(traits.OnOffApi_ServiceDesc, onoff.NewModelServer(defaultModel)))
	if err != nil {
		t.Fatalf("failed to add route: %v", err)
	}

	conn := NewLoopback(r)
	onOffClient := traits.NewOnOffApiClient(conn)
	occupancyClient := traits.NewOccupancySensorApiClient(conn)
	airQualityClient := traits.NewAirQualitySensorApiClient(conn)
	modeClient := traits.NewModeApiClient(conn)
	// "foo" should route to the fooModel
	res, err := onOffClient.GetOnOff(context.Background(), &traits.GetOnOffRequest{Name: "foo"})
	if err != nil {
		t.Errorf("failed to get onoff for foo: %v", err)
	} else if res.State != traits.OnOff_OFF {
		t.Errorf("expected OFF for foo, got %v", res.State)
	}
	// "bar" (or anything that's not "foo") should route to the defaultModel
	res, err = onOffClient.GetOnOff(context.Background(), &traits.GetOnOffRequest{Name: "bar"})
	if err != nil {
		t.Errorf("failed to get onoff for bar: %v", err)
	} else if res.State != traits.OnOff_ON {
		t.Errorf("expected ON for bar, got %v", res.State)
	}
	// "foo" for the occupancy service should route to the occupancyModel
	res2, err := occupancyClient.GetOccupancy(context.Background(), &traits.GetOccupancyRequest{Name: "foo"})
	if err != nil {
		t.Errorf("failed to get occupancy for foo: %v", err)
	} else if res2.State != traits.Occupancy_OCCUPIED {
		t.Errorf("expected OCCUPIED for foo, got %v", res2.State)
	}
	// "bar" for the occupancy service should fail to resolve
	_, err = occupancyClient.GetOccupancy(context.Background(), &traits.GetOccupancyRequest{Name: "bar"})
	if statusErr, _ := status.FromError(err); statusErr.Code() != codes.NotFound {
		t.Errorf("expected NotFound for bar, got %v", statusErr)
	}
	// there are no matching routes registered for the air quality service on device "bar", so it should fail to resolve
	_, err = airQualityClient.GetAirQuality(context.Background(), &traits.GetAirQualityRequest{Name: "bar"})
	if statusErr, _ := status.FromError(err); statusErr.Code() != codes.NotFound {
		t.Errorf("expected NotFound for air quality, got %v", statusErr)
	}
	// the mode service isn't registered on the router so this should fail, even though there is an all-service route
	// for "foo"
	_, err = modeClient.GetModeValues(context.Background(), &traits.GetModeValuesRequest{Name: "foo"})
	if statusErr, _ := status.FromError(err); statusErr.Code() != codes.Unimplemented {
		t.Errorf("expected Unimplemented for mode, got %v", statusErr)
	}

}

func routedRegistryService(t *testing.T, serviceName, keyName string) *Service {
	t.Helper()
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(serviceName))
	if err != nil {
		t.Fatalf("descriptor for service %q not in registry: %v", serviceName, err)
	}
	servDesc, ok := desc.(protoreflect.ServiceDescriptor)
	if !ok {
		t.Fatalf("%q is not a service", serviceName)
	}
	s, err := NewRoutedService(servDesc, keyName)
	if err != nil {
		t.Fatalf("failed to create routed service: %v", err)
	}
	return s
}
