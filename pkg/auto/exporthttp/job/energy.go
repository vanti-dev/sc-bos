package job

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/exporthttp/types"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// EnergyJob gets energy consumed over the previous execution interval (typically 24 hours)
type EnergyJob struct {
	BaseJob
	client     gen.MeterHistoryClient
	infoClient gen.MeterInfoClient
	Meters     []string
}

func (e *EnergyJob) GetName() string {
	return "energy"
}

func (e *EnergyJob) Do(ctx context.Context, sendFn sender) error {
	consumption := float32(.0)

	now := time.Now().UTC()
	filterTime := now.Sub(e.PreviousExecution.UTC())

	for _, meter := range e.Meters {
		cctx, cancel := context.WithTimeout(ctx, e.Timeout.Or(defaultTimeout))

		multiplier, err := e.getUnitMultiplier(cctx, meter)

		if err != nil {
			e.Logger.Error("getting unit multiplier", zap.String("meter", meter), zap.Error(err))
		}

		earliest, latest, err := getRecordsByTime(cctx, e.client.ListMeterReadingHistory, meter, now, filterTime)

		cancel()

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
			Value: float32(math.Floor(float64(consumption))),
			Units: "kWh",
		},
	}

	bytes, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = sendFn(ctx, e.GetUrl(), bytes)

	e.PreviousExecution = time.Now().UTC()

	return err
}

func (e *EnergyJob) getUnitMultiplier(ctx context.Context, meter string) (float32, error) {
	infoResp, err := e.infoClient.DescribeMeterReading(ctx, &gen.DescribeMeterReadingRequest{Name: meter})

	if err != nil {
		return 1, err
	}

	//  convert reading to kWh
	var multiplier float32

	switch strings.ToLower(infoResp.GetUsageUnit()) {
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
