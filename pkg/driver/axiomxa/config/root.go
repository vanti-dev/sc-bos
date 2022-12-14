package config

import (
	"encoding/json"
	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type Root struct {
	driver.BaseConfig
	HTTP         *HTTP         `json:"http,omitempty"`
	MessagePorts *MessagePorts `json:"messagePorts,omitempty"`
	Database     *Database     `json:"database,omitempty"`
}

func ReadBytes(data []byte) (root Root, err error) {
	err = json.Unmarshal(data, &root)
	return
}

type HTTP struct {
	BaseURL string `json:"baseUrl,omitempty"`
}

type Database struct {
	DSN          string `json:"dsn,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}
