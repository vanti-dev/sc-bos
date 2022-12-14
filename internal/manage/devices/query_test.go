package devices

import (
	"fmt"
	"testing"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func Test_isMessageValueEqualString(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		value string
		msg   proto.Message
		want  bool
	}{
		{"nil msg", "foo.bar", "any", nil, false},
		{"root string", "name", "foo", &traits.Metadata{Name: "foo"}, true},
		{"root string not equal", "name", "foo", &traits.Metadata{Name: "bar"}, false},
		{"root string absent", "name", "foo", &traits.Metadata{}, false},
		{"root map", "more.val", "foo", &traits.Metadata{More: map[string]string{"val": "foo"}}, true},
		{"root map not equal", "more.val", "foo", &traits.Metadata{More: map[string]string{"val": "bar"}}, false},
		{"root map nil", "more.val", "foo", &traits.Metadata{}, false},
		{"root map absent", "more.val", "foo", &traits.Metadata{More: map[string]string{}}, false},
		{"nested string", "id.bacnet", "1234", &traits.Metadata{Id: &traits.Metadata_ID{Bacnet: "1234"}}, true},
		{"nested string not equal", "id.bacnet", "1234", &traits.Metadata{Id: &traits.Metadata_ID{Bacnet: "not 1234"}}, false},
		{"nested string absent prop", "id.bacnet", "1234", &traits.Metadata{Id: &traits.Metadata_ID{}}, false},
		{"nested string absent message", "id.bacnet", "1234", &traits.Metadata{}, false},
		{"nested map", "id.more.foo", "1234", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, true},
		{"trailing .", "id.", "1234", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, false},
		{"leading .", ".id", "1234", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, false},
		{"property of scalar", "name.foo", "1234", &traits.Metadata{Name: "1234"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMessageValueEqualString(tt.path, tt.value, tt.msg); got != tt.want {
				t.Errorf("isMessageValueEqualString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Example_isMessageValueEqualString() {
	msg := &gen.Device{
		Name: "MyDevice",
		Metadata: &traits.Metadata{
			Membership: &traits.Metadata_Membership{
				Subsystem: "Lighting",
			},
		},
	}

	member := isMessageValueEqualString("metadata.membership.subsystem", "Lighting", msg)
	fmt.Println(member)
	// Output: true
}
