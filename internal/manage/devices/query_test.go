package devices

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func Test_getMessageString(t *testing.T) {
	tests := []struct {
		name string
		path string
		msg  proto.Message
		want []string
	}{
		{"nil msg", "foo.bar", nil, nil},
		{"root string", "name", &traits.Metadata{Name: "foo"}, []string{"foo"}},
		{"root string absent", "name", &traits.Metadata{}, []string{}},
		{"root map", "more.val", &traits.Metadata{More: map[string]string{"val": "foo"}}, []string{"foo"}},
		{"root map nil", "more.val", &traits.Metadata{}, []string{}},
		{"root map absent", "more.val", &traits.Metadata{More: map[string]string{}}, []string{}},
		{"nested string", "id.bacnet", &traits.Metadata{Id: &traits.Metadata_ID{Bacnet: "1234"}}, []string{"1234"}},
		{"nested string absent prop", "id.bacnet", &traits.Metadata{Id: &traits.Metadata_ID{}}, []string{}},
		{"nested string absent message", "id.bacnet", &traits.Metadata{}, []string{}},
		{"nested map", "id.more.foo", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, []string{"1234"}},
		{"trailing .", "id.", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, []string{}},
		{"leading .", ".id", &traits.Metadata{Id: &traits.Metadata_ID{More: map[string]string{"foo": "1234"}}}, []string{}},
		{"property of scalar", "name.foo", &traits.Metadata{Name: "1234"}, []string{}},
		{"match all (one) in array", "traits.name", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}}}, []string{"foo"}},
		{"match all in array", "traits.name", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}, {Name: "bar"}}}, []string{"foo", "bar"}},
		{"match in array with Index", "traits[0].name", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}, {Name: "bar"}}}, []string{"foo"}},
		{"match in array with Index", "traits[1].name", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}, {Name: "bar"}}}, []string{"bar"}},
		{"match in array doesn't exist", "traits[1].name", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "foo"}}}, []string{}},
		{"match in array negative", "traits[-1].name", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo"}}}, []string{"foo"}},
		{"match nested in array", "traits.more.units", &traits.Metadata{Traits: []*traits.TraitMetadata{{More: map[string]string{"units": "dogs"}}, {More: map[string]string{"units": "cats"}}}}, []string{"dogs", "cats"}},
		{"match all in array with primitive", "dns", &traits.Metadata_NIC{Dns: []string{"foo", "bar"}}, []string{"foo", "bar"}},
		{"match in array with primitive[0]", "dns[0]", &traits.Metadata_NIC{Dns: []string{"foo", "bar"}}, []string{"foo"}},
		{"match in array with primitive[1]", "dns[1]", &traits.Metadata_NIC{Dns: []string{"foo", "bar"}}, []string{"bar"}},
		{"match in array with primitive[-1]", "dns[-1]", &traits.Metadata_NIC{Dns: []string{"foo", "bar"}}, []string{"bar"}},
		{"match in array with primitive[-2]", "dns[-2]", &traits.Metadata_NIC{Dns: []string{"foo", "bar"}}, []string{"foo"}},
		{"match nested in array with wrong Index", "traits[0].more.units", &traits.Metadata{Traits: []*traits.TraitMetadata{{Name: "bar"}, {Name: "foo", More: map[string]string{"units": "cats"}}}}, []string{}},
	}
	cmpStr := func(a, b string) bool { return a < b }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(getMessageString(tt.path, tt.msg))
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty(), cmpopts.SortSlices(cmpStr)); diff != "" {
				t.Errorf("getMessageString() -want +got:\n%s", diff)
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
