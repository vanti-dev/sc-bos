package hpd3

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMotionServer_GetMotionDetection(t *testing.T) {
	type testCase struct {
		name   string
		points PointData
		expect *traits.MotionDetection
	}

	cases := []testCase{
		{
			name:   "false",
			points: PointData{Motion1: false},
			expect: &traits.MotionDetection{State: traits.MotionDetection_NOT_DETECTED},
		},
		{
			name:   "true",
			points: PointData{Motion1: true},
			expect: &traits.MotionDetection{State: traits.MotionDetection_DETECTED},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := &memoryClient{points: c.points}

			server := &motionServer{
				client: client,
				logger: testLogger(),
			}
			ctx, cancel := testCtx(t)
			defer cancel()

			res, err := server.GetMotionDetection(ctx, &traits.GetMotionDetectionRequest{})
			if err != nil {
				t.Error(err)
			}

			diff := cmp.Diff(c.expect, res, protocmp.Transform())
			if diff != "" {
				t.Errorf("incorrect Motion result (-want +got):\n%s\n", diff)
			}
		})
	}

}
