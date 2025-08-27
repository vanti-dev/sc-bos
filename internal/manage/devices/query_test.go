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
		leaf   leaf
		want   string
		wantOk bool
	}{
		{leaf: leaf{fd("double_val"), protoreflect.ValueOfFloat64(1.1)}, want: "1.1", wantOk: true},
		{leaf: leaf{fd("float_val"), protoreflect.ValueOfFloat32(2.1)}, want: "2.1", wantOk: true},
		{leaf: leaf{fd("int32_val"), protoreflect.ValueOfInt32(-10)}, want: "-10", wantOk: true},
		{leaf: leaf{fd("int64_val"), protoreflect.ValueOfInt64(-20)}, want: "-20", wantOk: true},
		{leaf: leaf{fd("uint32_val"), protoreflect.ValueOfUint32(30)}, want: "30", wantOk: true},
		{leaf: leaf{fd("uint64_val"), protoreflect.ValueOfUint64(40)}, want: "40", wantOk: true},
		{leaf: leaf{fd("sint32_val"), protoreflect.ValueOfInt32(-50)}, want: "-50", wantOk: true},
		{leaf: leaf{fd("sint64_val"), protoreflect.ValueOfInt64(-60)}, want: "-60", wantOk: true},
		{leaf: leaf{fd("fixed32_val"), protoreflect.ValueOfUint32(70)}, want: "70", wantOk: true},
		{leaf: leaf{fd("fixed64_val"), protoreflect.ValueOfUint64(80)}, want: "80", wantOk: true},
		{leaf: leaf{fd("sfixed32_val"), protoreflect.ValueOfInt32(-90)}, want: "-90", wantOk: true},
		{leaf: leaf{fd("sfixed64_val"), protoreflect.ValueOfInt64(-100)}, want: "-100", wantOk: true},
		{leaf: leaf{fd("bool_val"), protoreflect.ValueOfBool(true)}, want: "true", wantOk: true},
		{leaf: leaf{fd("enum_val"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_A.Number())}, want: "RESULT_ENUM_A", wantOk: true},
		{leaf: leaf{fd("duration_val"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 10).ProtoReflect())}, want: "10s", wantOk: true},
		{leaf: leaf{fd("timestamp_val"), protoreflect.ValueOfMessage(timestamppb.New(time.Date(2025, 8, 10, 12, 5, 40, 111, time.UTC)).ProtoReflect())}, want: "2025-08-10T12:05:40.000000111Z", wantOk: true},
	}
	for _, tt := range tests {
		t.Run(string(tt.leaf.fd.Name()), func(t *testing.T) {
			gotS, gotOk := tt.leaf.toString()
			if gotS != tt.want || gotOk != tt.wantOk {
				t.Errorf("leaf.toString() = %q, %v; want %q, %v", gotS, gotOk, tt.want, tt.wantOk)
			}
		})
	}
}

func Test_rangeLeafs(t *testing.T) {
	// m is the data structure we test against
	m := testData()

	// paths are all leaf paths in m
	paths := resultPaths("", nil)
	leafs := resultLeafs(nil)
	paths = nestedResultPaths("", paths)
	leafs = nestedResultLeafs(leafs)
	for _, nestingPath := range nestedMessagePaths {
		paths = nestedResultPaths(nestingPath, paths)
		leafs = nestedResultLeafs(leafs)
	}
	if len(leafs) != len(paths) {
		t.Fatalf("expected leafs and paths to have the same length, got %d and %d", len(leafs), len(paths))
	}

	type testCase struct {
		path      string
		wantLeafs []leaf
		wantErr   bool
	}
	// newTestCase is a helper for making test cases.
	newTestCase := func(path string, wantIndexes ...int) testCase {
		t.Helper()
		if len(wantIndexes) == 0 {
			t.Fatalf("wantIndexes must not be empty for path %q", path)
		}
		wantLeafs := make([]leaf, len(wantIndexes))
		for i, idx := range wantIndexes {
			if idx < 0 || idx >= len(leafs) {
				t.Fatalf("index %d out of range for path %q", idx, path)
			}
			wantLeafs[i] = leafs[idx]
		}
		return testCase{
			path:      path,
			wantLeafs: wantLeafs,
		}
	}

	tests := []testCase{
		// no path means all paths
		{wantLeafs: leafs},
	}

	// each path should resolve to exactly the wanted value
	for i, path := range paths {
		tests = append(tests, newTestCase(path, i))
	}

	// repeated values should be optional in the path
	for i, scalarType := range scalarTypes {
		// each scalar type ahs the non-repeated value followed by two repeated values
		leafIdx := i*3 + 1
		tests = append(tests, newTestCase(fmt.Sprintf("r_%s", scalarType), leafIdx, leafIdx+1))
	}

	// negative indexes should work for repeated values
	for i, scalarType := range scalarTypes {
		// each scalar type ahs the non-repeated value followed by two repeated values
		leafIdx := i*3 + 1
		tests = append(tests, newTestCase(fmt.Sprintf("r_%s[-1]", scalarType), leafIdx+1))
		tests = append(tests, newTestCase(fmt.Sprintf("r_%s[-2]", scalarType), leafIdx))
	}

	// path doesn't match any value, return empty iterator
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
		// path isn't a leaf (or repeated leaf)
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
		// fields of leaf messages
		"timestamp_val.seconds",
		"r_timestamp[0].seconds",
		"m_string_timestamp.a.seconds",
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
			got, err := rangeLeafsOptions{Stable: true}.Range(tt.path, m)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("rangeLeafs() error = %v, wantErr %v. Got %v", err, tt.wantErr, slices.Collect(got))
				}
				return
			}
			var wantFds []protoreflect.FieldDescriptor
			var wantValues []protoreflect.Value
			for _, l := range tt.wantLeafs {
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
				t.Errorf("rangeLeafs() field descriptors (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(wantValues, gotValues, transformReflectValue, protocmp.Transform()); diff != "" {
				t.Errorf("rangeLeafs() values (-want,+got):\n%s", diff)
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

// testData returns a *querypb.Result with all fields populated.
// Nested structures are populated to a depth of 3: root.results.result.string_val == "a".
// Each nesting has a non-nil Result with all default values: r.r_result[2].string_val == "", r.m_string_result["d"].string_val == "".
func testData() *querypb.Result {
	// newResult makes a fully populated querypb.Result
	newResult := func() *querypb.Result {
		return &querypb.Result{
			DoubleVal:        1.1,
			RDouble:          []float64{1.2, 1.3},
			FloatVal:         2.1,
			RFloat:           []float32{2.2, 2.3},
			Int32Val:         -10,
			RInt32:           []int32{-11, -12},
			Int64Val:         -20,
			RInt64:           []int64{-21, -22},
			Uint32Val:        30,
			RUint32:          []uint32{31, 32},
			Uint64Val:        40,
			RUint64:          []uint64{41, 42},
			Sint32Val:        -50,
			RSint32:          []int32{-51, -52},
			Sint64Val:        -60,
			RSint64:          []int64{-61, -62},
			Fixed32Val:       70,
			RFixed32:         []uint32{71, 72},
			Fixed64Val:       80,
			RFixed64:         []uint64{81, 82},
			Sfixed32Val:      -90,
			RSfixed32:        []int32{-91, -92},
			Sfixed64Val:      -100,
			RSfixed64:        []int64{-101, -102},
			BoolVal:          true,
			RBool:            []bool{false, true},
			StringVal:        "a",
			RString:          []string{"b", "c"},
			BytesVal:         []byte("abc"),
			RBytes:           [][]byte{[]byte("def"), []byte("ghi")},
			EnumVal:          querypb.ResultEnum_RESULT_ENUM_A,
			REnum:            []querypb.ResultEnum{querypb.ResultEnum_RESULT_ENUM_B, querypb.ResultEnum_RESULT_ENUM_C},
			DurationVal:      durationpb.New(time.Second * 10),
			RDuration:        []*durationpb.Duration{durationpb.New(time.Second * 20), durationpb.New(time.Second * 30)},
			TimestampVal:     timestamppb.New(time.Unix(100, 10)),
			RTimestamp:       []*timestamppb.Timestamp{timestamppb.New(time.Unix(200, 20)), timestamppb.New(time.Unix(300, 30))},
			MInt32String:     map[int32]string{-11: "a", -10: "b"},
			MInt64String:     map[int64]string{-21: "a", -20: "b"},
			MUint32String:    map[uint32]string{30: "a", 31: "b"},
			MUint64String:    map[uint64]string{40: "a", 41: "b"},
			MSint32String:    map[int32]string{-51: "a", -50: "b"},
			MSint64String:    map[int64]string{-61: "a", -60: "b"},
			MFixed32String:   map[uint32]string{70: "a", 71: "b"},
			MFixed64String:   map[uint64]string{80: "a", 81: "b"},
			MSfixed32String:  map[int32]string{-91: "a", -90: "b"},
			MSfixed64String:  map[int64]string{-101: "a", -100: "b"},
			MBoolString:      map[bool]string{false: "a", true: "b"},
			MStringString:    map[string]string{"a": "a", "b": "b"},
			MStringEnum:      map[string]querypb.ResultEnum{"a": querypb.ResultEnum_RESULT_ENUM_A, "b": querypb.ResultEnum_RESULT_ENUM_B},
			MStringDuration:  map[string]*durationpb.Duration{"a": durationpb.New(time.Second * 10), "b": durationpb.New(time.Second * 20)},
			MStringTimestamp: map[string]*timestamppb.Timestamp{"a": timestamppb.New(time.Unix(100, 10)), "b": timestamppb.New(time.Unix(200, 20))},
		}
	}
	fillNestedResult := func(r *querypb.Result) {
		r.Result = newResult()
		r.RResult = []*querypb.Result{newResult(), newResult()}
		r.MStringResult = map[string]*querypb.Result{
			"a": newResult(),
			"b": newResult(),
		}
	}

	res := newResult()
	fillNestedResult(res)
	fillNestedResult(res.Result)
	for _, result := range res.RResult {
		fillNestedResult(result)
	}
	for _, result := range res.MStringResult {
		fillNestedResult(result)
	}
	res.RResult = append(res.RResult, new(querypb.Result)) // r_result[2] is empty, so it has no values
	res.MStringResult["c"] = new(querypb.Result)           // c is empty, so it has no values
	return res
}

var scalarTypes = []string{
	"double",
	"float",
	"int32",
	"int64",
	"uint32",
	"uint64",
	"sint32",
	"sint64",
	"fixed32",
	"fixed64",
	"sfixed32",
	"sfixed64",
	"bool",
	"string",
	"bytes",
	"enum",
	"duration",
	"timestamp",
}

// resultPaths adds to paths all populated paths in the testData Result.
func resultPaths(prefix string, paths []string) []string {
	mapKeyTypes := [][]string{
		{"int32", "-11", "-10"},
		{"int64", "-21", "-20"},
		{"uint32", "30", "31"},
		{"uint64", "40", "41"},
		{"sint32", "-51", "-50"},
		{"sint64", "-61", "-60"},
		{"fixed32", "70", "71"},
		{"fixed64", "80", "81"},
		{"sfixed32", "-91", "-90"},
		{"sfixed64", "-101", "-100"},
		{"bool", "false", "true"},
		{"string", "a", "b"},
	}

	prefixStart := len(paths)
	for _, t := range scalarTypes {
		paths = append(paths, fmt.Sprintf("%s_val", t))
		paths = append(paths, fmt.Sprintf("r_%s[0]", t), fmt.Sprintf("r_%s[1]", t))
	}
	for _, t := range mapKeyTypes {
		ts := t[0]
		for _, v := range t[1:] {
			paths = append(paths, fmt.Sprintf("m_%s_string.%s", ts, v))
		}
	}
	paths = append(paths, "m_string_enum.a", "m_string_enum.b")
	paths = append(paths, "m_string_duration.a", "m_string_duration.b")
	paths = append(paths, "m_string_timestamp.a", "m_string_timestamp.b")

	// prepend the prefix
	needsPrefix := paths[prefixStart:]
	for i, path := range needsPrefix {
		needsPrefix[i] = prefix + path
	}
	return paths
}

var nestedMessagePaths = []string{
	"result.",
	"r_result[0].",
	"r_result[1].",
	"m_string_result.a.",
	"m_string_result.b.",
}

func nestedResultPaths(prefix string, paths []string) []string {
	for _, p := range nestedMessagePaths {
		paths = resultPaths(prefix+p, paths)
	}
	return paths
}

func resultLeafs(leafs []leaf) []leaf {
	valFd := func(field string) protoreflect.FieldDescriptor {
		md := (&querypb.Result{}).ProtoReflect().Descriptor()
		return md.Fields().ByName(protoreflect.Name(field))
	}
	mapFd := func(field string) protoreflect.FieldDescriptor {
		return valFd(field).MapValue()
	}

	return append(leafs, []leaf{
		{valFd("double_val"), protoreflect.ValueOfFloat64(1.1)},
		{valFd("r_double"), protoreflect.ValueOfFloat64(1.2)},
		{valFd("r_double"), protoreflect.ValueOfFloat64(1.3)},
		{valFd("float_val"), protoreflect.ValueOfFloat32(2.1)},
		{valFd("r_float"), protoreflect.ValueOfFloat32(2.2)},
		{valFd("r_float"), protoreflect.ValueOfFloat32(2.3)},
		{valFd("int32_val"), protoreflect.ValueOfInt32(-10)},
		{valFd("r_int32"), protoreflect.ValueOfInt32(-11)},
		{valFd("r_int32"), protoreflect.ValueOfInt32(-12)},
		{valFd("int64_val"), protoreflect.ValueOfInt64(-20)},
		{valFd("r_int64"), protoreflect.ValueOfInt64(-21)},
		{valFd("r_int64"), protoreflect.ValueOfInt64(-22)},
		{valFd("uint32_val"), protoreflect.ValueOfUint32(30)},
		{valFd("r_uint32"), protoreflect.ValueOfUint32(31)},
		{valFd("r_uint32"), protoreflect.ValueOfUint32(32)},
		{valFd("uint64_val"), protoreflect.ValueOfUint64(40)},
		{valFd("r_uint64"), protoreflect.ValueOfUint64(41)},
		{valFd("r_uint64"), protoreflect.ValueOfUint64(42)},
		{valFd("sint32_val"), protoreflect.ValueOfInt32(-50)},
		{valFd("r_sint32"), protoreflect.ValueOfInt32(-51)},
		{valFd("r_sint32"), protoreflect.ValueOfInt32(-52)},
		{valFd("sint64_val"), protoreflect.ValueOfInt64(-60)},
		{valFd("r_sint64"), protoreflect.ValueOfInt64(-61)},
		{valFd("r_sint64"), protoreflect.ValueOfInt64(-62)},
		{valFd("fixed32_val"), protoreflect.ValueOfUint32(70)},
		{valFd("r_fixed32"), protoreflect.ValueOfUint32(71)},
		{valFd("r_fixed32"), protoreflect.ValueOfUint32(72)},
		{valFd("fixed64_val"), protoreflect.ValueOfUint64(80)},
		{valFd("r_fixed64"), protoreflect.ValueOfUint64(81)},
		{valFd("r_fixed64"), protoreflect.ValueOfUint64(82)},
		{valFd("sfixed32_val"), protoreflect.ValueOfInt32(-90)},
		{valFd("r_sfixed32"), protoreflect.ValueOfInt32(-91)},
		{valFd("r_sfixed32"), protoreflect.ValueOfInt32(-92)},
		{valFd("sfixed64_val"), protoreflect.ValueOfInt64(-100)},
		{valFd("r_sfixed64"), protoreflect.ValueOfInt64(-101)},
		{valFd("r_sfixed64"), protoreflect.ValueOfInt64(-102)},
		{valFd("bool_val"), protoreflect.ValueOfBool(true)},
		{valFd("r_bool"), protoreflect.ValueOfBool(false)},
		{valFd("r_bool"), protoreflect.ValueOfBool(true)},
		{valFd("string_val"), protoreflect.ValueOfString("a")},
		{valFd("r_string"), protoreflect.ValueOfString("b")},
		{valFd("r_string"), protoreflect.ValueOfString("c")},
		{valFd("bytes_val"), protoreflect.ValueOfBytes([]byte("abc"))},
		{valFd("r_bytes"), protoreflect.ValueOfBytes([]byte("def"))},
		{valFd("r_bytes"), protoreflect.ValueOfBytes([]byte("ghi"))},
		{valFd("enum_val"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_A.Number())},
		{valFd("r_enum"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_B.Number())},
		{valFd("r_enum"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_C.Number())},
		{valFd("duration_val"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 10).ProtoReflect())},
		{valFd("r_duration"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 20).ProtoReflect())},
		{valFd("r_duration"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 30).ProtoReflect())},
		{valFd("timestamp_val"), protoreflect.ValueOfMessage(timestamppb.New(time.Unix(100, 10)).ProtoReflect())},
		{valFd("r_timestamp"), protoreflect.ValueOfMessage(timestamppb.New(time.Unix(200, 20)).ProtoReflect())},
		{valFd("r_timestamp"), protoreflect.ValueOfMessage(timestamppb.New(time.Unix(300, 30)).ProtoReflect())},
		{mapFd("m_int32_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_int32_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_int64_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_int64_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_uint32_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_uint32_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_uint64_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_uint64_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_sint32_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_sint32_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_sint64_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_sint64_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_fixed32_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_fixed32_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_fixed64_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_fixed64_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_sfixed32_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_sfixed32_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_sfixed64_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_sfixed64_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_bool_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_bool_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_string_string"), protoreflect.ValueOfString("a")},
		{mapFd("m_string_string"), protoreflect.ValueOfString("b")},
		{mapFd("m_string_enum"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_A.Number())},
		{mapFd("m_string_enum"), protoreflect.ValueOfEnum(querypb.ResultEnum_RESULT_ENUM_B.Number())},
		{mapFd("m_string_duration"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 10).ProtoReflect())},
		{mapFd("m_string_duration"), protoreflect.ValueOfMessage(durationpb.New(time.Second * 20).ProtoReflect())},
		{mapFd("m_string_timestamp"), protoreflect.ValueOfMessage(timestamppb.New(time.Unix(100, 10)).ProtoReflect())},
		{mapFd("m_string_timestamp"), protoreflect.ValueOfMessage(timestamppb.New(time.Unix(200, 20)).ProtoReflect())},
	}...)
}

func nestedResultLeafs(leafs []leaf) []leaf {
	// add all nested result leafs
	for range nestedMessagePaths {
		leafs = resultLeafs(leafs)
	}
	return leafs
}

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
		got := cmpFunc(leaf{})
		if got {
			t.Errorf("expected nil condition value to return false for any leaf, got true")
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

		strLeaf := func(val string) leaf {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return leaf{
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

		timestampLeaf := func(val time.Time) leaf {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return leaf{
				fd: md.Fields().ByName("timestamp_val"),
				v:  protoreflect.ValueOfMessage(timestamppb.New(val).ProtoReflect()),
			}
		}

		t.Run("not a timestamp", func(t *testing.T) {
			l := leaf{
				fd: (&querypb.Result{}).ProtoReflect().Descriptor().Fields().ByName("string_val"),
				v:  protoreflect.ValueOfString("not a timestamp"),
			}
			cmpFunc := conditionToCmpFunc(&gen.Device_Query_Condition{Value: &gen.Device_Query_Condition_TimestampEqual{TimestampEqual: timestamppb.New(now)}})
			if cmpFunc(l) {
				t.Errorf("expected condition to not match non-timestamp leaf, got true")
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

		nameLeaf := func(val string) leaf {
			md := (&querypb.Result{}).ProtoReflect().Descriptor()
			return leaf{
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
}
