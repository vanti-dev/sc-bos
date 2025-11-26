package sccexporter

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

// DataFetcher is a function that fetches device data for a specific trait
type DataFetcher func(ctx context.Context) (string, error)

// MeterReadingWithUnits extends gen.MeterReading with unit information.
// This allows us to efficiently add usageUnit and producedUnit fields
// without multiple marshal/unmarshal cycles.
type MeterReadingWithUnits struct {
	*gen.MeterReading
	UsageUnit    string `json:"usageUnit,omitempty"`
	ProducedUnit string `json:"producedUnit,omitempty"`
}

type device struct {
	name     string
	logger   *zap.Logger
	traits   map[trait.Name]DataFetcher
	info     map[string]string
	metaData *traits.Metadata
}

func newDevice(name string, logger *zap.Logger, metaData *traits.Metadata) *device {
	d := &device{
		name:     name,
		logger:   logger,
		traits:   make(map[trait.Name]DataFetcher),
		info:     make(map[string]string),
		metaData: metaData,
	}
	return d
}

func (d *device) getMeterData(ctx context.Context, meterClient gen.MeterApiClient) (string, error) {
	reading, err := meterClient.GetMeterReading(ctx, &gen.GetMeterReadingRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get meter reading", zap.String("meter", d.name), zap.Error(err))
		return "", err
	}

	// Use struct embedding to efficiently add unit fields in a single marshal operation
	data := &MeterReadingWithUnits{
		MeterReading: reading,
		UsageUnit:    d.info["usageUnit"],
		ProducedUnit: d.info["producedUnit"],
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		d.logger.Error("failed to marshal meter reading", zap.String("meter", d.name), zap.Error(err))
		return "", err
	}

	return string(bytes), nil
}

func (d *device) getAirQualityData(ctx context.Context, airQualityClient traits.AirQualitySensorApiClient) (string, error) {
	airQuality, err := airQualityClient.GetAirQuality(ctx, &traits.GetAirQualityRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get air quality", zap.String("airQualitySensor", d.name), zap.Error(err))
		return "", err
	}

	aq, err := json.Marshal(airQuality)
	if err != nil {
		d.logger.Error("failed to marshal air quality", zap.String("airQualitySensor", d.name), zap.Error(err))
		return "", err
	}

	return string(aq), nil
}

func (d *device) getOccupancyData(ctx context.Context, occupancyClient traits.OccupancySensorApiClient) (string, error) {
	occupancy, err := occupancyClient.GetOccupancy(ctx, &traits.GetOccupancyRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get occupancy", zap.String("occupancySensor", d.name), zap.Error(err))
		return "", err
	}

	o, err := json.Marshal(occupancy)
	if err != nil {
		d.logger.Error("failed to marshal occupancy", zap.String("occupancySensor", d.name), zap.Error(err))
		return "", err
	}

	return string(o), nil
}

func (d *device) getAirTemperatureData(ctx context.Context, client traits.AirTemperatureApiClient) (string, error) {
	airTemperature, err := client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{
		Name: d.name,
	})
	if err != nil {
		d.logger.Error("failed to get air temperature", zap.String("airTemperatureSensor", d.name), zap.Error(err))
		return "", err
	}

	at, err := json.Marshal(airTemperature)
	if err != nil {
		d.logger.Error("failed to marshal air temperature", zap.String("airTemperatureSensor", d.name), zap.Error(err))
		return "", err
	}

	return string(at), nil
}
