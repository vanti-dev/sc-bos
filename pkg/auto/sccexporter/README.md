## sccexporter automation 

The sccexporter automation provides a quick way of exporting trait data from devices to SCC via MQTT.
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
    Name string            `json:"name"`
    Data map[string]string `json:"data,omitempty"`
}
```

- **Name** (`string`): The unique Smart Core name of the device
- **Data** (`map[string]string`): Map where keys are trait names and values are JSON-encoded trait data or metadata
  - Keys are trait names (e.g., "smartcore.bos.Meter", "smartcore.traits.AirQualitySensor")
  - Values are JSON-encoded trait readings
  - Multiple traits can be included in a single message
  - Metadata is included every 100 intervals

### Data Messages

The `Data` map contains entries where each key is a trait name and the value is a JSON string with that trait's current readings. 

- **Multiple Traits**: A single message can contain data for multiple traits that the device implements
- **Meter Traits**: The JSON includes additional fields from trait info (usageUnit, producedUnit)
- **Metadata**: Included using the special "metadata" key at a configurable interval (default: every 100 data collection cycles)

### Example: Message

```json
{
  "agent": "van/uk/brum/ugs/building/scc-exporter",
  "device": {
    "name": "van/uk/brum/ugs/sensors/multi-sensor-01",
    "data": {
      "smartcore.traits.AirQualitySensor": "{\"carbonDioxideLevel\":450.5,\"score\":75.5}",
      "smartcore.traits.AirTemperature": "{\"ambientTemperature\":{\"valueCelsius\":22.5}}",
      "smartcore.traits.OccupancySensor": "{\"state\":\"OCCUPIED\",\"peopleCount\":5}"
    }
  },
  "timestamp": "2025-11-25T10:30:00Z"
}
```

### Notes

- The `Data` field is a map where keys are trait names (or "metadata" for device metadata)
- Values in the `Data` map are JSON-encoded strings (double-encoded in the final JSON payload)
- **Multiple traits from the same device are combined into a single message** - all trait data appears in the same `Data` map
- Trait names like "smartcore.bos.Meter" or "smartcore.traits.AirQualitySensor" are used as keys
- Meter readings automatically include `usageUnit` and `producedUnit` fields from the meter's trait info
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
    "clientCert": "/path/to/client.crt",
    "clientKey": "/path/to/client.key",
    "caCert": "/path/to/ca.crt",
    "sendInterval": "*/15 * * * *",
    "metadataInterval": 100,
    "publishTimeout": "5s",
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
- **clientCert** (string, required): Path to client certificate file for TLS
- **clientKey** (string, required): Path to client key file for TLS
- **caCert** (string, required): Path to CA certificate file for TLS
- **sendInterval** (schedule, optional): Cron schedule for data collection, default: "*/15 * * * *" (every 15 minutes)
- **metadataInterval** (int, optional): Include metadata every N data collection cycles, default: 100
- **publishTimeout** (duration, optional): Timeout for publishing to MQTT, default: "5s"
- **qos** (int, optional): MQTT Quality of Service level (0, 1, or 2), default: 1



