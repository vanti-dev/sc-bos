package config

import (
	"encoding/json"
	"fmt"

	"github.com/vanti-dev/sc-bos/pkg/driver"
)

const DefaultPort = 60001

type Root struct {
	driver.BaseConfig
	HTTP         *HTTP        `json:"http,omitempty"`
	MessagePorts MessagePorts `json:"messagePorts,omitempty"`
	Database     *Database    `json:"database,omitempty"`

	Devices []Device `json:"devices,omitempty"`
}

func ReadBytes(data []byte) (root Root, err error) {
	err = json.Unmarshal(data, &root)
	if err != nil {
		return
	}
	if root.MessagePorts.Bind == "" {
		root.MessagePorts.Bind = fmt.Sprintf(":%d", DefaultPort)
	}
	return
}

type Device struct {
	Name string // Smart Core name
	// NetworkDesc is the human specified name identifying an Axiom controller.
	NetworkDesc string
	// DeviceDesc is the human specified name identifying a card reader managed by NetworkDesc.
	DeviceDesc string
	// UDMITopicPrefix is used for telemetry and config when using UDMI.
	// Defaults to Name.
	UDMITopicPrefix string
}

type HTTP struct {
	BaseURL string `json:"baseUrl,omitempty"`
}

type Database struct {
	DSN          string `json:"dsn,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}
