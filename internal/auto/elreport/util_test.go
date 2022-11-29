package elreport

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
)

func TestSortDeduplicateFaults(t *testing.T) {
	type testCase struct {
		input  []gen.EmergencyLightFault
		expect []gen.EmergencyLightFault
	}

	cases := map[string]testCase{
		"Empty": {
			input:  nil,
			expect: nil,
		},
		"Single": {
			input:  []gen.EmergencyLightFault{gen.EmergencyLightFault_LAMP_FAULT},
			expect: []gen.EmergencyLightFault{gen.EmergencyLightFault_LAMP_FAULT},
		},
		"Single_Duplicate": {
			input:  []gen.EmergencyLightFault{gen.EmergencyLightFault_LAMP_FAULT, gen.EmergencyLightFault_LAMP_FAULT},
			expect: []gen.EmergencyLightFault{gen.EmergencyLightFault_LAMP_FAULT},
		},
		"Two_Reorder": {
			input:  []gen.EmergencyLightFault{gen.EmergencyLightFault_OTHER_FAULT, gen.EmergencyLightFault_BATTERY_FAULT},
			expect: []gen.EmergencyLightFault{gen.EmergencyLightFault_BATTERY_FAULT, gen.EmergencyLightFault_OTHER_FAULT},
		},
		"Two_Duplicate": {
			input: []gen.EmergencyLightFault{
				gen.EmergencyLightFault_OTHER_FAULT,
				gen.EmergencyLightFault_OTHER_FAULT,
				gen.EmergencyLightFault_BATTERY_FAULT,
			},
			expect: []gen.EmergencyLightFault{
				gen.EmergencyLightFault_BATTERY_FAULT,
				gen.EmergencyLightFault_OTHER_FAULT,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// this mutates c.input, but that is OK as we only use it once
			actual := sortDeduplicateFaults(c.input)
			if diff := cmp.Diff(c.expect, actual, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s\n", diff)
			}
		})
	}
}

func TestFaultsEquivalent_Binary(t *testing.T) {
	type binaryCase struct {
		x, y   []gen.EmergencyLightFault
		expect bool
	}

	binaryCases := map[string]binaryCase{
		"Single_Different": {
			x:      []gen.EmergencyLightFault{gen.EmergencyLightFault_COMMUNICATION_FAILURE},
			y:      []gen.EmergencyLightFault{gen.EmergencyLightFault_LAMP_FAULT},
			expect: false,
		},
		"Duplicate": {
			x: []gen.EmergencyLightFault{
				gen.EmergencyLightFault_DURATION_TEST_FAILED,
				gen.EmergencyLightFault_DURATION_TEST_FAILED,
			},
			y: []gen.EmergencyLightFault{
				gen.EmergencyLightFault_DURATION_TEST_FAILED,
			},
			expect: true,
		},
		"Reversed": {
			x: []gen.EmergencyLightFault{
				gen.EmergencyLightFault_COMMUNICATION_FAILURE,
				gen.EmergencyLightFault_FUNCTION_TEST_FAILED,
			},
			y: []gen.EmergencyLightFault{
				gen.EmergencyLightFault_FUNCTION_TEST_FAILED,
				gen.EmergencyLightFault_COMMUNICATION_FAILURE,
			},
			expect: true,
		},
	}

	for name, c := range binaryCases {
		t.Run(name, func(t *testing.T) {
			// the function is supposed to be commutative so test both ways round
			actualForward := faultsEquivalent(c.x, c.y)
			actualReverse := faultsEquivalent(c.y, c.x)
			if actualForward != c.expect {
				t.Errorf("Expected faultsEquivalent(x, y) == %v, but got %v", c.expect, actualForward)
			}
			if actualReverse != c.expect {
				t.Errorf("Expected faultsEquivalent(y, x) == %v, but got %v", c.expect, actualReverse)
			}
		})
	}

}

func TestFaultsEquivalent_Reflexive(t *testing.T) {
	reflexiveCases := map[string][]gen.EmergencyLightFault{
		"Empty": nil,
		"Single": {
			gen.EmergencyLightFault_BATTERY_FAULT,
		},
		"Two": {
			gen.EmergencyLightFault_COMMUNICATION_FAILURE,
			gen.EmergencyLightFault_OTHER_FAULT,
		},
	}

	for name, c := range reflexiveCases {
		t.Run(name, func(t *testing.T) {
			if !faultsEquivalent(c, c) {
				t.Errorf("faultsEquivalent is not reflexive on %v", c)
			}
		})
	}
}
