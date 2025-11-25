package job

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

func Test_processMeter(t *testing.T) {
	type args struct {
		historyFn  listMeterReadingFn
		meter      string
		now        time.Time
		filterTime time.Duration
	}

	logger := zap.NewNop()

	start := time.Time{}.Add(time.Nanosecond)
	now := time.Now()
	tests := []struct {
		name         string
		args         args
		wantEarliest *gen.MeterReadingRecord
		wantLatest   *gen.MeterReadingRecord
		wantErr      error
		consumption  float32
	}{
		{
			name: "happy path",
			args: args{
				historyFn: func(ctx context.Context, in *gen.ListMeterReadingHistoryRequest, opts ...grpc.CallOption) (*gen.ListMeterReadingHistoryResponse, error) {
					return twoMeterReadingPages(ctx, start, now, in, opts...)
				},
				now:        now,
				filterTime: -24 * time.Hour,
			},
			wantEarliest: &gen.MeterReadingRecord{
				MeterReading: &gen.MeterReading{
					Usage:     0,
					StartTime: timestamppb.New(start),
					EndTime:   timestamppb.New(start.Add(time.Minute)),
				},
				RecordTime: timestamppb.New(start.Add(time.Minute)),
			},
			wantLatest: &gen.MeterReadingRecord{
				MeterReading: &gen.MeterReading{
					Usage:     45,
					StartTime: timestamppb.New(start.Add(18 * time.Minute)),
					EndTime:   timestamppb.New(start.Add(30 * time.Minute)),
				},
				RecordTime: timestamppb.New(start.Add(30 * time.Minute)),
			},
			wantErr:     nil,
			consumption: 45,
		},
		{
			name: "error path",
			args: args{
				historyFn: func(ctx context.Context, in *gen.ListMeterReadingHistoryRequest, opts ...grpc.CallOption) (*gen.ListMeterReadingHistoryResponse, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			wantEarliest: nil,
			wantLatest:   nil,
			wantErr:      fmt.Errorf("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			earliest, latest, err := getRecordsByTime(context.Background(), logger, tt.args.historyFn, tt.args.meter, tt.args.now, tt.args.filterTime)

			assert.Equal(t, tt.wantEarliest, earliest, tt.name)
			assert.Equal(t, tt.wantLatest, latest, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)

			gotConsumption := processMeterRecords(1, earliest, latest)

			assert.Equal(t, tt.consumption, gotConsumption, tt.name)
		})
	}
}

func twoMeterReadingPages(_ context.Context, start, now time.Time, in *gen.ListMeterReadingHistoryRequest, _ ...grpc.CallOption) (*gen.ListMeterReadingHistoryResponse, error) {
	if in.GetPeriod().EndTime.AsTime().Equal(now) {
		return &gen.ListMeterReadingHistoryResponse{
			MeterReadingRecords: []*gen.MeterReadingRecord{
				{
					MeterReading: &gen.MeterReading{
						Usage:     45,
						StartTime: timestamppb.New(start.Add(18 * time.Minute)),
						EndTime:   timestamppb.New(start.Add(30 * time.Minute)),
					},
					RecordTime: timestamppb.New(start.Add(30 * time.Minute)),
				},
			},
		}, nil
	}
	return &gen.ListMeterReadingHistoryResponse{
		MeterReadingRecords: []*gen.MeterReadingRecord{
			{
				MeterReading: &gen.MeterReading{
					Usage:     0,
					StartTime: timestamppb.New(start),
					EndTime:   timestamppb.New(start.Add(time.Minute)),
				},
				RecordTime: timestamppb.New(start.Add(time.Minute)),
			},
		},
	}, nil
}
