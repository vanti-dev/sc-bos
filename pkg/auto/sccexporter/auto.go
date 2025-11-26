// Package sccexporter exports Device to scc from an on-premise smart core instance.
// Given a list of traits, the exporter will discover all devices which implement that trait and export
// the Device on a scheduled interval to an MQTT broker in the format expected by scc.
package sccexporter

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/auto"
	"github.com/smart-core-os/sc-bos/pkg/auto/sccexporter/config"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	meterpb "github.com/smart-core-os/sc-bos/pkg/gentrait/meter"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-bos/pkg/task/service"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

const AutoName = "sccexporter"

var Factory auto.Factory = factory{}

type factory struct{}

type AutoImpl struct {
	*service.Service[config.Root]
	auto.Services

	airQualityClient     traits.AirQualitySensorApiClient
	airTemperatureClient traits.AirTemperatureApiClient
	metadataClient       traits.MetadataApiClient
	meterClient          gen.MeterApiClient
	meterInfoClient      gen.MeterInfoClient
	occupancyClient      traits.OccupancySensorApiClient
}

func (f factory) New(services auto.Services) service.Lifecycle {
	a := &AutoImpl{
		Services: services,
	}
	a.Service = service.New(service.MonoApply(a.applyConfig), service.WithParser(config.ParseConfig))
	a.Logger = a.Logger.Named(AutoName)
	return a
}

func (a *AutoImpl) initialiseClients(n *node.Node) {
	a.airQualityClient = traits.NewAirQualitySensorApiClient(n.ClientConn())
	a.airTemperatureClient = traits.NewAirTemperatureApiClient(n.ClientConn())
	a.metadataClient = traits.NewMetadataApiClient(n.ClientConn())
	a.meterClient = gen.NewMeterApiClient(n.ClientConn())
	a.meterInfoClient = gen.NewMeterInfoClient(n.ClientConn())
	a.occupancyClient = traits.NewOccupancySensorApiClient(n.ClientConn())
}

func (a *AutoImpl) applyConfig(ctx context.Context, cfg config.Root) error {

	a.initialiseClients(a.Node)

	grp, autoCtx := errgroup.WithContext(ctx)
	allDevices := make(map[string]*device)

	// discover all the devices which implement the configured traits and set up the allDevices map
	for _, traitName := range cfg.Traits {
		err := a.getAllTraitImplementors(autoCtx, trait.Name(traitName), allDevices)
		if err != nil {
			a.Logger.Error("failed to get devices for trait", zap.String("trait", traitName), zap.Error(err))
			return err
		}

		switch traitName {
		case string(meterpb.TraitName):
			// grab the trait info for all meters first and save it in the device so we can push it
			// only supports the Meter info, think it's the only one we really need for data...
			a.getMeterInfo(autoCtx, trait.Name(traitName), allDevices)
		}
	}

	mqttClient, err := newMqttClient(cfg.Mqtt)
	if err != nil {
		a.Logger.Error("failed to create mqtt client", zap.Error(err))
		return err
	}

	s := newSccConnector(a.Logger, cfg.Mqtt, mqttClient)

	grp.Go(func() error {
		return s.publishToScc(autoCtx)
	})

	t := time.Now()
	iterationCount := 0
	grp.Go(func() error {
		for {
			next := cfg.Mqtt.SendInterval.Next(t)
			select {
			case <-autoCtx.Done():
				return nil
			case <-time.After(time.Until(next)):
				t = next
			}

			// send the metadata on first run and then every now and again.
			// it is never actually refreshed during the automation. Don't think any site actually modifies it at runtime.
			publishMetadata := (iterationCount % *cfg.Mqtt.MetadataInterval) == 0
			iterationCount++

			// limit concurrent device data fetches to 100, an arbitrary number but seems sensible
			sem := make(chan struct{}, 100)
			var wg errgroup.Group

			for _, dev := range allDevices {
				dev := dev

				sem <- struct{}{}
				wg.Go(func() error {
					defer func() { <-sem }()
					a.fetchAndPublishDeviceData(autoCtx, dev, cfg.Mqtt.Agent, s.messagesCh, publishMetadata, cfg.FetchTimeout.Duration)
					return nil
				})
			}

			// Wait for all device fetches to complete before next interval, should be fine as long as the interval is sensible
			if err := wg.Wait(); err != nil {
				a.Logger.Error("error fetching device data", zap.Error(err))
			}
		}
	})

	// applyConfig returns immediately, background tasks run until ctx is cancelled
	// When ctx is cancelled (on reconfigure or stop), cleanup happens after all goroutines complete
	go func() {
		err := grp.Wait()
		mqttClient.Disconnect(500)
		close(s.messagesCh)
		if err != nil && err != context.Canceled {
			a.Logger.Error("sccexporter automation stopped with error", zap.Error(err))
		}
	}()

	return nil

}

// getAllTraitImplementors populates the devices map with devices that have the given trait
func (a *AutoImpl) getAllTraitImplementors(ctx context.Context, traitName trait.Name, devices map[string]*device) error {
	resp, err := a.Services.Devices.ListDevices(ctx, &gen.ListDevicesRequest{
		Query: &gen.Device_Query{
			Conditions: []*gen.Device_Query_Condition{
				{
					Field: "metadata.traits.name",
					Value: &gen.Device_Query_Condition_StringEqual{
						StringEqual: string(traitName),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	for _, deviceInfo := range resp.Devices {
		if _, ok := devices[deviceInfo.Name]; !ok {
			devices[deviceInfo.Name] = newDevice(deviceInfo.Name, a.Logger, deviceInfo.Metadata)
			switch traitName {
			case trait.AirTemperature:
				devices[deviceInfo.Name].traits[trait.AirTemperature] = func(ctx context.Context) (string, error) {
					return devices[deviceInfo.Name].getAirTemperatureData(ctx, a.airTemperatureClient)
				}
			case trait.AirQualitySensor:
				devices[deviceInfo.Name].traits[trait.AirQualitySensor] = func(ctx context.Context) (string, error) {
					return devices[deviceInfo.Name].getAirQualityData(ctx, a.airQualityClient)
				}
			case meterpb.TraitName:
				devices[deviceInfo.Name].traits[meterpb.TraitName] = func(ctx context.Context) (string, error) {
					return devices[deviceInfo.Name].getMeterData(ctx, a.meterClient)
				}
			case trait.OccupancySensor:
				devices[deviceInfo.Name].traits[trait.OccupancySensor] = func(ctx context.Context) (string, error) {
					return devices[deviceInfo.Name].getOccupancyData(ctx, a.occupancyClient)
				}
			default:
				a.Logger.Warn("trait is configured but not supported",
					zap.String("trait", string(traitName)), zap.String("device", deviceInfo.Name))
			}
		}
	}
	return nil
}

// fetchAndPublishDeviceData fetches data for all traits of a device and sends it to the messages channel.
// If includeMetadata is true, device metadata is also added to the Data map.
// The fetchTimeout parameter limits how long we wait for each device's data to prevent slow devices from blocking.
func (a *AutoImpl) fetchAndPublishDeviceData(ctx context.Context, dev *device, agent string, messagesCh chan message, includeMetadata bool, fetchTimeout time.Duration) {
	// Add per-device timeout to prevent slow/hanging devices from blocking the entire collection cycle
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	toSend := message{
		Agent: agent,
		Device: Device{
			Name: dev.name,
			Data: make(map[trait.Name]string),
		},
		Timestamp: time.Now(),
	}

	// fetch the data for each trait the device has and stick it in the same message
	// so we just have one message per device per interval
	for traitName, fetcher := range dev.traits {
		data, err := fetcher(ctx)
		if err != nil {
			// Check if it's a timeout error for better logging
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				a.Logger.Warn("timeout fetching trait data",
					zap.String("device", dev.name),
					zap.String("trait", string(traitName)),
					zap.Duration("timeout", fetchTimeout))
			} else {
				a.Logger.Error("failed to fetch trait data",
					zap.String("device", dev.name),
					zap.String("trait", string(traitName)),
					zap.Error(err))
			}
			continue
		}
		toSend.Device.Data[traitName] = data
	}

	// Include metadata if requested
	if includeMetadata && dev.metaData != nil {
		metadata, err := json.Marshal(dev.metaData)
		if err != nil {
			a.Logger.Error("failed to marshal device metadata",
				zap.String("device", dev.name),
				zap.Error(err))
		} else {
			toSend.Device.Data[trait.Metadata] = string(metadata)
		}
	}

	// Send the message if we have any data
	if len(toSend.Device.Data) > 0 {
		messagesCh <- toSend
	}
}

// getMeterInfo populates the device info map with Meter usageUnit and producedUnit.
func (a *AutoImpl) getMeterInfo(ctx context.Context, traitName trait.Name, devices map[string]*device) {
	for deviceName, dev := range devices {

		if _, hasMeterTrait := dev.traits[traitName]; !hasMeterTrait {
			continue
		}

		support, err := a.meterInfoClient.DescribeMeterReading(ctx, &gen.DescribeMeterReadingRequest{
			Name: deviceName,
		})
		if err != nil {
			a.Logger.Warn("failed to get meter info",
				zap.String("device", deviceName),
				zap.Error(err))
			continue
		}

		if support.UsageUnit != "" {
			dev.info["usageUnit"] = support.UsageUnit
		}
		if support.ProducedUnit != "" {
			dev.info["producedUnit"] = support.ProducedUnit
		}
	}
}
