package job

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress/types"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// EnergyJob gets energy consumed over the previous interval (typically 24 hours)
type EnergyJob struct {
	BaseJob
	client     gen.MeterHistoryClient
	infoClient gen.MeterInfoClient
	Meters     []string
	Interval   time.Duration
}

func (e *EnergyJob) GetName() string {
	return "energy"
}

func (e *EnergyJob) GetClients() []any {
	return []any{&e.client, &e.infoClient}
}

func (e *EnergyJob) GetInterval() time.Duration {
	return e.Interval
}

func (e *EnergyJob) Do(ctx context.Context, sendFn sender) error {
	consumption := float32(.0)

	now := time.Now()

	for _, meter := range e.Meters {
		multiplier, err := e.getUnitMultiplier(ctx, meter)

		if err != nil {
			e.Logger.Error("getting unit multiplier", zap.String("meter", meter), zap.Error(err))
			continue
		}

		earliest, latest, err := getRecordsByTime(ctx, e.client.ListMeterReadingHistory, meter, now, e.GetInterval())

		if err != nil {
			e.Logger.Error("getting records by time", zap.String("meter", meter), zap.Error(err))
			continue
		}

		consumption += processMeterRecords(multiplier, earliest, latest)

	}

	body := &types.EnergyConsumption{
		Meta: types.Meta{
			Timestamp: now,
			Site:      e.GetSite(),
		},
		TodaysEnergyConsumption: types.Float32Measure{
			Value: consumption,
			Units: "kWh",
		},
	}

	bytes, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return sendFn(ctx, e.GetUrl(), bytes)
}

func (e *EnergyJob) getUnitMultiplier(ctx context.Context, meter string) (float32, error) {
	infoResp, err := e.infoClient.DescribeMeterReading(ctx, &gen.DescribeMeterReadingRequest{Name: meter})

	if err != nil {
		return 1, err
	}

	//  convert reading to kWh
	var multiplier float32

	switch strings.ToLower(infoResp.GetUnit()) {
	case "wh":
		multiplier = 1 / 1_000
	case "mwh":
		multiplier = 1_000
	case "kwh":
		fallthrough
	default:
		multiplier = 1
	}

	return multiplier, nil
}
