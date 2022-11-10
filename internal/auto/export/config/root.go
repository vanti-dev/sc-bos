package config

import (
	"encoding/json"
	"github.com/vanti-dev/bsp-ew/internal/auto"
)

type Root struct {
	auto.Config

	// Broker configures an MQTT broker to export data to.
	Broker *MQTTBroker `json:"broker,omitempty"`

	Sources []RawSource `json:"sources,omitempty"`
}

type Source struct {
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	TopicPrefix string `json:"topicPrefix,omitempty"`
}

type RawSource struct {
	Source
	Raw json.RawMessage `json:"-"`
}

func (r *RawSource) UnmarshalJSON(buf []byte) error {
	_ = json.Unmarshal(buf, &r.Raw)
	return json.Unmarshal(buf, &r.Source)
}
