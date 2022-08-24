package rpcutil

import (
	"sort"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func MaskContains(mask *fieldmaskpb.FieldMask, field string) bool {
	if mask == nil {
		// a nil field mask means that same as a field mask including all fields
		return true
	}

	mask = proto.Clone(mask).(*fieldmaskpb.FieldMask)
	mask.Normalize()
	sort.Strings(mask.Paths)
	i := sort.SearchStrings(mask.Paths, field)
	return i < len(mask.Paths) && mask.Paths[i] == field
}
