package hpd3

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestOccupancyServer_GetOccupancy(t *testing.T) {
	type testCase struct {
		name   string
		points PointData
		expect *traits.Occupancy
	}

	cases := []testCase{
		{
			name: "none",
			points: PointData{
				Presence1:           false,
				NumberOfPeopleTotal: 0,
			},
			expect: &traits.Occupancy{
				State:       traits.Occupancy_UNOCCUPIED,
				PeopleCount: 0,
			},
		},
		{
			name: "some",
			points: PointData{
				Presence1:           true,
				NumberOfPeopleTotal: 3,
			},
			expect: &traits.Occupancy{
				State:       traits.Occupancy_OCCUPIED,
				PeopleCount: 3,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := &memoryClient{points: c.points}

			server := &occupancyServer{
				client: client,
				logger: testLogger(),
			}
			ctx, cancel := testCtx(t)
			defer cancel()

			res, err := server.GetOccupancy(ctx, &traits.GetOccupancyRequest{})
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
