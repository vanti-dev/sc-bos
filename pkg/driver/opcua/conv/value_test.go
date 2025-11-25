package conv

import (
	"testing"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

func Test_ToTraitEnum(t *testing.T) {

	s := map[string]string{
		"0": "OPEN",
		"1": "OPENING",
		"2": "CLOSED",
		"3": "CLOSING",
	}
	status, err := ToTraitEnum[gen.Transport_Door_DoorStatus]("1", s, gen.Transport_Door_DoorStatus_value)
	if err != nil {
		t.Errorf("ToTraitEnum failed: %v", err)
	}
	if status != gen.Transport_Door_OPENING {
		t.Errorf("ToTraitEnum failed: %v", status)
	}
}
