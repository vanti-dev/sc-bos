// Package meteremail provides an automation that collects the instantaneous meter readings for a set of given devices
// The automation uses the Meter API to fetch meter readings on a configurable fixed date
// formats an email using html/template, and sends it to some recipients using smtp.
// Test program for meter reading automation is in cmd/tools/test-meteremail/main.go
package meteremail

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/auto/meteremail/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/task"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"time"
)

const AutoName = "meteremail"

var Factory auto.Factory = factory{}

type factory struct{}

type autoImpl struct {
	*service.Service[config.Root]
	auto.Services
}

func (f factory) New(services auto.Services) service.Lifecycle {
	a := &autoImpl{Services: services}
	a.Service = service.New(service.MonoApply(a.applyConfig), service.WithParser(config.ReadBytes))
	a.Logger = a.Logger.Named(AutoName)
	return a
}

func (a *autoImpl) getMeterReadingAndSource(ctx context.Context, conn *grpc.ClientConn, meterName string, meterType MeterType) (*config.Source, *MeterReading, error) {
	meterClient := gen.NewMeterApiClient(conn)
	metaDataClient := traits.NewMetadataApiClient(conn)

	meterReq := &gen.GetMeterReadingRequest{
		Name: meterName,
	}

	meterRes, err := retryT(ctx, func(ctx context.Context) (*gen.MeterReading, error) {
		return meterClient.GetMeterReading(ctx, meterReq)
	})
	if err != nil {
		a.Logger.Warn("failed to fetch meter readings", zap.Error(err))
		return nil, nil, err
	}

	metaDataReq := &traits.GetMetadataRequest{
		Name: meterName,
	}

	metaDataRes, err := metaDataClient.GetMetadata(ctx, metaDataReq)
	if err != nil {
		// not a major problem if we can't get metadata as we still have the name to ID the meter
		a.Logger.Warn("failed to fetch meta data", zap.Error(err))
	}

	meterReading := &MeterReading{MeterType: meterType, Date: time.Now(), Reading: meterRes.Usage}
	source := &config.Source{Name: meterName}

	if metaDataRes.Location != nil {
		source.Floor = metaDataRes.Location.Floor
		source.Zone = metaDataRes.Location.Zone
	}

	return source, meterReading, nil
}

func (a *autoImpl) applyConfig(ctx context.Context, cfg config.Root) error {
	logger := a.Logger
	logger = logger.With(zap.String("snmp.host", cfg.Destination.Host), zap.Int("snmp.port", cfg.Destination.Port))

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true // TODO set from config
	conn, err := grpc.Dial(cfg.ServerAddr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: can't connect: %s\n", err.Error())
		os.Exit(1)
	}

	sendTime := cfg.Destination.SendTime
	now := cfg.Now
	if now == nil {
		now = a.Now
	}
	if now == nil {
		now = time.Now
	}
	go func() {
		t := now()
		for {
			next := sendTime.Next(t)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Until(next)):
				// Use the time we were planning on running instead of the current time.
				// We do this to make output more predictable
				t = next
			}

			attrs := Attrs{
				Now:   t,
				Stats: []Stats{},
			}

			for _, meterName := range cfg.ElectricMeters {
				source, reading, err := a.getMeterReadingAndSource(ctx, conn, meterName, MeterTypeElectric)
				if err == nil {
					attrs.Stats = append(attrs.Stats, Stats{Source: *source, MeterReading: *reading})
				}
			}

			for _, meterName := range cfg.WaterMeters {
				source, reading, err := a.getMeterReadingAndSource(ctx, conn, meterName, MeterTypeWater)
				if err == nil {
					attrs.Stats = append(attrs.Stats, Stats{Source: *source, MeterReading: *reading})
				}
			}

			err := cfg.Destination.AttachFile(".testattach.csv")

			if err != nil {
				logger.Warn("failed to add attachment", zap.Error(err))
			}

			err = retry(ctx, func(ctx context.Context) error {
				return sendEmail(cfg.Destination, attrs)
			})
			if err != nil {
				logger.Warn("failed to send email", zap.Error(err))
			} else {
				logger.Info("email sent")
			}
		}
	}()

	return nil
}

func retry(ctx context.Context, f func(context.Context) error) error {
	return task.Run(ctx, func(ctx context.Context) (task.Next, error) {
		return 0, f(ctx)
	}, task.WithBackoff(10*time.Second, 10*time.Minute), task.WithRetry(40))
}

func retryT[T any](ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	var t T
	err := retry(ctx, func(ctx context.Context) error {
		var err error
		t, err = f(ctx)
		return err
	})
	return t, err
}
