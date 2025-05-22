package config

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/internal/util/pgxutil"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	auto.Config
	Source  *Source  `json:"source,omitempty"`
	Storage *Storage `json:"storage,omitempty"`
}

type Source struct {
	Name  string     `json:"name,omitempty"`
	Trait trait.Name `json:"trait,omitempty"`
	// ReadMask instructs the history service to only read the specified fields
	ReadMask *FieldMask `json:"readMask,omitempty"`
	// Sample controls whether periodic downsampling is applied to the data source.
	// It is intended for high-frequency sources where storing every data point is excessive.
	// When Sample is set, the source is polled using a Cron schedule, and only values
	// that differ from the previous sample are recorded — reducing storage by skipping redundant data.
	// If Sample is not set (the default) or disabled, the source will be treated as event-driven:
	// all changes will be recorded at their native frequency, with no sampling applied.
	Sample *Sampling `json:"sample,omitempty"`
}

func (s Source) SourceName() string {
	return fmt.Sprintf("%s[%s]", s.Name, s.Trait)
}

type FieldMask fieldmaskpb.FieldMask

func (f *FieldMask) PB() *fieldmaskpb.FieldMask {
	if f == nil {
		return nil
	}
	return (*fieldmaskpb.FieldMask)(f)
}

func (f *FieldMask) UnmarshalJSON(bytes []byte) error {
	return protojson.Unmarshal(bytes, (*fieldmaskpb.FieldMask)(f))
}

func (f *FieldMask) MarshalJSON() ([]byte, error) {
	return protojson.Marshal((*fieldmaskpb.FieldMask)(f))
}

type Sampling struct {
	// Schedule is the time schedule between samples being recorded.
	// The default is every 1 hour at minute 0 < 0 * * * * >.
	Schedule *jsontypes.Schedule `json:"schedule,omitempty"`
	// Disabled is the flag used to disable sampling.
	Disabled bool `json:"disabled,omitempty"`
}

type Storage struct {
	Type string `json:"type,omitempty"`
	pgxutil.ConnectConfig
	Name string `json:"name,omitempty"`
	// TTL is the time-to-live for records. Zero-value (not-specified) means "forever".
	TTL *TTL `json:"ttl,omitempty"`
}

type TTL struct {
	MaxAge   jsontypes.Duration `json:"maxAge,omitempty"`
	MaxCount int64              `json:"maxCount,omitempty"`
}

func ParseConfig(data []byte) (Root, error) {
	var cfg Root
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.Source.Sample != nil && cfg.Source.Sample.Schedule == nil {
		cfg.Source.Sample.Schedule = jsontypes.MustParseSchedule("0 * * * *") // every hour at minute 0
	}
	return cfg, nil
}
