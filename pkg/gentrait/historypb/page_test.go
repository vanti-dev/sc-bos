package historypb

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	timepb "github.com/smart-core-os/sc-api/go/types/time"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history/memstore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_pageReader_listRecords(t *testing.T) {
	s := memstore.New()
	now := time.UnixMilli(0)
	t.Cleanup(memstore.SetNow(s, func() time.Time {
		return now
	}))
	ctx := context.Background()

	for _, i := range []int64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90} {
		now = time.UnixMilli(i)
		data, err := proto.Marshal(&traits.Occupancy{StateChangeTime: timestamppb.New(now)})
		if err != nil {
			t.Fatal(err)
		}
		_, err = s.Append(ctx, data)
		if err != nil {
			t.Fatal(err)
		}
	}

	pr := occupancyPager
	page, size, nextToken, err := pr.listRecords(context.Background(), s, &timepb.Period{}, 5, "")
	if err != nil {
		t.Fatal(err)
	}
	if size != 10 {
		t.Fatalf("size want %v, got %v", 10, size)
	}
	if nextToken == "" {
		t.Fatalf("nextPage want something, got nothing")
	}
	wantPage := make([]*gen.OccupancyRecord, 5)
	for i, t := range []int64{0, 10, 20, 30, 40} {
		wantPage[i] = &gen.OccupancyRecord{
			RecordTime: timestamppb.New(time.UnixMilli(t)),
			Occupancy:  &traits.Occupancy{StateChangeTime: timestamppb.New(time.UnixMilli(t))},
		}
	}
	if diff := cmp.Diff(wantPage, page, protocmp.Transform()); diff != "" {
		t.Fatalf("page (-want,+got)\n%s", diff)
	}

	page, size, nextToken, err = pr.listRecords(context.Background(), s, &timepb.Period{}, 5, nextToken)
	if err != nil {
		t.Fatal(err)
	}
	if size != 10 {
		t.Fatalf("size want %v, got %v", 10, size)
	}
	if nextToken != "" {
		t.Fatalf("nextPage want empty, got %v", nextToken)
	}
	for i, t := range []int64{50, 60, 70, 80, 90} {
		wantPage[i] = &gen.OccupancyRecord{
			RecordTime: timestamppb.New(time.UnixMilli(t)),
			Occupancy:  &traits.Occupancy{StateChangeTime: timestamppb.New(time.UnixMilli(t))},
		}
	}
	if diff := cmp.Diff(wantPage, page, protocmp.Transform()); diff != "" {
		t.Fatalf("page (-want,+got)\n%s", diff)
	}
}

func Test_pageReader_listRecords_period(t *testing.T) {
	s := memstore.New()
	now := time.UnixMilli(0)
	t.Cleanup(memstore.SetNow(s, func() time.Time {
		return now
	}))
	ctx := context.Background()

	for _, i := range []int64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90} {
		now = time.UnixMilli(i)
		data, err := proto.Marshal(&traits.Occupancy{StateChangeTime: timestamppb.New(now)})
		if err != nil {
			t.Fatal(err)
		}
		_, err = s.Append(ctx, data)
		if err != nil {
			t.Fatal(err)
		}
	}

	pr := occupancyPager
	page, size, nextToken, err := pr.listRecords(context.Background(), s, &timepb.Period{
		StartTime: timestamppb.New(time.UnixMilli(30)),
		EndTime:   timestamppb.New(time.UnixMilli(70)),
	}, 5, "")
	if err != nil {
		t.Fatal(err)
	}
	if size != 4 {
		t.Fatalf("size want %v, got %v", 10, size)
	}
	if nextToken != "" {
		t.Fatalf("nextPage want nothing, got %s", nextToken)
	}
	wantPage := make([]*gen.OccupancyRecord, 4)
	for i, t := range []int64{30, 40, 50, 60} {
		wantPage[i] = &gen.OccupancyRecord{
			RecordTime: timestamppb.New(time.UnixMilli(t)),
			Occupancy:  &traits.Occupancy{StateChangeTime: timestamppb.New(time.UnixMilli(t))},
		}
	}
	if diff := cmp.Diff(wantPage, page, protocmp.Transform()); diff != "" {
		t.Fatalf("page (-want,+got)\n%s", diff)
	}
}
