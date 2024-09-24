package job

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
)

// TemperatureJob gets average air temperature at current point in time
type TemperatureJob struct {
	BaseJob
	client  traits.AirTemperatureApiClient
	Sensors []string
}

func (t *TemperatureJob) GetName() string {
	return "temperature"
}

func (t *TemperatureJob) GetClients() []any {
	return []any{&t.client}
}

func (t *TemperatureJob) Do(ctx context.Context, sendFn sender) error {
	sum := .0
	count := 0

	for _, sensor := range t.Sensors {
		resp, err := t.client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{Name: sensor})

		if err != nil {
			t.Logger.Error("getting air temperature", zap.String("sensor", sensor), zap.Error(err))
			continue
		}
		count++

		sum += resp.GetAmbientTemperature().GetValueCelsius()
	}

	average := sum / float64(count)

	body := &AverageTemperature{
		Meta: Meta{
			Site:      t.GetSite(),
			Timestamp: time.Now(),
		},
		AverageTemperature: Float64Measure{
			Value: average,
			Units: "°C",
		},
	}

	bytes, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return sendFn(ctx, t.GetUrl(), bytes)
}
