package wiser_knx

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

func DefaultConfig() Config {
	return Config{
		Poll:         10 * time.Second,
		Username:     "remote",
		PasswordFile: "/run/secrets/wiser-knx-password",
	}
}

func ParseConfig(raw []byte) (Config, error) {
	parsed := DefaultConfig()
	err := json.Unmarshal(raw, &parsed)
	return parsed, err
}

// Config is the configuration for the Wiser for KNX driver.
type Config struct {
	driver.BaseConfig

	// The IP address of the Wiser for KNX controller.
	Host string `json:"host"`

	// The username to use when connecting to the Wiser for KNX controller (default: "remote").
	Username string `json:"username,omitempty"`

	// The password to use when connecting to the Wiser for KNX controller.
	Password string `json:"password,omitempty"`

	// Path to a secret containing the password to use when connecting to the Wiser for KNX controller.
	PasswordFile string `json:"passwordFile,omitempty"`

	// The poll interval to use when polling the Wiser for KNX controller (default: 10 seconds).
	Poll time.Duration `json:"poll,omitempty"`

	// The list of exported objects on the Wiser for KNX controller.
	Devices []Device `json:"devices,omitempty"`
}

func (c Config) LoadPassword() (string, error) {
	if c.Password != "" {
		return c.Password, nil
	}
	p, err := os.ReadFile(c.PasswordFile)
	return strings.TrimSpace(string(p)), err
}

// Device config for an object on the Wiser for KNX controller.
// Note: only 1 of Address or Addresses should be specified.
type Device struct {
	// The device name (e.g. "inf/sc-01/lights/flex-room-01")
	Name string `json:"name"`
	// The address of the object on the Wiser for KNX controller (e.g. "1/1/1")
	Address string `json:"address,omitempty"`
	// Map of device component (light, override) to address (e.g. {"light": "1/1/1"})
	Addresses map[string]string `json:"addresses,omitempty"`
	// The metadata associated with the device.
	Metadata *traits.Metadata `json:"metadata,omitempty"`
}
