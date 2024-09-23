package job

import (
	"context"
	"slices"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	sctime "github.com/smart-core-os/sc-api/go/types/time"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type listMeterReadingFn func(ctx context.Context, in *gen.ListMeterReadingHistoryRequest, opts ...grpc.CallOption) (*gen.ListMeterReadingHistoryResponse, error)

func getAllRecords(ctx context.Context, historyFn listMeterReadingFn, meter string, now time.Time, filterTime time.Duration) ([]*gen.MeterReadingRecord, error) {
	var (
		pageToken string
		records   []*gen.MeterReadingRecord
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
			return nil, err
		}

		records = append(records, resp.GetMeterReadingRecords()...)

		if resp.GetNextPageToken() == "" {
			break
		}

		pageToken = resp.GetNextPageToken()

	}

	return records, nil
}

func processMeterRecords(multiplier float32, records []*gen.MeterReadingRecord) float32 {
	if len(records) < 2 {
		return 0
	}

	slices.SortFunc(records, func(i, j *gen.MeterReadingRecord) int {
		timeI := i.GetRecordTime().AsTime()
		timeJ := j.GetRecordTime().AsTime()

		if timeI.Before(timeJ) {
			return 1
		}
		if timeI.After(timeJ) {
			return -1
		}
		return 0 // If times are equal
	})

	latest := records[0]
	earliest := records[len(records)-1]

	return multiplier * (latest.GetMeterReading().GetUsage() - earliest.GetMeterReading().GetUsage())

}
