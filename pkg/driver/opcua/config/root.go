package config

import (
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

const (
	PointsEventTopicSuffix = "/event/pointset"
)

// OpcUaConfig config related to communicating with the OPC UA server.
type OpcUaConfig struct {
	// Endpoint is the OPC UA server endpoint.
	Endpoint string `json:"endpoint,omitempty"`
	// SubscriptionInterval for OPC UA subscription, defaults to 5s if not set.
	SubscriptionInterval *jsontypes.Duration `json:"subscriptionInterval,omitempty"`
}

// Variable is an OPC UA VariableNode, which is essentially a data point which we can read/write to (with permission).
type Variable struct {
	// NodeId identifies the VariableNode in the OPC UA server.
	NodeId string `json:"nodeId,omitempty"`
}

// Device represents a smart core device.
type Device struct {
	// Name the Smart Core device name
	Name string `json:"name,omitempty"`
	// Meta the Smart Core device metadata
	Meta *traits.Metadata `json:"meta,omitempty"`
	// Variables a list of OPC variables the device has
	Variables []Variable `json:"variables,omitempty"`
	// Traits a map Smart Core traits the device implements
	Traits []RawTrait `json:"traits,omitempty"`
}

type Timing struct {
	Timeout      jsontypes.Duration `json:"timeout,omitempty"`
	BackoffStart jsontypes.Duration `json:"backoffStart,omitempty"`
	BackoffMax   jsontypes.Duration `json:"backoffMax,omitempty"`
}

type Root struct {
	driver.BaseConfig

	Meta        *traits.Metadata `json:"meta,omitempty"`
	OpcUaConfig OpcUaConfig      `json:"opcUaConfig,omitempty"`
	Devices     []Device         `json:"devices,omitempty"`
	Timing      Timing           `json:"Timing,omitempty"`
}

func SetDefaults(cfg *Root) {
	if cfg.OpcUaConfig.SubscriptionInterval == nil {
		cfg.OpcUaConfig.SubscriptionInterval = &jsontypes.Duration{Duration: 5 * time.Second}
	}
	if cfg.Timing.Timeout.Duration == 0 {
		cfg.Timing.Timeout = jsontypes.Duration{Duration: 10 * time.Second}
	}
	if cfg.Timing.BackoffStart.Duration == 0 {
		cfg.Timing.BackoffStart = jsontypes.Duration{Duration: 2 * time.Second}
	}
	if cfg.Timing.BackoffMax.Duration == 0 {
		cfg.Timing.BackoffMax = jsontypes.Duration{Duration: 30 * time.Second}
	}
	if cfg.Timing.BackoffMax.Duration < cfg.Timing.BackoffStart.Duration {
		cfg.Timing.BackoffMax = cfg.Timing.BackoffStart
	}
}
