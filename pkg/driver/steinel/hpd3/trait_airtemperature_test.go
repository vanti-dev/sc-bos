package hpd3

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAirTemperatureServer_GetAirTemperature(t *testing.T) {
	type testCase struct {
		name   string
		points PointData
		expect *traits.AirTemperature
	}

	cases := []testCase{
		{
			name: "basic",
			points: PointData{
				Temperature: 25.5,
				Humidity:    50.5,
			},
			expect: &traits.AirTemperature{
				AmbientTemperature: &types.Temperature{ValueCelsius: 25.5},
				AmbientHumidity:    ptr[float32](50.5),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := &memoryClient{points: c.points}
			server := &airTemperatureServer{
				client: client,
				logger: testLogger(),
			}
			ctx, cancel := testCtx(t)
			defer cancel()

			res, err := server.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{})
			if err != nil {
				t.Error(err)
			}

			diff := cmp.Diff(c.expect, res, protocmp.Transform(), cmpopts.EquateApprox(0.0001, 0))
			if diff != "" {
				t.Errorf("incorrect Air Temperature result (-want +got):\n%s\n", diff)
			}
		})
	}
}
