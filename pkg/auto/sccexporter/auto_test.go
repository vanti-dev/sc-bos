package sccexporter

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-bos/internal/manage/devices"
	"github.com/smart-core-os/sc-bos/pkg/auto"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	meterpb "github.com/smart-core-os/sc-bos/pkg/gentrait/meter"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensorpb"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperaturepb"
	"github.com/smart-core-os/sc-golang/pkg/trait/metadatapb"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
)

func TestMetadata(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	root := node.New("metadata")

	sccexporter := &AutoImpl{
		Services: auto.Services{
			Logger:  logger,
			Node:    root,
			Devices: gen.NewDevicesApiClient(root.ClientConn()),
		},
	}

	metadata := &traits.Metadata{
		Name: "foo",
		Appearance: &traits.Metadata_Appearance{
			Title:       "Foo Device",
			Description: "A device for testing metadata",
		},
		Location: &traits.Metadata_Location{
			Floor: "1",
			Zone:  "bar",
		},
	}

	metaModel := metadatapb.NewModel(resource.WithInitialValue(metadata))
	modelServer := metadatapb.NewModelServer(metaModel)
	metaClient := node.WithClients(metadatapb.WrapApi(modelServer))
	root.Announce("foo", node.HasTrait(trait.Metadata, metaClient))

	sccexporter.initialiseClients(root)

	// Create a test device
	dev := newDevice("foo", logger, metadata)

	// Create a channel to receive messages
	messagesCh := make(chan message, 1)

	// Test fetchAndPublishDeviceData with metadata
	agent := "test-agent"

	sccexporter.fetchAndPublishDeviceData(context.Background(), dev, agent, messagesCh, true, 30*time.Second)

	// Verify the message was sent
	require.Len(t, messagesCh, 1)

	// Read the message from the channel
	msg := <-messagesCh

	// Verify message structure
	require.Equal(t, agent, msg.Agent)
	require.Equal(t, "foo", msg.Device.Name)
	require.NotEmpty(t, msg.Device.Data)
	require.Contains(t, msg.Device.Data, trait.Metadata)

	// Verify the metadata JSON can be unmarshalled and contains expected data
	var receivedMetadata traits.Metadata
	err = json.Unmarshal([]byte(msg.Device.Data[trait.Name(trait.Metadata)]), &receivedMetadata)
	require.NoError(t, err)

	require.Equal(t, "foo", receivedMetadata.Name)
	require.Equal(t, "Foo Device", receivedMetadata.Appearance.Title)
	require.Equal(t, "A device for testing metadata", receivedMetadata.Appearance.Description)
	require.Equal(t, "1", receivedMetadata.Location.Floor)
	require.Equal(t, "bar", receivedMetadata.Location.Zone)
}

func TestFetchAndPublishDeviceData(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	ctx := context.Background()
	agent := "test-agent"

	t.Run("single trait with data", func(t *testing.T) {
		// Create a device with a mock trait fetcher
		dev := newDevice("test-device", logger, nil)

		// Use JSON string for test data
		testDataJSON := `{
			"value1": 42.5,
			"value2": "test-value"
		}`

		dev.traits[trait.Name("test-trait")] = func(ctx context.Context) (string, error) {
			return testDataJSON, nil
		}

		// Create message channel
		messagesCh := make(chan message, 1)

		// Create AutoImpl instance
		a := &AutoImpl{
			Services: auto.Services{
				Logger: logger,
			},
		}

		// Call the function
		a.fetchAndPublishDeviceData(ctx, dev, agent, messagesCh, false, 30*time.Second)

		// Verify message was sent
		require.Len(t, messagesCh, 1)
		msg := <-messagesCh

		// Verify message structure
		require.Equal(t, agent, msg.Agent)
		require.Equal(t, "test-device", msg.Device.Name)
		require.NotEmpty(t, msg.Device.Data)
		require.Contains(t, msg.Device.Data, trait.Name("test-trait"))

		// Verify trait data by unmarshaling JSON
		var data map[string]any
		err = json.Unmarshal([]byte(msg.Device.Data[trait.Name("test-trait")]), &data)
		require.NoError(t, err)
		require.Equal(t, 42.5, data["value1"])
		require.Equal(t, "test-value", data["value2"])
	})

	t.Run("multiple traits with data", func(t *testing.T) {
		// Create a device with multiple mock trait fetchers
		dev := newDevice("multi-trait-device", logger, nil)

		// Add first trait
		trait1JSON := `{
			"measurement": 100.0,
			"status": "active"
		}`
		dev.traits[trait.Name("trait1")] = func(ctx context.Context) (string, error) {
			return trait1JSON, nil
		}

		// Add second trait
		trait2JSON := `{
			"temperature": 22.5,
			"humidity": 45.0
		}`
		dev.traits[trait.Name("trait2")] = func(ctx context.Context) (string, error) {
			return trait2JSON, nil
		}

		// Add third trait
		trait3JSON := `{
			"value": 42,
			"unit": "kWh"
		}`
		dev.traits[trait.Name("trait3")] = func(ctx context.Context) (string, error) {
			return trait3JSON, nil
		}

		// Create message channel
		messagesCh := make(chan message, 1)

		// Create AutoImpl instance
		a := &AutoImpl{
			Services: auto.Services{
				Logger: logger,
			},
		}

		// Call the function
		a.fetchAndPublishDeviceData(ctx, dev, agent, messagesCh, false, 30*time.Second)

		// Verify message was sent
		require.Len(t, messagesCh, 1)
		msg := <-messagesCh

		// Verify all three traits are present in the data map
		require.Len(t, msg.Device.Data, 3)
		require.Contains(t, msg.Device.Data, trait.Name("trait1"))
		require.Contains(t, msg.Device.Data, trait.Name("trait2"))
		require.Contains(t, msg.Device.Data, trait.Name("trait3"))

		// Verify trait1 data
		var data1 map[string]any
		err = json.Unmarshal([]byte(msg.Device.Data[trait.Name("trait1")]), &data1)
		require.NoError(t, err)
		require.Equal(t, float64(100.0), data1["measurement"])
		require.Equal(t, "active", data1["status"])

		// Verify trait2 data
		var data2 map[string]any
		err = json.Unmarshal([]byte(msg.Device.Data[trait.Name("trait2")]), &data2)
		require.NoError(t, err)
		require.Equal(t, float64(22.5), data2["temperature"])
		require.Equal(t, float64(45.0), data2["humidity"])

		// Verify trait3 data
		var data3 map[string]any
		err = json.Unmarshal([]byte(msg.Device.Data[trait.Name("trait3")]), &data3)
		require.NoError(t, err)
		require.Equal(t, float64(42), data3["value"])
		require.Equal(t, "kWh", data3["unit"])
	})

	t.Run("trait fetcher returns error", func(t *testing.T) {
		// Create a device with a failing trait fetcher
		dev := newDevice("error-device", logger, nil)
		dev.traits[trait.Name("failing-trait")] = func(ctx context.Context) (string, error) {
			return "", context.DeadlineExceeded
		}

		// Use JSON string for working trait
		workingDataJSON := `{
			"value": 123
		}`

		dev.traits[trait.Name("working-trait")] = func(ctx context.Context) (string, error) {
			return workingDataJSON, nil
		}

		// Create message channel
		messagesCh := make(chan message, 1)

		// Create AutoImpl instance
		a := &AutoImpl{
			Services: auto.Services{
				Logger: logger,
			},
		}

		// Call the function
		a.fetchAndPublishDeviceData(ctx, dev, agent, messagesCh, false, 30*time.Second)

		// Verify message was sent (only with working trait data)
		require.Len(t, messagesCh, 1)
		msg := <-messagesCh

		// Verify data is present from the working trait
		require.NotEmpty(t, msg.Device.Data)
		require.Contains(t, msg.Device.Data, trait.Name("working-trait"))

		// Verify the data by unmarshaling JSON
		var data map[string]any
		err = json.Unmarshal([]byte(msg.Device.Data[trait.Name("working-trait")]), &data)
		require.NoError(t, err)
		require.Equal(t, float64(123), data["value"]) // JSON numbers are float64
	})

	t.Run("all traits fail - no message sent", func(t *testing.T) {
		// Create a device where all traits fail
		dev := newDevice("all-fail-device", logger, nil)
		dev.traits[trait.Name("trait1")] = func(ctx context.Context) (string, error) {
			return "", context.DeadlineExceeded
		}
		dev.traits[trait.Name("trait2")] = func(ctx context.Context) (string, error) {
			return "", context.Canceled
		}

		// Create message channel
		messagesCh := make(chan message, 1)

		// Create AutoImpl instance
		a := &AutoImpl{
			Services: auto.Services{
				Logger: logger,
			},
		}

		// Call the function
		a.fetchAndPublishDeviceData(ctx, dev, agent, messagesCh, false, 30*time.Second)

		// Verify no message was sent
		require.Len(t, messagesCh, 0)
	})

	t.Run("device with no traits", func(t *testing.T) {
		// Create a device with no traits
		dev := newDevice("empty-device", logger, nil)

		// Create message channel
		messagesCh := make(chan message, 1)

		// Create AutoImpl instance
		a := &AutoImpl{
			Services: auto.Services{
				Logger: logger,
			},
		}

		// Call the function
		a.fetchAndPublishDeviceData(ctx, dev, agent, messagesCh, false, 30*time.Second)

		// Verify no message was sent
		require.Len(t, messagesCh, 0)
	})

	t.Run("timeout on slow device", func(t *testing.T) {
		// Create a device with a slow trait fetcher
		dev := newDevice("slow-device", logger, nil)
		dev.traits[trait.Name("slow-trait")] = func(ctx context.Context) (string, error) {
			// Simulate a slow device that takes longer than timeout
			select {
			case <-time.After(2 * time.Second):
				return `{"value": "should-not-see-this"}`, nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		// Create message channel
		messagesCh := make(chan message, 1)

		// Create AutoImpl instance
		a := &AutoImpl{
			Services: auto.Services{
				Logger: logger,
			},
		}

		// Call the function with short timeout (100ms)
		a.fetchAndPublishDeviceData(ctx, dev, agent, messagesCh, false, 100*time.Millisecond)

		// Verify no message was sent due to timeout
		require.Len(t, messagesCh, 0)
	})
}

func TestGetMeterDeviceAndData(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	root := node.New("metadata")

	startTime := time.Now().Add(-time.Hour)
	endTime := time.Now()

	meterReading := &gen.MeterReading{
		Usage:     123.45,
		StartTime: timestamppb.New(startTime),
		EndTime:   timestamppb.New(endTime),
		Produced:  67.89,
	}

	devicesApi := devices.NewServer(root)
	meterModel := meterpb.NewModel(resource.WithInitialValue(meterReading))
	modelServer := meterpb.NewModelServer(meterModel)
	meterClient := node.WithClients(gen.WrapMeterApi(modelServer))
	root.Announce("foo",
		node.HasTrait(meterpb.TraitName, meterClient),
		node.HasServices(root.ClientConn(), gen.DevicesApi_ServiceDesc),
	)

	sccexporter := &AutoImpl{
		Services: auto.Services{
			Logger:  logger,
			Node:    root,
			Devices: gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, devicesApi)),
		},
	}
	sccexporter.initialiseClients(root)

	allDevices := make(map[string]*device)
	err = sccexporter.getAllTraitImplementors(context.Background(), meterpb.TraitName, allDevices)
	require.NoError(t, err)

	require.Len(t, allDevices, 1)
	dev, exists := allDevices["foo"]
	require.True(t, exists)
	require.Equal(t, "foo", dev.name)

	res := allDevices["foo"].traits
	require.Len(t, res, 1)
	traitData, err := res[meterpb.TraitName](context.Background())
	require.NoError(t, err)

	reading := gen.MeterReading{}
	err = json.Unmarshal([]byte(traitData), &reading)
	require.NoError(t, err)

	require.Equal(t, meterReading.Usage, reading.Usage)
	require.Equal(t, meterReading.Produced, reading.Produced)
	require.Equal(t, meterReading.StartTime.AsTime(), reading.StartTime.AsTime())
	require.Equal(t, meterReading.EndTime.AsTime(), reading.EndTime.AsTime())
}

func TestGetMeterDeviceAndDataWithInfo(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	root := node.New("meter")

	startTime := time.Now().Add(-time.Hour)
	endTime := time.Now()

	meterReading := &gen.MeterReading{
		Usage:     123.45,
		StartTime: timestamppb.New(startTime),
		EndTime:   timestamppb.New(endTime),
		Produced:  67.89,
	}

	meterInfo := &gen.MeterReadingSupport{
		UsageUnit:    "kWh",
		ProducedUnit: "kWh",
	}

	devicesApi := devices.NewServer(root)
	meterModel := meterpb.NewModel(resource.WithInitialValue(meterReading))
	modelServer := meterpb.NewModelServer(meterModel)
	infoServer := &meterpb.InfoServer{MeterReading: meterInfo}
	meterClient := node.WithClients(
		gen.WrapMeterApi(modelServer),
		gen.WrapMeterInfo(infoServer),
	)
	root.Announce("foo",
		node.HasTrait(meterpb.TraitName, meterClient),
		node.HasServices(root.ClientConn(), gen.DevicesApi_ServiceDesc),
	)

	sccexporter := &AutoImpl{
		Services: auto.Services{
			Logger:  logger,
			Node:    root,
			Devices: gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, devicesApi)),
		},
	}
	sccexporter.initialiseClients(root)

	allDevices := make(map[string]*device)
	err = sccexporter.getAllTraitImplementors(context.Background(), meterpb.TraitName, allDevices)
	require.NoError(t, err)

	require.Len(t, allDevices, 1)
	dev, exists := allDevices["foo"]
	require.True(t, exists)
	require.Equal(t, "foo", dev.name)

	// Call getMeterInfo to populate usageUnit and producedUnit
	sccexporter.getMeterInfo(context.Background(), meterpb.TraitName, allDevices)

	// Verify that usageUnit and producedUnit are populated
	require.Equal(t, "kWh", dev.info["usageUnit"])
	require.Equal(t, "kWh", dev.info["producedUnit"])

	// Now fetch the meter data and verify it includes the units
	res := allDevices["foo"].traits
	require.Len(t, res, 1)
	traitData, err := res[meterpb.TraitName](context.Background())
	require.NoError(t, err)

	// Unmarshal to map to check for added fields
	var readingMap map[string]any
	err = json.Unmarshal([]byte(traitData), &readingMap)
	require.NoError(t, err)

	// Verify the meter data includes the reading values
	require.Equal(t, meterReading.Usage, float32(readingMap["usage"].(float64)))
	require.Equal(t, meterReading.Produced, float32(readingMap["produced"].(float64)))

	// Verify the units were added
	require.Equal(t, "kWh", readingMap["usageUnit"])
	require.Equal(t, "kWh", readingMap["producedUnit"])
}

func TestGetAirQualityDeviceAndData(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	root := node.New("airquality")

	co2Level := float32(450.5)
	score := float32(75.5)

	airQuality := &traits.AirQuality{
		CarbonDioxideLevel: &co2Level,
		Score:              &score,
	}

	devicesApi := devices.NewServer(root)
	airQualityModel := airqualitysensorpb.NewModel()
	_, err = airQualityModel.UpdateAirQuality(airQuality)
	require.NoError(t, err)
	modelServer := airqualitysensorpb.NewModelServer(airQualityModel)
	airQualityClient := node.WithClients(airqualitysensorpb.WrapApi(modelServer))
	root.Announce("foo",
		node.HasTrait(trait.AirQualitySensor, airQualityClient),
		node.HasServices(root.ClientConn(), gen.DevicesApi_ServiceDesc),
	)

	sccexporter := &AutoImpl{
		Services: auto.Services{
			Logger:  logger,
			Node:    root,
			Devices: gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, devicesApi)),
		},
	}
	sccexporter.initialiseClients(root)

	allDevices := make(map[string]*device)
	err = sccexporter.getAllTraitImplementors(context.Background(), trait.AirQualitySensor, allDevices)
	require.NoError(t, err)

	require.Len(t, allDevices, 1)
	dev, exists := allDevices["foo"]
	require.True(t, exists)
	require.Equal(t, "foo", dev.name)

	res := allDevices["foo"].traits
	require.Len(t, res, 1)
	traitData, err := res[trait.AirQualitySensor](context.Background())
	require.NoError(t, err)

	receivedAirQuality := traits.AirQuality{}
	err = json.Unmarshal([]byte(traitData), &receivedAirQuality)
	require.NoError(t, err)

	require.Equal(t, *airQuality.CarbonDioxideLevel, *receivedAirQuality.CarbonDioxideLevel)
	require.Equal(t, *airQuality.Score, *receivedAirQuality.Score)
}

func TestGetAirTemperatureDeviceAndData(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	root := node.New("airtemperature")

	celsius := 22.5

	airTemperature := &traits.AirTemperature{
		AmbientTemperature: &types.Temperature{ValueCelsius: celsius},
	}

	devicesApi := devices.NewServer(root)
	airTemperatureModel := airtemperaturepb.NewModel()
	_, err = airTemperatureModel.UpdateAirTemperature(airTemperature)
	require.NoError(t, err)
	modelServer := airtemperaturepb.NewModelServer(airTemperatureModel)
	airTemperatureClient := node.WithClients(airtemperaturepb.WrapApi(modelServer))
	root.Announce("foo",
		node.HasTrait(trait.AirTemperature, airTemperatureClient),
		node.HasServices(root.ClientConn(), gen.DevicesApi_ServiceDesc),
	)

	sccexporter := &AutoImpl{
		Services: auto.Services{
			Logger:  logger,
			Node:    root,
			Devices: gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, devicesApi)),
		},
	}
	sccexporter.initialiseClients(root)

	allDevices := make(map[string]*device)
	err = sccexporter.getAllTraitImplementors(context.Background(), trait.AirTemperature, allDevices)
	require.NoError(t, err)

	require.Len(t, allDevices, 1)
	dev, exists := allDevices["foo"]
	require.True(t, exists)
	require.Equal(t, "foo", dev.name)

	res := allDevices["foo"].traits
	require.Len(t, res, 1)
	traitData, err := res[trait.AirTemperature](context.Background())
	require.NoError(t, err)

	receivedAirTemperature := traits.AirTemperature{}
	err = json.Unmarshal([]byte(traitData), &receivedAirTemperature)
	require.NoError(t, err)

	require.Equal(t, airTemperature.AmbientTemperature.ValueCelsius, receivedAirTemperature.AmbientTemperature.ValueCelsius)
}

func TestGetOccupancyDeviceAndData(t *testing.T) {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	root := node.New("occupancy")

	stateChangeTime := time.Now().Add(-5 * time.Minute)

	occupancy := &traits.Occupancy{
		State:           traits.Occupancy_OCCUPIED,
		PeopleCount:     5,
		StateChangeTime: timestamppb.New(stateChangeTime),
	}

	devicesApi := devices.NewServer(root)
	occupancyModel := occupancysensorpb.NewModel()
	_, err = occupancyModel.SetOccupancy(occupancy)
	require.NoError(t, err)
	modelServer := occupancysensorpb.NewModelServer(occupancyModel)
	occupancyClient := node.WithClients(occupancysensorpb.WrapApi(modelServer))
	root.Announce("foo",
		node.HasTrait(trait.OccupancySensor, occupancyClient),
		node.HasServices(root.ClientConn(), gen.DevicesApi_ServiceDesc),
	)

	sccexporter := &AutoImpl{
		Services: auto.Services{
			Logger:  logger,
			Node:    root,
			Devices: gen.NewDevicesApiClient(wrap.ServerToClient(gen.DevicesApi_ServiceDesc, devicesApi)),
		},
	}
	sccexporter.initialiseClients(root)

	allDevices := make(map[string]*device)
	err = sccexporter.getAllTraitImplementors(context.Background(), trait.OccupancySensor, allDevices)
	require.NoError(t, err)

	require.Len(t, allDevices, 1)
	dev, exists := allDevices["foo"]
	require.True(t, exists)
	require.Equal(t, "foo", dev.name)

	res := allDevices["foo"].traits
	require.Len(t, res, 1)
	traitData, err := res[trait.OccupancySensor](context.Background())
	require.NoError(t, err)

	receivedOccupancy := traits.Occupancy{}
	err = json.Unmarshal([]byte(traitData), &receivedOccupancy)
	require.NoError(t, err)

	require.Equal(t, occupancy.State, receivedOccupancy.State)
	require.Equal(t, occupancy.PeopleCount, receivedOccupancy.PeopleCount)
	require.Equal(t, occupancy.StateChangeTime.AsTime(), receivedOccupancy.StateChangeTime.AsTime())
}
