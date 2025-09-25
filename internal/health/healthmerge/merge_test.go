package healthmerge

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name           string
		src, dst, want *gen.HealthCheck
		ignore         string
	}{
		{
			name: "oneof changes",
			src: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						CurrentValue: intValue(1),
					},
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Faults_{
					Faults: &gen.HealthCheck_Faults{},
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						CurrentValue: intValue(1),
					},
				},
			},
		},
		{
			name: "oneof changes reverse",
			src: &gen.HealthCheck{
				Check: &gen.HealthCheck_Faults_{
					Faults: &gen.HealthCheck_Faults{},
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						CurrentValue: intValue(1),
					},
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Faults_{
					Faults: &gen.HealthCheck_Faults{},
				},
			},
		},
		{
			name: "bounds.normal_values",
			src: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						Expected: intNormalValues(1, 2),
					},
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						Expected: intNormalValues(10, 20),
					},
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						Expected: intNormalValues(1, 2),
					},
				},
			},
		},
		{
			name: "bounds.abnormal_values",
			src: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						Expected: intAbnormalValues(1, 2),
					},
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						Expected: intAbnormalValues(10, 20),
					},
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Bounds_{
					Bounds: &gen.HealthCheck_Bounds{
						Expected: intAbnormalValues(1, 2),
					},
				},
			},
		},
		{
			name: "compliance_impact",
			src: &gen.HealthCheck{
				ComplianceImpacts: []*gen.HealthCheck_ComplianceImpact{
					{Contribution: gen.HealthCheck_ComplianceImpact_FAIL},
				},
			},
			dst: &gen.HealthCheck{
				ComplianceImpacts: []*gen.HealthCheck_ComplianceImpact{
					{Contribution: gen.HealthCheck_ComplianceImpact_RATING},
				},
			},
			want: &gen.HealthCheck{
				ComplianceImpacts: []*gen.HealthCheck_ComplianceImpact{
					{Contribution: gen.HealthCheck_ComplianceImpact_FAIL},
				},
			},
		},
		{
			ignore: "consistency with other merge functions",
			name:   "timestamp",
			src: &gen.HealthCheck{
				ToHealthyAck: &gen.HealthCheck_Ack{
					AckTime: timestamppb.New(time.Unix(100, 0)),
				},
			},
			dst: &gen.HealthCheck{
				ToHealthyAck: &gen.HealthCheck_Ack{
					AckTime: timestamppb.New(time.Unix(10, 100)),
				},
			},
			want: &gen.HealthCheck{
				ToHealthyAck: &gen.HealthCheck_Ack{
					AckTime: timestamppb.New(time.Unix(100, 0)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ignore != "" {
				t.Skipf("skipping test %q: %s", tt.name, tt.ignore)
			}
			Check(proto.Merge, tt.dst, tt.src)
			if diff := cmp.Diff(tt.want, tt.dst, protocmp.Transform()); diff != "" {
				t.Errorf("mergeCheck() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func intValue(value int) *gen.HealthCheck_Value {
	return &gen.HealthCheck_Value{
		Value: &gen.HealthCheck_Value_IntValue{IntValue: int64(value)},
	}
}

func intNormalValue(value int) *gen.HealthCheck_Bounds_NormalValue {
	return &gen.HealthCheck_Bounds_NormalValue{
		NormalValue: intValue(value),
	}
}
func intNormalRange(low, high int) *gen.HealthCheck_Bounds_NormalRange {
	return &gen.HealthCheck_Bounds_NormalRange{
		NormalRange: &gen.HealthCheck_ValueRange{
			Low:  intValue(low),
			High: intValue(high),
		},
	}
}
func intNormalValues(values ...int) *gen.HealthCheck_Bounds_NormalValues {
	valuespb := make([]*gen.HealthCheck_Value, len(values))
	for i, v := range values {
		valuespb[i] = intValue(v)
	}
	return &gen.HealthCheck_Bounds_NormalValues{
		NormalValues: &gen.HealthCheck_Values{Values: valuespb},
	}
}

func intAbnormalValues(values ...int) *gen.HealthCheck_Bounds_AbnormalValues {
	valuespb := make([]*gen.HealthCheck_Value, len(values))
	for i, v := range values {
		valuespb[i] = intValue(v)
	}
	return &gen.HealthCheck_Bounds_AbnormalValues{
		AbnormalValues: &gen.HealthCheck_Values{Values: valuespb},
	}
}
