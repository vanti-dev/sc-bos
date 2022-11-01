package bridge

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
)

func TestNotification_Decode(t *testing.T) {
	type testCase struct {
		input     notification
		expectErr error
		expect    []dali.InputEvent
	}

	cases := map[string]testCase{
		"Invalid": {
			input: notification{
				Valid:          false,
				Sequence:       0,
				NumInputEvents: 0,
				InputEvents:    nil,
			},
			expectErr: ErrInvalid,
			expect:    nil,
		},
		"Empty": {
			input: notification{
				Valid:          true,
				Sequence:       0,
				NumInputEvents: 0,
				InputEvents:    make([]inputEvent, 32),
			},
			expect: nil,
		},
		"Single": {
			input: notification{
				Valid:          true,
				Sequence:       123,
				NumInputEvents: 1,
				InputEvents: []inputEvent{
					{
						Parameters: dali.InputEventParametersForInstance(0x12, 0x34),
						Error:      false,
						Status:     0,
						Message:    "",
						Data:       1234,
					},
					{
						Parameters: dali.InputEventParametersForInstance(0x12, 0x34),
						Error:      false,
						Status:     0,
						Message:    "",
						Data:       5678,
					},
				},
			},
			expect: []dali.InputEvent{
				{
					InputEventParameters: dali.InputEventParametersForInstance(0x12, 0x34),
					Err:                  nil,
					Data:                 1234,
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := c.input.Decode()

			if diff := cmp.Diff(c.expectErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(c.expect, actual, cmpopts.EquateErrors(), cmpopts.EquateEmpty()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestInputEvent_Decode(t *testing.T) {
	type testCase struct {
		input  inputEvent
		expect dali.InputEvent
	}

	cases := map[string]testCase{
		"Error": {
			input: inputEvent{
				Parameters: dali.InputEventParametersForInstance(0x12, 0x34),
				Error:      true,
				Status:     123,
				Message:    "Message",
				Data:       1234,
			},
			expect: dali.InputEvent{
				InputEventParameters: dali.InputEventParametersForInstance(0x12, 0x34),
				Err:                  cmpopts.AnyError,
				Data:                 1234,
			},
		},
		"OK": {
			input: inputEvent{
				Parameters: dali.InputEventParametersForInstance(0x12, 0x34),
				Error:      false,
				Status:     0,
				Message:    "Message",
				Data:       1234,
			},
			expect: dali.InputEvent{
				InputEventParameters: dali.InputEventParametersForInstance(0x12, 0x34),
				Err:                  nil,
				Data:                 1234,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := c.input.Decode()

			if diff := cmp.Diff(c.expect, actual, cmpopts.EquateErrors()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
