package devices

import (
	"iter"
	"log"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

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

func conditionMatches(cond *gen.Device_Query_Condition, device *gen.Device) bool {
	// everything is a string comparison, for now. Can rework this later if that no longer is the case
	var cmp func(v string) bool
	switch c := cond.Value.(type) {
	case *gen.Device_Query_Condition_StringEqual:
		cmp = func(v string) bool {
			return v == c.StringEqual
		}
	case *gen.Device_Query_Condition_StringEqualFold:
		cmp = func(v string) bool {
			return strings.EqualFold(v, c.StringEqualFold)
		}
	case *gen.Device_Query_Condition_StringContains:
		cmp = func(v string) bool {
			return strings.Contains(v, c.StringContains)
		}
	case *gen.Device_Query_Condition_StringContainsFold:
		ls := strings.ToLower(c.StringContainsFold)
		cmp = func(v string) bool {
			return strings.Contains(strings.ToLower(v), ls)
		}
	default:
		return false
	}

	if cond.Field == "" {
		// any field
		return messageHasValueStringFunc(device, cmp)
	}

	for val := range getMessageString(cond.Field, device) {
		if cmp(val) {
			return true
		}
	}
	return false
}

// messageHasValueStringFunc returns whether any value in msg returns true when converted to a string and passed to f.
// See valueString for the string conversion mechanism.
func messageHasValueStringFunc(msg proto.Message, f func(v string) bool) bool {
	if msg == nil {
		return false
	}
	var match bool
	err := protorange.Range(msg.ProtoReflect(), func(values protopath.Values) error {
		last := values.Index(-1)

		var fd protoreflect.FieldDescriptor
		switch last.Step.Kind() {
		case protopath.FieldAccessStep:
			fd = last.Step.FieldDescriptor()
		case protopath.MapIndexStep:
			fd = values.Index(-2).Step.FieldDescriptor().MapValue()
		case protopath.ListIndexStep:
			fd = values.Index(-2).Step.FieldDescriptor()
		}

		str, ok := valueString(fd, last.Value)
		if !ok {
			return nil
		}
		if f(str) {
			match = true
			return protorange.Terminate
		}
		return nil
	})
	if err != nil {
		// this shouldn't happen as our Range func doesn't return unexpected errors
		log.Printf("Unexpected error during device query processing: %v", err)
		return false
	}
	return match
}

// getMessageString returns an iterator over the string values identified by path from msg.
// See valueString for details of string conversion.
func getMessageString(path string, msg proto.Message) iter.Seq[string] {
	if msg == nil {
		return func(yield func(string) bool) {}
	}
	segments, err := parsePath(path)
	if err != nil || len(segments) == 0 {
		return func(yield func(string) bool) {}
	}
	return getMessageValue(segments, msg.ProtoReflect())
}

// getMessageValue returns an iterator over the string values identified by path.
// The first path segment should refer to a field in msg.
func getMessageValue(path []pathSegment, msg protoreflect.Message) iter.Seq[string] {
	head := path[0]
	fieldDesc := msg.Descriptor().Fields().ByName(protoreflect.Name(head.Name))
	if fieldDesc == nil {
		return func(yield func(string) bool) {}
	}

	if head.IsIndex {
		return func(yield func(string) bool) {}
	}
	val := msg.Get(fieldDesc)
	if len(path) == 1 {
		// end of the path
		if fieldDesc.IsList() {
			return getListValue(path, fieldDesc, val.List())
		}
		return func(yield func(string) bool) {
			if str, got := valueString(fieldDesc, val); got && str != "" {
				yield(str)
			}
		}
	}

	return nextValue(path[1:], fieldDesc, val)
}

// getMapValue returns an iterator over the string values identified by path in the map m.
// The first path segment should refer to a key in the map.
func getMapValue(path []pathSegment, keyDesc, valueDesc protoreflect.FieldDescriptor, m protoreflect.Map) iter.Seq[string] {
	head := path[0]
	key, ok := parseMapKey(head.Name, keyDesc)
	if !ok {
		return func(yield func(string) bool) {}
	}
	value := m.Get(key)
	if !value.IsValid() { // means the key doesn't exist in the map
		return func(yield func(string) bool) {}
	}

	if len(path) == 1 {
		return func(yield func(string) bool) {
			if str, got := valueString(valueDesc, value); got && str != "" {
				yield(str)
			}
		}
	}

	return nextValue(path[1:], valueDesc, value)
}

// getListValue returns an iterator over the string values identified by path in the list l.
// The first path segment can either be an index segment, in which case it refers to a specific element in the list,
// or a non-index segment in which case it is ignored and all elements in the list are matched.
func getListValue(path []pathSegment, entryDesc protoreflect.FieldDescriptor, l protoreflect.List) iter.Seq[string] {
	head := path[0]
	if head.IsIndex {
		// search for a specific element
		idx := head.Index
		if idx < 0 {
			// [-1] is the last element
			idx = l.Len() + idx
		}
		if idx < 0 || idx >= l.Len() {
			// index out of range
			return func(yield func(string) bool) {}
		}
		item := l.Get(idx)
		return func(yield func(string) bool) {
			if str, got := valueString(entryDesc, item); got {
				if str != "" {
					yield(str)
				}
				return
			}
			for v := range getMessageValue(path[1:], item.Message()) {
				if !yield(v) {
					return
				}
			}
		}
	}
	// search all elements
	return func(yield func(string) bool) {
		for i := 0; i < l.Len(); i++ {
			item := l.Get(i)
			// this handles all scalars
			if str, got := valueString(entryDesc, item); got {
				if str != "" {
					if !yield(str) {
						return
					}
				}
				continue
			}
			// this deals with groups/messages
			for v := range getMessageValue(path, item.Message()) {
				if !yield(v) {
					return
				}
			}
		}
	}

}

// nextValue calls the correct getXxxValue func for the given field descriptor.
func nextValue(path []pathSegment, fieldDesc protoreflect.FieldDescriptor, val protoreflect.Value) iter.Seq[string] {
	switch {
	case fieldDesc.IsMap():
		return getMapValue(path, fieldDesc.MapKey(), fieldDesc.MapValue(), val.Map())
	case fieldDesc.IsList():
		return getListValue(path, fieldDesc, val.List())
	case fieldDesc.Message() != nil: // note this is true for map types, so check that first
		return getMessageValue(path, val.Message())
	default:
		return func(yield func(string) bool) {} // there's more to the path but the value has no properties
	}
}

// valueString converts a protoreflect.Value into a string ready for comparison to another string.
// Unlike v.String() this converts enum values to their enum name where available,
// otherwise converts them to a string representation of the enum number.
// Bytes are converted to string.
func valueString(fd protoreflect.FieldDescriptor, v protoreflect.Value) (string, bool) {
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
		// MessageKind, GroupKind
		return "", false
	}
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
