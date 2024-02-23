package config

import (
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTBroker struct {
	Host         string `json:"host,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}

func (b MQTTBroker) ClientOptions() (*mqtt.ClientOptions, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(b.Host)
	password := b.Password
	if password == "" {
		if b.PasswordFile != "" {
			var passFileBody []byte
			passFileBody, err := os.ReadFile(b.PasswordFile)
			if err != nil {
				return nil, err
			} else {
				password = strings.TrimSpace(string(passFileBody))
			}
		}
	}

	if password != "" { // allow connection without password if no password provided
		opts.SetPassword(password)
	}

	opts.SetUsername(b.Username)
	opts.SetOrderMatters(false)
	return opts, nil
}
