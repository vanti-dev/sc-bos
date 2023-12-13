// Command list-occupancy-history reads records from [gen.OccupancySensorHistoryClient] and writes them to a CSV file.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	timepb "github.com/smart-core-os/sc-api/go/types/time"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func main() {
	conn, err := grpc.DialContext(context.Background(), "10.1.104.3:23557", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	now, err := time.Parse(time.DateTime, "2023-11-23 00:00:00")
	if err != nil {
		panic(err)
	}
	now = now.Truncate(time.Second)
	from, to := now.Add(-7*24*time.Hour), now

	client := gen.NewOccupancySensorHistoryClient(conn)
	req := &gen.ListOccupancyHistoryRequest{
		Name: "enterprisewharf.co.uk/zones/building",
		Period: &timepb.Period{
			StartTime: timestamppb.New(from),
			EndTime:   timestamppb.New(to),
		},
		PageSize: 1000,
	}
	var records []*gen.OccupancyRecord
	for {
		resp, err := client.ListOccupancyHistory(context.Background(), req)
		if err != nil {
			panic(err)
		}
		records = append(records, resp.OccupancyRecords...)
		if resp.NextPageToken == "" {
			break
		}
		req.PageToken = resp.NextPageToken
	}

	// output CSV
	for _, r := range records {
		fmt.Printf("%s,%d\n", r.RecordTime.AsTime().Format(time.RFC3339), r.Occupancy.PeopleCount)
	}
}
