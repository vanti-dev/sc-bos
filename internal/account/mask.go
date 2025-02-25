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

func maskContains(mask fmutils.NestedMask, field protoreflect.Name) bool {
	if _, ok := mask[string(field)]; ok {
		return true
	}
	return false
}

func maskOrDefault(m proto.Message, mask *fieldmaskpb.FieldMask) fmutils.NestedMask {
	if mask == nil {
		return maskForAllFields(m)
	}
	return fmutils.NestedMaskFromPaths(mask.GetPaths())
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
		if r.Has(descriptor) {
			setFields = append(setFields, string(descriptor.Name()))
		}
		return true
	})
	return fmutils.NestedMaskFromPaths(setFields)
}

func maskForAllFields(m proto.Message) fmutils.NestedMask {
	var allFields []string
	m.ProtoReflect().Range(func(descriptor protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		allFields = append(allFields, string(descriptor.Name()))
		return true
	})
	return fmutils.NestedMaskFromPaths(allFields)
}

func fieldsToUpdate(old, new proto.Message, mask fmutils.NestedMask) ([]protoreflect.Name, error) {
	if _, ok := mask["*"]; ok {
		if len(mask) != 1 {
			// a wildcard should be the only field in the mask
			return nil, status.Errorf(codes.InvalidArgument, "mask containing wildcard and other fields is invalid")
		}
		mask = maskForAllFields(new)
	}
	maskFields := make(map[protoreflect.Name]bool)
	for path := range mask {
		maskFields[protoreflect.Name(path)] = true
	}

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
		if !maskFields[field.Name()] {
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
