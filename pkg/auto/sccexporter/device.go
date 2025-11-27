package sccexporter

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	meterpb "github.com/smart-core-os/sc-bos/pkg/gentrait/meter"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

// DataFetcher is a function that fetches device data for a specific trait
type DataFetcher func(ctx context.Context) ([]byte, error)

// MeterReadingWithUnits extends gen.MeterReading with gen.MeterReadingSupport information.
type MeterReadingWithUnits struct {
	*gen.MeterReading
	*gen.MeterReadingSupport
}

type device struct {
	name     string
	logger   *zap.Logger
	traits   map[trait.Name]DataFetcher
	info     map[trait.Name]proto.Message
	metaData *traits.Metadata
}

func newDevice(name string, logger *zap.Logger, metaData *traits.Metadata) *device {
	d := &device{
		name:     name,
		logger:   logger,
		traits:   make(map[trait.Name]DataFetcher),
		info:     make(map[trait.Name]proto.Message),
		metaData: metaData,
	}
	return d
}

func (d *device) getMeterData(ctx context.Context, meterClient gen.MeterApiClient) ([]byte, error) {
	reading, err := meterClient.GetMeterReading(ctx, &gen.GetMeterReadingRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get meter reading", zap.String("meter", d.name), zap.Error(err))
		return nil, err
	}

	// Marshal the meter reading to JSON using protojson
	readingBytes, err := protojson.Marshal(reading)
	if err != nil {
		d.logger.Error("failed to marshal meter reading", zap.String("meter", d.name), zap.Error(err))
		return nil, err
	}

	// If we don't have meter support info, just return the reading
	info, ok := d.info[meterpb.TraitName].(*gen.MeterReadingSupport)
	if !ok || info == nil {
		return readingBytes, nil
	}

	// Unmarshal to a map so we can add the unit fields
	var readingMap map[string]any
	if err := json.Unmarshal(readingBytes, &readingMap); err != nil {
		d.logger.Error("failed to unmarshal meter reading", zap.String("meter", d.name), zap.Error(err))
		return nil, err
	}

	// Add the unit fields from MeterReadingSupport
	if info.UsageUnit != "" {
		readingMap["usageUnit"] = info.UsageUnit
	}
	if info.ProducedUnit != "" {
		readingMap["producedUnit"] = info.ProducedUnit
	}

	// Marshal back to JSON with the added fields
	bytes, err := json.Marshal(readingMap)
	if err != nil {
		d.logger.Error("failed to marshal meter reading with units", zap.String("meter", d.name), zap.Error(err))
		return nil, err
	}

	return bytes, nil
}

func (d *device) getAirQualityData(ctx context.Context, airQualityClient traits.AirQualitySensorApiClient) ([]byte, error) {
	airQuality, err := airQualityClient.GetAirQuality(ctx, &traits.GetAirQualityRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get air quality", zap.String("airQualitySensor", d.name), zap.Error(err))
		return nil, err
	}

	aq, err := protojson.Marshal(airQuality)
	if err != nil {
		d.logger.Error("failed to marshal air quality", zap.String("airQualitySensor", d.name), zap.Error(err))
		return nil, err
	}

	return aq, nil
}

func (d *device) getOccupancyData(ctx context.Context, occupancyClient traits.OccupancySensorApiClient) ([]byte, error) {
	occupancy, err := occupancyClient.GetOccupancy(ctx, &traits.GetOccupancyRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get occupancy", zap.String("occupancySensor", d.name), zap.Error(err))
		return nil, err
	}

	o, err := protojson.Marshal(occupancy)
	if err != nil {
		d.logger.Error("failed to marshal occupancy", zap.String("occupancySensor", d.name), zap.Error(err))
		return nil, err
	}

	return o, nil
}

func (d *device) getAirTemperatureData(ctx context.Context, client traits.AirTemperatureApiClient) ([]byte, error) {
	airTemperature, err := client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get air temperature", zap.String("airTemperatureSensor", d.name), zap.Error(err))
		return nil, err
	}

	at, err := protojson.Marshal(airTemperature)
	if err != nil {
		d.logger.Error("failed to marshal air temperature", zap.String("airTemperatureSensor", d.name), zap.Error(err))
		return nil, err
	}

	return at, nil
}
