package tc3dali

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNetID_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`"1.2.3.4.5.6"`)
	var netID NetID
	err := json.Unmarshal(jsonData, &netID)
	if err != nil {
		t.Error(err)
	}
	expect := NetID{1, 2, 3, 4, 5, 6}
	if diff := cmp.Diff(expect, netID); diff != "" {
		t.Errorf("mismatched NetID (-want +got):\n%s", diff)
	}
}

func TestNetID_MarshalJSON(t *testing.T) {
	netID := NetID{6, 5, 4, 3, 2, 1}
	got, err := json.Marshal(netID)
	if err != nil {
		t.Error(err)
	}
	expect := []byte(`"6.5.4.3.2.1"`)
	if diff := cmp.Diff(expect, got); diff != "" {
		t.Errorf("mismatched JSON (-want +got):\n%s", diff)
	}
}
