package job

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress/types"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// WaterJob gets water consumed over the previous interval (typically 24 hours)
type WaterJob struct {
	BaseJob
	client     gen.MeterHistoryClient
	infoClient gen.MeterInfoClient
	Meters     []string
	Interval   time.Duration
}

func (w *WaterJob) GetName() string {
	return "water"
}

func (w *WaterJob) GetClients() []any {
	return []any{&w.client, &w.infoClient}
}

func (w *WaterJob) GetInterval() time.Duration {
	return w.Interval
}

func (w *WaterJob) Do(ctx context.Context, sendFn sender) error {
	consumption := float32(.0)

	now := time.Now()

	for _, meter := range w.Meters {
		multiplier, err := w.getUnitMultiplier(ctx, meter)

		if err != nil {
			w.Logger.Error("getting unit multiplier", zap.String("meter", meter), zap.Error(err))
			continue
		}

		records, err := getRecordsByTime(ctx, w.client.ListMeterReadingHistory, meter, now, w.GetInterval())

		if err != nil {
			w.Logger.Error("getting records by time", zap.String("meter", meter), zap.Error(err))
			continue
		}

		consumption += processMeterRecords(multiplier, records)

	}

	body := &types.WaterConsumption{
		Meta: types.Meta{
			Timestamp: now,
			Site:      w.GetSite(),
		},
		TodaysWaterConsumption: types.IntMeasure{
			Value: int32(consumption),
			Units: "litres",
		},
	}

	bytes, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return sendFn(ctx, w.GetUrl(), bytes)
}

func (w *WaterJob) getUnitMultiplier(ctx context.Context, meter string) (float32, error) {
	infoResp, err := w.infoClient.DescribeMeterReading(ctx, &gen.DescribeMeterReadingRequest{Name: meter})

	if err != nil {
		return 1, err
	}

	// convert reading to litres
	var multiplier float32

	switch infoResp.GetUnit() {
	case "cm3": // TODO: these strings may need correcting I tried guessing them
		multiplier = 1 / 1_000_000
	case "m3":
		fallthrough
	case "litres":
		fallthrough
	default:
		multiplier = 1
	}

	return multiplier, nil
}
