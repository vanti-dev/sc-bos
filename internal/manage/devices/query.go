package devices

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	return isMessageValueStringFunc(cond.Field, device, cmp)
}

// isMessageValueStringFunc returns whether the value identified by path in msg returns true when passed to f.
// The path argument looks like `some.prop.path` and navigates the message fields.
func isMessageValueStringFunc(path string, msg proto.Message, f func(v string) bool) bool {
	if msg == nil {
		return false
	}
	fd, v, ok := getMessageValue(path, msg.ProtoReflect())
	if !ok {
		return false
	}
	vs, ok := valueString(fd, v)
	if !ok {
		return false
	}

	return f(vs)
}

// messageHasValueStringFunc returns whether any value in msg returns true when converted to a string and passed to f.
// See valueString for the string conversion mechanism.
func messageHasValueStringFunc(msg proto.Message, f func(v string) bool) bool {
	if msg == nil {
		return false
	}
	var match bool
	err := protorange.Range(msg.ProtoReflect(), func(values protopath.Values) error {
		out := values.Index(-1)
		str, ok := valueString(out.Step.FieldDescriptor(), out.Value)
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

// getMessageString returns the property identified by path from msg as a string.
// Returns false if the path does not match any property, or that property cannot be represented as a string.
// See valueString for details of string conversion.
func getMessageString(path string, msg proto.Message) (string, bool) {
	if msg == nil {
		return "", false
	}
	fd, v, ok := getMessageValue(path, msg.ProtoReflect())
	if !ok {
		return "", false
	}
	vs, ok := valueString(fd, v)
	if !ok {
		return "", false
	}

	return vs, true
}

// getMessageValue returns the protoreflect.Value identified by path in msg.
// Returns false if the path can't be resolved.
func getMessageValue(path string, msg protoreflect.Message) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	prop, rest, found := strings.Cut(path, ".")
	fieldDesc := msg.Descriptor().Fields().ByName(protoreflect.Name(prop))
	if fieldDesc == nil {
		return fieldDesc, protoreflect.ValueOf(nil), false
	}
	val := msg.Get(fieldDesc)
	if !found {
		// end of the path
		return fieldDesc, val, true
	}

	return nextValue(rest, fieldDesc, val)
}

// getMapValue returns the protoreflect.Value identified by path in the map m.
// Returns false if the path can't be resolved.
func getMapValue(path string, keyDesc, valueDesc protoreflect.FieldDescriptor, m protoreflect.Map) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	prop, rest, found := strings.Cut(path, ".")
	key, ok := parseMapKey(prop, keyDesc)
	if !ok {
		return nil, protoreflect.Value{}, false
	}
	value := m.Get(key)
	if !value.IsValid() { // means the key doesn't exist in the map
		return nil, protoreflect.Value{}, false
	}

	if !found {
		return valueDesc, value, true
	}

	return nextValue(rest, valueDesc, value)
}

// getListValue returns the protoreflect.Value identified by path in the list l.
// Returns false if the path can't be resolved.
func getListValue(path string, entryDesc protoreflect.FieldDescriptor, l protoreflect.List) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	// todo: support list values
	// I guess it'd be nice to support "any in list match" but also "index in list matches" semantics.
	// In the first case we'd need to refactor this whole things to combine value lookup and condition checking.
	// In the second case we'd need to make our path parsing logic more capable/complicated.
	return nil, protoreflect.Value{}, false
}

// nextValue calls the correct getXxxValue func for the given field descriptor.
func nextValue(rest string, fieldDesc protoreflect.FieldDescriptor, val protoreflect.Value) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	switch {
	case fieldDesc.IsMap():
		return getMapValue(rest, fieldDesc.MapKey(), fieldDesc.MapValue(), val.Map())
	case fieldDesc.IsList():
		return getListValue(rest, fieldDesc, val.List())
	case fieldDesc.Message() != nil: // note this is true for map types, so check that first
		return getMessageValue(rest, val.Message())
	default:
		return fieldDesc, val, false // there's more to the path but the value has no properties
	}
}

// valueString converts a protoreflect.Value into a string ready for comparison to another string.
// Unlike v.String() this converts enum values to their enum name where available,
// otherwise converts them to a string representation of the enum number.
// Bytes are converted to string.
func valueString(fd protoreflect.FieldDescriptor, v protoreflect.Value) (string, bool) {
	if fd == nil {
		return valueToString(v)
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

func valueToString(v protoreflect.Value) (string, bool) {
	return fmt.Sprintf("%v", v.Interface()), true
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
