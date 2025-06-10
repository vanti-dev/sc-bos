// App dbadd-alerts is a tool that creates alerts via the AlertAdminApi.
// The app will seed the database with a number of alerts then sit in the background sending new updates and resolving old alerts.
// Alerts will be associated with devices that already exist on the server.
//
// At a high level this tool:
//  1. Connects to a SC server
//  2. Collects all the devices the server has - via the DevicesApi
//  3. Every few seconds, picks a random device and generates a random alert for that device
//  4. Sometimes alerts associated with that random device are resolved instead of being created.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var (
	AlertRate = rate.Every(5 * time.Second)
	MinAlerts = 100
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	err := run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	host := "localhost:23557"
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}

	devicesClient := gen.NewDevicesApiClient(conn)
	alertClient := gen.NewAlertApiClient(conn)
	alertAdminClient := gen.NewAlertAdminApiClient(conn)

	// get the devices that will be used as sources for the alerts
	devices, err := listDevices(ctx, devicesClient)
	if err != nil {
		return fmt.Errorf("listDevices: %w", err)
	}
	if len(devices) == 0 {
		return fmt.Errorf("no devices found, nothing to act as alert sources")
	}
	randDevice := func() *gen.Device {
		return devices[rand.Intn(len(devices))]
	}
	writeAlert := func() (*gen.Alert, bool, error) {
		d := randDevice()
		a := genAlert(d)
		// randomly resolve some alerts instead of updating them
		if rand.Intn(10) == 0 {
			_, err := alertAdminClient.ResolveAlert(ctx, &gen.ResolveAlertRequest{Alert: a, AllowMissing: true})
			if err != nil {
				return nil, false, fmt.Errorf("resolve alert: %w", err)
			}
			return a, true, nil
		}
		_, err := alertAdminClient.CreateAlert(ctx, &gen.CreateAlertRequest{Alert: a, MergeSource: true})
		if err != nil {
			return nil, false, fmt.Errorf("add alert: %w", err)
		}
		return a, false, nil
	}

	// work out what the current state of things to play nicely with existing databases
	_, alertsNeeded, err := getSeedState(ctx, alertClient)
	if err != nil {
		return fmt.Errorf("getSeedState: %w", err)
	}

	// prefill the db with alerts up to a set threshold
	if alertsNeeded > 0 {
		fmt.Println("Seeding alerts:", alertsNeeded)
	}
	for i := 0; i < alertsNeeded; i++ {
		if _, _, err := writeAlert(); err != nil {
			return fmt.Errorf("seed: %w", err)
		}
	}

	limiter := rate.NewLimiter(AlertRate, 1)
	for {
		if err := limiter.Wait(ctx); err != nil {
			return nil // ctx is done, stop adding alerts
		}

		a, resolved, err := writeAlert()
		if err != nil {
			return fmt.Errorf("write alert: %w", err)
		}
		if resolved {
			fmt.Println("Resolved alert:", a.Source)
		} else {
			fmt.Println("Added alert:", a.Source, a.Severity)
		}
	}
}

func listDevices(ctx context.Context, devicesClient gen.DevicesApiClient) ([]*gen.Device, error) {
	req := &gen.ListDevicesRequest{}
	var res []*gen.Device

	for {
		resp, err := devicesClient.ListDevices(ctx, req)
		if err != nil {
			return res, err
		}
		res = append(res, resp.Devices...)
		if resp.NextPageToken == "" {
			break
		}
		req.PageToken = resp.NextPageToken
	}
	return res, nil
}

func getSeedState(ctx context.Context, alertClient gen.AlertApiClient) (startTime time.Time, alertsNeeded int, err error) {
	fail := func(err error) (time.Time, int, error) {
		return time.Time{}, 0, err
	}
	alerts, err := alertClient.ListAlerts(ctx, &gen.ListAlertsRequest{PageSize: 1})
	if err != nil {
		return fail(err)
	}

	startTime = time.Now().Add(-50 * time.Hour)
	if len(alerts.Alerts) > 0 {
		startTime = alerts.Alerts[0].CreateTime.AsTime()
	}

	metadata, err := alertClient.GetAlertMetadata(ctx, &gen.GetAlertMetadataRequest{})
	if err != nil {
		return fail(err)
	}
	alertsNeeded = MinAlerts - int(metadata.TotalCount)
	if alertsNeeded <= 0 {
		alertsNeeded = 0
	}
	return startTime, alertsNeeded, nil
}

func genAlert(d *gen.Device) *gen.Alert {
	return &gen.Alert{
		Description: fmt.Sprintf("Something happened to %s", d.Name),
		CreateTime:  timestamppb.Now(),
		Severity:    randSeverity(),
		Floor:       d.GetMetadata().GetLocation().GetFloor(),
		Zone:        d.GetMetadata().GetLocation().GetZone(),
		Source:      d.Name,
		Subsystem:   d.GetMetadata().GetMembership().GetSubsystem(),
	}
}

var severities = maps.Keys(gen.Alert_Severity_name)

func randSeverity() gen.Alert_Severity {
	// don't include unspecified
	return gen.Alert_Severity(severities[rand.Intn(len(severities)-1)+1])
}
