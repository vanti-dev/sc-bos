// Package azureiot implements data upload to Azure IoT Hub.
// The driver pulls (or polls) for data from configured Smart Core devices, encodes the data as JSON, and uploads it to Azure IoT Hub via it's MQTT API.
// Published data is the protojson encoding of the PullTraitResponse message type, for example [traits.PullAirQualityResponse].
// See [traits.go] for supported traits.
//
// Azure IoT Hub has it's own concept of devices that are separate from Smart Core names.
// When configuring the driver you start by defining the IoT Hub device to publish to,
// then define which Smart Core devices publish data to that IoT hub device.
// Finally you define which traits to pull from each Smart Core device.
//
// There are two ways to configure the IoT Hub device: manually via the Hub console, or dynamically via the Device Provisioning Service (DPS).
// This driver supports both methods.
// The usual way to configure the IoT Hub device is via a Connection String which can be retrieved from the Hub console.
// See [config.go] for details.
//
// To reduce configuration overhead in both Smart Core and Azure IoT Hub, it is recommended to setup a single IoT Hub
// device for each Smart Core node and publish multiple Smart Core devices data to it.
package azureiot
