package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress/types"
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
		cctx, cancel := context.WithTimeout(ctx, 5*time.Second)

		resp, err := t.client.GetAirTemperature(cctx, &traits.GetAirTemperatureRequest{Name: sensor})
		cancel()

		if err != nil {
			t.Logger.Error("getting air temperature", zap.String("sensor", sensor), zap.Error(err))
			continue
		}
		count++

		sum += resp.GetAmbientTemperature().GetValueCelsius()
	}

	if count == 0 {
		return errors.Wrap(errNoSensorsRetrieved, "getting air temperature")
	}

	average := sum / float64(count)

	body := &types.AverageTemperature{
		Meta: types.Meta{
			Site:      t.GetSite(),
			Timestamp: time.Now(),
		},
		AverageTemperature: types.Float64Measure{
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
