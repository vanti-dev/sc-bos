package config

import (
	"os"
	"strings"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/driver"
)

type Root struct {
	driver.BaseConfig

	// Smart core metadata associated with this device.
	Metadata *traits.Metadata `json:"metadata,omitempty"`

	IpAddress    string `json:"ipAddress"`
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`

	// to support UDMI/MQTT automation
	UDMITopicPrefix string `json:"udmiTopicPrefix,omitempty"`
}

func (c *Root) LoadPassword() (string, error) {
	if c.Password != "" {
		return c.Password, nil
	}
	bs, err := os.ReadFile(c.PasswordFile)
	return strings.TrimSpace(string(bs)), err
}
