package node

import (
	"context"
	"strings"
	"testing"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/onoff"
)

func TestNode_Client(t *testing.T) {
	t.Run("get client", func(t *testing.T) {
		n := New("Test")
		n.Support(Clients(onoff.WrapApi(onoff.NewModelServer(onoff.NewModel(onoff.WithInitialOnOff(&traits.OnOff{State: traits.OnOff_ON}))))))

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
		want := "traits.OnOffApiClient"
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Expecting error to mention the type we're missing: should contain %s, got %s", want, err)
		}
	})
}
