package azureiot

import (
	"encoding/json"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Config struct {
	GroupKey     string             `json:"groupKey"`
	GroupKeyFile string             `json:"groupKeyFile"`
	IDScope      string             `json:"idScope"`
	Devices      []DeviceConfig     `json:"devices"`
	PollInterval jsontypes.Duration `json:"pollInterval"`
}

type DeviceConfig struct {
	Name             string `json:"name"`             // Smart Core name to monitor
	RegistrationID   string `json:"registrationID"`   // Device Provisioning Service registration ID - usually becomes the device name in the cloud
	ConnectionString string `json:"connectionString"` // Use a connection string to directly connect to IoT Hub, bypassing Device Provisioning Service
}

func ParseConfig(jsonBytes []byte) (Config, error) {
	parsed := DefaultConfig()
	err := json.Unmarshal(jsonBytes, &parsed)
	return parsed, err
}

func DefaultConfig() Config {
	return Config{
		PollInterval: jsontypes.Duration{Duration: 5 * time.Minute},
	}
}
