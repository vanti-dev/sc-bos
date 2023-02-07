package config

import (
	"github.com/vanti-dev/sc-bos/pkg/auto"
)

type Root struct {
	auto.Config

	// Broker configures an MQTT broker to export data to, and subscribe to topics on.
	Broker *MQTTBroker `json:"broker,omitempty"`

	// the names to use for rpc requests to UdmiService
	Sources []string `json:"sources,omitempty"`
}
