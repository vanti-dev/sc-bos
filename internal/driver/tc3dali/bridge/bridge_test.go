package bridge

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestAdsResponse_AsError(t *testing.T) {
	type testCase struct {
		input  bridgeResponse
		expect error
	}

	cases := map[string]testCase{
		"Invalid": {
			input: bridgeResponse{
				Valid: false,
			},
			expect: ErrInvalid,
		},
		"Not_IsError_and_zero": {
			// IsError false, zero status
			input: bridgeResponse{
				Valid:   true,
				IsError: false,
				Status:  0,
				Message: "foo",
			},
			expect: nil,
		},
		"Not_IsError_and_nonzero": {
			// IsError false, nonzero status
			input: bridgeResponse{
				Valid:   true,
				IsError: false,
				Status:  123,
				Message: "foo",
			},
			expect: nil,
		},
		"IsError_and_zero": {
			// IsError true, zero status
			// The bridge should never return this, but it's considered a non-error state because status==0
			// means 'OK'
			input: bridgeResponse{
				Valid:   true,
				IsError: true,
				Status:  0,
				Message: "foo",
			},
			expect: nil,
		},
		"IsError_and_nonzero": {
			// IsError true, nonzero status
			input: bridgeResponse{
				Valid:   true,
				IsError: true,
				Status:  123,
				Message: "foo",
			},
			expect: Error{
				Status:  123,
				Message: "foo",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := c.input.AsError()
			if diff := cmp.Diff(c.expect, actual, cmpopts.EquateErrors()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
