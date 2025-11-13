package job

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	sctime "github.com/smart-core-os/sc-api/go/types/time"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type listMeterReadingFn func(ctx context.Context, in *gen.ListMeterReadingHistoryRequest, opts ...grpc.CallOption) (*gen.ListMeterReadingHistoryResponse, error)

func getRecordsByTime(ctx context.Context, logger *zap.Logger, historyFn listMeterReadingFn, meter string, now time.Time, filterTime time.Duration) (earliest, latest *gen.MeterReadingRecord, err error) {
	var resp *gen.ListMeterReadingHistoryResponse

	start := now.Add(-filterTime)

	resp, err = getLastReadingBefore(ctx, meter, start, historyFn)

	if err != nil {
		return nil, nil, err
	}

	if len(resp.GetMeterReadingRecords()) == 0 {
		logger.Error("no records found in earliest", zap.String("meter", meter), zap.Time("start", start))
		return earliest, latest, nil
	}

	earliest = resp.GetMeterReadingRecords()[0]

	resp, err = getLastReadingBefore(ctx, meter, now, historyFn)

	if err != nil {
		return nil, nil, err
	}

	if len(resp.GetMeterReadingRecords()) == 0 {
		logger.Error("no records found in latest", zap.String("meter", meter), zap.Time("end", now))
		return earliest, earliest, nil // make sure this resolves consumption to 0 by returning < earliest, earliest, nil >
	}

	latest = resp.GetMeterReadingRecords()[0]

	return earliest, latest, nil
}

func processMeterRecords(multiplier float32, earliest, latest *gen.MeterReadingRecord) float32 {
	return multiplier * (latest.GetMeterReading().GetUsage() - earliest.GetMeterReading().GetUsage())
}

func getLastReadingBefore(ctx context.Context, meter string, t time.Time, historyFn listMeterReadingFn) (*gen.ListMeterReadingHistoryResponse, error) {
	return historyFn(ctx, &gen.ListMeterReadingHistoryRequest{
		Name: meter,
		Period: &sctime.Period{
			EndTime: timestamppb.New(t),
		},
		PageSize: 1,
		OrderBy:  "record_time DESC",
	})
}
