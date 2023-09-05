package hpd3

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAirQualityServer_GetAirQuality(t *testing.T) {
	client := &memoryClient{
		points: PointData{
			CO2: 123,
			VOC: 456,
		},
	}

	server := &airQualityServer{
		client: client,
		logger: testLogger(),
	}
	ctx, cancel := testCtx(t)
	defer cancel()

	res, err := server.GetAirQuality(ctx, &traits.GetAirQualityRequest{})
	if err != nil {
		t.Error(err)
	}

	expected := &traits.AirQuality{
		CarbonDioxideLevel:       ptr[float32](123.0),
		VolatileOrganicCompounds: ptr[float32](456.0 / 1000),
	}
	diff := cmp.Diff(expected, res, protocmp.Transform(), cmpopts.EquateApprox(0.0001, 0))
	if diff != "" {
		t.Errorf("incorrect Air Quality result (-want +got):\n%s\n", diff)
	}
}
