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
	found := false
	switch c := cond.Value.(type) {
	case *gen.Device_Query_Condition_StringEqual:
		cmp = func(v string) bool {
			if v == c.StringEqual {
				found = true
			}
			return found
		}
	case *gen.Device_Query_Condition_StringEqualFold:
		cmp = func(v string) bool {
			if strings.EqualFold(v, c.StringEqualFold) {
				found = true
			}
			return found
		}
	case *gen.Device_Query_Condition_StringContains:
		cmp = func(v string) bool {
			if strings.Contains(v, c.StringContains) {
				found = true
			}
			return found
		}
	case *gen.Device_Query_Condition_StringContainsFold:
		ls := strings.ToLower(c.StringContainsFold)
		cmp = func(v string) bool {
			if strings.Contains(strings.ToLower(v), ls) {
				found = true
			}
			return found
		}
	default:
		return false
	}

	if cond.Field == "" {
		// any field
		return messageHasValueStringFunc(device, cmp)
	}
	itr := getMessageString(cond.Field, device)

	itr(cmp)

	return found
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

// getMessageString returns the property identified by path from msg as a string.
// Returns false if the path does not match any property, or that property cannot be represented as a string,
// or if the comparator func parameter doesn't equate the string value Found
// See valueString for details of string conversion.
func getMessageString(path string, msg proto.Message) iter.Seq[string] {
	if msg == nil {
		return func(yield func(string) bool) {}
	}
	return getMessageValue(path, msg.ProtoReflect())
}

// getMessageValue returns the protoreflect.Value identified by path in msg.
// Returns false if the path can't be resolved, or if the comparator func parameter doesn't equate the value Found
func getMessageValue(path string, msg protoreflect.Message) iter.Seq[string] {
	deconstructedPath := depath(path)
	fieldDesc := msg.Descriptor().Fields().ByName(protoreflect.Name(deconstructedPath.Before))
	if fieldDesc == nil {
		return func(yield func(string) bool) {}
	}

	if deconstructedPath.Found && deconstructedPath.Index < 0 {
		return func(yield func(string) bool) {}
	}
	val := msg.Get(fieldDesc)
	if deconstructedPath.After == "" {
		str, got := valueString(fieldDesc, val)
		// end of the path
		return func(yield func(string) bool) {
			if !(got && yield(str)) {
				return
			}
		}
	}

	return nextValue(deconstructedPath.After, fieldDesc, val)
}

// getMapValue returns the protoreflect.Value identified by path in the map m.
// Returns false if the path can't be resolved, or if the comparator func parameter doesn't equate the value Found
func getMapValue(path string, keyDesc, valueDesc protoreflect.FieldDescriptor, m protoreflect.Map) iter.Seq[string] {
	prop, rest, found := strings.Cut(path, ".")
	key, ok := parseMapKey(prop, keyDesc)
	if !ok {
		return func(yield func(string) bool) {}
	}
	value := m.Get(key)
	if !value.IsValid() { // means the key doesn't exist in the map
		return func(yield func(string) bool) {}
	}

	if !found {
		return func(yield func(string) bool) {
			str, got := valueString(valueDesc, value)

			if !(got && yield(str)) {
				return
			}
		}
	}

	return nextValue(rest, valueDesc, value)
}

// getListValue returns the protoreflect.Value identified by path in the list l.
// Returns false if the path can't be resolved, or if the comparator func parameter doesn't equate the value Found
func getListValue(path string, entryDesc protoreflect.FieldDescriptor, l protoreflect.List) iter.Seq[string] {
	deconstructedPath := depath(path)
	if deconstructedPath.Index < 0 {
		// don't permit negative Index
		if deconstructedPath.Found {
			return func(yield func(string) bool) {}
		}

		var itr []iter.Seq[string]
		// search all elements
		for i := 0; i < l.Len(); i++ {
			val := l.Get(i)

			if !val.Message().IsValid() {
				continue
			}

			desc := val.Message().Descriptor().Fields().ByName(protoreflect.Name(deconstructedPath.Before))

			if desc == nil {
				return func(yield func(string) bool) {}
			}

			var it iter.Seq[string]
			switch {
			case desc.IsMap(), desc.IsList(), desc.Message() != nil:
				it = nextValue(deconstructedPath.After, desc, val.Message().Get(desc))
			default:
				str, got := valueString(desc, val.Message().Get(desc))

				it = func(yield func(string) bool) {
					if !(got && yield(str)) {
						return
					}
				}
			}

			itr = append(itr, it)

		}

		return func(yield func(string) bool) {
			for _, it := range itr {
				it(yield)
			}
		}
	}

	if deconstructedPath.Index >= l.Len() {
		return func(yield func(string) bool) {}
	}

	val := l.Get(deconstructedPath.Index)

	str, got := valueString(entryDesc, val)

	if got {
		return func(yield func(string) bool) {
			if yield(str) {
				return
			}
		}
	}

	if !val.Message().IsValid() {
		return func(yield func(string) bool) {}
	}

	desc := val.Message().Descriptor().Fields().ByName(protoreflect.Name(deconstructedPath.After))

	if desc == nil {
		return func(yield func(string) bool) {}
	}

	str, got = valueString(desc, val.Message().Get(desc))

	if got {
		return func(yield func(string) bool) {
			if yield(str) {
				return
			}
		}
	}

	return nextValue(deconstructedPath.After, desc, val.Message().Get(desc))
}

// nextValue calls the correct getXxxValue func for the given field descriptor.
func nextValue(rest string, fieldDesc protoreflect.FieldDescriptor, val protoreflect.Value) iter.Seq[string] {
	switch {
	case fieldDesc.IsMap():
		return getMapValue(rest, fieldDesc.MapKey(), fieldDesc.MapValue(), val.Map())
	case fieldDesc.IsList():
		return getListValue(rest, fieldDesc, val.List())
	case fieldDesc.Message() != nil: // note this is true for map types, so check that first
		return getMessageValue(rest, val.Message())
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
