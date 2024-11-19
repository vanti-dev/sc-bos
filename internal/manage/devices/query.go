package devices

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var arrIndexRegx = regexp.MustCompile("\\[(-?[0-9]+)]")

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
	_, res := compareMessageString(cond.Field, device, cmp)
	return res
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

// compareMessageString returns the property identified by path from msg as a string.
// Returns false if the path does not match any property, or that property cannot be represented as a string,
// or if the comparator func parameter doesn't equate the string value found
// See valueString for details of string conversion.
func compareMessageString(path string, msg proto.Message, f func(v string) bool) (string, bool) {
	if msg == nil {
		return "", false
	}
	fd, v, ok := compareMessageValue(path, msg.ProtoReflect(), f)
	if !ok {
		return "", false
	}
	vs, ok := valueString(fd, v)
	if !ok {
		return "", false
	}

	return vs, f(vs)
}

// compareMessageValue returns the protoreflect.Value identified by path in msg.
// Returns false if the path can't be resolved, or if the comparator func parameter doesn't equate the value found
func compareMessageValue(path string, msg protoreflect.Message, f func(v string) bool) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	prop, rest, found, index := depath(path)
	fieldDesc := msg.Descriptor().Fields().ByName(protoreflect.Name(prop))
	if fieldDesc == nil {
		return fieldDesc, protoreflect.ValueOf(nil), false
	}

	if found && index < 0 {
		return fieldDesc, protoreflect.ValueOf(nil), false
	}
	val := msg.Get(fieldDesc)
	if rest == "" {
		// end of the path
		return fieldDesc, val, true
	}

	return nextValue(rest, fieldDesc, val, f)
}

// compareMapValue returns the protoreflect.Value identified by path in the map m.
// Returns false if the path can't be resolved, or if the comparator func parameter doesn't equate the value found
func compareMapValue(path string, keyDesc, valueDesc protoreflect.FieldDescriptor, m protoreflect.Map, f func(v string) bool) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
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

	return nextValue(rest, valueDesc, value, f)
}

func compareNestedElement(restPath string, desc protoreflect.FieldDescriptor, val protoreflect.Value, f func(v string) bool) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	if val.IsValid() {
		str, got := valueString(desc, val.Message().Get(desc))
		if got && f(str) {
			return desc, val.Message().Get(desc), true
		}
		descriptor, value, found := nextValue(restPath, desc, val.Message().Get(desc), f)

		if found {
			if _, rest, found := strings.Cut(restPath, "."); found {
				return compareNestedElement(rest, descriptor, value, f)
			}

			str, got = valueString(descriptor, value)

			return descriptor, value, f(str) && got
		}
	}

	return nil, protoreflect.Value{}, false
}

// compareListValue returns the protoreflect.Value identified by path in the list l.
// Returns false if the path can't be resolved, or if the comparator func parameter doesn't equate the value found
func compareListValue(path string, entryDesc protoreflect.FieldDescriptor, l protoreflect.List, f func(v string) bool) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	prop, rest, found, index := depath(path)

	if index < 0 {
		// don't permit negative index
		if found {
			return nil, protoreflect.Value{}, false
		}

		// search all elements
		for i := 0; i < l.Len(); i++ {
			val := l.Get(i)

			if !val.Message().IsValid() {
				continue
			}

			desc := val.Message().Descriptor().Fields().ByName(protoreflect.Name(prop))

			if descriptor, value, found := compareNestedElement(rest, desc, val, f); found {
				return descriptor, value, true
			}
		}

		return nil, protoreflect.Value{}, false
	}

	if index >= l.Len() {
		return nil, protoreflect.Value{}, false
	}

	val := l.Get(index)

	switch entryDesc.Kind() {
	case protoreflect.StringKind:
		return entryDesc, val, f(val.String())
	case protoreflect.MessageKind:
		break
	default: // we don't support other primitives in f func
		return nil, protoreflect.Value{}, false
	}

	if !val.Message().IsValid() {
		return nil, protoreflect.Value{}, false
	}

	desc := val.Message().Descriptor().Fields().ByName(protoreflect.Name(rest))

	if desc == nil {
		return nil, protoreflect.Value{}, false
	}

	if _, _, found := compareNestedElement(rest, desc, val, f); found {
		return desc, val.Message().Get(desc), true
	}

	return nextValue(rest, desc, val, f)
}

// nextValue calls the correct getXxxValue func for the given field descriptor.
func nextValue(rest string, fieldDesc protoreflect.FieldDescriptor, val protoreflect.Value, f func(v string) bool) (protoreflect.FieldDescriptor, protoreflect.Value, bool) {
	switch {
	case fieldDesc.IsMap():
		return compareMapValue(rest, fieldDesc.MapKey(), fieldDesc.MapValue(), val.Map(), f)
	case fieldDesc.IsList():
		return compareListValue(rest, fieldDesc, val.List(), f)
	case fieldDesc.Message() != nil: // note this is true for map types, so check that first
		return compareMessageValue(rest, val.Message(), f)
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

// depath deconstructs a path to extract an array index if present
// returns true for found if a possible integer is found wrapped in [ ]
func depath(path string) (before, after string, found bool, index int) {
	before, after, _ = strings.Cut(path, ".")

	if arrIndexRegx.MatchString(path) {
		matches := arrIndexRegx.FindStringSubmatch(path)

		if len(matches) < 2 {
			return before, after, true, -1
		}

		index, err := strconv.ParseInt(matches[1], 10, 32)
		if err == nil && index > -1 {
			matchedIndices := arrIndexRegx.FindStringIndex(path)

			if matchedIndices == nil || len(matchedIndices) < 2 {
				return before, after, true, -1
			}
			if matchedIndices[0] == 0 {
				// An index is found at the start of path
				// return before,after only
				return before, after, true, int(index)
			}
			return arrIndexRegx.ReplaceAllString(before, ""), fmt.Sprintf("[%d].%s", index, after), true, int(index)
		}

		return arrIndexRegx.ReplaceAllString(before, ""), after, true, -1
	}

	return before, after, false, -1
}
