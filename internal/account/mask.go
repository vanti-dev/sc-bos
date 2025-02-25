package account

import (
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

func maskContains(mask fmutils.NestedMask, field protoreflect.Name) bool {
	if _, ok := mask[string(field)]; ok {
		return true
	}
	return false
}

func resolveMask(m proto.Message, mask *fieldmaskpb.FieldMask, ignore ...protoreflect.Name) (fmutils.NestedMask, error) {
	if mask == nil {
		return maskForSetFields(m, ignore...), nil
	}
	if slices.Contains(mask.GetPaths(), "*") {
		if len(mask.GetPaths()) != 1 {
			return nil, ErrInvalidWildcardMask
		}
		return maskForAllFields(m, ignore...), nil
	}

	// check that all fields specified in the mask exist on the message descriptor
	desc := m.ProtoReflect().Descriptor()
	paths := make([]string, 0, len(mask.GetPaths()))
	for _, path := range mask.GetPaths() {
		field := desc.Fields().ByName(protoreflect.Name(path))
		if field == nil {
			return nil, status.Errorf(codes.InvalidArgument, "unknown field %q in mask", path)
		}
		if !slices.Contains(ignore, field.Name()) {
			paths = append(paths, path)
		}
	}
	return fmutils.NestedMaskFromPaths(paths), nil
}

// returns a FieldMask with the fields that are set in the message
// shallow - only the top level fields are considered
// except for the fields in the ignore list
func maskForSetFields(m proto.Message, ignore ...protoreflect.Name) fmutils.NestedMask {
	var setFields []string
	r := m.ProtoReflect()
	r.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		if slices.Contains(ignore, descriptor.Name()) {
			return true
		}
		setFields = append(setFields, string(descriptor.Name()))
		return true
	})
	return fmutils.NestedMaskFromPaths(setFields)
}

func maskForAllFields(m proto.Message, ignore ...protoreflect.Name) fmutils.NestedMask {
	var allFields []string
	fieldDescs := m.ProtoReflect().Descriptor().Fields()
	for i := 0; i < fieldDescs.Len(); i++ {
		fieldDesc := fieldDescs.Get(i)
		name := fieldDesc.Name()
		if slices.Contains(ignore, name) {
			continue
		}
		allFields = append(allFields, string(name))
	}
	return fmutils.NestedMaskFromPaths(allFields)
}

func fieldsToUpdate(old, new proto.Message, mask fmutils.NestedMask) ([]protoreflect.Name, error) {
	var fields []protoreflect.Name
	newR, oldR := new.ProtoReflect(), old.ProtoReflect()
	newDesc, oldDesc := newR.Descriptor(), oldR.Descriptor()
	if newDesc != oldDesc {
		return nil, fmt.Errorf("messages are of different types")
	}
	for path := range mask {
		field := newR.Descriptor().Fields().ByName(protoreflect.Name(path))
		if field == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid field %q in mask", path)
		}
		if _, ok := mask[string(field.Name())]; !ok {
			continue
		}

		newValue := newR.Get(field)
		oldValue := oldR.Get(field)
		if newValue.Equal(oldValue) {
			continue
		}

		fields = append(fields, field.Name())
	}
	return fields, nil
}
