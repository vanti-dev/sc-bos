package healthpb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestFaultCheck_AddOrUpdateFault(t *testing.T) {
	tests := map[string]struct {
		initial []*gen.HealthCheck_Error
		new     *gen.HealthCheck_Error
		want    []*gen.HealthCheck_Error
	}{
		"nil initial": {
			initial: nil,
			new:     newFault("", "", "summary1", "desc1"),
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1")},
		},
		"nil new": {
			initial: []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1")},
			new:     nil,
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1")},
		},
		"add new to end": {
			initial: []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1")},
			new:     newFault("", "", "summary2", "desc2"),
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1"), newFault("", "", "summary2", "desc2")},
		},
		"add new to start": {
			initial: []*gen.HealthCheck_Error{newFault("", "", "summary2", "desc2")},
			new:     newFault("", "", "summary1", "desc1"),
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1"), newFault("", "", "summary2", "desc2")},
		},
		"add new in middle": {
			initial: []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1"), newFault("", "", "summary3", "desc3")},
			new:     newFault("", "", "summary2", "desc2"),
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1"), newFault("", "", "summary2", "desc2"), newFault("", "", "summary3", "desc3")},
		},
		"replace existing by summary": {
			initial: []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc1")},
			new:     newFault("", "", "summary1", "desc2"),
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary1", "desc2")},
		},
		"replace existing by system/code": {
			initial: []*gen.HealthCheck_Error{newFault("sys1", "code1", "summary1", "desc1")},
			new:     newFault("sys1", "code1", "summary2", "desc2"),
			want:    []*gen.HealthCheck_Error{newFault("sys1", "code1", "summary2", "desc2")},
		},
		"add new with different system/code": {
			initial: []*gen.HealthCheck_Error{newFault("sys1", "code1", "summary1", "desc1")},
			new:     newFault("sys2", "code2", "summary1", "desc2"), // same summary
			want:    []*gen.HealthCheck_Error{newFault("sys1", "code1", "summary1", "desc1"), newFault("sys2", "code2", "summary1", "desc2")},
		},
		"replace existing with system/code, add new by summary": {
			initial: []*gen.HealthCheck_Error{newFault("sys1", "code1", "summary1", "desc1")},
			new:     newFault("", "", "summary2", "desc2"),
			want:    []*gen.HealthCheck_Error{newFault("", "", "summary2", "desc2"), newFault("sys1", "code1", "summary1", "desc1")},
		},
		"multiple initial, replace one": {
			initial: []*gen.HealthCheck_Error{
				newFault("", "", "summary1", "desc1"),
				newFault("sys1", "code1", "summary2", "desc2"),
				newFault("", "", "summary3", "desc3"),
			},
			new: newFault("sys1", "code1", "summary2", "desc2-updated"),
			want: []*gen.HealthCheck_Error{
				newFault("", "", "summary1", "desc1"),
				newFault("sys1", "code1", "summary2", "desc2-updated"),
				newFault("", "", "summary3", "desc3"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := &gen.HealthCheck{
				Check: &gen.HealthCheck_Faults_{Faults: &gen.HealthCheck_Faults{
					CurrentFaults: tt.initial,
				}},
			}
			fc, err := newFaultCheck(c)
			if err != nil {
				t.Fatalf("newFaultCheck() error = %v", err)
			}
			fc.AddOrUpdateFault(tt.new)
			if diff := cmp.Diff(tt.want, fc.check.GetFaults().GetCurrentFaults(), protocmp.Transform()); diff != "" {
				t.Errorf("AddOrUpdateFault() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func newFault(system, code, summary, desc string) *gen.HealthCheck_Error {
	res := &gen.HealthCheck_Error{
		SummaryText: summary,
		DetailsText: desc,
	}
	if system != "" || code != "" {
		res.Code = &gen.HealthCheck_Error_Code{
			System: system,
			Code:   code,
		}
	}
	return res
}
