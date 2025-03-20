package account

import (
	"errors"
	"fmt"

	"github.com/mennanov/fmutils"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var ErrInvalidWildcardMask = status.Error(codes.InvalidArgument, "field mask cannot contain wildcard and field names")

func resolveMask(m proto.Message, mask *fieldmaskpb.FieldMask, ignorePaths ...string) (fmutils.NestedMask, error) {
	ignoreMask := fmutils.NestedMaskFromPaths(ignorePaths)
	if mask == nil {
		return maskForSetFields(m, ignoreMask), nil
	}
	if slices.Contains(mask.GetPaths(), "*") {
		if len(mask.GetPaths()) != 1 {
			return nil, ErrInvalidWildcardMask
		}
		return maskForAllFields(m, ignoreMask), nil
	}

	// check that all fields specified in the mask exist on the message descriptor
	nestedMask := fmutils.NestedMaskFromPaths(mask.GetPaths())
	err := validateMask(m.ProtoReflect(), "", nestedMask)
	return nestedMask, err
}

func validateMask(m protoreflect.Message, prefix string, mask fmutils.NestedMask) error {
	desc := m.Descriptor()
	for path, submask := range mask {
		fieldDesc := desc.Fields().ByName(protoreflect.Name(path))
		if fieldDesc == nil {
			return status.Errorf(codes.InvalidArgument, "invalid mask: unknown field %q", prefix+path)
		}
		field := m.Get(fieldDesc)

		if len(submask) > 0 && fieldDesc.Kind() != protoreflect.MessageKind {
			return status.Errorf(codes.InvalidArgument, "invalid mask: field %q is not a message", prefix+path)
		}

		if fieldDesc.Kind() == protoreflect.MessageKind && field.IsValid() {
			if err := validateMask(field.Message(), prefix+path+".", submask); err != nil {
				return err
			}
		}
	}
	return nil
}

// returns a FieldMask with the fields that are set in the message
// shallow - only the top level fields are considered
// except for the fields in the ignore list
func maskForSetFields(m proto.Message, ignore fmutils.NestedMask) fmutils.NestedMask {
	return fmutils.NestedMaskFromPaths(setFieldNames(m.ProtoReflect(), "", ignore))
}

// returns a list of leaf (non-message) field paths that are set in the given message
// field paths are prefixed with the given prefix
// ignore fields will not be added; only simple field names can be ignored, not paths containing a '.'
func setFieldNames(m protoreflect.Message, prefix string, ignore fmutils.NestedMask) []string {
	var setFields []string
	m.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		subignore, hasSubignore := ignore[string(descriptor.Name())]
		if hasSubignore && len(subignore) == 0 {
			// this field is ignored
			return true
		}
		if descriptor.Kind() == protoreflect.MessageKind {
			prefix := prefix + string(descriptor.Name()) + "."
			setFields = append(setFields, setFieldNames(value.Message(), prefix, subignore)...)
		} else {
			setFields = append(setFields, prefix+string(descriptor.Name()))
		}
		return true
	})
	return setFields
}

func maskForAllFields(m proto.Message, ignore fmutils.NestedMask) fmutils.NestedMask {
	return fmutils.NestedMaskFromPaths(allFieldPaths(m.ProtoReflect().Descriptor(), "", ignore))
}

func allFieldPaths(d protoreflect.MessageDescriptor, prefix string, ignore fmutils.NestedMask) []string {
	var allFields []string
	for i := range d.Fields().Len() {
		descriptor := d.Fields().Get(i)
		subignore, hasSubignore := ignore[string(descriptor.Name())]
		if hasSubignore && len(subignore) == 0 {
			// this field is ignored
			continue
		}
		allFields = append(allFields, prefix+string(descriptor.Name()))
		if descriptor.Kind() == protoreflect.MessageKind {
			prefix := prefix + string(descriptor.Name()) + "."
			allFields = append(allFields, allFieldPaths(descriptor.Message(), prefix, subignore)...)
		}
	}
	return allFields
}

func fieldsToUpdate(old, new proto.Message, mask fmutils.NestedMask) ([]string, error) {
	return diffMessages(old.ProtoReflect(), new.ProtoReflect(), "", mask)
}

func diffMessages(old, new protoreflect.Message, prefix string, mask fmutils.NestedMask) (paths []string, err error) {
	oldDesc, newDesc := old.Descriptor(), new.Descriptor()
	if oldDesc != newDesc {
		return nil, errors.New("messages are of different types")
	}

	for path, subMask := range mask {
		field := newDesc.Fields().ByName(protoreflect.Name(path))
		if field == nil {
			return nil, fmt.Errorf("unknown field %q in mask", path)
		}
		newValue := new.Get(field)
		oldValue := old.Get(field)
		if newValue.Equal(oldValue) {
			continue
		}

		// no child fields, so this mask entry represents the entire sub-message
		if len(subMask) == 0 {
			paths = append(paths, prefix+string(field.Name()))
			continue
		}

		if field.Kind() == protoreflect.MessageKind {
			prefix := prefix + string(field.Name()) + "."
			subPaths, err := diffMessages(oldValue.Message(), newValue.Message(), prefix, subMask)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		} else {
			paths = append(paths, prefix+string(field.Name()))
		}
	}
	return paths, nil
}
