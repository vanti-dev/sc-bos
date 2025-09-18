package config

import (
	"encoding/json"

	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Root struct {
	auto.Config
	Devices []*gen.Device_Query_Condition `json:"devices"`
	Check   *gen.HealthCheck              `json:"check"`
	Source  Source                        `json:"source"`
}

// Source configures which property of a device is checked by a health check.
type Source struct {
	Trait    trait.Name `json:"trait"`
	Resource Resource   `json:"resource,omitempty"`
	Value    Value      `json:"value,omitempty"`
}

// Resource is the name of a trait resource.
// For example "OnOff" or "Brightness", for which there would be GetOnOff and/or PullOnOff rpc methods in the trait API.
// When empty, the first declared resource in the trait is used.
type Resource string

func (r Resource) String() string {
	return string(r)
}

// Value is a dot-separated path to a field in a trait resource specified in Source.
type Value string

func (v Value) String() string {
	return string(v)
}

func (v Value) ToFieldMask() *fieldmaskpb.FieldMask {
	if v == "" {
		return nil
	}
	// convert camelCase to snake_case for protobuf field names
	ps := strcase.SnakeCase(string(v)) // ascii only
	return &fieldmaskpb.FieldMask{Paths: []string{ps}}
}

func Read(data []byte) (Root, error) {
	var cfg Root
	err := json.Unmarshal(data, &cfg)
	return cfg, err
}
