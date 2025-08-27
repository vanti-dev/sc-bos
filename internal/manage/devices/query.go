package devices

import (
	"bytes"
	"fmt"
	"iter"
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

// deviceMatchesQuery returns true if all of the conditions in query match fields of device.
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

// conditionMatches returns true if the condition matches leaf values in device.
// If the condition has a path, only leafs matching the path are considered,
// otherwise all leafs in device are considered.
func conditionMatches(cond *gen.Device_Query_Condition, device *gen.Device) bool {
	cmp := conditionToCmpFunc(cond)
	leafs, err := rangeLeafs(cond.Field, device)
	if err != nil {
		return false
	}
	for v := range leafs {
		if cmp(v) {
			return true
		}
	}
	return false
}

// conditionToCmpFunc converts a Device_Query_Condition into a function that checks if leaf values match the condition.
func conditionToCmpFunc(cond *gen.Device_Query_Condition) func(leaf) bool {
	strCmp := func(f func(string) bool) func(v leaf) bool {
		return func(v leaf) bool {
			s, ok := v.toString()
			if !ok {
				return false
			}
			return f(s)
		}
	}
	timestampCmp := func(f func(*timestamppb.Timestamp) bool) func(leaf) bool {
		return func(v leaf) bool {
			t, ok := v.toTimestamp()
			if !ok {
				return false
			}
			return f(t)
		}
	}
	descendantCmp := func(f func(string) bool) func(v leaf) bool {
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
	}

	return func(v leaf) bool {
		return false // no condition matches, return false
	}
}

// getMessageString returns an iterator over the string values identified by path from msg.
// See valueString for details of string conversion.
func getMessageString(path string, msg proto.Message) iter.Seq[string] {
	return func(yield func(string) bool) {
		leafs, err := rangeLeafs(path, msg)
		if err != nil {
			return
		}
		for v := range leafs {
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

// leaf is a searchable value in a proto message.
// Searchable values are scalars and scalar-like messages including google.protobuf.Timestamp and google.protobuf.Duration.
type leaf struct {
	fd protoreflect.FieldDescriptor
	v  protoreflect.Value
}

// rangeLeafsOptions defines options for the rangeLeafs function.
// Useful for testing where stable order is important, search matching doesn't need to be stable.
type rangeLeafsOptions struct {
	// If true, the iterator will yield values in the same order as a stable protorange.Range.
	Stable bool
}

var emptyLeafs = func(yield func(leaf) bool) {}

func (opts rangeLeafsOptions) Range(path string, msg proto.Message) (iter.Seq[leaf], error) {
	if msg == nil {
		return emptyLeafs, nil
	}

	if path == "" {
		rangeOpts := protorange.Options{
			Stable: opts.Stable,
		}
		return func(yield func(leaf) bool) {
			// we never return an error from the Range, so we can ignore the error return value
			_ = rangeOpts.Range(msg.ProtoReflect(), func(values protopath.Values) error {
				fd, isLeaf := getLeafDescriptor(values)
				if !isLeaf {
					return nil
				}
				if !yield(leaf{fd, values.Index(-1).Value}) {
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
	return scanMessage(segments, msg.ProtoReflect()), nil
}

// rangeLeafs returns an iterator over leaf values of msg matching path.
// If path is empty, it returns all leaf values in msg.
// See isLeaf for details on what a leaf is.
// If path is not empty and resolves to a non-leaf value, it returns an empty iterator.
// An error is returned if the path is invalid.
func rangeLeafs(path string, msg proto.Message) (iter.Seq[leaf], error) {
	return rangeLeafsOptions{}.Range(path, msg)
}

// scanMessage returns all leafs identified by path in msg.
// The first entry of path should be a field name in msg, otherwise an empty iterator is returned.
func scanMessage(path []pathSegment, msg protoreflect.Message) iter.Seq[leaf] {
	if len(path) == 0 {
		return emptyLeafs // callers should handle empty paths
	}
	head := path[0]
	if head.IsIndex {
		return emptyLeafs // path doesn't match the message structure
	}

	fd := msg.Descriptor().Fields().ByName(protoreflect.Name(head.Name))
	if fd == nil {
		return emptyLeafs // field doesn't exist in the message
	}

	v := msg.Get(fd)
	if v.Equal(fd.Default()) {
		return emptyLeafs // field exists but has no value
	}

	switch {
	case fd.IsMap():
		return scanMap(path[1:], fd, v.Map())
	case fd.IsList():
		return scanList(path[1:], fd, v.List())
	case fd.Message() != nil:
		vm := v.Message()
		if !vm.IsValid() {
			return emptyLeafs // field exists but has no value
		}
		if isLeaf(fd) {
			if len(path) == 1 {
				// end of the path, return the leaf
				return func(yield func(leaf) bool) {
					yield(leaf{fd, v})
				}
			}
			// path has more segments, but the value is a leaf
			return emptyLeafs
		}
		return scanMessage(path[1:], vm)
	case isLeaf(fd) && len(path) == 1:
		return func(yield func(leaf) bool) {
			yield(leaf{fd, v})
		}
	default:
		return emptyLeafs // the field value is of an unknown type

	}
}

// scanMap returns all leafs identified by path in the map m.
func scanMap(path []pathSegment, fd protoreflect.FieldDescriptor, m protoreflect.Map) iter.Seq[leaf] {
	if len(path) == 0 {
		return emptyLeafs // no more path segments, nothing to return
	}

	head := path[0]
	if head.IsIndex {
		return emptyLeafs // path expects a list, got a map
	}

	k, ok := parseMapKey(head.Name, fd.MapKey())
	if !ok {
		return emptyLeafs // invalid map key
	}
	v := m.Get(k)
	if !v.IsValid() || v.Equal(fd.MapValue().Default()) {
		return emptyLeafs // key doesn't exist in the map
	}

	// more to the path
	if len(path) > 1 {
		if isLeaf(fd.MapValue()) {
			return emptyLeafs // path has more segments, but the map value is a leaf
		}
		if md := fd.MapValue().Message(); md != nil {
			return scanMessage(path[1:], v.Message())
		}
		return emptyLeafs // there's more to the path but the value has no properties
	}

	// end of the path
	if !isLeaf(fd.MapValue()) {
		return emptyLeafs // path ends at a map value, but the map value is not a leaf
	}
	return func(yield func(leaf) bool) {
		yield(leaf{fd.MapValue(), v})
	}
}

// scanList returns all leafs identified by path in the list l.
func scanList(path []pathSegment, fd protoreflect.FieldDescriptor, l protoreflect.List) iter.Seq[leaf] {
	if len(path) == 0 {
		if !isLeaf(fd) {
			return emptyLeafs // path ends at this list, but the list contains non-leaf values
		}
		return func(yield func(leaf) bool) {
			for i := 0; i < l.Len(); i++ {
				if !yield(leaf{fd, l.Get(i)}) {
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
			return emptyLeafs // index out of range
		}
		v := l.Get(i)
		if isLeaf(fd) {
			if len(path) == 1 {
				// end of the path, return the leaf
				return func(yield func(leaf) bool) {
					yield(leaf{fd, v})
				}
			}
			// path has more segments, but the value is a leaf
			return emptyLeafs
		}

		if fd.Message() == nil {
			return emptyLeafs // it's not a leaf, but also not a message. Not sure what it is.
		}

		return scanMessage(path[1:], v.Message())
	}

	if isLeaf(fd) {
		return emptyLeafs // path has more segments, but the list contains leaf values
	}

	return func(yield func(leaf) bool) {
		for i := 0; i < l.Len(); i++ {
			v := l.Get(i)
			for leaf := range scanMessage(path, v.Message()) {
				if !yield(leaf) {
					return
				}
			}
		}
	}
}

// getLeafDescriptor returns the most relevant FieldDescriptor for values, and whether values is a leaf or not.
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
		// check the parent isn't a leaf
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

// isLeaf returns true if fd is a leaf.
// A leaf is a scalar value, or a well known value type.
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

// toString converts the leaf into a string ready for comparison to another string.
// Unlike l.v.String() this converts enum values to their enum name where available,
// otherwise converts them to a string representation of the enum number.
// Bytes are converted to string.
func (l leaf) toString() (string, bool) {
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
		// for leaf messages, we rely on the json representation to turn them into strings.
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

// toTimestamp converts the leaf into a *timestamppb.Timestamp if it is a valid timestamp.
func (l leaf) toTimestamp() (*timestamppb.Timestamp, bool) {
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
