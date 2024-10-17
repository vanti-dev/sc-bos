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

func getRecordsByTime(ctx context.Context, historyFn listMeterReadingFn, meter string, now time.Time, filterTime time.Duration) ([2]*gen.MeterReadingRecord, error) {
	var (
		pageToken string
		latest    = &gen.MeterReadingRecord{RecordTime: timestamppb.New(time.Time{})}
		earliest  = &gen.MeterReadingRecord{RecordTime: timestamppb.New(now)}
	)

	for {
		resp, err := historyFn(ctx, &gen.ListMeterReadingHistoryRequest{
			Name: meter,
			Period: &sctime.Period{
				StartTime: timestamppb.New(now.Add(-filterTime - time.Second)),
				EndTime:   timestamppb.New(now),
			},
			PageToken: pageToken,
		})

		if err != nil {
			return [2]*gen.MeterReadingRecord{}, err
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

	return [2]*gen.MeterReadingRecord{earliest, latest}, nil
}

func processMeterRecords(multiplier float32, records [2]*gen.MeterReadingRecord) float32 {
	latest := records[1]
	earliest := records[0]

	return multiplier * (latest.GetMeterReading().GetUsage() - earliest.GetMeterReading().GetUsage())
}
