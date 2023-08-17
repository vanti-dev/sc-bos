package azureiot

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

type Config struct {
	GroupKey     string         `json:"groupKey"`
	GroupKeyFile string         `json:"groupKeyFile"`
	Devices      []DeviceConfig `json:"devices"`
}

type DeviceConfig struct {
	Name             string       `json:"name"`             // Smart Core name to monitor
	RegistrationID   string       `json:"registrationID"`   // Device Provisioning Service registration ID - usually becomes the device name in the cloud
	ConnectionString string       `json:"connectionString"` // Use a connection string to directly connect to IoT Hub, bypassing Device Provisioning Service
	Traits           []trait.Name `json:"traits"`           // Which traits to poll and report data on
}
