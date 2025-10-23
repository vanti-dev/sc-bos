package devices

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/internal/manage/devices/testdata/testproto/querypb"
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

func Test_leaf_toString(t *testing.T) {
	fd := func(field string) protoreflect.FieldDescriptor {
		md := (&querypb.Result{}).ProtoReflect().Descriptor()
		return md.Fields().ByName(protoreflect.Name(field))
	}

	tests := []struct {
		leaf   value
		want   string
		wantOk bool
	}{
		{leaf: value{fd("double_val"), protoreflect.ValueOfFloat64(1.1)}, want: "1.1", wantOk: true},
		{leaf: value{fd("float_val"), protoreflect.ValueOfFloat32(2.1)}, want: "2.1", wantOk: true},
		{leaf: value{fd("int32_val"), protoreflect.ValueOfInt32(-10)}, want: "-10", wantOk: true},
		{leaf: value{fd("int64_val"), protoreflect.ValueOfInt64(-20)}, want: "-20", wantOk: true},
		{leaf: value{fd("uint32_val"), protoreflect.ValueOfUint32(30)}, want: "30", wantOk: true},
		{leaf: value{fd("uint64_val"), protoreflect.ValueOfUint64(40)}, want: "40", wantOk: true},
		{leaf: value{fd("sint32_val"), protoreflect.ValueOfInt32(-50)}, want: "-50", wantOk: true},
		{leaf: value{fd("sint64_val"), protoreflect.ValueOfInt64(-60)}, want: "-60", wantOk: true},
		{leaf: value{fd("fixed32_val"), protoreflect.ValueOfUint32(70)}, want: "70", wantOk: true},
		{leaf: value{fd("fixed64_val"), protoreflect.ValueOfUint64(80)}, want: "80", wantOk: true},
		{leaf: value{fd("sfixed32_val"), protoreflect.ValueOfInt32(-90)}, want: "-90", wantOk: true},
		{leaf: value{fd("sfixed64_val"), protoreflect.ValueOfInt64(-100)}, want: "-100", wantOk: true},
		{leaf: value{fd("bool_val"), protoreflect.ValueOfBool(true)}, want: "true", wantOk: true},
		{leaf: value{fd("enum_val"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_A.Number())}, want: "RESULT_ENUM_A", wantOk: true},
		{leaf: value{fd("duration_val"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 10).ProtoReflect())}, want: "10s", wantOk: true},
		{leaf: value{fd("timestamp_val"), protoreflect.ValueOfMessage(timestamppb.New(time.Date(2025, 8, 10, 12, 5, 40, 111, time.UTC)).ProtoReflect())}, want: "2025-08-10T12:05:40.000000111Z", wantOk: true},
	}
	for _, tt := range tests {
		t.Run(string(tt.leaf.fd.Name()), func(t *testing.T) {
			gotS, gotOk := tt.leaf.toString()
			if gotS != tt.want || gotOk != tt.wantOk {
				t.Errorf("value.toString() = %q, %v; want %q, %v", gotS, gotOk, tt.want, tt.wantOk)
			}
		})
	}
}

func Test_rangeValues(t *testing.T) {
	// m is the data structure we test against
	m := testData()

	// leafPaths are all leaf paths in m
	leafPaths := testDataLeafPaths()
	leafValues := testDataLeafValues()
	if len(leafValues) != len(leafPaths) {
		t.Fatalf("expected leafValues and leafPaths to have the same length, got %d and %d", len(leafValues), len(leafPaths))
	}

	type testCase struct {
		path     string
		wantVals []value
		wantErr  bool
	}
	tests := []testCase{
		// no path means all leafPaths
		{wantVals: leafValues},
	}

	// newLeafCase returns a testCase matching path against indexes in leafValues.
	newLeafCase := func(path string, wantIndexes ...int) testCase {
		t.Helper()
		if len(wantIndexes) == 0 {
			t.Fatalf("wantIndexes must not be empty for path %q", path)
		}
		wantLeafs := make([]value, len(wantIndexes))
		for i, idx := range wantIndexes {
			if idx < 0 || idx >= len(leafValues) {
				t.Fatalf("index %d out of range for path %q", idx, path)
			}
			wantLeafs[i] = leafValues[idx]
		}
		return testCase{
			path:     path,
			wantVals: wantLeafs,
		}
	}

	// each path should resolve to exactly the wanted value
	for i, path := range leafPaths {
		tests = append(tests, newLeafCase(path, i))
	}

	// repeated values should be optional in the path
	for i, scalarType := range scalarTypes {
		// in leafValues, each scalar type has the non-repeated value followed by two repeated values
		leafIdx := i*3 + 1
		tests = append(tests, newLeafCase(fmt.Sprintf("r_%s", scalarType), leafIdx, leafIdx+1))
	}

	// negative indexes should work for repeated values
	for i, scalarType := range scalarTypes {
		// in leafValues, each scalar type has the non-repeated value followed by two repeated values
		leafIdx := i*3 + 1
		tests = append(tests, newLeafCase(fmt.Sprintf("r_%s[-1]", scalarType), leafIdx+1))
		tests = append(tests, newLeafCase(fmt.Sprintf("r_%s[-2]", scalarType), leafIdx))
	}

	// tests for when path doesn't match any value, return empty iterator
	for _, path := range []string{
		// value doesn't exist, or is default
		"r_result[2].bool_val",
		"r_result[2].int32_val",
		"r_result[2].string_val",
		"r_result[2].timestamp_val",
		"m_string_range.c.bool_val",
		"m_string_range.c.int32_val",
		"m_string_range.c.string_val",
		"m_string_range.c.timestamp_val",
		"m_string_range.c.result.bool_val",
		"m_string_range.c.result.int32_val",
		"m_string_range.c.result.string_val",
		"m_string_range.c.result.timestamp_val",
		// path isn't a value (or repeated value)
		"result",
		"m_string_result",
		"m_string_string",
		"r_result",
		"r_result.result",
		"r_result[0].result",
		// message fields that don't exist
		"not_found",
		"r_string.not_found",
		"r_string[0].not_found",
		"result.not_found",
		"r_result[0].not_found",
		"m_string_result.a.not_found",
		// fields of non-message types
		"double_val.not_found",
		"string_val.not_found",
		"result.string_val.not_found",
		// repeated field indexes that don't exist
		"r_string[2]",
		"r_string[-3]",
		"r_result[0].r_string[2]",
		"r_result[3].string_val",
		// map keys that don't exist
		"m_int32_string.42",
		"m_string_string.not_found",
		"range.m_string_string.not_found",
		"r_range[0].m_string_string.not_found",
		"m_string_range.a.m_string_string.not_found",
		"m_string_range.not_found.string_val",
		// attempting to access a non-repeating field by index
		"string_val[0]",
		"range[0].string_val",
		"m_string_string[0]",
		"m_string_result[0]",
		// fields of value messages
		"timestamp_val.seconds",
		"r_timestamp[0].seconds",
		"m_string_timestamp.a.seconds",
		// hidden fields of map entries
		"m_int32_string.key",
		"m_int32_string.value",
		// bad map keys
		"m_int32_string.not_int",
		"m_int64_string.not_int",
		"m_uint32_string.not_uint",
		"m_uint64_string.not_uint",
		"m_sint32_string.not_sint",
		"m_sint64_string.not_sint",
		"m_fixed32_string.not_fixed",
		"m_fixed64_string.not_fixed",
		"m_sfixed32_string.not_sfixed",
		"m_sfixed64_string.not_sfixed",
		"m_bool_string.not_bool",
	} {
		tests = append(tests, testCase{path: path})
	}

	// invalid paths return an error, most cases are covered by the parsePath tests
	tests = append(tests, testCase{path: ".", wantErr: true})

	for _, tt := range tests {
		name := tt.path
		if name == "" {
			name = "all paths"
		}
		t.Run(name, func(t *testing.T) {
			got, err := rangeValuesOptions{Stable: true}.Range(tt.path, m)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("rangeValues() error = %v, wantErr %v. Got %v", err, tt.wantErr, slices.Collect(got))
				}
				return
			}
			var wantFds []protoreflect.FieldDescriptor
			var wantValues []protoreflect.Value
			for _, l := range tt.wantVals {
				wantFds = append(wantFds, l.fd)
				wantValues = append(wantValues, l.v)
			}

			var gotFds []protoreflect.FieldDescriptor
			var gotValues []protoreflect.Value
			for l := range got {
				gotFds = append(gotFds, l.fd)
				gotValues = append(gotValues, l.v)
			}
			if diff := cmp.Diff(wantFds, gotFds, transformReflectFieldDescriptor, protocmp.Transform()); diff != "" {
				t.Errorf("rangeValues() field descriptors (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(wantValues, gotValues, transformReflectValue, protocmp.Transform()); diff != "" {
				t.Errorf("rangeValues() values (-want,+got):\n%s", diff)
			}
		})
	}
}

var transformReflectFieldDescriptor = cmp.Transformer("", func(fd protoreflect.FieldDescriptor) any {
	// we only care about the name and number of the field descriptor
	return fmt.Sprintf("%s:%d", fd.FullName(), fd.Number())
})

// copied from the protorange package
var transformReflectValue = cmp.Transformer("", func(v protoreflect.Value) any {
	switch v := v.Interface().(type) {
	case protoreflect.Message:
		return v.Interface()
	case protoreflect.Map:
		ms := map[any]protoreflect.Value{}
		v.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			ms[k.Interface()] = v
			return true
		})
		return ms
	case protoreflect.List:
		ls := make([]protoreflect.Value, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			ls = append(ls, v.Get(i))
		}
		return ls
	default:
		return v
	}
})

func Test_conditionToCmpFunc(t *testing.T) {
	// helper for making StringList values
	sl := func(s ...string) *gen.Device_Query_StringList {
		return &gen.Device_Query_StringList{Strings: s}
	}

	t.Run("nil", func(t *testing.T) {
		cmpFunc := conditionToCmpFunc(&gen.Device_Query_Condition{})
		if cmpFunc == nil {
			t.Errorf("expected nil condition value to return non-nil cmpFunc, got %T", cmpFunc)
		}
		got := cmpFunc(value{})
		if got {
			t.Errorf("expected nil condition value to return false for any value, got true")
		}
	})

	condTestName := func(cond *gen.Device_Query_Condition) string {
		name := fmt.Sprintf("%T", cond.Value)
		name = strings.TrimPrefix(name, "*gen.Device_Query_Condition_")
		return name
	}

	t.Run("strings", func(t *testing.T) {
		tests := []struct {
			cond               *gen.Device_Query_Condition
			positive, negative []string
		}{
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_StringEqual{StringEqual: "foo"}},
				positive: []string{"foo"},
				negative: []string{"", "bar", "FOO", "foO", "fooo", "-foo", "10"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_StringEqualFold{StringEqualFold: "fOo"}},
				positive: []string{"fOo", "FOO", "foo"},
				negative: []string{"", "bar", "fOoO", "-fOo", "10"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_StringContains{StringContains: "apple"}},
				positive: []string{"apple pie", "apple", "pineapple", "apple123"},
				negative: []string{"", "pple", "appl", "Apple"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_StringContainsFold{StringContainsFold: "aPPle"}},
				positive: []string{"aPPle pie", "apple", "pineappLE", "aPPle123"},
				negative: []string{"", "pple", "appl", "banana"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_StringIn{StringIn: sl("foo", "BAR")}},
				positive: []string{"foo", "BAR"},
				negative: []string{"", "FOO", "bar", "foO", "fooo", "-foo", "10"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_StringInFold{StringInFold: sl("foo", "BAR")}},
				positive: []string{"foo", "BAR", "FOO", "bar", "foO"},
				negative: []string{"", "fooo", "-foo", "10"},
			},
		}

		strLeaf := func(val string) value {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return value{
				fd: md.Fields().ByName("string_val"),
				v:  protoreflect.ValueOfString(val),
			}
		}

		for _, tt := range tests {
			t.Run(condTestName(tt.cond), func(t *testing.T) {
				cmpFunc := conditionToCmpFunc(tt.cond)
				for _, str := range tt.positive {
					if !cmpFunc(strLeaf(str)) {
						t.Errorf("expected %q to match condition %s", str, tt.cond)
					}
				}
				for _, str := range tt.negative {
					if cmpFunc(strLeaf(str)) {
						t.Errorf("expected %q to not match condition %s", str, tt.cond)
					}
				}
			})
		}
	})

	t.Run("timestamps", func(t *testing.T) {
		now := time.Date(2025, 8, 10, 12, 5, 40, 111, time.UTC)
		var zero time.Time // zero time is used to test empty timestamps
		tests := []struct {
			cond               *gen.Device_Query_Condition
			positive, negative []time.Time
		}{
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampEqual{TimestampEqual: timestamppb.New(now)}},
				positive: []time.Time{now},
				negative: []time.Time{zero, now.Add(1), now.Add(-1)},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampLt{TimestampLt: timestamppb.New(now)}},
				positive: []time.Time{now.Add(-1), now.Add(-10 * time.Second)},
				negative: []time.Time{zero, now, now.Add(1), now.Add(10 * time.Second)},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampLte{TimestampLte: timestamppb.New(now)}},
				positive: []time.Time{now, now.Add(-1), now.Add(-10 * time.Second)},
				negative: []time.Time{zero, now.Add(1), now.Add(10 * time.Second)},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampGt{TimestampGt: timestamppb.New(now)}},
				positive: []time.Time{now.Add(1), now.Add(10 * time.Second)},
				negative: []time.Time{zero, now, now.Add(-1), now.Add(-10 * time.Second)},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampGte{TimestampGte: timestamppb.New(now)}},
				positive: []time.Time{now, now.Add(1), now.Add(10 * time.Second)},
				negative: []time.Time{zero, now.Add(-1), now.Add(-10 * time.Second)},
			},
		}

		timestampLeaf := func(val time.Time) value {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return value{
				fd: md.Fields().ByName("timestamp_val"),
				v:  protoreflect.ValueOfMessage(timestamppb.New(val).ProtoReflect()),
			}
		}

		t.Run("not a timestamp", func(t *testing.T) {
			l := value{
				fd: (&querypb.Result{}).ProtoReflect().Descriptor().Fields().ByName("string_val"),
				v:  protoreflect.ValueOfString("not a timestamp"),
			}
			cmpFunc := conditionToCmpFunc(&gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampEqual{TimestampEqual: timestamppb.New(now)}})
			if cmpFunc(l) {
				t.Errorf("expected condition to not match non-timestamp value, got true")
			}
		})

		for _, tt := range tests {
			t.Run(condTestName(tt.cond), func(t *testing.T) {
				cmpFunc := conditionToCmpFunc(tt.cond)
				for _, ts := range tt.positive {
					if !cmpFunc(timestampLeaf(ts)) {
						t.Errorf("expected %v to match condition %s", ts, tt.cond)
					}
				}
				for _, ts := range tt.negative {
					if cmpFunc(timestampLeaf(ts)) {
						t.Errorf("expected %v to not match condition %s", ts, tt.cond)
					}
				}
			})
		}
	})

	t.Run("names", func(t *testing.T) {
		tests := []struct {
			cond     *gen.Device_Query_Condition
			positive []string
			negative []string
		}{
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_NameDescendant{NameDescendant: "a/b"}},
				positive: []string{"a/b/c", "a/b/c/d"},
				negative: []string{"", "a", "a/", "a/b", "a/bc", "a/b/", "x/b/c"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_NameDescendantInc{NameDescendantInc: "a/b"}},
				positive: []string{"a/b", "a/b/c", "a/b/c/d"},
				negative: []string{"", "a", "a/", "a/bc", "a/b/", "x/b/c"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_NameDescendantIn{NameDescendantIn: sl("a/b", "1/2")}},
				positive: []string{"a/b/c", "a/b/c/d", "1/2/3", "1/2/3/4"},
				negative: []string{"", "a", "a/", "a/b", "a/bc", "a/b/", "x/b/c", "1", "1/2"},
			},
			{
				cond:     &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_NameDescendantIncIn{NameDescendantIncIn: sl("a/b", "1/2")}},
				positive: []string{"a/b", "a/b/c", "a/b/c/d", "1/2", "1/2/3", "1/2/3/4"},
				negative: []string{"", "a", "a/", "a/bc", "a/b/", "x/b/c", "1", "2"},
			},
		}

		nameLeaf := func(val string) value {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return value{
				fd: md.Fields().ByName("string_val"),
				v:  protoreflect.ValueOfString(val),
			}
		}

		for _, tt := range tests {
			t.Run(condTestName(tt.cond), func(t *testing.T) {
				cmpFunc := conditionToCmpFunc(tt.cond)
				for _, str := range tt.positive {
					if !cmpFunc(nameLeaf(str)) {
						t.Errorf("expected %q to match condition %s", str, tt.cond)
					}
				}
				for _, str := range tt.negative {
					if cmpFunc(nameLeaf(str)) {
						t.Errorf("expected %q to not match condition %s", str, tt.cond)
					}
				}
			})
		}
	})

	t.Run("presence", func(t *testing.T) {
		mkValue := func(field string, val protoreflect.Value) value {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return value{
				fd: md.Fields().ByName(protoreflect.Name(field)),
				v:  val,
			}
		}
		tests := []struct {
			cond     *gen.Device_Query_Condition
			positive []value
			negative []value
		}{
			{
				cond: &gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_Present{Present: &emptypb.Empty{}}},
				positive: []value{
					mkValue("string_val", protoreflect.ValueOfString("a")),
					mkValue("r_string", protoreflect.ValueOfString("b")),
					mkValue("enum_val", protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_A.Number())),
				},
				negative: []value{
					mkValue("string_val", protoreflect.Value{}),
					mkValue("r_string", protoreflect.Value{}),
					mkValue("enum_val", protoreflect.Value{}),
				},
			},
		}
		for _, tt := range tests {
			t.Run(condTestName(tt.cond), func(t *testing.T) {
				cmpFunc := conditionToCmpFunc(tt.cond)
				for _, l := range tt.positive {
					if !cmpFunc(l) {
						t.Errorf("expected %v to match condition %s", l, tt.cond)
					}
				}
				for _, l := range tt.negative {
					if cmpFunc(l) {
						t.Errorf("expected %v to not match condition %s", l, tt.cond)
					}
				}
			})
		}
	})
}
