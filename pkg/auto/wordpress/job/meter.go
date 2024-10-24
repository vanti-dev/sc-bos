package job

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	sctime "github.com/smart-core-os/sc-api/go/types/time"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type listMeterReadingFn func(ctx context.Context, in *gen.ListMeterReadingHistoryRequest, opts ...grpc.CallOption) (*gen.ListMeterReadingHistoryResponse, error)

func getRecordsByTime(ctx context.Context, historyFn listMeterReadingFn, meter string, now time.Time, filterTime time.Duration) (earliest, latest *gen.MeterReadingRecord, err error) {
	var (
		pageToken string
		resp      *gen.ListMeterReadingHistoryResponse
	)

	latest = &gen.MeterReadingRecord{RecordTime: timestamppb.New(time.Time{})}
	earliest = &gen.MeterReadingRecord{RecordTime: timestamppb.New(now)}

	for {
		resp, err = historyFn(ctx, &gen.ListMeterReadingHistoryRequest{
			Name: meter,
			Period: &sctime.Period{
				StartTime: timestamppb.New(now.Add(-filterTime - time.Second)),
				EndTime:   timestamppb.New(now),
			},
			PageToken: pageToken,
		})

		if err != nil {
			return nil, nil, err
		}

		for _, record := range resp.GetMeterReadingRecords() {
			if record.GetRecordTime().AsTime().Before(earliest.GetRecordTime().AsTime()) {
				earliest = record
			}

			if record.GetRecordTime().AsTime().After(latest.GetRecordTime().AsTime()) {
				latest = record
			}
		}

		if resp.GetNextPageToken() == "" {
			break
		}

		pageToken = resp.GetNextPageToken()
	}

	return earliest, latest, nil
}

func processMeterRecords(multiplier float32, earliest, latest *gen.MeterReadingRecord) float32 {
	return multiplier * (latest.GetMeterReading().GetUsage() - earliest.GetMeterReading().GetUsage())
}
