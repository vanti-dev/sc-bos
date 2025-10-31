package devices

import (
	"bytes"
	"fmt"
	"iter"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// deviceMatchesQuery returns true if all the conditions in query match fields of device.
func deviceMatchesQuery(query *gen.Device_Query, device *gen.Device) bool {
	if query == nil {
		return true
	}
	for _, condition := range query.Conditions {
		if !conditionMatches(condition, device) {
			return false
		}
	}

	// this means a query with no conditions always returns true
	return true
}

// conditionMatches returns true if the condition matches values in device.
// If the condition has a path, only values matching the path are considered,
// otherwise all leafs in device are considered.
func conditionMatches(cond *gen.Device_Query_Condition, msg proto.Message) bool {
	cmp := conditionToCmpFunc(cond)
	values, err := rangeValues(cond.Field, msg)
	if err != nil {
		return false
	}
	for v := range values {
		if cmp(v) {
			return true
		}
	}
	return false
}

// conditionToCmpFunc converts a Device_Query_Condition into a function that checks if value values match the condition.
func conditionToCmpFunc(cond *gen.Device_Query_Condition) func(value) bool {
	strCmp := func(f func(string) bool) func(v value) bool {
		return func(v value) bool {
			s, ok := v.toString()
			if !ok {
				return false
			}
			return f(s)
		}
	}
	timestampCmp := func(f func(*timestamppb.Timestamp) bool) func(value) bool {
		return func(v value) bool {
			t, ok := v.toTimestamp()
			if !ok {
				return false
			}
			return f(t)
		}
	}
	descendantCmp := func(f func(string) bool) func(v value) bool {
		return strCmp(func(s string) bool {
			if strings.HasSuffix(s, "/") {
				return false // trailing /'s don't match, they'd result in an empty segment
			}
			return f(s)
		})
	}

	switch c := cond.Value.(type) {
	case *gen.Device_Query_Condition_StringEqual:
		return strCmp(func(v string) bool {
			return v == c.StringEqual
		})
	case *gen.Device_Query_Condition_StringEqualFold:
		return strCmp(func(v string) bool {
			return strings.EqualFold(v, c.StringEqualFold)
		})
	case *gen.Device_Query_Condition_StringContains:
		return strCmp(func(v string) bool {
			return strings.Contains(v, c.StringContains)
		})
	case *gen.Device_Query_Condition_StringContainsFold:
		ls := strings.ToLower(c.StringContainsFold)
		return strCmp(func(v string) bool {
			return strings.Contains(strings.ToLower(v), ls)
		})
	case *gen.Device_Query_Condition_StringIn:
		set := make(map[string]struct{}, len(c.StringIn.Strings))
		for _, s := range c.StringIn.Strings {
			set[s] = struct{}{}
		}
		return strCmp(func(v string) bool {
			_, ok := set[v]
			return ok
		})
	case *gen.Device_Query_Condition_StringInFold:
		set := make(map[string]struct{}, len(c.StringInFold.Strings))
		for _, s := range c.StringInFold.Strings {
			set[strings.ToLower(s)] = struct{}{}
		}
		return strCmp(func(v string) bool {
			_, ok := set[strings.ToLower(v)]
			return ok
		})

	case *gen.Device_Query_Condition_TimestampEqual:
		return timestampCmp(func(t *timestamppb.Timestamp) bool {
			return t.AsTime().Equal(c.TimestampEqual.AsTime())
		})
	case *gen.Device_Query_Condition_TimestampGt:
		return timestampCmp(func(t *timestamppb.Timestamp) bool {
			return t.AsTime().After(c.TimestampGt.AsTime())
		})
	case *gen.Device_Query_Condition_TimestampGte:
		return timestampCmp(func(t *timestamppb.Timestamp) bool {
			return !t.AsTime().Before(c.TimestampGte.AsTime())
		})
	case *gen.Device_Query_Condition_TimestampLt:
		return timestampCmp(func(t *timestamppb.Timestamp) bool {
			// zero times shouldn't match
			goT := t.AsTime()
			if goT.IsZero() {
				return false
			}
			return goT.Before(c.TimestampLt.AsTime())
		})
	case *gen.Device_Query_Condition_TimestampLte:
		return timestampCmp(func(t *timestamppb.Timestamp) bool {
			// zero times shouldn't match
			goT := t.AsTime()
			if goT.IsZero() {
				return false
			}
			return !goT.After(c.TimestampLte.AsTime())
		})

	case *gen.Device_Query_Condition_NameDescendant:
		return descendantCmp(func(v string) bool {
			return strings.HasPrefix(v, c.NameDescendant+"/")
		})
	case *gen.Device_Query_Condition_NameDescendantInc:
		return descendantCmp(func(v string) bool {
			return v == c.NameDescendantInc || strings.HasPrefix(v, c.NameDescendantInc+"/")
		})
	case *gen.Device_Query_Condition_NameDescendantIn:
		tree := newTreeFromPaths(c.NameDescendantIn.Strings...)
		return descendantCmp(func(v string) bool {
			n, self := tree.matchDescendant(v)
			return n != nil && !self
		})
	case *gen.Device_Query_Condition_NameDescendantIncIn:
		tree := newTreeFromPaths(c.NameDescendantIncIn.Strings...)
		return descendantCmp(func(v string) bool {
			n, _ := tree.matchDescendant(v)
			return n != nil
		})

	case *gen.Device_Query_Condition_Present:
		return func(v value) bool {
			return v.v.IsValid()
		}

	case *gen.Device_Query_Condition_Matches:
		return func(v value) bool {
			if v.fd.Message() == nil {
				return false
			}
			msg := v.v.Message()
			if !msg.IsValid() {
				return false
			}
			for _, condition := range c.Matches.GetConditions() {
				if !conditionMatches(condition, msg.Interface()) {
					return false
				}
			}
			return true // implies any present 'matches' field without conditions will match
		}
	}

	return func(v value) bool {
		return false // no condition matches, return false
	}
}

// getMessageString returns an iterator over the string values identified by path from msg.
// See valueString for details of string conversion.
func getMessageString(path string, msg proto.Message) iter.Seq[string] {
	return func(yield func(string) bool) {
		values, err := rangeValues(path, msg)
		if err != nil {
			return
		}
		for v := range values {
			str, got := v.toString()
			if !got || str == "" {
				continue // not a string or empty
			}
			if !yield(str) {
				return // stop iterating
			}
		}
	}
}

// value is a searchable value in a proto message.
// A searchable value can be matched against a query condition.
type value struct {
	fd protoreflect.FieldDescriptor
	v  protoreflect.Value
}

// rangeValuesOptions defines options for the rangeValues function.
// Useful for testing where stable order is important, search matching doesn't need to be stable.
type rangeValuesOptions struct {
	// If true, the iterator will yield values in the same order as a stable protorange.Range.
	Stable bool
}

var emptyValues = func(yield func(value) bool) {}

func (opts rangeValuesOptions) Range(path string, msg proto.Message) (iter.Seq[value], error) {
	if msg == nil {
		return emptyValues, nil
	}

	if path == "" {
		rangeOpts := protorange.Options{
			Stable: opts.Stable,
		}
		return func(yield func(value) bool) {
			// we never return an error from the Range, so we can ignore the error return value
			_ = rangeOpts.Range(msg.ProtoReflect(), func(values protopath.Values) error {
				fd, isLeaf := getLeafDescriptor(values)
				if !isLeaf {
					return nil
				}
				if !yield(value{fd, values.Index(-1).Value}) {
					return protorange.Terminate
				}
				return nil
			}, nil)
		}, nil
	}

	segments, err := parsePath(path)
	if err != nil {
		return nil, fmt.Errorf("path: %w", err)
	}
	return opts.scanMessage(segments, msg.ProtoReflect()), nil
}

// rangeValues returns an iterator over values of msg matching path.
// If path is empty, it returns all leaf values in msg,
// see isLeaf for details on what a leaf value is.
// An error is returned if the path is invalid.
func rangeValues(path string, msg proto.Message) (iter.Seq[value], error) {
	return rangeValuesOptions{}.Range(path, msg)
}

// scanMessage returns all values identified by path in msg.
// The first entry of path should be a field name in msg, otherwise an empty iterator is returned.
func (opts rangeValuesOptions) scanMessage(path []pathSegment, msg protoreflect.Message) iter.Seq[value] {
	if len(path) == 0 {
		return emptyValues // callers should handle empty paths
	}
	head := path[0]
	if head.IsIndex {
		return emptyValues // path doesn't match the message structure
	}

	fd := msg.Descriptor().Fields().ByName(protoreflect.Name(head.Name))
	if fd == nil {
		return emptyValues // field doesn't exist in the message
	}

	v := msg.Get(fd)
	if v.Equal(fd.Default()) {
		return emptyValues // field exists but has no value
	}

	switch {
	case fd.IsMap():
		return opts.scanMap(path[1:], fd, v.Map())
	case fd.IsList():
		return opts.scanList(path[1:], fd, v.List())
	case fd.Message() != nil:
		vm := v.Message()
		if !vm.IsValid() {
			return emptyValues // field exists but has no value
		}
		if isLeaf(fd) && len(path) > 1 {
			return emptyValues // path has more segments, but the value is a value
		}
		if len(path) == 1 {
			// end of the path, return the value
			return func(yield func(value) bool) {
				yield(value{fd, v})
			}
		}
		return opts.scanMessage(path[1:], vm)
	case isLeaf(fd) && len(path) == 1:
		return func(yield func(value) bool) {
			yield(value{fd, v})
		}
	default:
		return emptyValues // the field value is of an unknown type

	}
}

// scanMap returns all leafs identified by path in the map m.
func (opts rangeValuesOptions) scanMap(path []pathSegment, fd protoreflect.FieldDescriptor, m protoreflect.Map) iter.Seq[value] {
	if len(path) == 0 {
		return func(yield func(value) bool) {
			if !opts.Stable {
				m.Range(func(key protoreflect.MapKey, v protoreflect.Value) bool {
					if !yield(value{fd.MapValue(), v}) {
						return false
					}
					return true
				})
			}

			// need stable order, capture and sort the entries
			type entry struct {
				k protoreflect.MapKey
				v protoreflect.Value
			}
			entries := make([]entry, 0, m.Len())
			m.Range(func(key protoreflect.MapKey, v protoreflect.Value) bool {
				entries = append(entries, entry{key, v})
				return true
			})
			// same logic that exists in the protorange package.
			// sorts false before true, numeric keys in ascending order,
			// and strings in lexicographical ordering according to UTF-8 codepoints.
			sort.Slice(entries, func(xe, ye int) bool {
				x := entries[xe].k.Value()
				y := entries[ye].k.Value()
				switch x.Interface().(type) {
				case bool:
					return !x.Bool() && y.Bool()
				case int32, int64:
					return x.Int() < y.Int()
				case uint32, uint64:
					return x.Uint() < y.Uint()
				case string:
					return x.String() < y.String()
				default:
					panic("invalid map key type")
				}
			})
			for _, e := range entries {
				if !yield(value{fd.MapValue(), e.v}) {
					return
				}
			}
		}
	}

	head := path[0]
	if head.IsIndex {
		return emptyValues // path expects a list, got a map
	}

	k, ok := parseMapKey(head.Name, fd.MapKey())
	if !ok {
		return emptyValues // invalid map key
	}
	v := m.Get(k)
	if !v.IsValid() || v.Equal(fd.MapValue().Default()) {
		return emptyValues // key doesn't exist in the map
	}

	if len(path) == 1 {
		// end of the path
		return func(yield func(value) bool) {
			yield(value{fd.MapValue(), v})
		}
	}

	// more to the path
	if isLeaf(fd.MapValue()) {
		return emptyValues // path has more segments, but the map value is a value
	}
	if md := fd.MapValue().Message(); md != nil {
		return opts.scanMessage(path[1:], v.Message())
	}
	return emptyValues // there's more to the path but the value has no properties
}

// scanList returns all leafs identified by path in the list l.
func (opts rangeValuesOptions) scanList(path []pathSegment, fd protoreflect.FieldDescriptor, l protoreflect.List) iter.Seq[value] {
	if len(path) == 0 {
		return func(yield func(value) bool) {
			for i := 0; i < l.Len(); i++ {
				if !yield(value{fd, l.Get(i)}) {
					return
				}
			}
		}
	}

	head := path[0]
	if head.IsIndex {
		i := head.Index
		if i < 0 {
			// negative index means counting from the end
			i = l.Len() + i
		}
		if i < 0 || i >= l.Len() {
			return emptyValues // index out of range
		}
		v := l.Get(i)
		if isLeaf(fd) && len(path) > 1 {
			return emptyValues // path has more segments, but the list contains values
		}
		if len(path) == 1 {
			// end of the path, return the value
			return func(yield func(value) bool) {
				yield(value{fd, v})
			}
		}

		if fd.Message() == nil {
			return emptyValues // it's not a value, but also not a message. Not sure what it is.
		}

		return opts.scanMessage(path[1:], v.Message())
	}

	// there's more to the path, but this is a list of leafs which can't have more path segments
	if isLeaf(fd) {
		return emptyValues // path has more segments, but the list contains value values
	}

	return func(yield func(value) bool) {
		for i := 0; i < l.Len(); i++ {
			v := l.Get(i)
			for leaf := range opts.scanMessage(path, v.Message()) {
				if !yield(leaf) {
					return
				}
			}
		}
	}
}

// getLeafDescriptor returns the most relevant FieldDescriptor for values, and whether values is a value or not.
func getLeafDescriptor(values protopath.Values) (protoreflect.FieldDescriptor, bool) {
	// getFd returns the most useful FieldDescriptor for values at index i.
	getFd := func(i int) protoreflect.FieldDescriptor {
		step := values.Index(i).Step
		switch step.Kind() {
		case protopath.FieldAccessStep:
			return step.FieldDescriptor()
		case protopath.MapIndexStep:
			return values.Index(i - 1).Step.FieldDescriptor().MapValue()
		case protopath.ListIndexStep:
			return values.Index(i - 1).Step.FieldDescriptor()
		default:
			return nil // unsupported step kind
		}
	}

	fd := getFd(-1)
	if fd == nil {
		return nil, false
	}

	switch values.Index(-1).Step.Kind() {
	case protopath.FieldAccessStep:
		// check the parent isn't a value
		if parentFd := getFd(-2); parentFd != nil && isLeaf(parentFd) {
			return nil, false
		}
		fd := getFd(-1)
		return fd, isLeaf(fd) && !fd.IsList()
	case protopath.MapIndexStep:
		fd := getFd(-1)
		return fd, isLeaf(fd)
	case protopath.ListIndexStep:
		fd := getFd(-1)
		return fd, isLeaf(fd)
	default:
		return nil, false
	}
}

// isLeaf returns true if fd is a value.
// A value is a scalar value, or a well known value type.
func isLeaf(fd protoreflect.FieldDescriptor) bool {
	if fd == nil {
		return false
	}
	if msg := fd.Message(); msg != nil {
		switch msg.FullName() {
		case "google.protobuf.Timestamp",
			"google.protobuf.Duration":
			return true // treat special messages as leafs
		default:
			return false
		}
	}
	return true // all other types are scalars, aka leafs
}

// toString converts the value into a string ready for comparison to another string.
// Unlike l.v.String() this converts enum values to their enum name where available,
// otherwise converts them to a string representation of the enum number.
// Bytes are converted to string.
func (l value) toString() (string, bool) {
	fd, v := l.fd, l.v
	if fd == nil {
		return "", false
	}
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return strconv.FormatBool(v.Bool()), true
	case protoreflect.FloatKind:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), true
	case protoreflect.DoubleKind:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), true
	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(v.Int(), 10), true
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(v.Uint(), 10), true
	case protoreflect.StringKind:
		return v.String(), true
	case protoreflect.EnumKind:
		enum := v.Enum()
		enumDesc := fd.Enum().Values().ByNumber(enum)
		if enumDesc == nil {
			// unknown enum
			return strconv.FormatInt(int64(enum), 10), true
		}
		return string(enumDesc.Name()), true
	case protoreflect.BytesKind:
		return string(v.Bytes()), true
	default:
		// for value messages, we rely on the json representation to turn them into strings.
		if msg := v.Message(); msg != nil {
			bs, err := protojson.Marshal(msg.Interface())
			if err != nil {
				// if we can't marshal the message, we can't convert it to a string
				return "", false
			}
			bs = bytes.Trim(bs, `"`) // remove quotes around the string
			return string(bs), true
		}
		// unsupported kinds
		return "", false
	}
}

// toTimestamp converts the value into a *timestamppb.Timestamp if it is a valid timestamp.
func (l value) toTimestamp() (*timestamppb.Timestamp, bool) {
	if l.fd == nil {
		return nil, false
	}
	if l.fd.Message() == nil || l.fd.Message().FullName() != "google.protobuf.Timestamp" {
		return nil, false
	}
	tm := l.v.Message()
	if !tm.IsValid() {
		return nil, false
	}
	t, ok := tm.Interface().(*timestamppb.Timestamp)
	if !ok {
		return nil, false
	}

	return t, true
}

// parseMapKey converts keyStr into a protoreflect.MapKey using fd to choose the correct conversion method to use.
func parseMapKey(keyStr string, fd protoreflect.FieldDescriptor) (protoreflect.MapKey, bool) {
	fail := func() (protoreflect.MapKey, bool) {
		return protoreflect.MapKey{}, false
	}
	switch fd.Kind() {
	case protoreflect.BoolKind:
		gt, err := strconv.ParseBool(keyStr)
		if err != nil {
			return fail()
		}
		return protoreflect.ValueOfBool(gt).MapKey(), true
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		gt, err := strconv.ParseInt(keyStr, 10, 32)
		if err != nil {
			return fail()
		}
		return protoreflect.ValueOfInt32(int32(gt)).MapKey(), true
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		gt, err := strconv.ParseInt(keyStr, 10, 64)
		if err != nil {
			return fail()
		}
		return protoreflect.ValueOfInt64(gt).MapKey(), true
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		gt, err := strconv.ParseUint(keyStr, 10, 32)
		if err != nil {
			return fail()
		}
		return protoreflect.ValueOfUint32(uint32(gt)).MapKey(), true
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		gt, err := strconv.ParseUint(keyStr, 10, 64)
		if err != nil {
			return fail()
		}
		return protoreflect.ValueOfUint64(gt).MapKey(), true
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(keyStr).MapKey(), true
	default:
		// unknown map key type!
		return fail()
	}
}
