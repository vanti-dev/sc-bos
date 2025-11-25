package config

import (
	"encoding/json"
	"strconv"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-bos/pkg/driver/opcua/conv"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

type Trait struct {
	Name     string           `json:"name,omitempty"`
	Kind     trait.Name       `json:"kind,omitempty"`
	Metadata *traits.Metadata `json:"metadata,omitempty"`
}

type RawTrait struct {
	Trait
	Raw json.RawMessage `json:"-"`
}

func (c *RawTrait) MarshalJSON() ([]byte, error) {
	return c.Raw, nil
}

func (c *RawTrait) UnmarshalJSON(buf []byte) error {
	if c == nil {
		*c = RawTrait{}
	}
	c.Raw = buf
	return json.Unmarshal(buf, &c.Trait)
}

// ValueSource configures a single Variable as the source of some trait value.
type ValueSource struct {
	NodeId string `json:"nodeId,omitempty"`
	Name   string `json:"name,omitempty"`
	// Description is a human-readable description of the source
	Description string `json:"description,omitempty"`
	// Optional. Used for converting simple units like kW -> W.
	// The value from the source will be multiplied by Scale when reading, and divided when writing.
	// For example if the trait is in watts and the device is in kW then Scale should be 1000 (aka kilo).
	Scale float64 `json:"scale,omitempty"`
	// Optional. Enum is a generic map to convert the OPC UA point value to something else.
	// For instance, converting the OCP UA value to an enum in a Smart Core trait which can be done by mapping the
	// OPC UA value as the key and the element from the generated <EnumName>_value field in the trait pb file.
	// The key needs to be an integer, it is defined as a string here for JSON marshaling.
	Enum map[string]string `json:"enum,omitempty"`
}

// GetValueFromIntKey get the value from the enum map given an integer OPC UA value
func (v ValueSource) GetValueFromIntKey(val any) any {
	if v.Enum != nil {
		i, err := conv.IntValue(val)
		if err == nil {
			if s, ok := v.Enum[strconv.Itoa(i)]; ok {
				return s
			}
		}
	}
	return val
}

// UdmiConfig is configured by a Device that wants to implement the UDMI trait.
type UdmiConfig struct {
	Trait
	// TopicPrefix is the prefix prepended to the topic in a gen.MqttMessage
	TopicPrefix string `json:"topicPrefix,omitempty"`
	// Points the points we want to send to the UDMI bus. point name -> point config (nodeId and optional enum)
	Points map[string]*ValueSource `json:"points"`
}

// MeterConfig is configured by a Device that wants to implement the Meter trait.
type MeterConfig struct {
	Trait
	Unit  string       `json:"unit,omitempty"`
	Usage *ValueSource `json:"usage,omitempty"`
}

type Door struct {
	Title  string       `json:"title,omitempty"`
	Deck   int          `json:"deck,omitempty"`
	Status *ValueSource `json:"status,omitempty"`
}

type LocationType string

const (
	// SingleFloor tells us that the OPC UA node represents a single floor.
	SingleFloor LocationType = "SingleFloor"
)

type Location struct {
	// Type tells us how to interpret the value source. It must be one of the defined LocationType values.
	// For example, the node in the value source could describe a single location,
	// with other nodes telling us about the next location.
	// Or it could contain an array that lists all the destinations. It is unclear what this needs to handle,
	// so it needs to be flexible enough and extensible to handle future integrations.
	Type   LocationType `json:"type,omitempty"`
	Source ValueSource  `json:"source,omitempty"`
}

// TransportConfig is configured by a Device that wants to implement the Transport trait.
type TransportConfig struct {
	Trait
	ActualPosition  *ValueSource `json:"actualPosition,omitempty"`
	Doors           []*Door      `json:"doors,omitempty"`
	Load            *ValueSource `json:"load,omitempty"`
	LoadUnit        string       `json:"loadUnit,omitempty"`
	MaxLoad         int32        `json:"maxLoad,omitempty"`
	MovingDirection *ValueSource `json:"movingDirection,omitempty"`
	// The OPC UA node(s) which tells us the where the transport is going to stop at next.
	// If the OPC UA server has more than one point which tells us about the next destinations,
	// this array should be ordered so that it matches the order of the physical transport stops.
	// i.e [0] = first stop, [1] = second stop, etc.
	NextDestinations []*Location  `json:"nextDestinations,omitempty"`
	OperatingMode    *ValueSource `json:"operatingMode,omitempty"`
	Speed            *ValueSource `json:"speed,omitempty"`
	SpeedUnit        string       `json:"speedUnit,omitempty"`
}

type ElectricConfig struct {
	Trait
	Demand *ElectricDemandConfig `json:"demand,omitempty"`
}

type ElectricDemandConfig struct {
	*ElectricPhaseConfig                        // single phase
	Phases               [3]ElectricPhaseConfig `json:"phases,omitempty"`
}

type ElectricPhaseConfig struct {
	Current *ValueSource `json:"current,omitempty"`
	Voltage *ValueSource `json:"voltage,omitempty"`
	Rating  *ValueSource `json:"rating,omitempty"`

	PowerFactor   *ValueSource `json:"powerFactor,omitempty"`
	RealPower     *ValueSource `json:"realPower,omitempty"`
	ApparentPower *ValueSource `json:"apparentPower,omitempty"`
	ReactivePower *ValueSource `json:"reactivePower,omitempty"`
}
