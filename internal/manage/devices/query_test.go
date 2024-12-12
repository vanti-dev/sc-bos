package devices

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func Test_getMessageString(t *testing.T) {
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
		{"match in array", "traits.name", "foo", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}}}, true},
		{"match in array with Index", "traits[0].name", "foo", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}}}, true},
		{"match in array doesn't exist", "traits[1].name", "foo", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}}}, false},
		{"match in array negative", "traits[-1].name", "foo", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}}}, false},
		{"match any in array", "traits.name", "foo", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo"}}}, true},
		{"match any in array - no matches", "traits.name", "baz", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo"}}}, false},
		{"match nested in array", "traits.more.units", "cats", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo", More: map[string]string{"units": "cats"}}}}, true},
		{"match in array with primitive", "dns[0]", "bar", &traits.Metadata_NIC{Dns: []string{"bar"}}, true},
		{"match nested in array with wrong Index", "traits[0].more.units", "cats", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo", More: map[string]string{"units": "cats"}}}}, false},
		{"match in array with primitive and Index - no matches", "dns[0]", "foo", &traits.Metadata_NIC{Dns: []string{"bar"}}, false},
		{"match nested in array - no matches", "traits.more.units", "dogs", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo", More: map[string]string{"units": "cats"}}}}, false},
		{"match in array with primitive - no matches", "dns[0]", "foo", &traits.Metadata_NIC{Dns: []string{"bar"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itr := getMessageString(tt.path, tt.msg)

			found := false

			f := func(v string) bool {
				if v == tt.value {
					found = true
				}
				return found
			}

			itr(f)

			if found != tt.want {
				t.Errorf("getMessageString() = %v, want %v", found, tt.want)
			}
		})
	}
}
func Test_messageHasValueStringFunc(t *testing.T) {
	tests := []struct {
		name  string
		value string
		msg   proto.Message
		want  bool
	}{
		{"nil msg", "any", nil, false},
		{"root string", "foo", &traits.Metadata{Name: "foo"}, true},
		{"root string not equal", "foo", &traits.Metadata{Name: "bar"}, false},
		{"root string absent", "foo", &traits.Metadata{}, false},
		{"root map", "foo", &traits.Metadata{More: map[string]string{"val": "foo"}}, true},
		{"root map not equal", "foo", &traits.Metadata{More: map[string]string{"val": "bar"}}, false},
		{"root map nil", "foo", &traits.Metadata{}, false},
		{"root map absent", "foo", &traits.Metadata{More: map[string]string{}}, false},
		{"root map no key match", "val", &traits.Metadata{More: map[string]string{"val": "foo"}}, false},
		{"nested string", "1234", &traits.Metadata{Id: &traits.Metadata_ID{Bacnet: "1234"}}, true},
		{"nested string absent prop", "1234", &traits.Metadata{Id: &traits.Metadata_ID{}}, false},
		{"nested string absent message", "1234", &traits.Metadata{}, false},
		{"nested map", "1234", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, true},
		{"list property", "1234", &traits.Metadata{Nics: []*traits.Metadata_NIC{{DisplayName: "1234"}}}, true},
		// There was a bug caused by incorrectly serialising messages to string before comparing,
		// i.e. comparing against `{ "foo" [] [] 0x9872000020 }`, which caused false matches
		{"special char", "{", &traits.Metadata{Name: "{foo}"}, true},
		{"bad string prop", "{", &traits.Metadata{Id: &traits.Metadata_ID{Bacnet: "1234"}}, false},
		{"bad string map", "{", &traits.Metadata{More: map[string]string{"val": "bar"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := messageHasValueStringFunc(tt.msg, func(v string) bool {
				return strings.Contains(v, tt.value)
			}); got != tt.want {
				t.Errorf("messageHasValueStringFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Example_isMessageValueStringFunc() {
	msg := &gen.Device{
		Name: "MyDevice",
		Metadata: &traits.Metadata{
			Membership: &traits.Metadata_Membership{
				Subsystem: "Lighting",
			},
		},
	}

	itr := getMessageString("metadata.membership.subsystem", msg)

	found := false

	itr(func(v string) bool {
		if v == "Lighting" {
			found = true
		}
		return found
	})

	fmt.Println(found)
	// Output: true
}
