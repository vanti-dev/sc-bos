package xovis

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

func DefaultConfig() DriverConfig {
	return DriverConfig{
		PasswordFile: "/run/secrets/xovis-password",
	}
}

func ParseConfig(raw []byte) (DriverConfig, error) {
	parsed := DefaultConfig()
	err := json.Unmarshal(raw, &parsed)
	return parsed, err
}

type DriverConfig struct {
	driver.BaseConfig
	MultiSensor  bool            `json:"multiSensor"`
	Host         string          `json:"host"`
	Username     string          `json:"username"`
	Password     string          `json:"password,omitempty"`
	PasswordFile string          `json:"passwordFile,omitempty"`
	DataPush     *DataPushConfig `json:"dataPush"`
	Devices      []DeviceConfig  `json:"devices,omitempty"`
}

func (c DriverConfig) LoadPassword() (string, error) {
	if c.Password != "" {
		return c.Password, nil
	}
	bs, err := os.ReadFile(c.PasswordFile)
	return strings.TrimSpace(string(bs)), err
}

type DeviceConfig struct {
	Name       string           `json:"name"`
	Occupancy  *LogicConfig     `json:"occupancy"`  // an Occupancy logic
	EnterLeave *LogicConfig     `json:"enterLeave"` // an In/Out logic
	Metadata   *traits.Metadata `json:"metadata,omitempty"`

	// to support UDMI/MQTT automation
	UDMITopicPrefix string `json:"udmiTopicPrefix,omitempty"`
}

type LogicConfig struct {
	ID int `json:"id"`
}

type DataPushConfig struct {
	HTTPPort    int    `json:"httpPort,omitempty"`
	WebhookPath string `json:"webhookPath"`
}
