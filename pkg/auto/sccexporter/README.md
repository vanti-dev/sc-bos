## sccexporter automation 

The sccexporter automation provides a quick way of exporting trait data from devices to Smart Core Connect (SCC) via MQTT.
For all the traits provided in the config, the automation will discover all the devices which implement those traits,
grab the data at a given interval, and publish the data to SCC via MQTT.

### Supported Traits

The automation currently supports the following traits:
- **Meter** (`smartcore.bos.Meter`) - Energy/resource metering data
- **AirQualitySensor** (`smartcore.traits.AirQualitySensor`) - Air quality measurements
- **AirTemperature** (`smartcore.traits.AirTemperature`) - Temperature data
- **OccupancySensor** (`smartcore.traits.OccupancySensor`) - Occupancy information

### Data Collection

- **Trait Data**: Collected and published at the configured interval for all devices
- **Trait Info**: For Meter traits, additional info (usageUnit, producedUnit) is fetched once during startup and included in each data payload
- **Device Metadata**: Collected on startup and included in the data payload at a configurable interval (default: every 100 data collection cycles)

## Message Structure

The `message` structure represents the data format published to the MQTT broker:

```go
type message struct {
    Agent     string    `json:"agent"`
    Device    Device    `json:"device"`
    Timestamp time.Time `json:"timestamp"`
}
```

### Fields

- **Agent** (`string`): Identifier for the agent/system sending the data (configured in automation settings)
- **Device** (`Device`): Contains device information and trait data or metadata
- **Timestamp** (`time.Time`): When the message was created

### Device Structure

```go
type Device struct {
    Name string                       `json:"name"`
    Data map[string]json.RawMessage   `json:"data,omitempty"`
}
```

- **Name** (`string`): The unique Smart Core name of the device
- **Data** (`map[string]json.RawMessage`): Map where keys are trait names and values are JSON-encoded trait data or metadata
  - Keys are trait names (e.g., "smartcore.bos.Meter", "smartcore.traits.AirQualitySensor")
  - Values are `json.RawMessage` (byte arrays) containing JSON-encoded trait data
  - Protobuf messages are serialized using `protojson.Marshal`, which produces JSON with camelCase field names
  - Multiple traits can be included in a single message
  - Metadata is included every 100 intervals using the "metadata" key

### Data Messages

The `Data` map contains entries where each key is a trait name and the value is a `json.RawMessage` (byte array) containing JSON-encoded trait data.

#### Encoding Details

- **Protobuf Messages**: Most trait data is encoded using `protojson.Marshal`, which produces JSON with:
  - camelCase field names (e.g., `carbonDioxideLevel` instead of `carbon_dioxide_level`)
  - RFC 3339 timestamps (e.g., `2025-11-27T10:30:00Z`)
  - Numeric enums as string names (e.g., `"OCCUPIED"` instead of `1`)

- **Meter Data**: Uses a hybrid approach:
  1. The `MeterReading` protobuf message is serialized with `protojson.Marshal`
  2. Additional fields (`usageUnit`, `producedUnit`) from `MeterReadingSupport` are merged into the JSON
  3. Final result is re-encoded as standard JSON to combine both sources

- **Multiple Traits**: A single message can contain data for multiple traits that the device implements
- **Metadata**: Device metadata is included using the special "metadata" key at a configurable interval (default: every 100 data collection cycles)

### Example: Message with Multiple Traits

```json
{
  "agent": "van/uk/brum/ugs/building/scc-exporter",
  "device": {
    "name": "van/uk/brum/ugs/sensors/multi-sensor-01",
    "data": {
      "smartcore.traits.AirQualitySensor": {
        "carbonDioxideLevel": 450.5,
        "score": 75.5
      },
      "smartcore.traits.AirTemperature": {
        "ambientTemperature": {
          "valueCelsius": 22.5
        }
      },
      "smartcore.traits.OccupancySensor": {
        "state": "OCCUPIED",
        "peopleCount": 5,
        "stateChangeTime": "2025-11-27T10:25:00Z"
      }
    }
  },
  "timestamp": "2025-11-27T10:30:00Z"
}
```

### Example: Meter Message with Units

```json
{
  "agent": "van/uk/brum/ugs/building/scc-exporter",
  "device": {
    "name": "van/uk/brum/ugs/meters/elec-main",
    "data": {
      "smartcore.bos.Meter": {
        "usage": 123.45,
        "produced": 67.89,
        "startTime": "2025-11-27T09:30:00Z",
        "endTime": "2025-11-27T10:30:00Z",
        "usageUnit": "kWh",
        "producedUnit": "kWh"
      }
    }
  },
  "timestamp": "2025-11-27T10:30:00Z"
}
```

### Notes

- The `Data` field is a map where keys are trait names (or "metadata" for device metadata)
- Values in the `Data` map are `json.RawMessage` (byte arrays) containing JSON-encoded data
- **Protobuf Encoding**: Trait data from protobuf messages uses `protojson.Marshal`, which produces:
  - camelCase field names (e.g., `carbonDioxideLevel`, `ambientTemperature`)
  - String representations for enums (e.g., `"OCCUPIED"`)
  - RFC 3339 formatted timestamps (e.g., `"2025-11-27T10:30:00Z"`)
- **Meter Data Special Handling**: 
  - Meter readings combine data from two protobuf messages (`MeterReading` and `MeterReadingSupport`)
  - The `usageUnit` and `producedUnit` fields from `MeterReadingSupport` are merged into the reading JSON
  - This allows meter data to include both the readings and their units in a single JSON object
- **Multiple traits from the same device are combined into a single message** - all trait data appears in the same `Data` map
- Trait names like "smartcore.bos.Meter" or "smartcore.traits.AirQualitySensor" are used as keys
- **Metadata is included in the same message as trait data** at a configurable interval (default: every 100 data collection cycles) using the "metadata" key
- This reduces MQTT traffic by consolidating all device information into fewer messages

## Configuration

### Example Configuration

```json
{
  "type": "sccexporter",
  "traits": ["smartcore.bos.Meter", "smartcore.traits.AirQualitySensor"],
  "fetchTimeout": "5s",
  "mqtt": {
    "agent": "van/uk/site/building/exporter",
    "host": "ssl://mqtt.example.com:8883",
    "topic": "scc/data",
    "clientId": "scc-exporter-1",
    "clientCertPath": "/path/to/client.crt",
    "clientKeyPath": "/path/to/client.key",
    "caCertPath": "/path/to/ca.crt",
    "connectTimeout": "5s",
    "publishTimeout": "5s",
    "sendInterval": "*/15 * * * *",
    "metadataInterval": 100,
    "qos": 1
  }
}
```

### Configuration Options

#### General Settings

- **fetchTimeout** (duration, optional): Maximum time to wait for a single device's trait data fetch. If a device takes longer than this, the fetch is cancelled and that device is skipped for this cycle. This prevents slow or hanging devices from blocking the entire collection cycle. Default: "5s"

#### Mqtt Settings

- **agent** (string, required): Identifier for the agent/system sending the data
- **host** (string, required): MQTT broker address (e.g., "ssl://mqtt.example.com:8883")
- **topic** (string, required): MQTT topic to publish messages to
- **clientId** (string, required): Unique client identifier for MQTT connection
- **clientCertPath** (string, required): Path to client certificate file for TLS
- **clientKeyPath** (string, required): Path to client key file for TLS
- **caCertPath** (string, required): Path to CA certificate file for TLS
- **connectTimeout** (duration, optional): Timeout for connecting to MQTT broker, default: "5s"
- **publishTimeout** (duration, optional): Timeout for publishing to MQTT, default: "5s"
- **sendInterval** (schedule, optional): Cron schedule for data collection, default: "*/15 * * * *" (every 15 minutes)
- **metadataInterval** (int, optional): Include metadata every N data collection cycles, default: 100
- **qos** (int, optional): MQTT Quality of Service level (0, 1, or 2), default: 1



