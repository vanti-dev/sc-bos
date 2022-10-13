package node

import (
	"context"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("get client", func(t *testing.T) {
		n := New("Test")
		n.Support(Clients(traits.NewOnOffApiClient(nil)))

		var client traits.OnOffApiClient
		var err error
		client, err = Client[traits.OnOffApiClient](n)
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		_ = client
	})
	t.Run("no type", func(t *testing.T) {
		n := New("Test")
		client, err := Client[traits.OnOffApiClient](n)
		if err == nil {
			t.Fatalf("Expecting error, got client %v", client)
		}
	})
}

func TestFindClient(t *testing.T) {
	t.Run("get client", func(t *testing.T) {
		n := New("Test")
		n.Support(Clients(traits.NewOnOffApiClient(nil)))

		var client traits.OnOffApiClient
		FindClient(n, &client)
		if client == nil {
			t.Fatalf("Expecting client but got none")
		}
	})
	t.Run("no client", func(t *testing.T) {
		n := New("Test")
		var client traits.OnOffApiClient
		FindClient(n, &client)
		if client != nil {
			t.Fatalf("Expecting no client, got %v", client)
		}
	})
}

func TestNode_Client(t *testing.T) {
	t.Run("get client", func(t *testing.T) {
		n := New("Test")
		n.Support(Clients(onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(traits.OnOff_ON)))))

		var client traits.OnOffApiClient
		err := n.Client(&client)
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		if client == nil {
			t.Fatalf("Expecting client but got none")
		}
		state, err := client.GetOnOff(context.Background(), &traits.GetOnOffRequest{})
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		if state.State != traits.OnOff_ON {
			t.Fatalf("State doesn't match, want %v, got %v", traits.OnOff_ON, state.State)
		}
	})
	t.Run("no client", func(t *testing.T) {
		n := New("Test")
		n.Support(Clients(traits.NewAirTemperatureApiClient(nil)))
		var client traits.OnOffApiClient
		err := n.Client(&client)
		if err == nil {
			t.Fatalf("Expected err, got none")
		}
	})
}
