package merge

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
				Check: &gen.HealthCheck_Check{
					Bounds: intNormalValue(1),
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intNormalRange(10, 100),
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intNormalValue(1),
				},
			},
		},
		{
			name: "check.normal_values",
			src: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intNormalValues(1, 2),
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intNormalValues(10, 20),
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intNormalValues(1, 2),
				},
			},
		},
		{
			name: "check.abnormal_values",
			src: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intAbnormalValues(1, 2),
				},
			},
			dst: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intAbnormalValues(10, 20),
				},
			},
			want: &gen.HealthCheck{
				Check: &gen.HealthCheck_Check{
					Bounds: intAbnormalValues(1, 2),
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

func intNormalValue(value int) *gen.HealthCheck_Check_NormalValue {
	return &gen.HealthCheck_Check_NormalValue{
		NormalValue: &gen.HealthCheck_Value{
			Value: &gen.HealthCheck_Value_IntValue{IntValue: int64(value)},
		},
	}
}
func intNormalRange(low, high int) *gen.HealthCheck_Check_NormalRange {
	return &gen.HealthCheck_Check_NormalRange{
		NormalRange: &gen.HealthCheck_ValueRange{
			Low:  &gen.HealthCheck_Value{Value: &gen.HealthCheck_Value_IntValue{IntValue: int64(low)}},
			High: &gen.HealthCheck_Value{Value: &gen.HealthCheck_Value_IntValue{IntValue: int64(high)}},
		},
	}
}
func intNormalValues(values ...int) *gen.HealthCheck_Check_NormalValues {
	valuespb := make([]*gen.HealthCheck_Value, len(values))
	for i, v := range values {
		valuespb[i] = &gen.HealthCheck_Value{
			Value: &gen.HealthCheck_Value_IntValue{IntValue: int64(v)},
		}
	}
	return &gen.HealthCheck_Check_NormalValues{
		NormalValues: &gen.HealthCheck_Values{Values: valuespb},
	}
}

func intAbnormalValues(values ...int) *gen.HealthCheck_Check_AbnormalValues {
	valuespb := make([]*gen.HealthCheck_Value, len(values))
	for i, v := range values {
		valuespb[i] = &gen.HealthCheck_Value{
			Value: &gen.HealthCheck_Value_IntValue{IntValue: int64(v)},
		}
	}
	return &gen.HealthCheck_Check_AbnormalValues{
		AbnormalValues: &gen.HealthCheck_Values{Values: valuespb},
	}
}
