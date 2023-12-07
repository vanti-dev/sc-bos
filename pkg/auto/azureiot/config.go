package azureiot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

const (
	DefaultPollInterval = 5 * time.Minute
	DefaultBacklogSize  = 1000 // in messages per device
)

type Config struct {
	// Azure IoT Hub dps credentials, used for devices where DeviceConfig.ConnectionString == ""
	GroupKey     string `json:"groupKey"`
	GroupKeyFile string `json:"groupKeyFile"`
	IDScope      string `json:"idScope"`

	Devices []DeviceConfig `json:"devices"`

	// When a device doesn't support pullTraits, use this poll interval.
	// Defaults to DefaultPollInterval.
	PollInterval *jsontypes.Duration `json:"pollInterval"`
}

// DeviceConfig represents an Azure IoT Hub device.
// This device may map to one or more Smart Core devices.
type DeviceConfig struct {
	SCDeviceConfig                  // config used when an IoT Hub device maps to one Smart Core device
	Children       []SCDeviceConfig `json:"children,omitempty"` // config used when an IoT Hub device maps to multiple Smart Core devices
	// todo: DiscoverDevices bool // if true, discover devices from Smart Core and use them as the Children list

	// Azure IoT Hub connection details, one of these must be provided
	RegistrationID   string `json:"registrationID,omitempty"`   // Device Provisioning Service registration ID - usually becomes the device name in the cloud
	ConnectionString string `json:"connectionString,omitempty"` // Use a connection string to directly connect to IoT Hub, bypassing Device Provisioning Service
}

// SCDeviceConfig represents a Smart Core device information will be pulled from.
type SCDeviceConfig struct {
	Name                string       `json:"name"`                // Smart Core name to monitor
	Traits              []trait.Name `json:"traits"`              // Traits to pullTraits from
	IgnoreUnknownTraits bool         `json:"ignoreUnknownTraits"` // Ignore unknown traits instead of erroring.
	// todo: empty traits list means the intersection of device and auto supported traits

	PollInterval *jsontypes.Duration `json:"pollInterval"` // Defaults to Config.PollInterval
}

func ParseConfig(jsonBytes []byte) (Config, error) {
	cfg := DefaultConfig()
	err := json.Unmarshal(jsonBytes, &cfg)

	if i := cfg.PollInterval.Or(DefaultPollInterval); i < minPollInterval {
		return cfg, fmt.Errorf("pollInterval %v must be at least %v", i, minPollInterval)
	}

	for i, deviceCfg := range cfg.Devices {
		if deviceCfg.ConnectionString == "" {
			// 	needsGroupKey = true

			// if the group key is used, then an ID scope is also mandatory
			if cfg.IDScope == "" {
				return cfg, fmt.Errorf("id scope is required when using group keys")
			}
		}
		if deviceCfg.PollInterval == nil {
			deviceCfg.PollInterval = cfg.PollInterval
		}
		for i, child := range deviceCfg.Children {
			if child.PollInterval == nil {
				child.PollInterval = deviceCfg.PollInterval
			}
			deviceCfg.Children[i] = child
		}

		cfg.Devices[i] = deviceCfg
	}

	return cfg, err
}

func DefaultConfig() Config {
	return Config{
		PollInterval: &jsontypes.Duration{Duration: DefaultPollInterval},
	}
}
