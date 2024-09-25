package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
)

// AirQualityJob gets average Co2 at current point in time
type AirQualityJob struct {
	BaseJob
	client  traits.AirQualitySensorApiClient
	Sensors []string
}

func (a *AirQualityJob) GetName() string {
	return "air_quality"
}

func (a *AirQualityJob) GetClients() []any {
	return []any{&a.client}
}

func (a *AirQualityJob) Do(ctx context.Context, sendFn sender) error {
	sum := float32(0)
	count := 0

	for _, sensor := range a.Sensors {
		resp, err := a.client.GetAirQuality(ctx, &traits.GetAirQualityRequest{Name: sensor})

		if err != nil {
			a.Logger.Error("getting air quality", zap.String("sensor", sensor), zap.Error(err))
			continue
		}

		count++
		sum += *resp.CarbonDioxideLevel
	}

	if count == 0 {
		return errors.Wrap(errNoSensorsRetrieved, "getting air quality")
	}

	average := sum / float32(count)

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
