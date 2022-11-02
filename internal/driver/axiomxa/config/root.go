package config

import (
	"encoding/json"
)

type Root struct {
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

type MessagePorts struct {
	LocalAddress string `json:"localAddress,omitempty"`
	PathPrefix   string `json:"pathPrefix,omitempty"`
}

type Database struct {
	DSN          string `json:"dsn,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}
