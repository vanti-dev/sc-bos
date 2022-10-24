package bacnet

import (
	"github.com/vanti-dev/gobacnet"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"testing"
)

func TestYabeRoom(t *testing.T) {
	client, err := gobacnet.NewClient("bridge100", 0)
	if err != nil {
		t.Fatalf("NewClient error %v", err)
	}
	t.Cleanup(client.Close)

	iAm, err := client.WhoIs(800_000, 900_000)
	if err != nil {
		t.Fatalf("WhoIs error %v", err)
	}

	t.Logf("IAm %+v", iAm)

	netAddr, err := iAm[0].Addr.UDPAddr()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Connecting to %v", netAddr)
	bacAddr := bactypes.UDPToAddress(&netAddr)
	t.Logf("After conversion: iAm Addr:%v, udp addr:%v, bac addr:%v", iAm[0].Addr, netAddr, bacAddr)
	devices, err := client.RemoteDevices(bacAddr, iAm[0].ID.Instance)
	if err != nil {
		t.Fatalf("RemoteDevices error %v", err)
	}
	t.Logf("RemoteDevices %+v", devices)
}
