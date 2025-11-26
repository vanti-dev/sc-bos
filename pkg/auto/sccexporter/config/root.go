package config

import (
	"encoding/json"
	"time"

	"github.com/smart-core-os/sc-bos/pkg/auto"
	"github.com/smart-core-os/sc-bos/pkg/util/jsontypes"
)

type Mqtt struct {
	Agent            string              `json:"agent"`
	Host             string              `json:"host"`
	Topic            string              `json:"topic"` // the topic to publish to
	ClientId         string              `json:"clientId"`
	ClientKey        string              `json:"clientKey"`
	ClientCert       string              `json:"clientCert"`
	CaCert           string              `json:"caCert"`
	PublishTimeout   *jsontypes.Duration `json:"publishTimeout,omitempty,omitzero"`   // timeout for publishing to MQTT, default to 5s
	Qos              *int                `json:"qos,omitempty"`                       // MQTT qos, default to 1
	SendInterval     *jsontypes.Schedule `json:"sendInterval,omitempty,omitzero"`     // time between sends, default to 15m
	MetadataInterval *int                `json:"metadataInterval,omitempty,omitzero"` // how often to include metadata (every N data sends), default to 100
}

type Root struct {
	auto.Config

	// A list of all traits we want to send data for
	Traits []string `json:"traits"`
	Mqtt   Mqtt     `json:"mqtt"`
	// FetchTimeout is the maximum time to wait for a single device's trait data fetch
	// If a device takes longer than this, the fetch is cancelled and the device is skipped
	// Default is 5 seconds
	FetchTimeout *jsontypes.Duration `json:"fetchTimeout,omitempty,omitzero"`
}

func ParseConfig(data []byte) (Root, error) {
	root := Root{}

	if err := json.Unmarshal(data, &root); err != nil {
		return Root{}, err
	}

	if root.Mqtt.SendInterval == nil {
		root.Mqtt.SendInterval = jsontypes.MustParseSchedule("*/15 * * * *")
	}
	if root.Mqtt.PublishTimeout == nil || root.Mqtt.PublishTimeout.Duration == 0 {
		root.Mqtt.PublishTimeout = &jsontypes.Duration{Duration: 5 * time.Second}
	}
	if root.Mqtt.Qos == nil {
		q := 1
		root.Mqtt.Qos = &q
	}
	if root.Mqtt.MetadataInterval == nil {
		interval := 100
		root.Mqtt.MetadataInterval = &interval
	}
	if root.FetchTimeout == nil || root.FetchTimeout.Duration == 0 {
		root.FetchTimeout = &jsontypes.Duration{Duration: 5 * time.Second}
	}

	return root, nil
}
