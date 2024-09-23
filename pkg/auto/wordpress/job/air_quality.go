package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
)

// AirQualityJob gets average Co2 at current point in time
type AirQualityJob struct {
	BaseJob
	client  traits.AirQualitySensorApiClient
	Sensors []string
}

func (b *AirQualityJob) GetName() string {
	return "air_quality"
}

func (a *AirQualityJob) GetClients() []any {
	return []any{&a.client}
}
func (a *AirQualityJob) Do(ctx context.Context, sendFn sender) error {
	if len(a.Sensors) < 1 {
		return nil
	}

	sum := float32(0)

	for _, sensor := range a.Sensors {
		resp, err := a.client.GetAirQuality(ctx, &traits.GetAirQualityRequest{Name: sensor})

		if err != nil {
			return err
		}

		sum += *resp.CarbonDioxideLevel
	}

	average := sum / float32(len(a.Sensors))

	body := &AverageCo2{
		Meta: Meta{
			Site:      a.GetSite(),
			Timestamp: time.Now(),
		},
		AverageCo2: Float32Measure{
			Value: average,
			Units: "ppm",
		},
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return sendFn(ctx, a.GetUrl(), bytes)
}
