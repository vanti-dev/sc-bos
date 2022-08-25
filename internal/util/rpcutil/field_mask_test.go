package rpcutil

import (
	"testing"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestMaskContains(t *testing.T) {
	type testCase struct {
		mask     *fieldmaskpb.FieldMask
		field    string
		expected bool
	}

	cases := map[string]testCase{
		"Nil": {
			mask:     nil,
			field:    "foo",
			expected: true,
		},
		"Empty": {
			mask:     &fieldmaskpb.FieldMask{},
			field:    "foo",
			expected: false,
		},
		"Single_Match": {
			mask:     &fieldmaskpb.FieldMask{Paths: []string{"foo"}},
			field:    "foo",
			expected: true,
		},
		"Single_No_Match": {
			mask:     &fieldmaskpb.FieldMask{Paths: []string{"foo"}},
			field:    "bar",
			expected: false,
		},
		"Multiple_Ordered_Match": {
			mask:     &fieldmaskpb.FieldMask{Paths: []string{"bar", "baz", "foo"}},
			field:    "baz",
			expected: true,
		},
		"Multiple_Reverse_Ordered_Match": {
			mask:     &fieldmaskpb.FieldMask{Paths: []string{"foo", "baz", "bar"}},
			field:    "baz",
			expected: true,
		},
		"Multiple_No_Match": {
			mask:     &fieldmaskpb.FieldMask{Paths: []string{"bar", "baz", "foo"}},
			field:    "zzz",
			expected: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := MaskContains(c.mask, c.field)
			if actual != c.expected {
				t.Errorf("expected MaskContains(...)=%v but got %v", c.expected, actual)
			}
		})
	}
}
